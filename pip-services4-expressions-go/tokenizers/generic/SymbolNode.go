package generic

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

// SymbolNode a <code>SymbolNode</code> object is a member of a tree that contains all possible prefixes
// of allowable symbols. Multi-character symbols appear in a <code>SymbolNode</code> tree
// with one node for each character.
// <p/>
// For example, the symbol <code>=:~</code> will appear in a tree as three nodes. The first
// node contains an equals sign, and has a child; that child contains a colon and has a child;
// this third child contains a tilde, and has no children of its own. If the colon node had
// another child for a dollar sign character, then the tree would contain the symbol <code>=:$</code>.
// <p/>
// A tree of <code>SymbolNode</code> objects collaborate to read a (potentially multi-character)
// symbol from an input stream. A root node with no character of its own finds an initial node
// that represents the first character in the input. This node looks to see if the next character
// in the stream matches one of its children. If so, the node delegates its reading task to its child.
// This approach walks down the tree, pulling symbols from the input that match the path down the tree.
// <p/>
// When a node does not have a child that matches the next character, we will have read the longest
// possible symbol prefix. This prefix may or may not be a valid symbol.
// Consider a tree that has had <code>=:~</code> added and has not had <code>=:</code> added.
// In this tree, of the three nodes that contain <code>=:~</code>, only the first and third contain
// complete symbols. If, say, the input contains <code>=:a</code>, the colon node will not have
// a child that matches the 'a' and so it will stop reading. The colon node has to "unread": it must
// push back its character, and ask its parent to unread. Unreading continues until it reaches
// an ancestor that represents a valid symbol.
type SymbolNode struct {
	parent    *SymbolNode
	character rune
	children  *utilities.CharReferenceMap
	tokenType int
	valid     bool
	ancestry  []rune
}

// NewSymbolNode constructs a SymbolNode with the given parent, representing the given character.
//
//	Parameters:
//		- parent: This node's parent
//		- character: This node's associated character.
func NewSymbolNode(parent *SymbolNode, character rune) *SymbolNode {
	return &SymbolNode{
		parent:    parent,
		character: character,
		tokenType: tokenizers.Unknown,
	}
}

// EnsureChildWithChar find or create a child for the given character.
func (c *SymbolNode) EnsureChildWithChar(value rune) *SymbolNode {
	if c.children == nil {
		c.children = utilities.NewCharReferenceMap()
	}

	childNode, _ := c.children.Lookup(value).(*SymbolNode)
	if childNode == nil {
		childNode = NewSymbolNode(c, value)
		c.children.AddInterval(value, value, childNode)
	}
	return childNode
}

// AddDescendantLine add a line of descendants that represent the characters in the given string.
func (c *SymbolNode) AddDescendantLine(value []rune, tokenType int) {
	if len(value) > 0 {
		childNode := c.EnsureChildWithChar(value[0])
		childNode.AddDescendantLine(value[1:], tokenType)
	} else {
		c.valid = true
		c.tokenType = tokenType
	}
}

// DeepestRead find the descendant that takes as many characters as possible from the input.
func (c *SymbolNode) DeepestRead(scanner io.IScanner) *SymbolNode {
	nextSymbol := scanner.Read()

	var childNode *SymbolNode
	if !utilities.CharValidator.IsEof(nextSymbol) {
		childNode = c.FindChildWithChar(nextSymbol)
	}
	if childNode == nil {
		scanner.Unread()
		return c
	}

	return childNode.DeepestRead(scanner)
}

// FindChildWithChar find a child with the given character.
func (c *SymbolNode) FindChildWithChar(value rune) *SymbolNode {
	if c.children == nil {
		return nil
	}
	result, _ := c.children.Lookup(value).(*SymbolNode)
	return result
}

// UnreadToValid find a descendant which is down the path the given string indicates.
//
//	func (c *SymbolNode) FindDescendant(value []rune) *SymbolNode {
//	   tempChar := CharValidator.Eof
//	   if len(value) > 0 {
//	   tempChar = value[0]
//	   }
//	   childNode := c.FindChildWithChar(tempChar)
//	   if !CharValidator.IsEof(tempChar) && childNode != nil && len(value) > 1 {
//	       childNode = childNode.FindDescendant(value[1:])
//	   }
//	   return childNode
//	}
//
// Unwind to a valid node; this node is "valid" if its ancestry represents a complete symbol.
// If this node is not valid, put back the character and ask the parent to unwind.
func (c *SymbolNode) UnreadToValid(scanner io.IScanner) *SymbolNode {
	if !c.valid && c.parent != nil {
		scanner.Unread()
		return c.parent.UnreadToValid(scanner)
	}
	return c
}

//internal SymbolNode Parent { get { return _parent; } }
//internal SymbolNode[] Children { get { return _children; } }
//internal char Character { get { return _character; } }

func (c *SymbolNode) Valid() bool {
	return c.valid
}

func (c *SymbolNode) SetValid(value bool) {
	c.valid = value
}

func (c *SymbolNode) TokenType() int {
	return c.tokenType
}

func (c *SymbolNode) SetTokenType(value int) {
	c.tokenType = value
}

// Ancestry show the symbol this node represents.
//
//	Returns: the symbol this node represents
func (c *SymbolNode) Ancestry() []rune {
	if c.ancestry == nil || len(c.ancestry) == 0 {
		if c.parent != nil {
			c.ancestry = c.parent.Ancestry()
		}
		if c.character != 0 {
			c.ancestry = append(c.ancestry, c.character)
		}
	}
	return c.ancestry
}
