package db

import (
	"database/sql"
	"log"
	"sync"
)

var (
	dbInstance *sql.DB
	once       sync.Once
	err        error
)

func OpenDB() *sql.DB {
	once.Do(func() {
		dbInstance, err = sql.Open("sqlite3", "file:chatterm.db?cache=shared&mode=rwc")
		if err != nil {
			log.Fatal(err)
		}
	})
	return dbInstance
}
