package db

import (
	"database/sql"
	"log"

	"github.com/zigzter/chatterm/types"
)

type ChatMessageRepository interface {
	Insert(msg types.InsertChat) error
	Search(query string) ([]ChatMessageRepo, error)
}

type ChatMessageRepo struct {
	db *sql.DB
}

func NewChatMessageRepository(db *sql.DB) *ChatMessageRepo {
	if db == nil {
		panic("chat message repository: missing db")
	}
	return &ChatMessageRepo{db: db}
}

func (c *ChatMessageRepo) Insert(msg types.InsertChat) {
	sqlStatement := `
        INSERT INTO chat_messages (username, user_id, channel, content, timestamp) 
        VALUES (?, ?, ?, ?, ?)`
	_, err := c.db.Exec(
		sqlStatement,
		msg.Username,
		msg.UserID,
		msg.Channel,
		msg.Content,
		msg.Timestamp,
	)
	if err != nil {
		log.Fatal("Cannot insert chat message:", err)
	}
}

func (c *ChatMessageRepo) Search(query string) ([]string, error) {
	stmt := `SELECT username, user_id, channel, content, timestamp
        FROM chat_messages WHERE chat_messages MATCH ?`
	rows, err := c.db.Query(stmt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var m string
		if err := rows.Scan(&m); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
