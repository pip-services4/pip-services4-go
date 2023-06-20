package test_random

import (
	"strings"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/random"
	"github.com/stretchr/testify/assert"
)

func TestPhrase(t *testing.T) {
	assert.True(t, random.Text.Phrase(-1, -2) == "")
	assert.True(t, random.Text.Phrase(-1, 0) == "")
	assert.True(t, random.Text.Phrase(-2, -1) == "")

	text := random.Text.Phrase(4, 4)
	assert.True(t, len(text) >= 4 && len(text) <= 10)
	text = random.Text.Phrase(4, 10)
	assert.True(t, len(text) >= 4)
}

func TestFullName(t *testing.T) {
	text := random.Text.FullName()
	assert.True(t, strings.Index(text, " ") >= 0)
}

func TestPhone(t *testing.T) {
	text := random.Text.Phone()
	assert.True(t, strings.Index(text, "(") >= 0)
	assert.True(t, strings.Index(text, ")") >= 0)
	assert.True(t, strings.Index(text, "-") >= 0)
}

func TestEmail(t *testing.T) {
	text := random.Text.Email()
	assert.True(t, strings.Index(text, "@") >= 0)
	assert.True(t, strings.Index(text, ".com") >= 0)
}
