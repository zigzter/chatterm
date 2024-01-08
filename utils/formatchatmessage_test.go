package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatChatMessage(t *testing.T) {
	got := FormatChatMessage(ParsedChatMessage)
	want := "[10:05]\x1b[32m\U000f04e5\x1b[38;2;0;248;255mgandalf\x1b[39m\x1b[0m: All we have to decide is what to do with the time that is given to us.\x1b[0m\n"

	assert.Equal(t, got, want)
}
