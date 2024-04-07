package request_photo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/jmoiron/sqlx"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
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

type docsDoc struct {
	File string `json:"file"`
}

// UploadPhotoAsFile загружает фотографию как документ
func (r *Repo) UploadPhotoAsFile(ctx context.Context, VK *api.VK, doc object.DocsDoc, vkID int) error {
	resp, err := http.Get(doc.URL)
	if err != nil {
		return fmt.Errorf("failed to retrieve document URL: %w", err)
	}
	defer resp.Body.Close()

	upload, err := VK.DocsGetMessagesUploadServer(api.Params{
		"type":    "doc",
		"peer_id": vkID,
	})
	if err != nil {
		return fmt.Errorf("failed to get upload server URL: %w", err)
	}

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read document content: %w", err)
	}
	fileBody := bytes.NewReader(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", doc.Title)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, fileBody)
	if err != nil {
		return fmt.Errorf("failed to copy file content to form file: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", upload.UploadURL, body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer response.Body.Close()

	docs := &docsDoc{}
	err = json.NewDecoder(response.Body).Decode(docs)
	if err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	log.Println(docs.File)

	savedDoc, err := VK.DocsSave(api.Params{
		"file":  docs.File,
		"title": doc.Title,
	})
	if err != nil {
		return fmt.Errorf("failed to save document: %w", err)
	}

	attachment := "doc" + strconv.Itoa(savedDoc.Doc.OwnerID) + "_" + strconv.Itoa(savedDoc.Doc.ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO request_photo(attachment, user_id) VALUES ($1, $2)", attachment, vkID)
	if err != nil {
		return fmt.Errorf("failed to insert into database: %w", err)
	}

	return nil
}

type photosPhoto struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

// UploadPhoto загружает фотографию
func (r *Repo) UploadPhoto(ctx context.Context, VK *api.VK, photo object.PhotosPhoto, vkID int) error {
	resp, err := http.Get(photo.Sizes[4].URL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	uploadServer, err := VK.PhotosGetMessagesUploadServer(api.Params{
		"peer_id": vkID,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	fileBody := bytes.NewReader(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = io.Copy(part, fileBody)
	if err != nil {
		log.Println(err)
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", uploadServer.UploadURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)

	uploadResult := &photosPhoto{}
	err = json.NewDecoder(response.Body).Decode(uploadResult)
	if err != nil {
		log.Fatal(err)
	}

	savedPhoto, err := VK.PhotosSaveMessagesPhoto(api.Params{
		"photo":  uploadResult.Photo,
		"server": uploadResult.Server,
		"hash":   uploadResult.Hash,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	attachment := "photo" + strconv.Itoa(savedPhoto[0].OwnerID) + "_" + strconv.Itoa(savedPhoto[0].ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO request_photo(attachment, user_id) VALUES ($1, $2)", attachment, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeletePhotoRequest удаляет заявку на добавление фотографии в альбом
func (r *Repo) DeletePhotoRequest(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM request_photo WHERE id = (SELECT id FROM request_photo WHERE user_id = $1 ORDER BY id DESC LIMIT 1)", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	return nil
}

// GetPhotoLastID возвращает ID последней заявки на добавление фото в альбом
func (r *Repo) GetPhotoLastID(ctx context.Context, vkID int) (int, error) {
	var photo ds.RequestPhoto
	err := r.db.GetContext(ctx, &photo, "SELECT id FROM request_photo WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return 0, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return photo.ID, nil
}

// UpdateYear добавляет год события для фотографии
func (r *Repo) UpdateYear(ctx context.Context, photoID int, year int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE request_photo SET year = $1 WHERE id = $2", year, photoID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateStudyProgram добавляет программу обучения для фотографии
func (r *Repo) UpdateStudyProgram(ctx context.Context, photoID int, studyProgram string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE request_photo SET study_program = $1 WHERE id = $2", studyProgram, photoID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// GetEventNames возвращает список названий событий для фотографии
func (r *Repo) GetEventNames() (string, error) {
	var output string

	rows, err := r.db.Query("SELECT CONCAT(id, ') ', name) AS formatted_string FROM events")
	if err != nil {
		return "", fmt.Errorf("[db.Query]: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var formattedString string
		err := rows.Scan(&formattedString)
		if err != nil {
			return "", fmt.Errorf("[db.Scan]: %w", err)
		}
		output += formattedString + "\n"
	}

	return output, nil
}

// GetEventMaxID возвращает максимальное ID из событий для фотографий
func (r *Repo) GetEventMaxID() (int, error) {
	var maxID int
	err := r.db.Get(&maxID, "SELECT MAX(id) FROM events")
	if err != nil {
		return 0, fmt.Errorf("[db.Get]: %w", err)
	}

	return maxID, nil
}

// UpdateEvent добавляет событие для фотографии
func (r *Repo) UpdateEvent(ctx context.Context, photoID, eventNumber int) error {
	var name string
	err := r.db.Get(&name, "SELECT name FROM events WHERE id = $1", eventNumber)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE request_photo SET (event, is_event_new) = ($1, $2) WHERE id = $3", name, false, photoID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateUserEvent добавляет пользовательское название события
func (r *Repo) UpdateUserEvent(ctx context.Context, photoID int, category string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE request_photo SET (event, is_event_new) = ($1, $2) WHERE id = $3", category, true, photoID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateDescription добавляет описание фотографии
func (r *Repo) UpdateDescription(ctx context.Context, photoID int, description string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE request_photo SET description = $1 WHERE id = $2", description, photoID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// CheckParams возвращает все параметры фотографии на загрузку в альбом
func (r *Repo) CheckParams(ctx context.Context, photoID int) (string, string, error) {
	sqlQuery := `
	SELECT 
    	'1) Год события: ' || COALESCE(CAST(year AS VARCHAR), 'Не указано') AS year,
    	'2) Программа обучения: ' || COALESCE(study_program, 'Не указано') AS studyProgram,
    	'3) Название события: ' || COALESCE(event, 'Не указано') AS event,
    	'4) Описание: ' || COALESCE(description, 'Не указано') AS description,
    	'5) Отмеченные люди: ' || COALESCE(array_to_string(marked_people, ', '), 'Не указано') AS markedPeople,
    	attachment
	FROM request_photo
	WHERE id = $1;`

	var (
		year         string
		studyProgram string
		event        string
		description  string
		markedPeople string
		attachment   string
	)

	err := r.db.QueryRow(sqlQuery, photoID).Scan(&year, &studyProgram, &event, &description, &markedPeople, &attachment)
	if err != nil {
		return "", "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	output := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", year, studyProgram, event, description, markedPeople)

	return output, attachment, nil
}
