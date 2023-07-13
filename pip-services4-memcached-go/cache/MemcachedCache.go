package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	ccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

/*
MemcachedCache are distributed cache that stores values in Memcaches caching service.

The current implementation does not support authentication.

Configuration parameters:

  - connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
  - options:
  - max_size:              maximum number of values stored in this cache (default: 1000)
  - max_key_size:          maximum key length (default: 250)
  - max_expiration:        maximum expiration duration in milliseconds (default: 2592000)
  - max_value:             maximum value length (default: 1048576)
  - pool_size:             pool size (default: 5)
  - reconnect:             reconnection timeout in milliseconds (default: 10 sec)
  - retries:               number of retries (default: 3)
  - timeout:               default caching timeout in milliseconds (default: 1 minute)
  - failures:              number of failures before stop retrying (default: 5)
  - retry:                 retry timeout in milliseconds (default: 30 sec)
  - idle:                  idle timeout before disconnect in milliseconds (default: 5 sec)

References:

- *:discovery:*:*:1.0    (optional) IDiscovery services to resolve connection

Example:

		ctx := context.Background()

	    cache := NewMemcachedCache[string]();
	    cache.Configure(ctx, cconf.NewConfigParamsFromTuples(
	      "host", "localhost",
	      "port", 11211,
	    ));

	    err := cache.Open(ctx)
	      ...


	    ret, err := cache.Store(ctx, "key1", []byte("ABC"))
	    if err != nil {
	    	...
	    }

	    res, err := cache.Retrive(ctx,"key1")
	    value, _ := res.([]byte)
	    fmt.Println(string(value))     // Result: "ABC"
*/
type MemcachedCache[T any] struct {
	connectionResolver *ccon.ConnectionResolver
	// maxKeySize         int
	// maxExpiration      int64
	// maxValue           int64
	// poolSize           int
	// reconnect          int
	timeout int
	// retries            int
	// failures           int
	// retry              int
	remove bool
	//idle   int
	client    *memcache.Client
	convertor cconv.IJSONEngine[T]
	logger    clog.CompositeLogger
}

// NewMemcachedCache method are creates a new instance of this cache.
func NewMemcachedCache[T any]() *MemcachedCache[T] {
	c := &MemcachedCache[T]{
		connectionResolver: ccon.NewEmptyConnectionResolver(),
		// maxKeySize:         250,
		// maxExpiration:      2592000,
		// maxValue:           1048576,
		// poolSize:           5,
		// reconnect:          10000,
		timeout: 5000,
		// retries:            5,
		// failures:           5,
		// retry:              30000,
		remove: false,
		//idle:   5000,
		client:    nil,
		convertor: cconv.NewDefaultCustomTypeJsonConvertor[T](),
		logger:    *clog.NewCompositeLogger(),
	}
	return c
}

// Configure method are configures component by passing configuration parameters.
//   - ctx context.Context
//   - config    configuration parameters to be set.
func (c *MemcachedCache[T]) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.connectionResolver.Configure(ctx, config)
	c.logger.Configure(ctx, config)

	// c.maxKeySize = config.GetAsIntegerWithDefault("options.max_key_size", c.maxKeySize)
	// c.maxExpiration = config.GetAsLongWithDefault("options.max_expiration", c.maxExpiration)
	// c.maxValue = config.GetAsLongWithDefault("options.max_value", c.maxValue)
	// c.poolSize = config.GetAsIntegerWithDefault("options.pool_size", c.poolSize)
	// c.reconnect = config.GetAsIntegerWithDefault("options.reconnect", c.reconnect)
	c.timeout = config.GetAsIntegerWithDefault("options.timeout", c.timeout)
	// c.retries = config.GetAsIntegerWithDefault("options.retries", c.retries)
	// c.failures = config.GetAsIntegerWithDefault("options.failures", c.failures)
	// c.retry = config.GetAsIntegerWithDefault("options.retry", c.retry)
	// c.remove = config.GetAsBooleanWithDefault("options.remove", c.remove)
	//c.idle = config.GetAsIntegerWithDefault("options.idle", c.idle)
}

// SetReferences are sets references to dependent components.
//   - ctx context.Context
//   - references 	references to locate the component dependencies.
func (c *MemcachedCache[T]) SetReferences(ctx context.Context, references cref.IReferences) {
	c.connectionResolver.SetReferences(ctx, references)
	c.logger.SetReferences(ctx, references)
}

