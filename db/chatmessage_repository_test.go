package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/zigzter/chatterm/types"
)

var db *sql.DB
var repo *ChatMessageRepo

func TestMain(m *testing.M) {
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS chat_messages USING fts5(
            username,
            user_id,
            channel,
            content,
            timestamp
        )`)
	if err != nil {
		panic(err)
	}
	repo = NewChatMessageRepository(db)
	code := m.Run()
	db.Close()
	os.Exit(code)
}

func TestChatMessageRepo(t *testing.T) {
	t.Run("Test Repo Insertion", func(t *testing.T) {
		repo.Insert(types.InsertChat{Username: "gandalf", Content: "some text"})
	})

	t.Run("Test Repo Query Build", func(t *testing.T) {
		got := repo.BuildQuery("all's well that ends better")
		want := "SELECT username, user_id, channel, content, timestamp FROM chat_messages WHERE content MATCH  all's well that ends better"
		assert.Equal(t, got, want)
	})
}
