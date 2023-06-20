package cache

import "context"

// NullCache Dummy cache implementation that doesn't do anything.
// It can be used in testing or in situations when cache is required but shall be disabled.
type NullCache[T any] struct{}

// 	NewNullCache creates a new instance of the cache.
//	Returns: *NullCache
func NewNullCache[T any]() *NullCache[T] {
	return &NullCache[T]{}
}

// Retrieve retrieves cached value from the cache using its key.
// If value is missing in the cache or expired it returns null.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique value key.
//	Returns: T, error
func (c *NullCache[T]) Retrieve(ctx context.Context, key string) (T, error) {
	var defaultValue T
	return defaultValue, nil
}

// Store value in the cache with expiration time, if success return stored value.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique value key.
//		- value T a value to store.
//		- timeout int64 expiration timeout in milliseconds.
//	Returns T, error
func (c *NullCache[T]) Store(ctx context.Context, key string, value T, timeout int64) (T, error) {
	return value, nil
}

// Remove a value from the cache by its key.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.\
//		- key string a unique value key.
//	Returns: error
func (c *NullCache[T]) Remove(ctx context.Context, key string) error {
	return nil
}

// Contains check is value contains in cache and time not expire.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.\
//		- key string a unique value key.
//	Returns: bool
func (c *NullCache[T]) Contains(ctx context.Context, key string) bool {
	return false
}
