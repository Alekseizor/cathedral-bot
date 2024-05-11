package ds

type SearchAlbum struct {
	ID           int     `db:"id"`
	Year         *int    `db:"year"`
	StudyProgram *string `db:"study_program"`
	Event        *string `db:"event"`
	Surname      *string `db:"surname"`
	Pointer      *int    `db:"pointer"`
	UserID       int     `db:"user_id"`
}