// IsOpen Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *MemcachedCache[T]) IsOpen() bool {
	return c.client != nil
}

// Open method are opens the component.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//
// Retruns: error or nil no errors occured.
func (c *MemcachedCache[T]) Open(ctx context.Context) error {
	connections, err := c.connectionResolver.ResolveAll(ctx)

	if err == nil && len(connections) == 0 {
		err = cerr.NewConfigError(cctx.GetTraceId(ctx), "NO_CONNECTION", "Connection is not configured")
	}

	if err != nil {
		return err
	}

	var servers []string = make([]string, 0)
	for _, connection := range connections {
		host := connection.Host()
		port := connection.Port()
		if port == 0 {
			port = 11211
		}

		servers = append(servers, host+":"+strconv.FormatInt(int64(port), 10))
	}

	// options = {
	//     maxKeySize: c.maxKeySize,
	//     maxExpiration: c.maxExpiration,
	//     maxValue: c.maxValue,
	//     poolSize: c.poolSize,
	//     reconnect: c.reconnect,
	//     timeout: c.timeout,
	//     retries: c.retries,
	//     failures: c.failures,
	//     retry: c.retry,
	//     remove: c.remove,
	//     idle: c.idle
	// };

	c.client = memcache.New(servers...)
	c.client.Timeout = time.Duration(c.timeout) * time.Millisecond
	//c.client.MaxIdleConns = c.idle

	return nil
}

// Close method are closes component and frees used resources.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//
// Retruns: error or nil no errors occured.
func (c *MemcachedCache[T]) Close(ctx context.Context) error {
	c.client = nil
	return nil
}

func (c *MemcachedCache[T]) checkOpened(traceId string) (state bool, err error) {
	if !c.IsOpen() {
		err = cerr.NewInvalidStateError(traceId, "NOT_OPENED", "Connection is not opened")
		return false, err
	}

	return true, nil
}

// Retrieve method are retrieves cached value from the cache using its key.
// If value is missing in the cache or expired it returns nil.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique value key.
//     Retruns: cached value or error.
func (c *MemcachedCache[T]) Retrieve(ctx context.Context, key string) (value T, err error) {
	var defaultValue T

	if state, err := c.checkOpened(cctx.GetTraceId(ctx)); !state {
		return defaultValue, err
	}
	item, err := c.client.Get(key)
	if err != nil && err == memcache.ErrCacheMiss {
		err = nil
	}
	if item != nil {
		value, err := c.convertor.FromJson(string(item.Value))
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, err
}

// Store method are stores value in the cache with expiration time.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique value key.
//   - value             a value to store.
//   - timeout           expiration timeout in milliseconds.
//
// Returns: error or nil for success
func (c *MemcachedCache[T]) Store(ctx context.Context, key string, value T, timeout int64) (result T, err error) {
	var defaultValue T

	if state, err := c.checkOpened(cctx.GetTraceId(ctx)); !state {
		return defaultValue, err
	}

	timeoutInSec := int32(timeout) / 1000

	jsonVal, err := c.convertor.ToJson(value)
	if err != nil {
		return defaultValue, err
	}

	item := memcache.Item{
		Key:        key,
		Value:      []byte(jsonVal),
		Expiration: timeoutInSec,
	}
	return value, c.client.Set(&item)
}

// Remove method are removes a value from the cache by its key.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique value key.
//
// Retruns: an error or nil for success
func (c *MemcachedCache[T]) Remove(ctx context.Context, key string) error {
	state, err := c.checkOpened(cctx.GetTraceId(ctx))

	if !state {
		return err
	}

	err = c.client.Delete(key)
	if err != nil && err == memcache.ErrCacheMiss {
		err = nil
	}
	return err
}

// Contains check is value stores
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique value key.
func (c *MemcachedCache[T]) Contains(ctx context.Context, key string) bool {
	state, err := c.checkOpened(cctx.GetTraceId(ctx))
	if !state {
		c.logger.Error(ctx, err, "Connection is not opened")
		return false
	}

	if _, err := c.client.Get(key); err != nil {
		return false
	}

	return true
}
