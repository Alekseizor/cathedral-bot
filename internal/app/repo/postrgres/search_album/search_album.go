package search_album

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

// Repo инстанс репо для работы с параметрами для поиска альбома
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с параметрами для поиска альбома
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// CreateSearchAlbum создает запись для поиска альбома
func (r *Repo) CreateSearchAlbum(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO search_album (user_id) VALUES ($1)", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeleteSearchAlbum удаляет запись для поиска альбома
func (r *Repo) DeleteSearchAlbum(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateYear добавляет год события для поиска альбома
func (r *Repo) UpdateYear(ctx context.Context, vkID int, year int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET year = $1 WHERE user_id = $2", year, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateStudyProgram добавляет программу обучения для поиска альбома
func (r *Repo) UpdateStudyProgram(ctx context.Context, vkID int, studyProgram string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET study_program = $1 WHERE user_id = $2", studyProgram, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateEvent добавляет название события для поиска альбома
func (r *Repo) UpdateEvent(ctx context.Context, vkID int, eventNumber int) error {
	var event string
	err := r.db.Get(&event, "SELECT name FROM events ORDER BY name OFFSET $1 LIMIT 1", eventNumber)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE search_album SET event = $1 WHERE user_id = $2", event, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeleteYear удаляет год события для поиска альбома
func (r *Repo) DeleteYear(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET year = $1 WHERE user_id = $2", nil, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeleteStudyProgram удаляет программу обучения для поиска альбома
func (r *Repo) DeleteStudyProgram(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET study_program = $1 WHERE user_id = $2", nil, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeleteEvent удаляет название события для поиска альбома
func (r *Repo) DeleteEvent(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET event = $1 WHERE user_id = $2", nil, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// CountAlbums возвращает количество найденных альбомов по году
func (r *Repo) CountAlbums(ctx context.Context, vkID int) (int, error) {
	var searchAlbum ds.SearchAlbum

	err := r.db.GetContext(ctx, &searchAlbum, "SELECT * FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return 0, fmt.Errorf("[db.GetContext]: %w", err)
	}

	var count int

	query := "SELECT COUNT(*) FROM student_albums"
	var args []interface{}

	var argIndex = 1

	if searchAlbum.Year != nil {
		query += " WHERE year = $" + strconv.Itoa(argIndex)
		args = append(args, *searchAlbum.Year)
		argIndex++
	}

	if searchAlbum.StudyProgram != nil {
		if len(args) == 0 {
			query += " WHERE"
		} else {
			query += " AND"
		}
		query += " study_program = $" + strconv.Itoa(argIndex)
		args = append(args, *searchAlbum.StudyProgram)
		argIndex++
	}

	if searchAlbum.Event != nil {
		if len(args) == 0 {
			query += " WHERE"
		} else {
			query += " AND"
		}
		query += " event = $" + strconv.Itoa(argIndex)
		args = append(args, *searchAlbum.Event)
		argIndex++
	}

	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return count, nil
}

// ShowList возвращает список найденных альбомов
func (r *Repo) ShowList(ctx context.Context, vkID int) (string, int, int, error) {
	var searchAlbum ds.SearchAlbum

	err := r.db.GetContext(ctx, &searchAlbum, "SELECT * FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return "", 0, 0, fmt.Errorf("[db.GetContext]: %w", err)
	}

	query := "SELECT * FROM student_albums"
	var args []interface{}

	var argIndex = 1

	if searchAlbum.Year != nil {
		query += " WHERE year = $" + strconv.Itoa(argIndex)
		args = append(args, *searchAlbum.Year)
		argIndex++
	}

	if searchAlbum.StudyProgram != nil {
		if len(args) == 0 {
			query += " WHERE"
		} else {
			query += " AND"
		}
		query += " study_program = $" + strconv.Itoa(argIndex)
		args = append(args, *searchAlbum.StudyProgram)
		argIndex++
	}

	if searchAlbum.Event != nil {
		if len(args) == 0 {
			query += " WHERE"
		} else {
			query += " AND"
		}
		query += " event = $" + strconv.Itoa(argIndex)
		args = append(args, *searchAlbum.Event)
		argIndex++
	}

	var albums []ds.StudentAlbum
	err = r.db.SelectContext(ctx, &albums, query, args...)
	if err != nil {
		return "", 0, 0, fmt.Errorf("[db.SelectContext]: %w", err)
	}
	count := len(albums)

	query += " ORDER BY year DESC, study_program, event" + " offset $" + strconv.Itoa(argIndex) + " limit 10"
	args = append(args, *searchAlbum.Pointer)

	err = r.db.SelectContext(ctx, &albums, query, args...)
	if err != nil {
		return "", 0, 0, fmt.Errorf("[db.SelectContext]: %w", err)
	}

	var result string

	if len(albums) == 1 {
		yearStr := strconv.Itoa(albums[0].Year)
		result += yearStr + " // " + albums[0].StudyProgram + " // " + albums[0].Event + "\n" + albums[0].URL + "\n"
	} else {
		for idx, album := range albums {
			yearStr := strconv.Itoa(album.Year)
			idxStr := strconv.Itoa(*searchAlbum.Pointer + idx + 1)
			result += idxStr + ") " + yearStr + " // " + album.StudyProgram + " // " + album.Event + "\n" + album.URL + "\n"
		}
	}

	return result, *searchAlbum.Pointer, count, nil
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

// GetEventNames возвращает список названий событий для фотографии
func (r *Repo) GetEventNames() (string, error) {
	var events []ds.Event
	err := r.db.Select(&events, "SELECT * FROM events ORDER BY name")
	if err != nil {
		return "", fmt.Errorf("[db.Select]: %w", err)
	}

	var result string
	for idx, event := range events {
		idxStr := strconv.Itoa(idx + 1)
		result += idxStr + ") " + event.Name + "\n"
	}

	return result, nil
}

// GetTeacherNames возвращает ФИО преподавателей
func (r *Repo) GetTeacherNames(vkID int) (string, int, int, error) {
	var searchAlbum ds.SearchAlbum
	err := r.db.Get(&searchAlbum, "SELECT * FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return "", 0, 0, fmt.Errorf("[db.Get]: %w", err)
	}

	var name string
	if searchAlbum.Surname != nil {
		name = *searchAlbum.Surname
	}

	var count int
	err = r.db.Get(&count, "SELECT COUNT(*) FROM teacher_albums WHERE LOWER(name) LIKE $1 || '%'", strings.ToLower(name))
	if err != nil {
		return "", 0, 0, fmt.Errorf("[db.Get]: %w", err)
	}

	var albums []ds.TeacherAlbum
	err = r.db.Select(&albums, "SELECT * FROM teacher_albums WHERE LOWER(name) LIKE $1 || '%' OFFSET $2 LIMIT 10", strings.ToLower(name), searchAlbum.Pointer)
	if err != nil {
		return "", 0, 0, fmt.Errorf("[db.Select]: %w", err)
	}

	var result string
	for idx, album := range albums {
		idxStr := strconv.Itoa(*searchAlbum.Pointer + idx + 1)
		result += idxStr + ") " + album.Name + "\n" + album.URL + "\n"
	}

	return result, *searchAlbum.Pointer, count, nil
}

// ChangePointerTeacher меняет указатель при поиске альбома преподавателя
func (r *Repo) ChangePointerTeacher(vkID int, flag bool) error {
	var pointer int

	err := r.db.Get(&pointer, "SELECT pointer FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	if !flag {
		pointer -= 10
	} else {
		pointer += 10
	}

	if pointer >= 0 {
		_, err = r.db.Exec("UPDATE search_album SET pointer = $1 WHERE user_id = $2", pointer, vkID)
		if err != nil {
			return fmt.Errorf("[db.Exec]: %w", err)
		}
	}

	return nil
}

// ChangePointerStudents меняет указатель при поиске альбома студентов
func (r *Repo) ChangePointerStudents(vkID int, flag bool) error {
	var pointer int

	err := r.db.Get(&pointer, "SELECT pointer FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	if !flag {
		pointer -= 10
	} else {
		pointer += 10
	}

	if pointer >= 0 {
		_, err = r.db.Exec("UPDATE search_album SET pointer = $1 WHERE user_id = $2", pointer, vkID)
		if err != nil {
			return fmt.Errorf("[db.Exec]: %w", err)
		}
	}

	return nil
}

// DeletePointer удаляет указатель
func (r *Repo) DeletePointer(vkID int) error {
	_, err := r.db.Exec("UPDATE search_album SET pointer = 0 WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Exec]: %w", err)
	}
	return nil
}

// DeleteSurname удаляет фамилию преподавателя из поиска альбома
func (r *Repo) DeleteSurname(vkID int) error {
	_, err := r.db.Exec("UPDATE search_album SET surname = '' WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Exec]: %w", err)
	}
	return nil
}

// UpdateName добавляет первые буквы фамилии преподавателя для поиска альбома
func (r *Repo) UpdateName(ctx context.Context, vkID int, name string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET surname = $1 WHERE user_id = $2", name, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// GetSearchParams возвращает параметры поиска
func (r *Repo) GetSearchParams(vkID int) (string, error) {
	var searchAlbum ds.SearchAlbum
	err := r.db.Get(&searchAlbum, "SELECT * FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return "", fmt.Errorf("[db.Get]: %w", err)
	}

	var year, studyProgram, event string
	if searchAlbum.Year == nil {
		year = "Не указано"
	} else {
		year = strconv.Itoa(*searchAlbum.Year)
	}

	if searchAlbum.StudyProgram == nil {
		studyProgram = "Не указано"
	} else {
		studyProgram = *searchAlbum.StudyProgram
	}

	if searchAlbum.Event == nil {
		event = "Не указано"
	} else {
		event = *searchAlbum.Event
	}

	var result string
	result += "Параметры поиска:" + "\n" + "1) Год события: " + year + "\n" + "2) Программа обучения: " + studyProgram + "\n" + "3) Название события: " + event

	return result, nil
}
