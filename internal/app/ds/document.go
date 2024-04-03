package ds

type Document struct {
	ID            int      `db:"id"`
	Title         string   `db:"title"`
	Author        string   `db:"author"`
	Year          int      `db:"year"`
	Category      string   `db:"category"`
	IsCategoryNew bool     `db:"is_category_new"`
	Hashtags      []string `db:"hashtags"`
	URL           string   `db:"url"`
	UserID        int      `db:"user_id"`
}
