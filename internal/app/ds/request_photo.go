package ds

type RequestPhoto struct {
	ID           int      `db:"id"`
	Year         int      `db:"year"`
	StudyProgram string   `db:"study_program"`
	Event        string   `db:"event"`
	IsEventNew   bool     `db:"is_event_new"`
	Description  bool     `db:"description"`
	MarkedPeople []string `db:"marked_people"`
	Attachment   string   `db:"attachment"`
	UserID       int      `db:"user_id"`
}
