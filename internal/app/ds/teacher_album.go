package ds

type TeacherAlbum struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	URL  string `db:"url"`
}
