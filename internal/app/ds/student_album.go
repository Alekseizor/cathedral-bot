package ds

type StudentAlbum struct {
	ID           int    `db:"id"`
	Year         int    `db:"year"`
	StudyProgram string `db:"study_program"`
	Event        string `db:"event"`
	Description  string `db:"description"`
	URL          string `db:"url"`
	VkID         int    `db:"vk_id"`
}
