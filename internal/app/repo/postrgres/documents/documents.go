package documents

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (r *Repo) IfDocumentsExists(vkID int) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM documents WHERE user_id = $1", vkID)
	if err != nil {
		return false, fmt.Errorf("[db.Get]: %w", err)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (r *Repo) UploadDocument(ctx context.Context, document ds.RequestDocument) error {
	query := `
        INSERT INTO documents (title, author, year, category, description, hashtags, attachment, user_id)
        VALUES (:title, :author, :year, :category, :description, :hashtags, :attachment, :user_id)
    `

	_, err := r.db.NamedExecContext(ctx, query, document)
	if err != nil {
		return fmt.Errorf("[db.NamedExecContext]: %w", err)
	}

	return nil
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

func (r *Repo) UpdateCategoryByCategoryName(ctx context.Context, documentID int, categoryName string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE documents SET category = $1 WHERE id = $2", categoryName, documentID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (r *Repo) NewCategory(ctx context.Context, category string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO categories(name) VALUES ($1)", category)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (r *Repo) SearchDocuments(ctx context.Context, params ds.SearchDocument, vkID int) (int, error) {
	var conditions []string
	var values []interface{}

	if params.Year.Valid {
		conditions = append(conditions, "year = $"+strconv.Itoa(len(values)+1))
		year, _ := params.Year.Value()
		values = append(values, year)
	} else if params.StartYear.Valid && params.EndYear.Valid {
		startYear, _ := params.StartYear.Value()
		endYear, _ := params.EndYear.Value()
		conditions = append(conditions, "year BETWEEN $"+strconv.Itoa(len(values)+1)+" AND $"+strconv.Itoa(len(values)+2))
		values = append(values, startYear, endYear)
	}

	if params.Title.Valid {
		title, _ := params.Title.Value()
		conditions = append(conditions, "Lower(title) like lower('%'||$"+strconv.Itoa(len(values)+1)+"||'%')")
		values = append(values, title)
	}

	if params.Author.Valid {
		author, _ := params.Author.Value()
		conditions = append(conditions, "Lower(author) like lower('%'||$"+strconv.Itoa(len(values)+1)+"||'%')")
		values = append(values, author)
	}

	if len(params.Categories) > 0 {
		placeholder := "$" + strconv.Itoa(len(values)+1)
		conditions = append(conditions, "category = ANY("+placeholder+")")
		values = append(values, pq.Array(params.Categories))
	}

	if len(params.Hashtags) > 0 {
		placeholder := "$" + strconv.Itoa(len(values)+1)
		conditions = append(conditions, placeholder+" && hashtags")
		values = append(values, pq.Array(params.Hashtags))
	}

	query := "SELECT id FROM documents"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var documentsID []int
	for rows.Next() {
		var docID int
		err := rows.Scan(&docID)
		if err != nil {
			return 0, err
		}
		documentsID = append(documentsID, docID)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE search_document SET documents = $1 WHERE user_id = $2", pq.Array(documentsID), vkID)
	if err != nil {
		return 0, fmt.Errorf("[db.ExecContext]: %w", err)
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return len(documentsID), nil
}

// GetSearchDocuments выводит найденные документы
func (r *Repo) GetSearchDocuments(ctx context.Context, vkID int) (string, error) {
	var params ds.SearchDocument
	err := r.db.QueryRowContext(ctx, "SELECT documents, pointer_doc FROM search_document WHERE user_id = $1", vkID).Scan(&params.Documents, &params.PointerDoc)

	query := `SELECT id, title, author, year, category
		FROM documents
		WHERE id = ANY ($1)`

	var documentIDs pq.Int64Array
	if len(params.Documents) >= params.PointerDoc+5 {
		documentIDs = params.Documents[params.PointerDoc : params.PointerDoc+5]
	} else {
		documentIDs = params.Documents[params.PointerDoc:]
	}

	rows, err := r.db.QueryContext(ctx, query, pq.Array(documentIDs))
	if err != nil {
		return "", fmt.Errorf("[db.QueryContext]: %w", err)
	}

	defer rows.Close()

	var output string
	for rows.Next() {
		var doc ds.Document
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.Author, &doc.Year, &doc.Category); err != nil {
			return "", err
		}
		index := params.PointerDoc + 1
		if doc.Author == nil {
			text := "Автор не указан"
			doc.Author = &text
		}
		if doc.Year != nil {
			output += fmt.Sprintf("[%d]. %s, %s, %d (%s)\n", index, *doc.Title, *doc.Author, *doc.Year, *doc.Category)
		} else {
			output += fmt.Sprintf("[%d]. %s, %s (%s)\n", index, *doc.Title, *doc.Author, *doc.Category)
		}
		params.PointerDoc++
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("[db.Scan]: %w", err)
	}

	return output, nil
}
