package search_document

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Repo инстанс репо для работы с параметрами для поиска документа
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с параметрами для поиска документа
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// CreateSearch создает запись для поиска
func (r *Repo) CreateSearch(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO search_document (user_id) VALUES ($1)", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeleteSearch удаляет запись для поиска
func (r *Repo) DeleteSearch(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM search_document WHERE user_id = ($1)", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	return nil
}

// UpdateNameSearch добавляет название документа для поиска
func (r *Repo) UpdateNameSearch(ctx context.Context, name string, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_document SET title = $1 WHERE user_id = $2", name, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateAuthorSearch добавляет ФИО автора документа для поиска
func (r *Repo) UpdateAuthorSearch(ctx context.Context, author string, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_document SET author = $1 WHERE user_id = $2", author, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateYearSearch добавляет год создания документа для поиска
func (r *Repo) UpdateYearSearch(ctx context.Context, year, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_document SET (year, start_year, end_year) = ($1, NULL, NULL) WHERE user_id = $2", year, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateYearRangeSearch добавляет интервал по годам создания документа для поиска
func (r *Repo) UpdateYearRangeSearch(ctx context.Context, startYear, endYear, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_document SET (year, start_year, end_year) = (NULL, $1, $2) WHERE user_id = $3", startYear, endYear, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateCategoriesSearch добавляет список категорий документа для поиска
func (r *Repo) UpdateCategoriesSearch(ctx context.Context, categories []int, vkID int) error {
	var categoriesNames pq.StringArray
	err := r.db.SelectContext(ctx, &categoriesNames, "SELECT name FROM categories WHERE id = ANY($1)", pq.Array(categories))
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE search_document SET categories = $1 WHERE user_id = $2", categoriesNames, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateHashtagsSearch добавляет список хештегов документа для поиска
func (r *Repo) UpdateHashtagsSearch(ctx context.Context, hashtags []string, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_document SET hashtags = $1 WHERE user_id = $2", pq.Array(hashtags), vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}
