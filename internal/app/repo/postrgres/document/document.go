package document

import (
	"context"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Repo инстанс репо для работы с документами пользователя
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с документами пользователя
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// InsertDocumentURL добавляет URL для скачивания документа и id пользователя
func (r *Repo) InsertDocumentURL(ctx context.Context, title, docURL string, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO documents(title, url, user_id) VALUES ($1, $2, $3)", title, docURL, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

// UpdateName добавляет название документа
func (r *Repo) UpdateName(ctx context.Context, vkID int, name string) error {
	var doc ds.Document
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE documents SET title = $1 WHERE id = $2", name, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateAuthor добавляет ФИО автора документа
func (r *Repo) UpdateAuthor(ctx context.Context, vkID int, author string) error {
	var doc ds.Document
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE documents SET author = $1 WHERE id = $2", author, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateYear добавляет год создания документа
func (r *Repo) UpdateYear(ctx context.Context, vkID int, year int) error {
	var doc ds.Document
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE documents SET year = $1 WHERE id = $2", year, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// GetCategoryNames возвращает список названий категорий
func (r *Repo) GetCategoryNames() (string, error) {
	var output string

	rows, err := r.db.Query("SELECT CONCAT(id, '. ', name) AS formatted_string FROM categories")
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

// GetCategoryMaxID возвращает максимальное ID из категорий
func (r *Repo) GetCategoryMaxID() (int, error) {
	var maxID int
	err := r.db.Get(&maxID, "SELECT MAX(id) FROM categories")
	if err != nil {
		return 0, fmt.Errorf("[db.Get]: %w", err)
	}

	return maxID, nil
}

// GetCategoryNameByID возвращает название категории по ее ID
//func (r *Repo) GetCategoryNameByID(id int) (string, error) {
//	var name string
//	err := r.db.Get(&name, "SELECT name FROM categories WHERE id = $1", id)
//	if err != nil {
//		return "", fmt.Errorf("[db.Get]: %w", err)
//	}
//
//	return name, nil
//}

// UpdateCategory добавляет категорию документа
func (r *Repo) UpdateCategory(ctx context.Context, vkID, categoryNumber int) error {
	var doc ds.Document
	var name string
	err := r.db.Get(&name, "SELECT name FROM categories WHERE id = $1", categoryNumber)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	err = r.db.GetContext(ctx, &doc, "SELECT id FROM documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE documents SET category = $1 WHERE id = $2", name, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateUserCategory добавляет пользовательскую категорию документа
func (r *Repo) UpdateUserCategory(ctx context.Context, vkID int, category string) error {
	var doc ds.Document

	err := r.db.GetContext(ctx, &doc, "SELECT id FROM documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE documents SET (category, is_category_new) = ($1, $2) WHERE id = $3", category, true, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateHashtags добавляет хештеги к документу
func (r *Repo) UpdateHashtags(ctx context.Context, vkID int, hashtags []string) error {
	var doc ds.Document
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE documents SET hashtags = $1 WHERE id = $2", pq.Array(hashtags), doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// CheckParams возвращает все параметры заявки на загрузку документа
func (r *Repo) CheckParams(ctx context.Context, vkID int) (string, string, error) {
	var doc ds.Document
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return "", "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	sqlQuery := `
	SELECT 
    	'1. Название: ' || COALESCE(title, 'Не указано') AS name,
    	'2. Автор: ' || COALESCE(author, 'Не указано') AS author,
    	'3. Год создания документа: ' || COALESCE(CAST(year AS VARCHAR), 'Не указано') AS year,
    	'4. Категория: ' || COALESCE(category, 'Не указано') AS category,
    	'5. Хэштеги: ' || COALESCE(array_to_string(hashtags, ', '), 'Не указано') AS hashtag
	FROM documents
	WHERE id = $1;`

	var (
		name       string
		author     string
		year       string
		category   string
		hashtag    string
		attachment string
	)

	err = r.db.QueryRow(sqlQuery, doc.ID).Scan(&name, &author, &year, &category, &hashtag)
	if err != nil {
		return "", "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	attachment = "doc" + "185404885" + "_" + "673328305" + "_" + "zdIf99RBZYxvX0EQRfa4drMLgNLgFRzLZRalqSbtyns"

	output := fmt.Sprintf("Ваша заявка на загрузку документа:\n %s\n%s\n%s\n%s\n%s\n", name, author, year, category, hashtag)

	return output, attachment, nil
}
