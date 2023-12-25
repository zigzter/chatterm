package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageParser(t *testing.T) {
	got := MessageParser(RawChatMessage)
	want := ChatMessage{
		Color:          "#00F8FF",
		DisplayName:    "zigzter",
		IsFirstMessage: false,
		IsMod:          true,
		IsVIP:          false,
		Message:        "merry crimus OFFLINECHAT",
		Timestamp:      "10:05",
	}
	assert.Equal(t, got, want)
}
