package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

var (
	dbInstance *sql.DB
	once       sync.Once
	err        error
)

func OpenDB(path ...string) *sql.DB {
	once.Do(func() {
		// HACK: This should only be run by the call in main.go, so we're using an array to make it optional
		// TODO: Find a better way to do this
		dbInstance, err = sql.Open("sqlite3", fmt.Sprintf("file:%s/chatterm.db?cache=shared&mode=rwc", path[0]))
		if err != nil {
			log.Fatal(err)
		}
	})
	return dbInstance
}
