package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigzter/chatterm/types"
)

func TestMessageParser(t *testing.T) {
	got := MessageParser(RawChatMessage)
	want := types.ChatMessage{
		Color:           "#00F8FF",
		DisplayName:     "gandalf",
		IsFirstMessage:  false,
		ChannelUserType: "moderator",
		Message:         "All we have to decide is what to do with the time that is given to us.",
		Timestamp:       "10:05",
		UserId:          "20816785",
	}
	assert.Equal(t, want, got)
}
