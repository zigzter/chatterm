package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigzter/chatterm/types"
)

func TestUsernoticeParser(t *testing.T) {
	t.Run("Test SubParser", func(t *testing.T) {
		got := UsernoticeParser(RegResubMessage)
		want := types.SubMessage{
			DisplayName: "gimli",
			Message:     "Nobody tosses a dwarf",
			Months:      "12",
			Streak:      "1",
		}

		assert.Equal(t, want, got)
	})

	t.Run("Test AnnouncementParser", func(t *testing.T) {
		got := UsernoticeParser(RawAnnouncementMessage)
		want := types.AnnouncementMessage{
			DisplayName: "gandalf",
			Message:     "You. Shall. Not. Pass!",
		}

		assert.Equal(t, want, got)
	})

	t.Run("Test RaidParser", func(t *testing.T) {
		got := UsernoticeParser(RawRaidMessage)
		want := types.RaidMessage{
			DisplayName: "gandalf",
			ViewerCount: "15",
		}

		assert.Equal(t, want, got)
	})
}
