package personal_account_photo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

// Repo инстанс репо для работы с фотографиями пользователя
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с фотографиями пользователя
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// GetRequestPhoto возвращает заявку на загрузку фото
func (r *Repo) GetRequestPhoto(vkID int) (string, string, int, int, error) {
	var pointer int
	err := r.db.Get(&pointer, "SELECT pointer FROM personal_account_photo WHERE user_id = $1", vkID)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("[db.Get]: %w", err)
	}

	var status int
	var attachment string
	err = r.db.QueryRow("SELECT status, attachment FROM request_photo WHERE user_id = $1 offset $2 limit 1", vkID, pointer).Scan(&status, &attachment)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("[db.QueryRow]: %w", err)
	}

	var count int
	err = r.db.Get(&count, "SELECT count(*) FROM request_photo WHERE user_id = $1", vkID)
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
	case 4:
		statusString = "Одобрена"
	}

	message := fmt.Sprintf("Заявка %s/%s\nСтатус данной заявки: %s", pointerString, countString, statusString)

	return message, attachment, pointer, count, nil
}

// CreatePersonalAccountPhoto создает запись в личном кабинете пользователя
func (r *Repo) CreatePersonalAccountPhoto(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO personal_account_photo (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// ChangePointer меняет указатель на заявку в личном кабинете
func (r *Repo) ChangePointer(vkID int, flag bool) error {
	var pointer int

	err := r.db.Get(&pointer, "SELECT pointer FROM personal_account_photo WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	if !flag {
		pointer--
	} else {
		pointer++
	}

	if pointer >= 0 {
		_, err = r.db.Exec("UPDATE personal_account_photo SET pointer = $1 WHERE user_id = $2", pointer, vkID)
		if err != nil {
			return fmt.Errorf("[db.Exec]: %w", err)
		}
	}

	return nil
}

// DeletePointer удаляет указатель на заявку в личном кабинете
func (r *Repo) DeletePointer(vkID int) error {
	_, err := r.db.Exec("UPDATE personal_account_photo SET pointer = 0 WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.Exec]: %w", err)
	}
	return nil
}
