package documents

import "github.com/jmoiron/sqlx"

// Repo инстанс репо для работы с загруженными документами
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с загруженными документами
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}
