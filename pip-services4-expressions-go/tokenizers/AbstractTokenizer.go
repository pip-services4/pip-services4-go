package tokenizers

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

type ITokenizerOverrides interface {
	ReadNextToken() *Token
}

// AbstractTokenizer implements an abstract tokenizer class.
type AbstractTokenizer struct {
	Overrides ITokenizerOverrides
	mp        *utilities.CharReferenceMap

	skipUnknown      bool
	skipWhitespaces  bool
	skipComments     bool
	skipEof          bool
	mergeWhitespaces bool
	unifyNumbers     bool
	decodeStrings    bool

	commentState    ICommentState
	numberState     INumberState
	quoteState      IQuoteState
	symbolState     ISymbolState
	whitespaceState IWhitespaceState
	wordState       IWordState

	Scanner        io.IScanner
	NextTokenValue *Token
	LastTokenType  int
}

func InheritAbstractTokenizer(overrides ITokenizerOverrides) *AbstractTokenizer {
	c := AbstractTokenizer{
		Overrides:     overrides,
		mp:            utilities.NewCharReferenceMap(),
		LastTokenType: Unknown,
	}
	return &c
}

func (c *AbstractTokenizer) SkipUnknown() bool {
	return c.skipUnknown
}

func (c *AbstractTokenizer) SetSkipUnknown(value bool) {
	c.skipUnknown = value
}

func (c *AbstractTokenizer) SkipWhitespaces() bool {
	return c.skipWhitespaces
}

func (c *AbstractTokenizer) SetSkipWhitespaces(value bool) {
	c.skipWhitespaces = value
}

func (c *AbstractTokenizer) SkipComments() bool {
	return c.skipComments
}

func (c *AbstractTokenizer) SetSkipComments(value bool) {
	c.skipComments = value
}

func (c *AbstractTokenizer) SkipEof() bool {
	return c.skipEof
}

func (c *AbstractTokenizer) SetSkipEof(value bool) {
	c.skipEof = value
}

func (c *AbstractTokenizer) MergeWhitespaces() bool {
	return c.mergeWhitespaces
}

func (c *AbstractTokenizer) SetMergeWhitespaces(value bool) {
	c.mergeWhitespaces = value
}

func (c *AbstractTokenizer) UnifyNumbers() bool {
	return c.unifyNumbers
}

func (c *AbstractTokenizer) SetUnifyNumbers(value bool) {
	c.unifyNumbers = value
}

func (c *AbstractTokenizer) DecodeStrings() bool {
	return c.decodeStrings
}

func (c *AbstractTokenizer) SetDecodeStrings(value bool) {
	c.decodeStrings = value
}

func (c *AbstractTokenizer) CommentState() ICommentState {
	return c.commentState
}

func (c *AbstractTokenizer) SetCommentState(value ICommentState) {
	c.commentState = value
}

func (c *AbstractTokenizer) NumberState() INumberState {
	return c.numberState
}

func (c *AbstractTokenizer) SetNumberState(value INumberState) {
	c.numberState = value
}

func (c *AbstractTokenizer) QuoteState() IQuoteState {
	return c.quoteState
}

func (c *AbstractTokenizer) SetQuoteState(value IQuoteState) {
	c.quoteState = value
}

func (c *AbstractTokenizer) SymbolState() ISymbolState {
	return c.symbolState
}

func (c *AbstractTokenizer) SetSymbolState(value ISymbolState) {
	c.symbolState = value
}

func (c *AbstractTokenizer) WhitespaceState() IWhitespaceState {
	return c.whitespaceState
}

func (c *AbstractTokenizer) SetWhitespaceState(value IWhitespaceState) {
	c.whitespaceState = value
}

func (c *AbstractTokenizer) WordState() IWordState {
	return c.wordState
}

func (c *AbstractTokenizer) SetWordState(value IWordState) {
	c.wordState = value
}

func (c *AbstractTokenizer) GetCharacterState(symbol rune) ITokenizerState {
	state, _ := c.mp.Lookup(symbol).(ITokenizerState)
	return state
}

