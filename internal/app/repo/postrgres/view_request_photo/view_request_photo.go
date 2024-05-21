package view_request_photo

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/jmoiron/sqlx"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// Repo инстанс репо для работы с фотографиями пользователя
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с фотографиями пользователя
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// GetRequestPhoto возвращает заявку на загрузку фото
func (r *Repo) GetRequestPhoto(vkID int) (string, string, int, int, error) {
	var pointer int
	err := r.db.Get(&pointer, "SELECT pointer FROM view_request_photo WHERE user_id = $1", vkID)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("[db.Get]: %w", err)
	}

	sqlQuery := `
	SELECT 
    	'3) Год события: ' || COALESCE(CAST(year AS VARCHAR), 'Не указано') AS year,
    	'4) Программа обучения: ' || COALESCE(study_program, 'Не указано') AS studyProgram,
    	'5) Название события: ' || COALESCE(event, 'Не указано') AS event,
    	'6) Описание: ' || COALESCE(description, 'Не указано') AS description,
    	'7) Отмеченные люди: ' || COALESCE(array_to_string(marked_people, ', '), 'Не указано') AS markedPeople,
    	'8) Преподаватели: ' || COALESCE(array_to_string(teachers, ', '), 'Не указано') AS teachers,
    	status,
    	attachment
	FROM request_photo
	ORDER BY id
	offset $1 limit 1;`

	var (
		year         string
		studyProgram string
		event        string
		description  string
		markedPeople string
		teachers     string
		status       int
		attachment   string
	)
	err = r.db.QueryRow(sqlQuery, pointer).Scan(&year, &studyProgram, &event, &description, &markedPeople, &teachers, &status, &attachment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", 0, 0, nil
		}
		return "", "", 0, 0, fmt.Errorf("[db.QueryRow]: %w", err)
	}

	var count int
	err = r.db.Get(&count, "SELECT count(*) FROM request_photo")
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("[db.Get]: %w", err)
	}

	countString := strconv.Itoa(count)
	pointerString := strconv.Itoa(pointer + 1)

	var statusString string

	switch status {
	case 1:
		statusString = "Ожидает рассмотрения"
	case 2:
		statusString = "На рассмотрении"
	case 3:
		statusString = "Отклонена"
	case 4:
		statusString = "Одобрена"
	}
	message := fmt.Sprintf("1) Заявка %s/%s\n2) Статус данной заявки: %s\n%s\n%s\n%s\n%s\n%s\n%s",
		pointerString, countString, statusString, year, studyProgram, event, description, markedPeople, teachers)

	return message, attachment, pointer, count, nil
}

// CreatePersonalAccountPhoto создает запись в личном кабинете пользователя
func (r *Repo) CreatePersonalAccountPhoto(vkID int) error {
	_, err := r.db.Exec("INSERT INTO view_request_photo (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING", vkID)
	if err != nil {
		return fmt.Errorf("[db.Exec]: %w", err)
	}

	return nil
}

// ChangePointer меняет указатель на заявку в личном кабинете
func (r *Repo) ChangePointer(vkID int, flag bool) error {
	var pointer int

	err := r.db.Get(&pointer, "SELECT pointer FROM view_request_photo WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	if !flag {
		pointer--
	} else {
		pointer++
	}

	if pointer >= 0 {
		_, err = r.db.Exec("UPDATE view_request_photo SET pointer = $1 WHERE user_id = $2", pointer, vkID)
		if err != nil {
			return fmt.Errorf("[db.Exec]: %w", err)
		}
	}

	return nil
}

// DeletePointer удаляет указатель на заявку в личном кабинете
func (r *Repo) DeletePointer(vkID int) error {
	_, err := r.db.Exec("UPDATE view_request_photo SET pointer = 0 WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Exec]: %w", err)
	}
	return nil
}

// ApprovePhoto одобряет фото
func (r *Repo) ApprovePhoto(vkID int, vkUser *api.VK, groupID int) (string, error) {
	var pointer int
	err := r.db.Get(&pointer, "SELECT pointer FROM view_request_photo WHERE user_id = $1", vkID)
	if err != nil {
		return "", fmt.Errorf("[db.Get]: %w", err)
	}

	var req ds.RequestPhoto
	err = r.db.QueryRow("SELECT year, study_program, event, description, marked_people, teachers, url FROM request_photo WHERE id = (SELECT id FROM request_photo ORDER BY id OFFSET $1 LIMIT 1)", pointer).
		Scan(&req.Year, &req.StudyProgram, &req.Event, &req.Description, &req.MarkedPeople, &req.Teachers, &req.URL)
	if err != nil {
		return "", fmt.Errorf("[db.QueryRow]: %w", err)
	}

	var comment string
	if req.Year == nil {
		comment = "Не указан год события"
		return comment, err
	}
	if req.StudyProgram == nil {
		comment = "Не указана программа обучения"
		return comment, err
	}
	if req.Event == nil {
		comment = "Не указано название события"
		return comment, err
	}

	var albumID *int
	err = r.db.Get(&albumID, "SELECT vk_id FROM student_albums WHERE year = $1 AND study_program = $2 AND event = $3", req.Year, req.StudyProgram, req.Event)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("[db.Get]: %w", err)
	}

	if albumID == nil {
		albumID, err = r.CreateStudentAlbum(vkUser, groupID, *req.Year, *req.StudyProgram, *req.Event)
		if err != nil {
			return "", fmt.Errorf("error create student album: %w", err)
		}
	}

	markPeople := strings.Join(req.MarkedPeople, ", ")
	var description string
	if req.Description != nil {
		description += fmt.Sprintf("Описание:\n%s\n", *req.Description)
	}
	if markPeople != "" {
		description += fmt.Sprintf("Отмеченные люди слева направо:\n%s", markPeople)
	}

	err = r.AddPhotoToAlbum(vkUser, *albumID, groupID, req.URL, description)
	if err != nil {
		return "", err
	}

	for _, teacher := range req.Teachers {
		albumID = nil
		err = r.db.Get(&albumID, "SELECT vk_id FROM teacher_albums WHERE name = $1", teacher)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("[db.Get]: %w", err)
		}

		if albumID == nil {
			albumID, err = r.CreateTeacherAlbum(vkUser, groupID, teacher)
			if err != nil {
				return "", fmt.Errorf("create student album: %w", err)
			}
		}

		err = r.AddPhotoToAlbum(vkUser, *albumID, groupID, req.URL, description)
		if err != nil {
			return "", err
		}
	}

	_, err = r.db.Exec(`UPDATE request_photo SET status = 4 WHERE id = (SELECT id FROM request_photo ORDER BY id OFFSET $1 LIMIT 1)`, pointer)
	if err != nil {
		return "", fmt.Errorf("[db.Exec]: %w", err)
	}

	return "", nil
}

