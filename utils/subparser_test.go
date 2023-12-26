package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigzter/chatterm/types"
)

func TestSubParser(t *testing.T) {
	got := SubParser(RegResubMessage)
	want := types.SubMessage{
		DisplayName: "vexthorne",
		Message:     "merry Christmas",
		Months:      "12",
		Streak:      "1",
	}

	assert.Equal(t, got, want)
}
