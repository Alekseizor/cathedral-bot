package ds

type TeacherAlbum struct {
	ID      int    `db:"id"`
	Teacher string `db:"teacher"`
	URL     string `db:"url"`
}
