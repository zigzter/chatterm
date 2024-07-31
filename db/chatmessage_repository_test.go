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
var repo ChatMessageRepository
var message1 = types.InsertChat{
	Username: "gandalf", UserID: "1", Channel: "MiddleEarth",
	Content: "alls well that ends better", Timestamp: "11:11",
}
var message2 = types.InsertChat{
	Username: "gimli", UserID: "2", Channel: "MiddleEarth",
	Content: "nobody tosses a dwarf", Timestamp: "11:12",
}
var message3 = types.InsertChat{
	Username: "gandalf", UserID: "1", Channel: "MiddleEarth",
	Content:   "a wizard is never late, nor is he early, he arrives precisely when he means to",
	Timestamp: "11:15",
}

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
		repo.Insert(message1)
		repo.Insert(message2)
		repo.Insert(message3)
	})

	t.Run("Test Repo Query Build", func(t *testing.T) {
		base := "SELECT username, user_id, channel, content, timestamp FROM chat_messages"
		got := repo.BuildQuery("alls well that ends better from:gandalf")
		want := base + " WHERE chat_messages MATCH 'username:gandalf alls well that ends better' ORDER BY rank"
		assert.Equal(t, want, got)
	})

	t.Run("Test Get User Messages", func(t *testing.T) {
		got, err := repo.GetUsersMessages("gandalf", "MiddleEarth")
		want := []types.InsertChat{message1, message3}
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Test Repo Search", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  []types.InsertChat
		}{
			{"Test username search", "from:gandalf", []types.InsertChat{message1, message3}},
			{"Test text search", "tosses", []types.InsertChat{message2}},
			{"Test channel and username search", "from:gimli channel:MiddleEarth", []types.InsertChat{message2}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := repo.Search(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
