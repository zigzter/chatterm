package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSubMessage(t *testing.T) {
	got := FormatSubMessage(ParsedSubMessage)
	want := "gimli subscribed for 12 months: Nobody tosses a dwarf\n"

	assert.Equal(t, got, want)
}
