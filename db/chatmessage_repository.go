package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/zigzter/chatterm/types"
)

type ChatMessageRepository interface {
	Insert(msg types.InsertChat) error
	BuildQuery(input string) string
	Search(query string) ([]types.InsertChat, error)
	GetUsersMessages(username, channel string) ([]types.InsertChat, error)
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
	filters := []string{}
	for _, word := range queryWords {
		splitWord := strings.SplitN(word, ":", 2)
		if len(splitWord) > 1 {
			filter := splitWord[0]
			text := splitWord[1]
			if filter == "from" {
				filters = append(filters, "username:"+text)
			}
			if filter == "channel" {
				filters = append(filters, "channel:"+text)
			}
		} else {
			searchText += " " + word
		}
	}
	query += " WHERE chat_messages MATCH " + "'" + strings.Join(filters, " ") + searchText + "'"
	return query + " ORDER BY rank"
}

func (c *ChatMessageRepo) GetResults(query string) ([]types.InsertChat, error) {
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

// Search allows the user to search the chat DB for messages
func (c *ChatMessageRepo) Search(input string) ([]types.InsertChat, error) {
	query := c.BuildQuery(input)
	return c.GetResults(query)
}

// GetUsersMessages retrieves all messages from the given user in the given channel.
// We want to sort this by the timestamp, not the rank, therefore it's a separate method.
func (c *ChatMessageRepo) GetUsersMessages(username, channel string) ([]types.InsertChat, error) {
	query := fmt.Sprintf(
		`SELECT username, user_id, channel, content, timestamp
        FROM chat_messages
        WHERE chat_messages MATCH 'username:%s channel:%s'
        ORDER BY timestamp ASC`,
		username,
		channel,
	)
	return c.GetResults(query)
}
