package test_generic

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"
	"github.com/stretchr/testify/assert"
)

func TestGenericCommentStateNextToken(t *testing.T) {
	state := generic.NewGenericCommentState()

	scanner := io.NewStringScanner("# Comment \r# Comment ")
	token := state.NextToken(scanner, nil)
	assert.Equal(t, "# Comment ", token.Value())
	assert.Equal(t, tokenizers.Comment, token.Type())

	scanner = io.NewStringScanner("# Comment \n# Comment ")
	token = state.NextToken(scanner, nil)
	assert.Equal(t, "# Comment ", token.Value())
	assert.Equal(t, tokenizers.Comment, token.Type())
}
