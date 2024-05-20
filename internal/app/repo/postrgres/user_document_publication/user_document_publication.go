package user_document_publication

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

// Repo инстанс репо для работы с заявками пользователя на загрузку документа
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с заявками пользователя на загрузку документа
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// GetRequestDocument возвращает заявку на загрузку документа
func (r *Repo) GetRequestDocument(vkID int) (string, string, int, int, error) {
	var pointer int
	err := r.db.Get(&pointer, "SELECT pointer FROM user_document_publication WHERE user_id = $1", vkID)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("[db.Get]: %w", err)
	}

	sqlQuery := `
	SELECT 
    	'1. Название: ' || COALESCE(title, 'Не указано') AS name,
    	'2. Автор: ' || COALESCE(author, 'Не указано') AS author,
    	'3. Год создания документа: ' || COALESCE(CAST(year AS VARCHAR), 'Не указано') AS year,
    	'4. Категория: ' || COALESCE(category, 'Не указано') AS category,
    	'5. Описание: ' || COALESCE(description, 'Не указано') AS description,
    	'6. Хэштеги: ' || COALESCE(array_to_string(hashtags, ', '), 'Не указано') AS hashtag,
    	attachment,
    	status
	FROM requests_documents
	WHERE user_id = $1
	OFFSET $2
	LIMIT 1;`

	var (
		name        string
		author      string
		year        string
		category    string
		description string
		hashtag     string
		attachment  string
		status      int
	)

	// достали заявку
	err = r.db.QueryRow(sqlQuery, vkID, pointer).Scan(&name, &author, &year, &category, &description, &hashtag, &attachment, &status)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("[db.QueryRowContext]: %w", err)
	}

	var count int
	err = r.db.Get(&count, "SELECT count(*) FROM requests_documents WHERE user_id = $1", vkID)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("[db.Get]: %w", err)
	}

	countString := strconv.Itoa(count)
	pointerString := strconv.Itoa(pointer + 1)

	var statusString string

	switch status {
	case 1:
		statusString = "Ожидает рассмотрения"
	case 2:
		statusString = "На рассмотрении"
	case 3:
		statusString = "Отклонена"
	}

	output := fmt.Sprintf("Заявка %s/%s на публикацию документа [%s]:\n %s\n%s\n%s\n%s\n%s\n%s\n", pointerString, countString, statusString, name, author, year, category, description, hashtag)

	return output, attachment, pointer, count, nil
}

// CreateUserDocumentPublication создает запись в личном кабинете пользователя
func (r *Repo) CreateUserDocumentPublication(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO user_document_publication (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// ChangePointer меняет указатель на заявку в личном кабинете
func (r *Repo) ChangePointer(vkID int, flag bool) error {
	var pointer int

	err := r.db.Get(&pointer, "SELECT pointer FROM user_document_publication WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	if !flag {
		pointer--
	} else {
		pointer++
	}

	if pointer >= 0 {
		_, err = r.db.Exec("UPDATE user_document_publication SET pointer = $1 WHERE user_id = $2", pointer, vkID)
		if err != nil {
			return fmt.Errorf("[db.Exec]: %w", err)
		}
	}

	return nil
}

// DeletePointer удаляет указатель на заявку в личном кабинете
func (r *Repo) DeletePointer(vkID int) error {
	_, err := r.db.Exec("UPDATE user_document_publication SET pointer = 0 WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Exec]: %w", err)
	}
	return nil
}
