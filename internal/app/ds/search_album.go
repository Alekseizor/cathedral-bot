package ds

type SearchAlbum struct {
	ID           int     `db:"id"`
	Year         *int    `db:"year"`
	StudyProgram *string `db:"study_program"`
	Event        *string `db:"event"`
	Teacher      *string `db:"teacher"`
	UserID       int     `db:"user_id"`
}