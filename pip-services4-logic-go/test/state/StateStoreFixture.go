package test_state

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-logic-go/state"
	"github.com/stretchr/testify/assert"
)

const KEY1 = "key1"
const KEY2 = "key2"

const VALUE1 = "value1"
const VALUE2 = "value2"

type StateStoreFixture struct {
	state state.IStateStore[any]
}

func NewStateStoreFixture(state state.IStateStore[any]) *StateStoreFixture {
	return &StateStoreFixture{
		state: state,
	}
}

func (c *StateStoreFixture) TestSaveAndLoad(t *testing.T) {
	c.state.Save(context.Background(), KEY1, VALUE1)
	c.state.Save(context.Background(), KEY2, VALUE2)

	val := c.state.Load(context.Background(), KEY1)
	assert.NotNil(t, val)
	assert.Equal(t, VALUE1, val)

	values := c.state.LoadBulk(context.Background(), []string{KEY2})
	assert.Len(t, values, 1)
	assert.Equal(t, KEY2, values[0].Key)
	assert.Equal(t, VALUE2, values[0].Value)
}

func (c *StateStoreFixture) TestDelete(t *testing.T) {
	c.state.Save(context.Background(), KEY1, VALUE1)

	c.state.Delete(context.Background(), KEY1)

	val := c.state.Load(context.Background(), KEY1)
	assert.Nil(t, val)
}
