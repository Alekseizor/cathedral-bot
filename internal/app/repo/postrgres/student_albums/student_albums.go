package student_albums

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
    	'1) Год события: ' || COALESCE(CAST(year AS VARCHAR), 'Не указано') AS year,
    	'2) Программа обучения: ' || COALESCE(study_program, 'Не указано') AS studyProgram,
    	'3) Название события: ' || COALESCE(event, 'Не указано') AS event,
    	'4) Описание: ' || COALESCE(description, 'Не указано') AS description
	FROM student_albums
	WHERE id = $1;`

	var (
		year         string
		studyProgram string
		event        string
		description  string
	)

	err := r.db.QueryRowContext(ctx, sqlQuery, albumID).Scan(&year, &studyProgram, &event, &description)
	if err != nil {
		return "", fmt.Errorf("[db.QueryRow]: %w", err)
	}

	output := fmt.Sprintf("%s\n%s\n%s\n%s\n", year, studyProgram, event, description)

	return output, nil
}

func (r *Repo) GetAllAlbumsOutput(ctx context.Context) (string, error) {
	sqlQuery := `
    SELECT 
        'ID: ' || CAST(id AS VARCHAR) AS id,
        CASE 
            WHEN year IS NULL AND study_program IS NULL AND event IS NULL THEN '2. Название альбома: Остальное' 
            ELSE 'Название альбома: ' || COALESCE(CAST(year AS VARCHAR), '---') ||' // '|| COALESCE(study_program, '---')||' // ' || COALESCE(event, '---') 
        END AS name
    FROM student_albums;`

	rows, err := r.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return "", fmt.Errorf("[db.QueryContext]: %w", err)
	}
	defer rows.Close()

	var output string

	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return "", fmt.Errorf("[rows.Scan]: %w", err)
		}
		output += fmt.Sprintf("%s\n%s\n\n", id, name)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("[rows.Err]: %w", err)
	}

	return output, nil
}

func (r *Repo) CheckExistence(ctx context.Context, albumID int) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS (SELECT 1 FROM student_albums WHERE id = $1)", albumID)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return exists, nil
}
