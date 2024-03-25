package state

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Alekseizor/cathedral-bot/internal/app/ds"
)

// Repo инстанс репо для работы с пользователями
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо, подключения к бд еще нет!
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// Get достает стейт пользователя из БД
// Если пользователя нет, вернет пустую структуру
func (r *Repo) Get(ctx context.Context, vkID int) (string, error) {
	var user ds.State
	err := r.db.GetContext(ctx, &user, "SELECT title FROM state WHERE vk_id = $1", vkID)
	if err != nil {
		return "", err
	}

	return user.Title, nil
}

// Insert создает стейт для пользователя в БД
func (r *Repo) Insert(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO state VALUES ($1, 'start')", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

// Update обновляет стейт пользователя
func (r *Repo) Update(ctx context.Context, vkID int, newState string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE state SET title = $1 WHERE vk_id = $2", newState, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}
