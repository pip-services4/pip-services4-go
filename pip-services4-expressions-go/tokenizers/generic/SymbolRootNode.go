package generic

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
)

// SymbolRootNode this class is a special case of a <code>SymbolNode</code>. A <code>SymbolRootNode</code>
// object has no symbol of its own, but has children that represent all possible symbols.
type SymbolRootNode struct {
	*SymbolNode
}

// NewSymbolRootNode creates and initializes a root node.
func NewSymbolRootNode() *SymbolRootNode {
	return &SymbolRootNode{
		SymbolNode: NewSymbolNode(nil, 0),
	}
}

// Add the given string as a symbol.
//
//	Parameters:
//		- value: the character sequence to add.
func (c *SymbolRootNode) Add(value string, tokenType int) {
	if value == "" {
		panic("Value must have at least 1 character")
	}

	v := []rune(value)
	childNode := c.EnsureChildWithChar(v[0])

	if childNode.TokenType() == tokenizers.Unknown {
		childNode.SetValid(true)
		childNode.SetTokenType(tokenizers.Symbol)
	}

	childNode.AddDescendantLine(v[1:], tokenType)
}

// NextToken return a symbol string from a scanner.
//
//	Parameters:
//		- scanner: A scanner to read from
//		- firstChar: The first character of this symbol, already read from the scanner.
//	Returns: A symbol string from a scanner
func (c *SymbolRootNode) NextToken(scanner io.IScanner) *tokenizers.Token {
	nextSymbol := scanner.Read()
	line := scanner.Line()
	column := scanner.Column()

	childNode := c.FindChildWithChar(nextSymbol)
	if childNode != nil {
		childNode = childNode.DeepestRead(scanner)
		childNode = childNode.UnreadToValid(scanner)
		childNodeValue := string(childNode.Ancestry())
		return tokenizers.NewToken(childNode.TokenType(), childNodeValue, line, column)
	} else {
		tokenValue := string([]rune{nextSymbol})
		return tokenizers.NewToken(tokenizers.Symbol, tokenValue, line, column)
	}
}
