package test_generic

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/csv"
	test_tokenizers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/test/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
)

func TestCsvTokenizerWithDefaultParameters(t *testing.T) {
	tokenString := "\n\r\"John \"\"Da Man\"\"\",Repici,120 Jefferson St.,Riverside, NJ,08075\r\n" +
		"Stephen,Tyler,\"7452 Terrace \"\"At the Plaza\"\" road\",SomeTown,SD, 91234\r" +
		",Blankman,,SomeTown, SD, 00298\n"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Eol, "\n\r", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "\"John \"\"Da Man\"\"\"", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Repici", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "120 Jefferson St.", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Riverside", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, " NJ", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "08075", 0, 0),
		tokenizers.NewToken(tokenizers.Eol, "\r\n", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Stephen", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Tyler", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "\"7452 Terrace \"\"At the Plaza\"\" road\"", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "SomeTown", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "SD", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, " 91234", 0, 0),
		tokenizers.NewToken(tokenizers.Eol, "\r", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Blankman", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "SomeTown", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, " SD", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, ",", 0, 0),
		tokenizers.NewToken(tokenizers.Word, " 00298", 0, 0),
		tokenizers.NewToken(tokenizers.Eol, "\n", 0, 0),
	}

	tokenizer := csv.NewCsvTokenizer()
	tokenizer.SetSkipEof(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}

func TestCsvTokenizerWithOverridenParameters(t *testing.T) {
	tokenString := "\n\r'John, ''Da Man'''\tRepici\t120 Jefferson St.\tRiverside\t NJ\t08075\r\n" +
		"Stephen\t\"Tyler\"\t'7452 \t\nTerrace ''At the Plaza'' road'\tSomeTown\tSD\t 91234\r" +
		"\tBlankman\t\tSomeTown 'xxx\t'\t SD\t 00298\n"
	expectedTokens := []*tokenizers.Token{
		tokenizers.NewToken(tokenizers.Eol, "\n\r", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "'John, ''Da Man'''", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Repici", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "120 Jefferson St.", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Riverside", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, " NJ", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "08075", 0, 0),
		tokenizers.NewToken(tokenizers.Eol, "\r\n", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Stephen", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "\"Tyler\"", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "'7452 \t\nTerrace ''At the Plaza'' road'", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "SomeTown", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "SD", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, " 91234", 0, 0),
		tokenizers.NewToken(tokenizers.Eol, "\r", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "Blankman", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, "SomeTown ", 0, 0),
		tokenizers.NewToken(tokenizers.Quoted, "'xxx\t'", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, " SD", 0, 0),
		tokenizers.NewToken(tokenizers.Symbol, "\t", 0, 0),
		tokenizers.NewToken(tokenizers.Word, " 00298", 0, 0),
		tokenizers.NewToken(tokenizers.Eol, "\n", 0, 0),
	}

	tokenizer := csv.NewCsvTokenizer()
	tokenizer.SetFieldSeparators([]rune{'\t'})
	tokenizer.SetQuoteSymbols([]rune{'\'', '"'})
	tokenizer.SetEndOfLine("\n")
	tokenizer.SetSkipEof(true)
	tokenList := tokenizer.TokenizeBuffer(tokenString)

	test_tokenizers.AssertAreEqualsTokenLists(t, expectedTokens, tokenList)
}
