package ds

type Photo struct {
	ID           int      `db:"id"`
	Year         int      `db:"year"`
	StudyProgram string   `db:"study_program"`
	Event        string   `db:"event"`
	IsEventNew   bool     `db:"is_event_new"`
	Description  bool     `db:"description"`
	MarkedPeople []string `db:"marked_people"`
	URL          string   `db:"url"`
}
