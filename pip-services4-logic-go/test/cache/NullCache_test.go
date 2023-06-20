package test_cache

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-logic-go/cache"
	"github.com/stretchr/testify/assert"
)

func TestNullCache(t *testing.T) {
	_cache := cache.NewNullCache[any]()

	value, err := _cache.Retrieve(context.Background(), "key1")
	assert.Nil(t, value)
	assert.Nil(t, err)

	value, err = _cache.Store(context.Background(), "key1", "value1", 0)
	assert.Equal(t, "value1", value)
	assert.Nil(t, err)
}
