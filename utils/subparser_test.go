package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubParser(t *testing.T) {
	got := SubParser(RegResubMessage)
	want := SubMessage{
		DisplayName: "vexthorne",
		Message:     "merry Christmas",
		Months:      "12",
		Streak:      "1",
	}

	assert.Equal(t, got, want)
}