func (c *AbstractTokenizer) SetCharacterState(fromSymbol rune, toSymbol rune, state ITokenizerState) {
	c.mp.AddInterval(fromSymbol, toSymbol, state)
}

func (c *AbstractTokenizer) ClearCharacterStates() {
	c.mp.Clear()
}

func (c *AbstractTokenizer) Reader() io.IScanner {
	return c.Scanner
}

func (c *AbstractTokenizer) SetReader(value io.IScanner) {
	c.Scanner = value
	c.NextTokenValue = nil
	c.LastTokenType = Unknown
}

func (c *AbstractTokenizer) HasNextToken() bool {
	if c.NextTokenValue == nil {
		c.NextTokenValue = c.Overrides.ReadNextToken()
	}
	return c.NextTokenValue != nil
}

func (c *AbstractTokenizer) NextToken() *Token {
	token := c.NextTokenValue
	if token == nil {
		token = c.Overrides.ReadNextToken()
	}
	c.NextTokenValue = nil
	return token
}

func (c *AbstractTokenizer) ReadNextToken() *Token {
	if c.Scanner == nil {
		return nil
	}

	line := c.Scanner.PeekLine()
	column := c.Scanner.PeekColumn()
	var token *Token = nil

	for true {
		// Read character
		nextChar := c.Scanner.Peek()

		// If reached Eof then exit
		if utilities.CharValidator.IsEof(nextChar) {
			token = nil
			break
		}

		// Get state for character
		state := c.GetCharacterState(nextChar)
		if state != nil {
			token = state.NextToken(c.Scanner, c)
		}

		// Check for unknown characters and endless loops...
		if token == nil || token.Value() == "" {
			chr := c.Scanner.Read()
			token = NewToken(Unknown, string(chr), line, column)
		}

		// Skip unknown characters if option set.
		if token.Type() == Unknown && c.skipUnknown {
			c.LastTokenType = token.Type()
			continue
		}

		// Decode strings is option set.
		if _, ok := state.(IQuoteState); ok && c.decodeStrings {
			token = NewToken(token.Type(), c.QuoteState().DecodeString(token.Value(), nextChar), line, column)
		}

		// Skips comments if option set.
		if token.Type() == Comment && c.skipComments {
			c.LastTokenType = token.Type()
			continue
		}

		// Skips whitespaces if option set.
		if token.Type() == Whitespace && c.LastTokenType == Whitespace && c.skipWhitespaces {
			c.LastTokenType = token.Type()
			continue
		}

		// Unifies whitespaces if option set.
		if token.Type() == Whitespace && c.mergeWhitespaces {
			token = NewToken(Whitespace, " ", line, column)
		}

		// Unifies numbers if option set.
		if c.unifyNumbers &&
			(token.Type() == Integer || token.Type() == Float || token.Type() == HexDecimal) {
			token = NewToken(Number, token.Value(), line, column)
		}

		break
	}

	// Adds an Eof if option is not set.
	if token == nil && c.LastTokenType != Eof && !c.skipEof {
		token = NewToken(Eof, "", line, column)
	}

	// Assigns the last token type
	c.LastTokenType = Eof
	if token != nil {
		c.LastTokenType = token.Type()
	}

	return token
}

func (c *AbstractTokenizer) TokenizeStream(scanner io.IScanner) []*Token {
	c.SetReader(scanner)
	tokenList := []*Token{}
	token := c.NextToken()

	for token != nil {
		tokenList = append(tokenList, token)
		token = c.NextToken()
	}

	return tokenList
}

func (c *AbstractTokenizer) TokenizeBuffer(buffer string) []*Token {
	scanner := io.NewStringScanner(buffer)
	return c.TokenizeStream(scanner)
}

func (c *AbstractTokenizer) TokenizeStreamToStrings(scanner io.IScanner) []string {
	c.SetReader(scanner)
	stringList := []string{}
	token := c.NextToken()

	for token != nil {
		stringList = append(stringList, token.Value())
		token = c.NextToken()
	}

	return stringList
}

func (c *AbstractTokenizer) TokenizeBufferToStrings(buffer string) []string {
	scanner := io.NewStringScanner(buffer)
	return c.TokenizeStreamToStrings(scanner)
}
