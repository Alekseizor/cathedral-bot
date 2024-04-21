package ds

type StudentAlbum struct {
	ID           int    `db:"id"`
	Year         int    `db:"year"`
	StudyProgram string `db:"study_program"`
	Event        string `db:"event"`
	URL          string `db:"url"`
}
