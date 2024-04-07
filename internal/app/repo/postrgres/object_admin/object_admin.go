package object_admin

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Insert(ctx context.Context, adminID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO object_admin VALUES ($1, 0)", adminID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, fileID int, adminID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE object_admin SET object_id = $1 WHERE admin_id = $2", fileID, adminID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (r *Repo) Get(ctx context.Context, adminID int) (int, error) {
	var fileID int
	err := r.db.GetContext(ctx, &fileID, "SELECT object_id FROM object_admin WHERE admin_id = $1", adminID)
	if err != nil {
		return 0, fmt.Errorf("[db.GetContext]: %w", err)
	}
	return fileID, nil
}
