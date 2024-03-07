package types

type InsertChat struct {
	Username  string `db:"username"`
	UserID    string `db:"user_id"`
	Channel   string `db:"channel"`
	Content   string `db:"content"`
	Timestamp string `db:"timestamp"`
}