// RejectPhoto отклоняет фото
func (r *Repo) RejectPhoto(vkID int) error {
	var pointer int
	err := r.db.Get(&pointer, "SELECT pointer FROM view_request_photo WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	_, err = r.db.Exec(`UPDATE request_photo SET status = 3 WHERE id = (SELECT id FROM request_photo ORDER BY id OFFSET $1 LIMIT 1)`, pointer)
	if err != nil {
		return fmt.Errorf("[db.Exec]: %w", err)
	}

	return nil
}

// CreateStudentAlbum создает альбом студентов
func (r *Repo) CreateStudentAlbum(vkUser *api.VK, groupID int, year, studyProgram, event string) (*int, error) {
	resultName := fmt.Sprintf("%s // %s // %s", year, studyProgram, event)

	album, err := vkUser.PhotosCreateAlbum(api.Params{
		"title":    resultName,
		"group_id": groupID,
	})

	albumURL := fmt.Sprintf("https://vk.com/album-%d_%d", groupID, album.ID)

	_, err = r.db.Exec("INSERT INTO student_albums (year, study_program, event, url, vk_id) VALUES ($1, $2, $3, $4, $5)", year, studyProgram, event, albumURL, album.ID)
	if err != nil {
		return nil, fmt.Errorf("[db.Exec]: %w", err)
	}

	return &album.ID, nil
}

// CreateTeacherAlbum создает альбом преподавателя
func (r *Repo) CreateTeacherAlbum(vkUser *api.VK, groupID int, nameAlbum string) (*int, error) {
	album, err := vkUser.PhotosCreateAlbum(api.Params{
		"title":    nameAlbum,
		"group_id": groupID,
	})

	albumURL := fmt.Sprintf("https://vk.com/album-%d_%d", groupID, album.ID)

	_, err = r.db.Exec("INSERT INTO teacher_albums (name, url, vk_id) VALUES ($1, $2, $3)", nameAlbum, albumURL, album.ID)
	if err != nil {
		return nil, fmt.Errorf("[db.Exec]: %w", err)
	}

	return &album.ID, nil
}

type UploadResponse struct {
	Server     int    `json:"server"`
	PhotosList string `json:"photos_list"`
	AID        int    `json:"aid"`
	Hash       string `json:"hash"`
}

func (r *Repo) AddPhotoToAlbum(vkUser *api.VK, albumID int, groupID int, photoURL string, description string) error {
	// Получаем URL сервера для загрузки фотографии
	uploadServer, err := vkUser.PhotosGetUploadServer(api.Params{"album_id": albumID, "group_id": groupID})
	if err != nil {
		return fmt.Errorf("ошибка при получении сервера загрузки: %w", err)
	}

	// Получаем содержимое фотографии по URL
	resp, err := http.Get(photoURL)
	if err != nil {
		return fmt.Errorf("ошибка при загрузке фотографии: %w", err)
	}
	defer resp.Body.Close()

	// Создаем тело запроса
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файл в тело запроса
	part, err := writer.CreateFormFile("file1", "photo.jpg")
	if err != nil {
		return fmt.Errorf("ошибка при создании формы для файла: %w", err)
	}
	_, err = io.Copy(part, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка при записи фотографии в форму: %w", err)
	}
	writer.Close()

	// Отправляем POST-запрос на сервер загрузки
	resp, err = http.Post(uploadServer.UploadURL, writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	// Декодируем ответ
	var uploadResponse UploadResponse
	err = json.NewDecoder(resp.Body).Decode(&uploadResponse)
	if err != nil {
		return fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	// Сохраняем фотографию
	if description != "" {
		_, err = vkUser.PhotosSave(api.Params{
			"album_id":    albumID,
			"group_id":    groupID,
			"server":      uploadResponse.Server,
			"photos_list": uploadResponse.PhotosList,
			"hash":        uploadResponse.Hash,
			"caption":     description,
		})
	} else {
		_, err = vkUser.PhotosSave(api.Params{
			"album_id":    albumID,
			"group_id":    groupID,
			"server":      uploadResponse.Server,
			"photos_list": uploadResponse.PhotosList,
			"hash":        uploadResponse.Hash,
		})
	}

	if err != nil {
		return fmt.Errorf("ошибка при сохранении фотографии: %w", err)
	}

	return nil
}
