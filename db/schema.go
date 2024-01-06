package db

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
	createChatMessagesTable := `
        CREATE TABLE IF NOT EXISTS chat_messages (
            id INTEGER PRIMARY KEY,
            username TEXT,
            user_id TEXT,
            content TEXT,
            timestamp DATETIME
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
