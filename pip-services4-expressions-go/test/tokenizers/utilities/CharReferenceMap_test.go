package test_utilities

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
	"github.com/stretchr/testify/assert"
)

func TestCharReferenceMapDefaultInterval(t *testing.T) {
	m := utilities.NewCharReferenceMap()
	assert.Nil(t, m.Lookup('A'))
	assert.Nil(t, m.Lookup(0x2045))

	m.AddDefaultInterval(true)
	assert.NotNil(t, m.Lookup('A'))
	assert.NotNil(t, m.Lookup(0x2045))

	m.Clear()
	assert.Nil(t, m.Lookup('A'))
	assert.Nil(t, m.Lookup(0x2045))
}

func TestCharReferenceMapInterval(t *testing.T) {
	m := utilities.NewCharReferenceMap()
	assert.Nil(t, m.Lookup('A'))
	assert.Nil(t, m.Lookup(0x2045))

	m.AddInterval('A', 'z', true)
	assert.NotNil(t, m.Lookup('A'))
	assert.Nil(t, m.Lookup(0x2045))

	m.AddInterval(0x2000, 0x20ff, true)
	assert.NotNil(t, m.Lookup('A'))
	assert.NotNil(t, m.Lookup(0x2045))

	m.Clear()
	assert.Nil(t, m.Lookup('A'))
	assert.Nil(t, m.Lookup(0x2045))

	m.AddInterval('A', 0x20ff, true)
	assert.NotNil(t, m.Lookup('A'))
	assert.NotNil(t, m.Lookup(0x2045))
}
