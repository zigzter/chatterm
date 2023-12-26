package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigzter/chatterm/types"
)

func TestMessageParser(t *testing.T) {
	got := MessageParser(RawChatMessage)
	want := types.ChatMessage{
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
