package ds

import "github.com/lib/pq"

const (
	StatusInProgress = iota
	StatusUserConfirmed
	StatusAdminWorking
	StatusAdminDeclined
)

type RequestDocument struct {
	ID            int            `db:"id"`
	Title         string         `db:"title"`
	Author        string         `db:"author"`
	Year          int            `db:"year"`
	Category      string         `db:"category"`
	IsCategoryNew bool           `db:"is_category_new"`
	Description   string         `db:"description"`
	Hashtags      pq.StringArray `db:"hashtags"`
	Attachment    string         `db:"attachment"`
	UserID        int            `db:"user_id"`
	Status        int            `db:"status"`
}

type Document struct {
	ID          int            `db:"id"`
	Title       *string        `db:"title"`
	Author      *string        `db:"author"`
	Year        *int           `db:"year"`
	Category    *string        `db:"category"`
	Description *string        `db:"description"`
	Hashtags    pq.StringArray `db:"hashtags"`
	Attachment  *string        `db:"attachment"`
	UserID      int            `db:"user_id"`
}
