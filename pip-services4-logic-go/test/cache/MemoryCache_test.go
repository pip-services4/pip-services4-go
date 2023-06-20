package test_cache

import (
	"context"
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-logic-go/cache"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {
	var _cache cache.ICache[any]
	_cache = cache.NewMemoryCache[any]()

	value, err := _cache.Retrieve(context.Background(), "key1")
	assert.Nil(t, value)
	assert.Nil(t, err)

	value, err = _cache.Store(context.Background(), "key1", "value1", 250)
	assert.Equal(t, "value1", value)
	assert.Nil(t, err)

	value, err = _cache.Retrieve(context.Background(), "key1")
	assert.Equal(t, "value1", value)
	assert.Nil(t, err)

	time.Sleep(500 * time.Millisecond)

	value, err = _cache.Retrieve(context.Background(), "key1")
	assert.Nil(t, value)
	assert.Nil(t, err)
}
