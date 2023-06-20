package test_io

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/stretchr/testify/assert"
)

type ScannerFixture struct {
	scanner io.IScanner
	content []rune
}

func NewScannerFixture(scanner io.IScanner, content string) *ScannerFixture {
	return &ScannerFixture{
		scanner: scanner,
		content: []rune(content),
	}
}

func (c *ScannerFixture) TestRead(t *testing.T) {
	c.scanner.Reset()

	for i := 0; i < len(c.content); i++ {
		chr := c.scanner.Read()
		assert.Equal(t, c.content[i], chr)
	}

	chr := c.scanner.Read()
	assert.Equal(t, rune(-1), chr)

	chr = c.scanner.Read()
	assert.Equal(t, rune(-1), chr)
}

func (c *ScannerFixture) TestUnread(t *testing.T) {
	c.scanner.Reset()

	chr := c.scanner.Peek()
	assert.Equal(t, c.content[0], chr)

	chr = c.scanner.Read()
	assert.Equal(t, c.content[0], chr)

	chr = c.scanner.Read()
	assert.Equal(t, c.content[1], chr)

	c.scanner.Unread()
	chr = c.scanner.Read()
	assert.Equal(t, c.content[1], chr)

	c.scanner.UnreadMany(2)
	chr = c.scanner.Read()
	assert.Equal(t, c.content[0], chr)
	chr = c.scanner.Read()
	assert.Equal(t, c.content[1], chr)
}

func (c *ScannerFixture) TestLineColumn(t *testing.T, position int, charAt rune, line int, column int) {
	c.scanner.Reset()

	// Get in position
	for position > 1 {
		c.scanner.Read()
		position--
	}

	// Test forward scanning
	chr := c.scanner.Read()
	assert.Equal(t, charAt, chr)
	ln := c.scanner.Line()
	assert.Equal(t, line, ln)
	col := c.scanner.Column()
	assert.Equal(t, column, col)

	// Moving backward
	chr = c.scanner.Read()
	if chr != rune(-1) {
		c.scanner.Unread()
	}
	c.scanner.Unread()

	// Test backward scanning
	chr = c.scanner.Read()
	assert.Equal(t, charAt, chr)
	ln = c.scanner.Line()
	assert.Equal(t, line, ln)
	col = c.scanner.Column()
	assert.Equal(t, column, col)
}
