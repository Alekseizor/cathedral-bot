package documents

import (
	"context"
	"fmt"

	"github.com/lib/pq"
)

// EditTitle изменяет название документа по ID заявки
func (r *Repo) EditTitle(ctx context.Context, name string, reqDocID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE documents SET title = $1 WHERE id = $2", name, reqDocID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// EditAuthor изменяет ФИО автора документа по ID заявки
func (r *Repo) EditAuthor(ctx context.Context, author string, reqDocID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE documents SET author = $1 WHERE id = $2", author, reqDocID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// EditYear изменяет год создания документа по ID заявки
func (r *Repo) EditYear(ctx context.Context, year, reqDocID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE documents SET year = $1 WHERE id = $2", year, reqDocID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// EditCategory изменяет категорию документа
func (r *Repo) EditCategory(ctx context.Context, categoryNumber, reqDocID int) error {
	var name string
	err := r.db.Get(&name, "SELECT name FROM categories WHERE id = $1", categoryNumber)
	if err != nil {
		return fmt.Errorf("[db.Get]: %w", err)
	}

	_, err = r.db.ExecContext(ctx, "UPDATE documents SET category = $1 WHERE id = $2", name, reqDocID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// EditDescription изменяет описание документа по ID заявки
func (r *Repo) EditDescription(ctx context.Context, description string, reqDocID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE documents SET description = $1 WHERE id = $2", description, reqDocID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

// EditHashtags изменяет хештеги документа по ID заявки
func (r *Repo) EditHashtags(ctx context.Context, hashtags []string, reqDocID int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE documents SET hashtags = $1 WHERE id = $2", pq.Array(hashtags), reqDocID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}
