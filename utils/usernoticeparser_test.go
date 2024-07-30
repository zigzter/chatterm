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
			Color:       "#008000",
			Message:     "Nobody tosses a dwarf",
			Months:      "12",
			Streak:      "1",
			Timestamp:   "11:28",
		}

		assert.Equal(t, want, got)
	})

	t.Run("Test AnnouncementParser", func(t *testing.T) {
		got := UsernoticeParser(RawAnnouncementMessage)
		want := types.AnnouncementMessage{
			DisplayName: "gandalf",
			Color:       "#54BC75",
			Message:     "You. Shall. Not. Pass!",
			Timestamp:   "17:35",
		}

		assert.Equal(t, want, got)
	})

	t.Run("Test RaidParser", func(t *testing.T) {
		got := UsernoticeParser(RawRaidMessage)
		want := types.RaidMessage{
			DisplayName: "gandalf",
			Color:       "#9ACD32",
			ViewerCount: "15",
			Timestamp:   "16:36",
		}

		assert.Equal(t, want, got)
	})
}
