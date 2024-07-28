package db

import (
	"database/sql"
	"strings"

	"github.com/zigzter/chatterm/types"
)

type ChatMessageRepository interface {
	Insert(msg types.InsertChat) error
	BuildQuery(input string) string
	Search(query string) ([]types.InsertChat, error)
}

type ChatMessageRepo struct {
	db *sql.DB
}

func NewChatMessageRepository(db *sql.DB) ChatMessageRepository {
	if db == nil {
		panic("chat message repository: missing db")
	}
	return &ChatMessageRepo{db: db}
}

func (c *ChatMessageRepo) Insert(msg types.InsertChat) error {
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
	return err
}

// BuildQuery takes a user's search input and turns it into a SELECT statement
func (c *ChatMessageRepo) BuildQuery(input string) string {
	query := "SELECT username, user_id, channel, content, timestamp FROM chat_messages"
	queryWords := strings.Split(input, " ")
	searchText := ""
	userFilter := ""
	for _, word := range queryWords {
		splitWord := strings.SplitN(word, ":", 2)
		if len(splitWord) > 1 {
			filter := splitWord[0]
			text := splitWord[1]
			if filter == "from" {
				userFilter = "username:" + text
			}
		} else {
			searchText += " " + word
		}
	}
	query += " WHERE chat_messages MATCH " + "'" + userFilter + searchText + "'"
	return query
}

func (c *ChatMessageRepo) Search(input string) ([]types.InsertChat, error) {
	query := c.BuildQuery(input)
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []types.InsertChat
	for rows.Next() {
		var m types.InsertChat
		if err := rows.Scan(&m.Username, &m.UserID, &m.Channel, &m.Content, &m.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if err = rows.Err(); err != nil {
		return messages, err
	}
	return messages, nil
}
