package cache

import (
	"context"
	"sync"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
)

// MemoryCache that stores values in the process memory.
//
//	Configuration parameters:
//		- options:
//			- timeout: default caching timeout in milliseconds (default: 1 minute)
//			- max_size: maximum number of values stored in this cache (default: 1000)
//
// see ICache
//
//	Example:
//		cache := NewMemoryCache[string]();
//		res, err := cache.Store(contex.Background(), "key1", "ABC", 10000);
type MemoryCache[T any] struct {
	cache     map[string]*CacheEntry[string]
	mtx       *sync.Mutex
	timeout   int64
	maxSize   int
	convertor convert.IJSONEngine[T]
}

// NewMemoryCache creates a new instance of the cache.
// Returns: *MemoryCache
func NewMemoryCache[T any]() *MemoryCache[T] {
	return &MemoryCache[T]{
		cache:     map[string]*CacheEntry[string]{},
		mtx:       &sync.Mutex{},
		timeout:   60000,
		maxSize:   1000,
		convertor: convert.NewDefaultCustomTypeJsonConvertor[T](),
	}
}

// NewMemoryCacheFromConfig creates a new instance of the cache.
//
//	Parameters: cfg *config.ConfigParams configuration parameters to be set.
//	Returns: *MemoryCache
func NewMemoryCacheFromConfig[T any](ctx context.Context, cfg *config.ConfigParams) *MemoryCache[T] {
	c := NewMemoryCache[T]()
	c.Configure(ctx, cfg)
	return c
}

// Configure configures component by passing configuration parameters.
//
//	Parameters: config *config.ConfigParams configuration parameters to be set.
func (c *MemoryCache[T]) Configure(ctx context.Context, cfg *config.ConfigParams) {
	c.timeout = cfg.GetAsLongWithDefault("timeout", c.timeout)
	c.maxSize = cfg.GetAsIntegerWithDefault("max_size", c.maxSize)
}

// Cleanup memory cache, public thread save method
func (c *MemoryCache[T]) Cleanup() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cleanup()
}

// Cleanup memory cache, not thread save
func (c *MemoryCache[T]) cleanup() {
	var oldest *CacheEntry[string]
	var keysToRemove = make([]string, 0)

	for key, value := range c.cache {
		if value.IsExpired() {
			keysToRemove = append(keysToRemove, key)
		}
		if oldest == nil || oldest.Expiration().After(value.Expiration()) {
			oldest = value
		}
	}

	for _, key := range keysToRemove {
		delete(c.cache, key)
	}

	if len(c.cache) > c.maxSize && oldest != nil {
		delete(c.cache, oldest.Key())
	}
}

// Retrieve cached value from the cache using its key.
// If value is missing in the cache or expired it returns null.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique value key.
//	Returns T, error
func (c *MemoryCache[T]) Retrieve(ctx context.Context, key string) (T, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var defaultValue T

	if key == "" {
		return defaultValue, errors.NewInvalidStateError(
			cctx.GetTraceId(ctx),
			"INVALID_KEY",
			"key can not be empty string",
		)
	}

	entry := c.cache[key]
	if entry != nil {
		if entry.IsExpired() {
			delete(c.cache, key)
			return defaultValue, nil
		}
		value, err := c.convertor.FromJson(entry.Value())
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}

// Store value in the cache with expiration time, if success return stored value.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique value key.
//		- value T a value to store.
//		- timeout int64 expiration timeout in milliseconds.
//	Returns T, error
func (c *MemoryCache[T]) Store(ctx context.Context,
	key string, value T, timeout int64) (T, error) {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	var defaultValue T

	if key == "" {
		return value, errors.NewInvalidStateError(
			cctx.GetTraceId(ctx),
			"INVALID_KEY",
			"key can not be empty string",
		)
	}

	entry := c.cache[key]
	if timeout <= 0 {
		timeout = c.timeout
	}

	jsonVal, err := c.convertor.ToJson(value)
	if err != nil {
		return defaultValue, err
	}

	if entry != nil {
		entry.SetValue(jsonVal, timeout)
	} else {
		c.cache[key] = NewCacheEntry[string](key, jsonVal, timeout)
	}

	// cleanup
	if c.maxSize > 0 && len(c.cache) > c.maxSize {
		c.cleanup()
	}

	return value, nil
}

// Remove a value from the cache by its key.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique value key.
//	Returns: error
func (c *MemoryCache[T]) Remove(ctx context.Context, key string) error {

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if key == "" {
		return errors.NewInvalidStateError(
			cctx.GetTraceId(ctx),
			"INVALID_KEY",
			"key can not be empty string",
		)
	}

	delete(c.cache, key)

	return nil
}

// Contains check is value contains in cache and time not expire.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique value key.
//	Returns: bool
func (c *MemoryCache[T]) Contains(ctx context.Context, key string) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if entry, ok := c.cache[key]; ok {
		if entry.IsExpired() {
			delete(c.cache, key)
			return false
		}
		return true
	}
	return false
}

// Clear a value from the cache.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
func (c *MemoryCache[T]) Clear(ctx context.Context) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cache = make(map[string]*CacheEntry[string])

	return nil
}
