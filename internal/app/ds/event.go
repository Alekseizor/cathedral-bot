package ds

type Event struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
