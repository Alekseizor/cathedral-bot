package admin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Repo инстанс репо для работы с администраторами
type Repo struct {
	db *sqlx.DB
}

// New - создаем новый объект репо для работы с таблицей администратора
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CheckExistence(ctx context.Context, vkID int) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, "SELECT * FROM admin WHERE vk_id = $1", vkID)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("[db.GetContext]: %w", err)
	}

	return exists, nil
}

// AddAlbumsAdmin создает нового администратора фотоархива
func (r *Repo) AddAlbumsAdmin(ctx context.Context, vkID int) error {
	exists, err := r.CheckExistence(ctx, vkID)
	if err != nil {
		return fmt.Errorf("[admin.CheckExistence]: %w", err)
	}

	if exists {
		_, err := r.db.ExecContext(ctx, "UPDATE admin SET albums = $1 WHERE vk_id = $2", true, vkID)
		if err != nil {
			return fmt.Errorf("[db.ExecContext]: %w", err)
		}
		return nil
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO admin VALUES ($1, false,true)", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (r *Repo) DeleteAlbumsAdmin(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE admin SET albums = $1 WHERE vk_id = $2", false, vkID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (r *Repo) GetAlbumsAdmins(ctx context.Context) ([]int, error) {
	var vkIDAdmins []int
	err := r.db.GetContext(ctx, &vkIDAdmins, "SELECT vk_id FROM admin WHERE albums = true")
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return vkIDAdmins, nil
}

// AddDocumentAdmin создает нового администратора документоархива
func (r *Repo) AddDocumentAdmin(ctx context.Context, vkID int) error {
	exists, err := r.CheckExistence(ctx, vkID)
	if err != nil {
		return fmt.Errorf("[admin.CheckExistence]: %w", err)
	}

	if exists {
		_, err := r.db.ExecContext(ctx, "UPDATE admin SET documents = $1 WHERE vk_id = $2", true, vkID)
		if err != nil {
			return fmt.Errorf("[db.ExecContext]: %w", err)
		}
		return nil
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO admin VALUES ($1, true,false)", vkID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (r *Repo) DeleteDocumentAdmin(ctx context.Context, vkID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE admin SET documents = $1 WHERE vk_id = $2", false, vkID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (r *Repo) GetDocumentsAdmins(ctx context.Context) ([]int, error) {
	var vkIDAdmins []int
	err := r.db.GetContext(ctx, &vkIDAdmins, "SELECT vk_id FROM admin WHERE documents = true")
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return vkIDAdmins, nil
}
