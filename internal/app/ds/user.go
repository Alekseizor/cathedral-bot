package ds

// State - cтруктура для работы с таблицей state
type State struct {
	VkID  int    `db:"vk_id,omitempty"`
	Title string `db:"title,omitempty"`
}
