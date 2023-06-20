package state

import "time"

// StateEntry data object to store state values with their keys used by MemoryStateEntry
type StateEntry[T any] struct {
	key            string `json:"key" bson:"key"`
	value          T      `json:"value" bson:"value"`
	lastUpdateTime int64  `json:"last_update_time" bson:"last_update_time"` // timestamp in ms
}

// NewStateEntry method creates a new instance of the state entry and assigns its values.
//	Parameters:
//		- key   a unique key to locate the value.
//		- value a value to be stored.
func NewStateEntry[T any](key string, value T) *StateEntry[T] {
	return &StateEntry[T]{
		key:            key,
		value:          value,
		lastUpdateTime: time.Now().UTC().UnixNano() / (int64)(1000),
	}
}

// GetKey method gets the key to locate the state value.
//	Returns the value key.
func (c *StateEntry[T]) GetKey() string {
	return c.key
}

// GetValue method gets the sstate value.
//	Returns the value object.
func (c *StateEntry[T]) GetValue() T {
	return c.value
}

// GetLastUpdateTime method gets the last update time.
//	Returns the timestamp when the value ware stored.
func (c *StateEntry[T]) GetLastUpdateTime() int64 {
	return c.lastUpdateTime
}

// SetValue method sets a new state value.
//	Parameters:
//		- value a new cached value.
func (c *StateEntry[T]) SetValue(value T) {
	c.value = value
	c.lastUpdateTime = time.Now().UTC().UnixNano() / (int64)(1000)
}
