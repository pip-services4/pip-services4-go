package test_data

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	"github.com/stretchr/testify/assert"
)

func TestNextShort(t *testing.T) {
	id1 := keys.IdGenerator.NextShort()
	assert.NotEmpty(t, id1)
	assert.True(t, len(id1) >= 9)

	id2 := keys.IdGenerator.NextShort()
	assert.NotEmpty(t, id2)
	assert.True(t, len(id2) >= 9)

	assert.NotEqual(t, id1, id2)
}

func TestNextLong(t *testing.T) {
	id1 := keys.IdGenerator.NextLong()
	assert.NotEmpty(t, id1)
	assert.Len(t, id1, 32)

	id2 := keys.IdGenerator.NextLong()
	assert.NotEmpty(t, id2)
	assert.Len(t, id2, 32)

	assert.NotEqual(t, id1, id2)
}
