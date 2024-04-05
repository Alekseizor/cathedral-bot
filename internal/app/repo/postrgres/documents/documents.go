package documents

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repo инстанс репо для работы с загруженными документами
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с загруженными документами
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CheckExistence(ctx context.Context, documentID int) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS (SELECT 1 FROM documents WHERE id = $1)", documentID)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return exists, nil
}

// GetOutput возвращает данные о документе в формате вывода
func (r *Repo) GetOutput(ctx context.Context, documentID int) (string, string, error) {
	sqlQuery := `
	SELECT 
    	'1. Название: ' || COALESCE(title, 'Не указано') AS name,
    	'2. Автор: ' || COALESCE(author, 'Не указано') AS author,
    	'3. Год создания документа: ' || COALESCE(CAST(year AS VARCHAR), 'Не указано') AS year,
    	'4. Категория: ' || COALESCE(category, 'Не указано') AS category,
    	'5. Описание: ' || COALESCE(description, 'Не указано') AS description,
    	'6. Хэштеги: ' || COALESCE(array_to_string(hashtags, ', '), 'Не указано') AS hashtag,
    	attachment
	FROM documents
	WHERE id = $1;`

	var (
		name        string
		author      string
		year        string
		category    string
		description string
		hashtag     string
		attachment  string
	)

	err := r.db.QueryRow(sqlQuery, documentID).Scan(&name, &author, &year, &category, &description, &hashtag, &attachment)
	if err != nil {
		return "", "", fmt.Errorf("[db.QueryRow]: %w", err)
	}

	output := fmt.Sprintf(" Документ:\n %s\n%s\n%s\n%s\n%s\n%s\n", name, author, year, category, description, hashtag)

	return output, attachment, nil
}

func (r *Repo) Delete(ctx context.Context, documentID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM documents WHERE id = $1", documentID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}
