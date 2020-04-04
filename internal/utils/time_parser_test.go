package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTime(t *testing.T) {
	ans := ParseTime("7:00")
	assert.Equal(t, 420, ans)

	ans = ParseTime("07:00")
	assert.Equal(t, 420, ans)

	ans = ParseTime("10:00")
	assert.Equal(t, 600, ans)

	ans = ParseTime("60")
	assert.Equal(t, 60, ans)

	ans = ParseTime("100")
	assert.Equal(t, 100, ans)

	assert.Panics(t, func() { ParseTime("-10") })
	assert.Panics(t, func() { ParseTime(":50") })
	assert.Panics(t, func() { ParseTime("-7:50") })
}
