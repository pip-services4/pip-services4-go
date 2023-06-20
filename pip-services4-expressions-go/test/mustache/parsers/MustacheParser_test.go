package test_mustache_parsers

import (
	"testing"

	mparsers "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache/parsers"
	"github.com/stretchr/testify/assert"
)

func TestMustacheParserLexicalAnalysis(t *testing.T) {
	parser := mparsers.NewMustacheParser()
	err := parser.SetTemplate("Hello, {{{NAME}}}{{ #if ESCLAMATION }}!{{/if}}{{{^ESCLAMATION}}}.{{{/ESCLAMATION}}}")
	assert.Nil(t, err)

	expectedTokens := []*mparsers.MustacheToken{
		mparsers.NewMustacheToken(mparsers.TokenValue, "Hello, ", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenEscapedVariable, "NAME", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenSection, "ESCLAMATION", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenValue, "!", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenSectionEnd, "", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenInvertedSection, "ESCLAMATION", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenValue, ".", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenSectionEnd, "ESCLAMATION", 0, 0),
	}

	tokens := parser.InitialTokens()
	assert.Equal(t, len(expectedTokens), len(tokens))

	for i := 0; i < len(tokens); i++ {
		assert.Equal(t, expectedTokens[i].Type(), tokens[i].Type())
		assert.Equal(t, expectedTokens[i].Value(), tokens[i].Value())
	}
}

func TestMustacheParserSyntaxAnalysis(t *testing.T) {
	parser := mparsers.NewMustacheParser()
	err := parser.SetTemplate("Hello, {{{NAME}}}{{ #if ESCLAMATION }}!{{/if}}{{{^ESCLAMATION}}}.{{{/ESCLAMATION}}}")
	assert.Nil(t, err)

	expectedTokens := []*mparsers.MustacheToken{
		mparsers.NewMustacheToken(mparsers.TokenValue, "Hello, ", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenEscapedVariable, "NAME", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenSection, "ESCLAMATION", 0, 0),
		mparsers.NewMustacheToken(mparsers.TokenInvertedSection, "ESCLAMATION", 0, 0),
	}

	tokens := parser.ResultTokens()
	assert.Equal(t, len(expectedTokens), len(tokens))

	for i := 0; i < len(tokens); i++ {
		assert.Equal(t, expectedTokens[i].Type(), tokens[i].Type())
		assert.Equal(t, expectedTokens[i].Value(), tokens[i].Value())
	}
}

func TestMustacheParserVariableNames(t *testing.T) {
	parser := mparsers.NewMustacheParser()
	err := parser.SetTemplate("Hello, {{{NAME}}}{{ #if ESCLAMATION }}!{{/if}}{{{^ESCLAMATION}}}.{{{/ESCLAMATION}}}")
	assert.Nil(t, err)

	assert.Equal(t, 2, len(parser.VariableNames()))
	assert.Equal(t, "NAME", parser.VariableNames()[0])
	assert.Equal(t, "ESCLAMATION", parser.VariableNames()[1])
}
