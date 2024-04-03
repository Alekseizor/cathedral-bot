package ds

type RequestDocument struct {
	ID            int      `db:"id"`
	Title         string   `db:"title"`
	Author        string   `db:"author"`
	Year          int      `db:"year"`
	Category      string   `db:"category"`
	IsCategoryNew bool     `db:"is_category_new"`
	Description   bool     `db:"description"`
	Hashtags      []string `db:"hashtags"`
	Attachment    string   `db:"attachment"`
	UserID        int      `db:"user_id"`
}
