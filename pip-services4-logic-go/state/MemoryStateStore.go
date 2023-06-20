package state

import (
	"context"
	"sync"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// MemoryStateStore is a state store that keeps states in the process memory.
//
//	Configuration parameters:
//		- options:
//		- timeout: default caching timeout in milliseconds (default: disabled)
//
//	Example:
//		store := NewMemoryStateStore[MyType]();
//		value := store.Load(context.Background(), "123", "key1");
//		...
//		store.Save(context.Background(), "123", "key1", MyType{});
type MemoryStateStore[T any] struct {
	states    map[string]*StateEntry[string]
	timeout   int64
	mtx       sync.Mutex
	convertor convert.IJSONEngine[T]
}

const StoreOptionsTimeoutConfigParameter = "options.timeout"

// NewEmptyMemoryStateStore creates a new instance of the state store.
func NewEmptyMemoryStateStore[T any]() *MemoryStateStore[T] {
	return &MemoryStateStore[T]{
		states:    make(map[string]*StateEntry[string]),
		timeout:   0,
		convertor: convert.NewDefaultCustomTypeJsonConvertor[T](),
	}
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config configuration parameters to be set.
func (c *MemoryStateStore[T]) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.timeout = config.GetAsLongWithDefault(StoreOptionsTimeoutConfigParameter, c.timeout)
}

// cleanup clears component state.
func (c *MemoryStateStore[T]) cleanup() {
	if c.timeout == 0 {
		return
	}

	cutOffTime := time.Now().UTC().UnixNano() - c.timeout
	// Cleanup obsolete entries
	for prop, entry := range c.states {
		// Remove obsolete entry
		if entry.GetLastUpdateTime() < cutOffTime {
			delete(c.states, prop)
		}
	}
}

// Load state from the store using its key.
// If value is missing in the store it returns nil.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key           a unique state key.
//	Returns: the state value or nil if value wasn't found.
func (c *MemoryStateStore[T]) Load(ctx context.Context, key string) T {
	if len(key) == 0 {
		panic(errors.NewError("Key cannot be empty"))
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	// Cleanup the stored states
	c.cleanup()
	// Get entry from the store
	var defaultValue T
	if entry, ok := c.states[key]; ok && entry != nil {
		if res, err := c.convertor.FromJson(entry.GetValue()); err == nil {
			return res
		}
	}

	return defaultValue
}

// LoadBulk loads an array of states from the store using their keys.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- keys          unique state keys.
//	Returns: an array with state values and their corresponding keys.
func (c *MemoryStateStore[T]) LoadBulk(ctx context.Context, keys []string) []StateValue[T] {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	// Cleanup the stored states
	c.cleanup()

	result := make([]StateValue[T], 0)
	for _, key := range keys {
		if entry, ok := c.states[key]; ok && entry != nil {
			if res, err := c.convertor.FromJson(entry.GetValue()); err == nil {
				result = append(result, StateValue[T]{Key: key, Value: res})
			}
		}
	}
	return result
}

// Save state into the store.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key           a unique state key.
//		- value         a state value.
//	Returns: the state that was stored in the store.
func (c *MemoryStateStore[T]) Save(ctx context.Context, key string, value T) T {
	if len(key) == 0 {
		panic(errors.NewError("Key cannot be empty"))
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	// Cleanup the stored states
	c.cleanup()

	// Get the entry
	var entry *StateEntry[string]
	if val, ok := c.states[key]; ok {
		entry = val
	}

	// Update the entry
	if entry != nil {
		if buf, err := c.convertor.ToJson(value); err == nil {
			entry.SetValue(buf)
		}
	} else { // Or create a new entry
		if buf, err := c.convertor.ToJson(value); err == nil {
			c.states[key] = NewStateEntry(key, buf)
		}
	}

	var defaultValue T
	if entry, ok := c.states[key]; ok && entry != nil {
		res, _ := c.convertor.FromJson(entry.GetValue())
		return res
	}

	return defaultValue
}

// Delete a state from the store by its key.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key           a unique value key.
//	Returns: the state that was deleted in the store.
func (c *MemoryStateStore[T]) Delete(ctx context.Context, key string) T {
	if len(key) == 0 {
		panic(errors.NewError("Key cannot be empty"))
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	// Cleanup the stored states
	c.cleanup()

	var defaultValue T

	// Get the entry
	if entry, ok := c.states[key]; ok {
		delete(c.states, key)
		if entry != nil {
			res, _ := c.convertor.FromJson(entry.GetValue())
			return res
		}
	}

	return defaultValue
}
