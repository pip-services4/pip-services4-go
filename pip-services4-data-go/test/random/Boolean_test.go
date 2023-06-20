package test_random

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/random"
	"github.com/stretchr/testify/assert"
)

func TestChance(t *testing.T) {
	value := random.Boolean.Chance(5, 10)
	assert.True(t, value || !value)

	value = random.Boolean.Chance(5, 5)
	assert.True(t, value || !value)

	value = random.Boolean.Chance(0, 0)
	assert.True(t, !value)

	value = random.Boolean.Chance(-1, 0)
	assert.True(t, !value)

	value = random.Boolean.Chance(-1, -1)
	assert.True(t, !value)
}
