package ds

import (
	"database/sql"
	"github.com/lib/pq"
)

type SearchDocument struct {
	ID         int            `db:"id"`
	Title      sql.NullString `db:"title"`
	Author     sql.NullString `db:"author"`
	Year       sql.NullInt64  `db:"year"`
	StartYear  sql.NullInt64  `db:"start_year"`
	EndYear    sql.NullInt64  `db:"end_year"`
	Categories pq.StringArray `db:"categories"`
	Hashtags   pq.StringArray `db:"hashtags"`
	UserID     int            `db:"user_id"`
}

type ParseSearchDocument struct {
	Title        string `db:"title"`
	Author       string `db:"author"`
	YearInterval string `db:"year_interval"`
	Categories   string `db:"categories"`
	Hashtags     string `db:"hashtags"`
}
