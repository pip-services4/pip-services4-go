package test_generic

import (
	"testing"

	test_tokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/test/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"
)

func TestGenericTokenizerExpression(t *testing.T) {
	tokenString := "A+B/123 - \t 'xyz'\n <>-10.11# This is a comment"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Word, "A", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "+", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "B", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "/", 0, 0),
		tokenizers.NewToken(tokenizers.Integer, "123", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "-", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " \t ", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "'xyz'", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, "\n ", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "<>", 0, 0),
		tokenizers.NewToken(tokenizers.Float, "-10.11", 0, 0),
		tokenizers.NewToken(tokenizers.Comment, "# This is a comment", 0, 0),
		tokenizers.NewToken(tokenizers.Eof, "", 0, 0),
	}

	tokenizer := generic.NewGenericTokenizer()
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}

func TestGenericTokenizerQuoteToken(t *testing.T) {
	tokenString := "A'xyz'\"abc\ndeg\" 'jkl\"def'"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Word, "A", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "xyz", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "abc\ndeg", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "jkl\"def", 0, 0),
	}

	tokenizer := generic.NewGenericTokenizer()
	tokenizer.SetSkipEof(true)
	tokenizer.SetDecodeStrings(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}

func TestGenericTokenizerWordToken(t *testing.T) {
	tokenString := "A'xyz'Ebf_2\n2x_2"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Word, "A", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "xyz", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Ebf_2", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, "\n", 0, 0),
		tokenizers.NewToken(tokenizers.Integer, "2", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "x_2", 0, 0),
	}

	tokenizer := generic.NewGenericTokenizer()
	tokenizer.SetSkipEof(true)
	tokenizer.SetDecodeStrings(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}

func TestGenericTokenizerNumberToken(t *testing.T) {
	tokenString := "123-321 .543-.76-. -123.456"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Integer, "123", 0, 0),
		tokenizers.NewToken(tokenizers.Integer, "-321", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Float, ".543", 0, 0),
		tokenizers.NewToken(tokenizers.Float, "-.76", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "-", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ".", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Float, "-123.456", 0, 0),
	}

	tokenizer := generic.NewGenericTokenizer()
	tokenizer.SetSkipEof(true)
	tokenizer.SetDecodeStrings(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}

func TestGenericTokenizerWrongToken(t *testing.T) {
	tokenString := "1>2"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Integer, "1", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ">", 0, 0),
		tokenizers.NewToken(tokenizers.Integer, "2", 0, 0),
	}

	tokenizer := generic.NewGenericTokenizer()
	tokenizer.SetSkipEof(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}
