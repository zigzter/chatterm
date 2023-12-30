package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSubMessage(t *testing.T) {
	got := FormatSubMessage(ParsedSubMessage)
	want := "vexthorne subscribed for 12 months: merry Christmas\n"

	assert.Equal(t, got, want)
}
