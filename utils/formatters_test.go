package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatters(t *testing.T) {
	t.Run("Test FormatChatMessage", func(t *testing.T) {
		got := FormatChatMessage(ParsedChatMessage, 100)
		want := "[10:05][\U000f04e5]gandalf: All we have to decide is what to do with the time that is given to us.\n"
		assert.Equal(t, want, got)
	})

	t.Run("Test FormatSubMessage", func(t *testing.T) {
		got := FormatSubMessage(ParsedSubMessage, 100)
		want := "gimli subscribed for 12 months: Nobody tosses a dwarf\n"
		assert.Equal(t, want, got)
	})
}
