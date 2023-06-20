package test_generic

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"
	"github.com/stretchr/testify/assert"
)

func TestCCommentStateNextToken(t *testing.T) {
	state := generic.NewCCommentState()

	scanner := io.NewStringScanner("// Comment \n Comment ")
	failed := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				failed = true
			}
		}()
		state.NextToken(scanner, nil)
	}()
	assert.True(t, failed)

	scanner = io.NewStringScanner("/* Comment \n Comment */#")
	token := state.NextToken(scanner, nil)
	assert.Equal(t, "/* Comment \n Comment */", token.Value())
	assert.Equal(t, tokenizers.Comment, token.Type())
}
