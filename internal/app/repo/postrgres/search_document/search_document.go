package search_document

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
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

// CheckSearchParams возвращает все параметры для поиска документа
func (r *Repo) CheckSearchParams(ctx context.Context, vkID int) (string, error) {

	sqlQuery := `
	SELECT
    	CONCAT('Название: ', COALESCE(title, 'Не указано')) AS title,
    	CONCAT('Автор: ', COALESCE(author, 'Не указано')) AS author,
    	CASE
        	WHEN year IS NOT NULL THEN CONCAT('Год создания документа: ', year)
        	WHEN start_year IS NOT NULL AND end_year IS NOT NULL THEN CONCAT('Временной интервал: ', start_year, '-', end_year)
        	ELSE 'Год создания/Временной интервал: Не указано'
    	END AS year_interval,
    	CONCAT('Список категорий: ', COALESCE(ARRAY_TO_STRING(categories, ', '), 'Не указано')) AS categories,
    	CONCAT('Хэштеги: ', COALESCE(ARRAY_TO_STRING(hashtags, ', '), 'Не указано')) AS hashtags
	FROM
    	search_document
	WHERE user_id = $1;`

	var searchParams ds.ParseSearchDocument

	err := r.db.QueryRowContext(ctx, sqlQuery, vkID).Scan(&searchParams.Title, &searchParams.Author, &searchParams.YearInterval, &searchParams.Categories, &searchParams.Hashtags)
	if err != nil {
		return "", fmt.Errorf("[db.QueryRowContext]: %w", err)
	}

	output := fmt.Sprintf("Ваши параметры для поиска:\n %s\n%s\n%s\n%s\n%s\n", searchParams.Title, searchParams.Author, searchParams.YearInterval, searchParams.Categories, searchParams.Hashtags)

	return output, nil
}

// ParseSearch парсит параметры поиска по id юзера
func (r *Repo) ParseSearch(ctx context.Context, vkID int) (ds.SearchDocument, error) {
	var doc ds.SearchDocument
	err := r.db.QueryRowContext(ctx, "SELECT id, title, author, year, start_year, end_year, categories, hashtags FROM search_document WHERE user_id = $1", vkID).Scan(&doc.ID, &doc.Title, &doc.Author, &doc.Year, &doc.StartYear, &doc.EndYear, &doc.Categories, &doc.Hashtags)
	if err != nil {
		return ds.SearchDocument{}, fmt.Errorf("[db.QueryRowContext]: %w", err)
	}

	return doc, nil
}

// UpdatePointer обновляет указатель на докумет
func (r *Repo) UpdatePointer(ctx context.Context, value, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_document SET pointer_doc =  pointer_doc + $1 WHERE user_id = $2", value, vkID)
	if err != nil {
		return fmt.Errorf("[db.QueryRowContext]: %w", err)
	}

	return nil
}

// ParseSearchButtons парсит параметры, необходимые для оторбражения кнопок листинга
func (r *Repo) ParseSearchButtons(ctx context.Context, vkID int) (ds.SearchDocument, error) {
	var doc ds.SearchDocument
	err := r.db.QueryRowContext(ctx, "SELECT documents, pointer_doc FROM search_document WHERE user_id = $1", vkID).Scan(&doc.Documents, &doc.PointerDoc)
	if err != nil {
		return ds.SearchDocument{}, fmt.Errorf("[db.QueryRowContext]: %w", err)
	}

	return doc, nil
}
