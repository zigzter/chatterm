package db

import (
	"database/sql"
	"log"
)

func OpenDB() *sql.DB {
	db, err := sql.Open("sqlite3", "file:chatterm.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
