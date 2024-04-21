package teacher_albums

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repo инстанс репо для работы с альбомами студентов
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с параметрами для альбомов студентов
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetAllAlbumsOutput(ctx context.Context) (string, error) {
	sqlQuery := `
    SELECT 
        'ID: ' || CAST(id AS VARCHAR) AS id,
        'Название альбома: ' || teacher AS name
    FROM teacher_albums;`

	rows, err := r.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return "", fmt.Errorf("[db.QueryContext]: %w", err)
	}
	defer rows.Close()

	var output string

	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return "", fmt.Errorf("[rows.Scan]: %w", err)
		}
		output += fmt.Sprintf("%s\n%s\n\n", id, name)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("[rows.Err]: %w", err)
	}

	return output, nil
}
