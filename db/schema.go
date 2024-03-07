package db

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
	createChatMessagesTable := `
        CREATE VIRTUAL TABLE IF NOT EXISTS chat_messages USING fts5(
            username,
            user_id,
            channel,
            content,
            timestamp
        );`

	createUserIdMapTable := `
        CREATE TABLE IF NOT EXISTS userid_map (
            username TEXT PRIMARY KEY,
            user_id TEXT
        );`

	_, err := db.Exec(createChatMessagesTable)
	if err != nil {
		log.Fatal("Cannot create chat_messages table:", err)
	}

	_, err = db.Exec(createUserIdMapTable)
	if err != nil {
		log.Fatal("Cannot create userid_map table:", err)
	}
}
