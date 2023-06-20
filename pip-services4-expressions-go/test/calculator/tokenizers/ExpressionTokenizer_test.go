package test_calculator_tokenizers

import (
	"testing"

	ctokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/tokenizers"
	test_tokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/test/tokenizers"
	tokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/stretchr/testify/assert"
)

func TestExpressionTokenizerQuoteToken(t *testing.T) {
	tokenString := "A'xyz'\"abc\ndeg\" 'jkl\"def'\"ab\"\"de\"'df''er'"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Word, "A", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "xyz", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "abc\ndeg", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "jkl\"def", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "ab\"de", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "df'er", 0, 0),
	}

	tokenizer := ctokenizers.NewExpressionTokenizer()
	tokenizer.SetSkipEof(true)
	tokenizer.SetDecodeStrings(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}

func TestExpressionTokenizerWordToken(t *testing.T) {
	tokenString := "A'xyz'Ebf_2\n2_2"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Word, "A", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "xyz", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Ebf_2", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, "\n", 0, 0),
		tokenizers.NewToken(tokenizers.Integer, "2", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "_2", 0, 0),
	}

	tokenizer := ctokenizers.NewExpressionTokenizer()
	tokenizer.SetSkipEof(true)
	tokenizer.SetDecodeStrings(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}

func TestExpressionTokenizerNumberToken(t *testing.T) {
	tokenString := "123-321 .543-.76-. 123.456 123e45 543.11E+43 1e 3E-"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Integer, "123", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "-", 0, 0),
		tokenizers.NewToken(tokenizers.Integer, "321", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Float, ".543", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "-", 0, 0),
		tokenizers.NewToken(tokenizers.Float, ".76", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "-", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ".", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Float, "123.456", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Float, "123e45", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Float, "543.11E+43", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Integer, "1", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "e", 0, 0),
		tokenizers.NewToken(tokenizers.Whitespace, " ", 0, 0),
		tokenizers.NewToken(tokenizers.Integer, "3", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "E", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "-", 0, 0),
	}

	tokenizer := ctokenizers.NewExpressionTokenizer()
	tokenizer.SetSkipEof(true)
	tokenizer.SetDecodeStrings(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}

func TestExpressionTokenizerExpressionToken(t *testing.T) {
	tokenString := "A + b / (3 - Max(-123, 1)*2)"

	tokenizer := ctokenizers.NewExpressionTokenizer()
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	assert.Len(t, tokenList, 25)
}
