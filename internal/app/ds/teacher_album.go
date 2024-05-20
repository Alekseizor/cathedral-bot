package ds

type TeacherAlbum struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	URL         string `db:"url"`
	VkID        int    `db:"vk_id"`
}
