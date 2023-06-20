package test_random

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/random"
	"github.com/stretchr/testify/assert"
)

func TestNextInteger(t *testing.T) {
	value := random.Integer.Next(0, 5)
	assert.True(t, value <= 5)

	value = random.Integer.Next(2, 5)
	assert.True(t, value <= 5 && value >= 2)
}

func TestUpdateInteger(t *testing.T) {
	value := random.Integer.Update(0, 5)
	assert.True(t, value <= 5 && value >= -5)

	value = random.Integer.Update(5, 0)

	value = random.Integer.Update(0, 0)
	assert.True(t, value == 0)
}

func TestIntegerSequence(t *testing.T) {
	list := random.Integer.Sequence(1, 5)
	assert.True(t, len(list) <= 5 && len(list) >= 1)

	list = random.Integer.Sequence(-1, 0)
	assert.True(t, len(list) == 0)

	list = random.Integer.Sequence(-1, -4)
	assert.True(t, len(list) == 0)

	list = random.Integer.Sequence(4, 4)
	assert.True(t, len(list) == 4)

	list = random.Integer.Sequence(5, 5)
	assert.True(t, len(list) == 5)
}
