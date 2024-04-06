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

//func (r *Repo) SearchDocuments(ctx context.Context, params ds.SearchDocument) (int, error) {
//	var conditions []string
//	var values []interface{}
//
//	if params.Year.Valid {
//		conditions = append(conditions, "year = ?")
//		year, _ := params.Year.Value()
//		values = append(values, year)
//	} else if params.StartYear.Valid && params.EndYear.Valid {
//		startYear, _ := params.StartYear.Value()
//		endYear, _ := params.EndYear.Value()
//		conditions = append(conditions, "year BETWEEN ? AND ?")
//		values = append(values, startYear, endYear)
//	}
//
//	if params.Title.Valid {
//		title, _ := params.Title.Value()
//		conditions = append(conditions, "title = ?")
//		values = append(values, title)
//	}
//
//	if params.Author.Valid {
//		author, _ := params.Author.Value()
//		conditions = append(conditions, "author = ?")
//		values = append(values, author)
//	}
//
//	if len(params.Categories) > 0 {
//		conditions = append(conditions, "category = ANY(?)")
//		values = append(values, pq.Array(params.Categories))
//	}
//
//	if len(params.Hashtags) > 0 {
//		conditions = append(conditions, "hashtags = ANY(?)")
//		values = append(values, pq.Array(params.Hashtags))
//	}
//
//	query := "SELECT * FROM documents"
//	if len(conditions) > 0 {
//		query += " WHERE " + strings.Join(conditions, " AND ")
//	}
//
//	query = "SELECT * FROM documents WHERE year BETWEEN $1 AND $2 AND category = ANY($3)"
//	rows, err := r.db.QueryContext(ctx, query, values...)
//	if err != nil {
//		return 0, err
//	}
//	defer rows.Close()
//
//	var documents []ds.Documents
//	for rows.Next() {
//		var doc ds.Documents
//		err := rows.Scan(&doc.ID, &doc.Title, &doc.Author, &doc.Year, pq.Array(&doc.Category), pq.Array(&doc.Hashtags))
//		if err != nil {
//			return 0, err
//		}
//		documents = append(documents, doc)
//	}
//
//	if err := rows.Err(); err != nil {
//		return 0, err
//	}
//
//	return len(documents), nil
//}

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
