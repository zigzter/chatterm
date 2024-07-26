package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

var filterMap = map[string]string{
	"from": "username",
	// TODO: Add before, after
}

func (c *ChatMessageRepo) BuildQuery(input string) string {
	query := "SELECT username, user_id, channel, content, timestamp FROM chat_messages"
	queryWords := strings.Split(input, " ")
	filters := map[string]string{}
	searchText := ""
	for _, word := range queryWords {
		splitWord := strings.SplitN(word, ":", 2)
		if len(splitWord) > 1 {
			filters[splitWord[0]] = splitWord[1]
		} else {
			searchText += " " + word
		}
	}
	query += " WHERE content MATCH " + searchText
	for filter, value := range filters {
		query += fmt.Sprintf(" WHERE %s MATCH %s", filterMap[filter], value)
	}
	return query
}

func (c *ChatMessageRepo) Search(input string) ([]string, error) {
	query := c.BuildQuery(input)
	rows, err := c.db.Query(query)
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
