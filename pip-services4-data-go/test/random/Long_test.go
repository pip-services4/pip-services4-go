package test_random

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/random"
	"github.com/stretchr/testify/assert"
)

func TestNextLong(t *testing.T) {
	value := random.Long.Next(0, 5)
	assert.True(t, value <= 5)

	value = random.Long.Next(2, 5)
	assert.True(t, value <= 5 && value >= 2)
}

func TestUpdateLong(t *testing.T) {
	value := random.Long.Update(0, 5)
	assert.True(t, value <= 5 && value >= -5)

	value = random.Long.Update(5, 0)

	value = random.Long.Update(0, 0)
	assert.True(t, value == 0)
}

func TestLongSequence(t *testing.T) {
	list := random.Long.Sequence(1, 5)
	assert.True(t, len(list) <= 5 && len(list) >= 1)

	list = random.Long.Sequence(-1, 0)
	assert.True(t, len(list) == 0)

	list = random.Long.Sequence(-1, -4)
	assert.True(t, len(list) == 0)

	list = random.Long.Sequence(4, 4)
	assert.True(t, len(list) == 4)

	list = random.Long.Sequence(5, 5)
	assert.True(t, len(list) == 5)
}
