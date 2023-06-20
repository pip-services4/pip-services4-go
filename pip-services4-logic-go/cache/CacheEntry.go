package cache

import (
	"time"
)

// CacheEntry Data object to store cached values with their keys used by MemoryCache
type CacheEntry[T any] struct {
	key        string
	value      T
	expiration time.Time
}

// NewCacheEntry creates a new instance of the cache entry and assigns its values.
//	Parameters:
//		- key string a unique key to locate the value.
//		- value T a value to be stored.
//		- timeout int64 expiration timeout in milliseconds.
//	Returns *CacheEntry
func NewCacheEntry[T any](key string, value T, timeout int64) *CacheEntry[T] {
	return &CacheEntry[T]{
		key:        key,
		value:      value,
		expiration: time.Now().Add(time.Duration(timeout) * time.Millisecond),
	}
}

// Key gets the key to locate the cached value.
//	Returns: string the value key.
func (c *CacheEntry[T]) Key() string {
	return c.key
}

// Value gets the cached value.
//	Returns: any the value object.
func (c *CacheEntry[T]) Value() T {
	return c.value
}

// Expiration gets the expiration timeout.
//	Returns time.Time the expiration timeout in milliseconds.
func (c *CacheEntry[T]) Expiration() time.Time {
	return c.expiration
}

// SetValue a new value and extends its expiration.
//	Parameters:
//		- value any a new cached value.
//		- timeout int64 an expiration timeout in milliseconds.
func (c *CacheEntry[T]) SetValue(value T, timeout int64) {
	c.value = value
	c.expiration = time.Now().Add(time.Duration(timeout) * time.Millisecond)
}

// IsExpired checks if this value already expired.
//	Returns: bool true if the value already expires and false otherwise.
func (c *CacheEntry[T]) IsExpired() bool {
	return time.Now().After(c.expiration)
}
