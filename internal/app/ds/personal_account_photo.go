package ds

type PersonalAccountPhoto struct {
	ID      int  `db:"id"`
	Pointer *int `db:"pointer"`
	UserID  int  `db:"user_id"`
}
