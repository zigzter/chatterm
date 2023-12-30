package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatChatMessage(t *testing.T) {
	got := FormatChatMessage(ParsedChatMessage)
	want := "[10:05]\x1b[32m\U000f04e5\x1b[38;2;0;248;255mzigzter\x1b[39m\x1b[0m: merry crimus OFFLINECHAT\x1b[0m\n"

	assert.Equal(t, got, want)
}
