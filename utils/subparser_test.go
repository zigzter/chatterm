package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigzter/chatterm/types"
)

func TestSubParser(t *testing.T) {
	got := SubParser(RegResubMessage)
	want := types.SubMessage{
		DisplayName: "gimli",
		Message:     "Nobody tosses a dwarf",
		Months:      "12",
		Streak:      "1",
	}

	assert.Equal(t, got, want)
}
