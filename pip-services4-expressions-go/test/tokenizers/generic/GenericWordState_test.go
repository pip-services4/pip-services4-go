package test_generic

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"
	"github.com/stretchr/testify/assert"
)

func TestGenericWordStateNextToken(t *testing.T) {
	state := generic.NewGenericWordState()

	scanner := io.NewStringScanner("AB_CD=")
	token := state.NextToken(scanner, nil)
	assert.Equal(t, "AB_CD", token.Value())
	assert.Equal(t, tokenizers.Word, token.Type())
}
