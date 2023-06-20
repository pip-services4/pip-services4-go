package test_mustache_tokenizers

import (
	"testing"

	mtokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache/tokenizers"
	test_tokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/test/tokenizers"
	tokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
)

func TestMustacheTokenizerTemplate1(t *testing.T) {
	tokenString := "Hello, {{ Name }}!"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Special, "Hello, ", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "{{", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Name", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "}}", 0, 0),
		tokenizers.NewToken(tokenizers.Special, "!", 0, 0),
	}

	tokenizer := mtokenizers.NewMustacheTokenizer()
	tokenizer.SetSkipEof(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}
