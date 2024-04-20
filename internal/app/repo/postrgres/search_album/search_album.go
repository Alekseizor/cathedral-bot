package search_album

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Repo инстанс репо для работы с параметрами для поиска альбома
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с параметрами для поиска альбома
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// CreateSearchAlbum создает запись для поиска альбома
func (r *Repo) CreateSearchAlbum(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO search_album (user_id) VALUES ($1)", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeleteSearchAlbum удаляет запись для поиска альбома
func (r *Repo) DeleteSearchAlbum(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// UpdateYear добавляет год события для поиска альбома
func (r *Repo) UpdateYear(ctx context.Context, vkID int, year int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET year = $1 WHERE user_id = $2", year, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// DeleteYear удаляет год события для поиска альбома
func (r *Repo) DeleteYear(ctx context.Context, vkID int) error {
	var year *int
	_, err := r.db.ExecContext(ctx, "UPDATE search_album SET year = $1 WHERE user_id = $2", year, vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// YearCountAlbums возвращает количество найденных альбомов по году
func (r *Repo) YearCountAlbums(ctx context.Context, vkID int) (int, error) {
	var year int
	err := r.db.GetContext(ctx, &year, "SELECT year FROM search_album WHERE user_id = $1", vkID)
	if err != nil {
		return 0, fmt.Errorf("[db.GetContext]: %w", err)
	}

	var count int
	err = r.db.GetContext(ctx, &count, "SELECT count(*) FROM albums WHERE year = $1", year)
	if err != nil {
		return 0, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return count, nil
}
