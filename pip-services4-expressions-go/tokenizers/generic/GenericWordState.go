package generic

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

// GenericWordState a wordState returns a word from a scanner. Like other states, a tokenizer transfers the job
// of reading to this state, depending on an initial character. Thus, the tokenizer decides
// which characters may begin a word, and this state determines which characters may appear
// as a second or later character in a word. These are typically different sets of characters;
// in particular, it is typical for digits to appear as parts of a word, but not
// as the initial character of a word.
// <p/>
// By default, the following characters may appear in a word.
// The method <code>setWordChars()</code> allows customizing this.
// <blockquote><pre>
// From    To
//
//	'a', 'z'
//	'A', 'Z'
//	'0', '9'
//
//	 as well as: minus sign, underscore, and apostrophe.
//
// </pre></blockquote>
type GenericWordState struct {
	mp *utilities.CharReferenceMap
}

// NewGenericWordState constructs a word state with a default idea of what characters
// are admissible inside a word (as described in the class comment).
func NewGenericWordState() *GenericWordState {
	c := &GenericWordState{
		mp: utilities.NewCharReferenceMap(),
	}

	c.SetWordChars('a', 'z', true)
	c.SetWordChars('A', 'Z', true)
	c.SetWordChars('0', '9', true)
	c.SetWordChars('-', '-', true)
	c.SetWordChars('_', '_', true)
	//c.SetWordChars(39, 39, true)
	c.SetWordChars(0x00c0, 0x00ff, true)
	c.SetWordChars(0x0100, 0xffff, true)

	return c
}

// NextToken ignore word (such as blanks and tabs), and return the tokenizer's next token.
//
//	Returns: the tokenizer's next token
func (c *GenericWordState) NextToken(
	scanner io.IScanner, tokenizer tokenizers.ITokenizer) *tokenizers.Token {

	tokenValue := strings.Builder{}
	nextSymbol := scanner.Read()
	line := scanner.Line()
	column := scanner.Column()

	for c.mp.Lookup(nextSymbol) != nil {
		tokenValue.WriteRune(nextSymbol)
		nextSymbol = scanner.Read()
	}

	if !utilities.CharValidator.IsEof(nextSymbol) {
		scanner.Unread()
	}

	return tokenizers.NewToken(tokenizers.Word, tokenValue.String(), line, column)
}

// SetWordChars establish characters in the given range as valid characters for part of a word after
// the first character. Note that the tokenizer must determine which characters are valid
// as the beginning character of a word.
//
//	Parameters:
//		- fromSymbol: First character index of the interval.
//		- toSymbol: Last character index of the interval.
//		- enable: <code>true</code> if this state should use characters in the given range.
func (c *GenericWordState) SetWordChars(fromSymbol rune, toSymbol rune, enable bool) {
	if enable {
		c.mp.AddInterval(fromSymbol, toSymbol, true)
	} else {
		c.mp.AddInterval(fromSymbol, toSymbol, nil)
	}
}

// ClearWordChars clears definitions of word chars.
func (c *GenericWordState) ClearWordChars() {
	c.mp.Clear()
}
