package documents

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strconv"
	"strings"
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

func (r *Repo) CheckExistence(ctx context.Context, vkID int) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS (SELECT 1 FROM documents WHERE id = $1)", vkID)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return exists, nil
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
		conditions = append(conditions, "title = $"+strconv.Itoa(len(values)+1))
		values = append(values, title)
	}

	if params.Author.Valid {
		author, _ := params.Author.Value()
		conditions = append(conditions, "author = $"+strconv.Itoa(len(values)+1))
		values = append(values, author)
	}

	if len(params.Categories) > 0 {
		placeholder := "$" + strconv.Itoa(len(values)+1)
		conditions = append(conditions, "category = ANY("+placeholder+")")
		values = append(values, pq.Array(params.Categories))
	}

	if len(params.Hashtags) > 0 {
		placeholder := "$" + strconv.Itoa(len(values)+1)
		conditions = append(conditions, "hashtags = ANY("+placeholder+")")
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

	query := `SELECT id, title, category
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
	//rows, err := r.db.QueryContext(ctx, query, pq.Array([]int{1, 11, 21, 31, 41}))
	//if err != nil {
	//	return "", fmt.Errorf("[db.QueryContext]: %w", err)
	//}
	defer rows.Close()

	var output string
	for rows.Next() {
		var doc ds.Documents
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.Category); err != nil {
			return "", err
		}
		index := params.PointerDoc + 1
		output += fmt.Sprintf("[%d]. Название: %s. Категория: %s\n", index, doc.Title, doc.Category)
		params.PointerDoc++
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("[db.Scan]: %w", err)
	}

	return output, nil
}
