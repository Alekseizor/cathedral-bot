package ds

import (
	"github.com/lib/pq"
)

type RequestPhoto struct {
	ID           int            `db:"id"`
	Year         *string        `db:"year"`
	StudyProgram *string        `db:"study_program"`
	Event        *string        `db:"event"`
	IsEventNew   bool           `db:"is_event_new"`
	Description  *string        `db:"description"`
	CountPeople  int            `db:"count_people"`
	MarkedPerson int            `db:"marked_person"`
	MarkedPeople pq.StringArray `db:"marked_people"`
	Teachers     pq.StringArray `db:"teachers"`
	Pointer      int            `db:"pointer"`
	Attachment   string         `db:"attachment"`
	Attachments  pq.StringArray `db:"attachments"`
	URL          string         `db:"url"`
	URLS         pq.StringArray `db:"urls"`
	UserID       int            `db:"user_id"`
	Status       int            `db:"status"`
}
