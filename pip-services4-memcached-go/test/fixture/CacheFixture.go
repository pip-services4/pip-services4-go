package test_fixture

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	ccache "github.com/pip-services4/pip-services4-go/pip-services4-logic-go/cache"
	"github.com/stretchr/testify/assert"
)

var (
	KEY1 string = "key1"
	KEY2 string = "key2"
	KEY3 string = "key3"
	KEY4 string = "key4"
	KEY5 string = "key5"
	KEY6 string = "key6"

	VALUE1 string                 = "value1"
	VALUE2 map[string]interface{} = map[string]interface{}{"val": "value2"}
	VALUE3 time.Time              = time.Now()
	VALUE4 []int                  = []int{1, 2, 3, 4}
	VALUE5 int                    = 12345
	VALUE6 interface{}            = nil
)

type CacheFixture struct {
	cache ccache.ICache[any]
}

func NewCacheFixture(cache ccache.ICache[any]) *CacheFixture {
	c := CacheFixture{}
	c.cache = cache
	return &c
}

func (c *CacheFixture) TestStoreAndRetrieve(t *testing.T) {
	ctx := context.Background()

	_, err := c.cache.Store(ctx, KEY1, VALUE1, 5000)
	assert.Nil(t, err)

	_, err = c.cache.Store(ctx, KEY2, VALUE2, 5000)
	assert.Nil(t, err)

	_, err = c.cache.Store(ctx, KEY3, VALUE3, 5000)
	assert.Nil(t, err)

	_, err = c.cache.Store(ctx, KEY4, VALUE4, 5000)
	assert.Nil(t, err)

	_, err = c.cache.Store(ctx, KEY5, VALUE5, 5000)
	assert.Nil(t, err)

	_, err = c.cache.Store(ctx, KEY6, VALUE6, 5000)
	assert.Nil(t, err)

	<-time.After(500 * time.Millisecond)

	val, err := c.cache.Retrieve(ctx, KEY1)
	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, VALUE1, val.(string))

	val, err = c.cache.Retrieve(ctx, KEY2)

	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, VALUE2, val.(map[string]interface{}))

	val, err = c.cache.Retrieve(ctx, KEY3)
	assert.Nil(t, err)

	marshalVal3, err := json.Marshal(VALUE3)
	assert.Nil(t, err)

	assert.NotNil(t, val)
	assert.Equal(t, string(marshalVal3), ("\"" + val.(string) + "\""))

	val, err = c.cache.Retrieve(ctx, KEY4)

	assert.Nil(t, err)
	assert.NotNil(t, val)

	intArr := make([]int, 0, 4)

	for _, v := range val.([]any) {
		intArr = append(intArr, int(v.(float64)))
	}

	assert.Equal(t, VALUE4, intArr)

	val, err = c.cache.Retrieve(ctx, KEY5)

	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, VALUE5, int(val.(float64)))

	val, err = c.cache.Retrieve(ctx, KEY6)

	assert.Nil(t, err)
	assert.Nil(t, val)
}

func (c *CacheFixture) TestRetrieveExpired(t *testing.T) {
	ctx := context.Background()

	_, err := c.cache.Store(ctx, KEY1, VALUE1, 1000)
	assert.Nil(t, err)

	<-time.After(1500 * time.Millisecond)

	val, err := c.cache.Retrieve(ctx, KEY1)
	assert.Nil(t, err)
	assert.Nil(t, val)
}

func (c *CacheFixture) TestRemove(t *testing.T) {
	ctx := context.Background()

	_, err := c.cache.Store(ctx, KEY1, VALUE1, 1000)
	assert.Nil(t, err)

	err = c.cache.Remove(ctx, KEY1)
	assert.Nil(t, err)

	val, err := c.cache.Retrieve(ctx, KEY1)
	assert.Nil(t, err)
	assert.Nil(t, val)
}
