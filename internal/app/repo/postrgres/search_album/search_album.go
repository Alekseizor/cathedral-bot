package search_album

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	"strconv"
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
	err := r.db.Get(&event, "SELECT name FROM events WHERE id = $1", eventNumber)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE search_album SET event = $1 WHERE user_id = $2", event, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateTeacher добавляет ФИО преподавателя для поиска альбома
func (r *Repo) UpdateTeacher(ctx context.Context, vkID int, teacherName string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET teacher = $1 WHERE user_id = $2", teacherName, vkID)
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

// DeleteTeacher удаляет ФИО преподавателя для поиска альбома
func (r *Repo) DeleteTeacher(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET teacher = $1 WHERE user_id = $2", nil, vkID)
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
func (r *Repo) ShowList(ctx context.Context, vkID int) (string, error) {
	var searchAlbum ds.SearchAlbum

	err := r.db.GetContext(ctx, &searchAlbum, "SELECT * FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return "", fmt.Errorf("[db.GetContext]: %w", err)
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

	query += " ORDER BY year" + " DESC"

	var albums []ds.StudentAlbum
	err = r.db.SelectContext(ctx, &albums, query, args...)
	if err != nil {
		return "", fmt.Errorf("[db.SelectContext]: %w", err)
	}

	var result string

	if len(albums) == 1 {
		yearStr := strconv.Itoa(albums[0].Year)
		result += yearStr + " // " + albums[0].StudyProgram + " // " + albums[0].Event + "\n" + albums[0].URL + "\n"
	} else {
		for idx, album := range albums {
			yearStr := strconv.Itoa(album.Year)
			idxStr := strconv.Itoa(idx + 1)
			result += idxStr + ") " + yearStr + " // " + album.StudyProgram + " // " + album.Event + "\n" + album.URL + "\n"
		}
	}

	return result, nil
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

// GetTeacherNames возвращает ФИО преподавателей
func (r *Repo) GetTeacherNames() (string, error) {
	var teacherNames string

	rows, err := r.db.Query("SELECT CONCAT(id, ') ', name) AS formatted_string FROM teachers")
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
		teacherNames += formattedString + "\n"
	}

	return teacherNames, nil
}

// GetTeacherMaxID возвращает максимальное ID из преподавателей
func (r *Repo) GetTeacherMaxID() (int, error) {
	var maxID int
	err := r.db.Get(&maxID, "SELECT MAX(id) FROM teachers")
	if err != nil {
		return 0, fmt.Errorf("[db.Get]: %w", err)
	}

	return maxID, nil
}

// GetTeacherName возвращает ФИО преподавателя
func (r *Repo) GetTeacherName(ctx context.Context, teacherID int) (string, error) {
	var name string
	err := r.db.Get(&name, "SELECT name FROM teachers WHERE id = $1", teacherID)
	if err != nil {
		return "", fmt.Errorf("[db.Get]: %w", err)
	}

	return name, nil
}

// ShowTeacher возвращает найденный альбом преподавателя
func (r *Repo) ShowTeacher(ctx context.Context, vkID int) (string, error) {
	var teacher string

	err := r.db.GetContext(ctx, &teacher, "SELECT teacher FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	var album ds.TeacherAlbum
	err = r.db.GetContext(ctx, &album, "SELECT * FROM teacher_albums WHERE teacher = $1", teacher)
	if err != nil {
		return "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	result := album.Teacher + "\n" + album.URL

	return result, nil
}
