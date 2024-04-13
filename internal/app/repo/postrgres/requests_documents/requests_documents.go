package requests_documents

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
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

type docsDoc struct {
	File string `json:"file"`
}

// UpdateStatus изменяет статус заявки на загрузку документа по ID заявки
func (r *Repo) UpdateStatus(ctx context.Context, status int, reqDocID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE requests_documents SET status = $1 WHERE id = $2", status, reqDocID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UploadDocument загружает документ
func (r *Repo) UploadDocument(ctx context.Context, VK *api.VK, doc object.DocsDoc, vkID int) error {
	resp, err := http.Get(doc.URL)
	if err != nil {
		log.Println(err)
	}
	upload, _ := VK.DocsGetMessagesUploadServer(api.Params{
		"type":    "doc",
		"peer_id": vkID,
	})
	file, err := io.ReadAll(resp.Body)
	fileBody := bytes.NewReader(file)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", doc.Title)
	io.Copy(part, fileBody)
	writer.Close()
	req, _ := http.NewRequest("POST", upload.UploadURL, bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, _ := client.Do(req)
	docs := &docsDoc{}
	json.NewDecoder(response.Body).Decode(docs)
	log.Println(docs.File)
	savedDoc, _ := VK.DocsSave(api.Params{
		"file":  docs.File,
		"title": doc.Title,
	})

	attachment := "doc" + strconv.Itoa(savedDoc.Doc.OwnerID) + "_" + strconv.Itoa(savedDoc.Doc.ID)

	_, err = r.db.ExecContext(ctx, "INSERT INTO requests_documents(title, attachment, user_id) VALUES ($1, $2, $3)", doc.Title, attachment, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeleteDocumentRequest удаляет заявку документа
func (r *Repo) DeleteDocumentRequest(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM requests_documents WHERE id = (SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1)", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	return nil
}

// GetDocumentAttachment возвращает attachment документа
func (r *Repo) GetDocumentAttachment(ctx context.Context, vkID int) (string, error) {
	var doc ds.RequestDocument
	err := r.db.GetContext(ctx, &doc, "SELECT attachment FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	return doc.Attachment, nil
}

// GetDocumentLastID возвращает ID последней заявки
func (r *Repo) GetDocumentLastID(ctx context.Context, vkID int) (int, error) {
	var doc ds.RequestDocument
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return 0, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return doc.ID, nil
}

// UpdateName добавляет название документа
func (r *Repo) UpdateName(ctx context.Context, vkID int, name string) error {
	var doc ds.RequestDocument
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE requests_documents SET title = $1 WHERE id = $2", name, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateAuthor добавляет ФИО автора документа
func (r *Repo) UpdateAuthor(ctx context.Context, vkID int, author string) error {
	var doc ds.RequestDocument
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE requests_documents SET author = $1 WHERE id = $2", author, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateYear добавляет год создания документа
func (r *Repo) UpdateYear(ctx context.Context, vkID int, year int) error {
	var doc ds.RequestDocument
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE requests_documents SET year = $1 WHERE id = $2", year, doc.ID)
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
func (r *Repo) GetCategoryNameByID(ctx context.Context, id int) (string, error) {
	var name string
	err := r.db.GetContext(ctx, &name, "SELECT name FROM categories WHERE id = $1", id)
	if err != nil {
		return "", fmt.Errorf("[db.Get]: %w", err)
	}

	return name, nil
}

func (r *Repo) CheckCategoryExistence(ctx context.Context, category string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT EXISTS (SELECT 1 FROM categories WHERE name = $1)", category)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return exists, nil
}

func (r *Repo) InsertCategory(ctx context.Context, category string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO categories(name) VALUES ($1)", category)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

// UpdateCategory добавляет категорию документа
func (r *Repo) UpdateCategory(ctx context.Context, vkID, categoryNumber int) error {
	var doc ds.RequestDocument
	var name string
	err := r.db.Get(&name, "SELECT name FROM categories WHERE id = $1", categoryNumber)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	err = r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE requests_documents SET (category, is_category_new) = ($1, $2) WHERE id = $3", name, false, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateUserCategory добавляет пользовательскую категорию документа
func (r *Repo) UpdateUserCategory(ctx context.Context, vkID int, category string) error {
	var doc ds.RequestDocument

	err := r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE requests_documents SET (category, is_category_new) = ($1, $2) WHERE id = $3", category, true, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateDescription добавляет описание документа
func (r *Repo) UpdateDescription(ctx context.Context, vkID int, description string) error {
	var doc ds.RequestDocument
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE requests_documents SET description = $1 WHERE id = $2", description, doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateHashtags добавляет хештеги к документу
func (r *Repo) UpdateHashtags(ctx context.Context, vkID int, hashtags []string) error {
	var doc ds.RequestDocument
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return fmt.Errorf("[db.GetContext]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE requests_documents SET hashtags = $1 WHERE id = $2", pq.Array(hashtags), doc.ID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// CheckParams возвращает все параметры заявки на загрузку документа
func (r *Repo) CheckParams(ctx context.Context, vkID int) (string, string, error) {
	var doc ds.RequestDocument
	err := r.db.GetContext(ctx, &doc, "SELECT id FROM requests_documents WHERE user_id = $1 ORDER BY id DESC LIMIT 1", vkID)
	if err != nil {
		return "", "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	sqlQuery := `
	SELECT 
    	'1. Название: ' || COALESCE(title, 'Не указано') AS name,
    	'2. Автор: ' || COALESCE(author, 'Не указано') AS author,
    	'3. Год создания документа: ' || COALESCE(CAST(year AS VARCHAR), 'Не указано') AS year,
    	'4. Категория: ' || COALESCE(category, 'Не указано') AS category,
    	'5. Описание: ' || COALESCE(description, 'Не указано') AS description,
    	'6. Хэштеги: ' || COALESCE(array_to_string(hashtags, ', '), 'Не указано') AS hashtag,
    	attachment
	FROM requests_documents
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

	err = r.db.QueryRow(sqlQuery, doc.ID).Scan(&name, &author, &year, &category, &description, &hashtag, &attachment)
	if err != nil {
		return "", "", fmt.Errorf("[db.GetContext]: %w", err)
	}

	output := fmt.Sprintf("Ваша заявка на загрузку:\n %s\n%s\n%s\n%s\n%s\n%s\n", name, author, year, category, description, hashtag)

	return output, attachment, nil
}
