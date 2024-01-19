package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatters(t *testing.T) {
	t.Run("Test FormatChatMessage", func(t *testing.T) {
		got := FormatChatMessage(ParsedChatMessage, 100)
		want := "[10:05]\x1b[32m\U000f04e5\x1b[38;2;0;248;255mgandalf\x1b[39m\x1b[0m: All we have to decide is what to do with the time that is given to us.\x1b[0m\n"

		assert.Equal(t, got, want)
	})

	t.Run("Test FormatSubMessage", func(t *testing.T) {
		got := FormatSubMessage(ParsedSubMessage, 100)
		want := "gimli subscribed for 12 months: Nobody tosses a dwarf\n"

		assert.Equal(t, got, want)
	})
}
