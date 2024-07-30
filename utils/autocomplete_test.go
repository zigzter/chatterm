package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ac *Trie
var names = []string{"legolas", "gimli", "gandalf", "aragorn", "saruman", "sauron", "bilbo", "samwise", "merry", "pippin", "boromir", "elrond", "galadriel"}

func TestMain(m *testing.M) {
	ac = NewTrie()
	code := m.Run()
	os.Exit(code)
}

func TestAutocomplete(t *testing.T) {
	t.Run("Test populate", func(t *testing.T) {
		ac.Populate(names)
		got := len(ac.Root.children)
		want := 8
		assert.Equal(t, want, got)
	})

	t.Run("Test search", func(t *testing.T) {
		ac.Populate(names)
		testCases := []struct {
			prefix string
			want   []string
		}{
			{"leg", []string{"legolas"}},
			{"g", []string{"gimli", "gandalf", "galadriel"}},
			{"sa", []string{"saruman", "sauron", "samwise"}},
			{"a", []string{"aragorn"}},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Prefix: %s", tc.prefix), func(t *testing.T) {
				got := ac.Search(tc.prefix)
				assert.ElementsMatch(t, tc.want, got)
			})
		}
	})

	t.Run("Test suggestion update", func(t *testing.T) {
		ac.Search("sa")
		got := ac.UpdateSuggestion("sa")
		want := "saruman"
		assert.Equal(t, want, got)
		got = ac.UpdateSuggestion("sa")
		want = "sauron"
		assert.Equal(t, want, got)
	})
}
