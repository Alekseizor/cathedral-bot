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

// DeleteYear удаляет год события для поиска альбома
func (r *Repo) DeleteYear(ctx context.Context, vkID int) error {
	var year *int
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET year = $1 WHERE user_id = $2", year, vkID)
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
		return "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	var result string

	for idx, album := range albums {
		yearStr := strconv.Itoa(album.Year)
		idxStr := strconv.Itoa(idx + 1)
		result += idxStr + ") " + yearStr + " // " + album.StudyProgram + " // " + album.Event + "\n" + album.URL + "\n"
	}

	return result, nil
}
