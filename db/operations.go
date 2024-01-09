package db

import (
	"database/sql"
	"log"

	"github.com/zigzter/chatterm/types"
)

func InsertChatMessage(db *sql.DB, msg types.InsertChat) {
	sqlStatement := `
        INSERT INTO chat_messages (username, user_id, content, timestamp) 
        VALUES (?, ?, ?, ?)`
	_, err := db.Exec(sqlStatement, msg.Username, msg.UserId, msg.Content, msg.Timestamp)
	if err != nil {
		log.Fatal("Cannot insert chat message:", err)
	}
}

func InsertUserMap(db *sql.DB, username, userID string) {
	sqlStatement := `
        INSERT INTO userid_map (username, user_id) VALUES (?, ?)
        ON CONFLICT(username) DO UPDATE SET user_id = excluded.user_id;`
	_, err := db.Exec(sqlStatement, username, userID)
	if err != nil {
		log.Fatal("Cannot insert or update user map:", err)
	}
}

func QueryChatMessages(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Cannot query chat messages:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, userID, content, timestamp string
		err = rows.Scan(&id, &username, &userID, &content, &timestamp)
		if err != nil {
			log.Fatal("Scan failed:", err)
		}
		// TODO: figure out best way to format and return this
		log.Printf("Chat Message: %d, %s, %s, %s, %s\n", id, username, userID, content, timestamp)
	}
}

func GetUserId(db *sql.DB, username string) (string, error) {
	var userId string
	sqlStatement := "SELECT user_id FROM userid_map WHERE username = ? LIMIT 1"
	row := db.QueryRow(sqlStatement, username)

	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return userId, nil
}
