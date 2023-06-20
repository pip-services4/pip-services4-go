package state

import "context"

// NullStateStore dummy state store implementation that doesn't do anything.
// It can be used in testing or in situations when state management is not required
// but shall be disabled.
type NullStateStore[T any] struct {
}

func NewEmptyNullStateStore[T any]() *NullStateStore[T] {
	return &NullStateStore[T]{}
}

// Load state from the store using its key.
// If value is missing in the store it returns nil.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key           a unique state key.
//	Returns: the state value or nil if value wasn't found.
func (c *NullStateStore[T]) Load(ctx context.Context, key string) T {
	var defaultValue T
	return defaultValue
}

// LoadBulk loads an array of states from the store using their keys.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- keys          unique state keys.
//	Returns: an array with state values and their corresponding keys.
func (c *NullStateStore[T]) LoadBulk(ctx context.Context, keys []string) []StateValue[T] {
	return []StateValue[T]{}
}

// Save state into the store.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key           a unique state key.
//		- value         a state value.
//	Returns: the state that was stored in the store.
func (c *NullStateStore[T]) Save(ctx context.Context, key string, value T) T {
	return value
}

// Delete a state from the store by its key.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key           a unique value key.
//	Returns: the state that was deleted in the store.
func (c *NullStateStore[T]) Delete(ctx context.Context, key string) T {
	var defaultValue T
	return defaultValue
}
