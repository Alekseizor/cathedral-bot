package teacher_albums

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repo инстанс репо для работы с альбомами студентов
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с параметрами для альбомов студентов
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetAlbum(ctx context.Context, albumID int) (string, error) {
	sqlQuery := `
	SELECT 
    	'1) Альбом про преподавателя: ' || COALESCE(name, 'Не указано') AS name,
    	'2) Описание: ' || COALESCE(description, 'Не указано') AS description,
    	'3) url: ' || COALESCE(url, 'Не указано') AS url
	FROM teacher_albums
	WHERE id = $1;`

	var (
		teacher     string
		description string
		url         string
	)

	err := r.db.QueryRowContext(ctx, sqlQuery, albumID).Scan(&teacher, &description, &url)
	if err != nil {
		return "", fmt.Errorf("[db.QueryRowContext]: %w", err)
	}

	output := fmt.Sprintf("%s\n%s\n%s", teacher, description, url)

	return output, nil
}

func (r *Repo) GetAllAlbumsOutput(ctx context.Context) ([]string, error) {
	sqlQuery := `
    SELECT 
        'ID: ' || CAST(id AS VARCHAR) AS id,
        'Название альбома: ' || name AS name,
        'url: ' || url AS url
    FROM teacher_albums;`

	rows, err := r.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("[db.QueryContext]: %w", err)
	}
	defer rows.Close()

	output := make([]string, 0)
	index := 0
	outputElem := ""

	for rows.Next() {
		index++
		var id, name, url string
		if err := rows.Scan(&id, &name, &url); err != nil {
			return nil, fmt.Errorf("[rows.Scan]: %w", err)
		}
		outputElem += fmt.Sprintf("%s\n%s\n%s\n\n", id, name, url)

		if index == 10 {
			output = append(output, outputElem)
			index = 0
			outputElem = ""
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[rows.Err]: %w", err)
	}

	return output, nil
}

func (r *Repo) CheckExistence(ctx context.Context, albumID int) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS (SELECT 1 FROM teacher_albums WHERE id = $1)", albumID)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return exists, nil
}

func (r *Repo) UpdateName(ctx context.Context, albumID int, eventNumber int) error {
	var name string
	err := r.db.Get(&name, "SELECT name FROM teachers WHERE id = $1", eventNumber)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE teacher_albums SET teacher = $1 WHERE id = $2", name, albumID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (r *Repo) UpdateNewName(ctx context.Context, albumID int, teacher string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO teachers (name) VALUES ($1)", teacher)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE teacher_albums SET teacher = $1 WHERE id = $2", teacher, albumID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (r *Repo) UpdateDescription(ctx context.Context, albumID int, description string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE teacher_albums SET description = $1 WHERE id = $2", description, albumID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (r *Repo) GetVKID(ctx context.Context, albumID int) (string, error) {
	var vkID string
	err := r.db.GetContext(ctx, &vkID, "SELECT vk_id FROM teacher_albums WHERE id = $1", albumID)
	if err != nil {
		return "", fmt.Errorf("[db.GetContext]: %w", err)
	}
	return vkID, nil
}

func (r *Repo) GetTitle(ctx context.Context, albumID int) (string, error) {
	var (
		year         int
		studyProgram string
		event        string
	)

	err := r.db.QueryRowContext(ctx, "SELECT year,study_program,event FROM teacher_albums WHERE id = $1", albumID).Scan(&year, &studyProgram, &event)
	if err != nil {
		return "", fmt.Errorf("[db.QueryRowContext]: %w", err)
	}

	return fmt.Sprintf("%d // %s // %s", year, studyProgram, event), nil
}

func (r *Repo) Delete(ctx context.Context, albumID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM teacher_albums WHERE id = $1", albumID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}
