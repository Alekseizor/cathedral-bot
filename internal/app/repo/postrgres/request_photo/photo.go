package request_photo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
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

// InsertPhotoURL добавляет URL для скачивания фото
func (r *Repo) InsertPhotoURL(ctx context.Context, title, docURL string, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO documents(title, url, user_id) VALUES ($1, $2, $3)", title, docURL, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}
