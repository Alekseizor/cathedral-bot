package search_document

import "github.com/jmoiron/sqlx"

// Repo инстанс репо для работы с параметрами для поиска документа
type Repo struct {
	db *sqlx.DB
}

// New - создаем новое объект репо для работы с параметрами для поиска документа
func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}
