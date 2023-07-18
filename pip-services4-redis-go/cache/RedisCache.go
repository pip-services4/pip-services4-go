package persistence

import (
	"context"

	"strconv"
	"time"

	"github.com/go-redis/redis"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	ccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

/*
Distributed cache that stores values in Redis in-memory database.

Configuration parameters:

  - connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
  - credential(s):
  - store_key:             key to retrieve parameters from credential store
  - username:              user name (currently is not used)
  - password:              user password
  - options:
  - retries:               number of retries (default: 3)
  - timeout:               default caching timeout in milliseconds (default: 1 minute)
  - db_num:                database number in Redis  (default 0)
  - max_size:            	 maximum number of values stored in this cache (default: 1000)
  - cluster:            	 enable redis cluster

References:

- *:discovery:*:*:1.0        (optional) IDiscovery services to resolve connection
- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credential

Example:

		ctx := context.Background()

	    cache = NewRedisCache[string]();
	    cache.Configure(ctx, cconf.NewConfigParamsFromTuples(
	      "host", "localhost",
	      "port", 6379,
	    ));

	    err = cache.Open(ctx, "123")
	      ...

	    ret, err := cache.Store(ctx, "123", "key1", []byte("ABC"))
	    if err != nil {
	    	...
	    }

	    res, err := cache.Retrieve(ctx, "123", "key1")
	    value, _ := res.([]byte)
	    fmt.Println(string(value))     // Result: "ABC"
*/
type RedisCache[T any] struct {
	connectionResolver *ccon.ConnectionResolver
	credentialResolver *cauth.CredentialResolver

	timeout   int
	retries   int
	dbNum     int
	isCluster bool

	client redis.UniversalClient
	logger clog.CompositeLogger

	convertor cconv.IJSONEngine[T]
}

// NewRedisCache method are creates a new instance of this cache.
func NewRedisCache[T any]() *RedisCache[T] {
	return &RedisCache[T]{
		logger:             *clog.NewCompositeLogger(),
		connectionResolver: ccon.NewEmptyConnectionResolver(),
		credentialResolver: cauth.NewEmptyCredentialResolver(),
		timeout:            30000,
		retries:            30000,
		dbNum:              3,
		convertor:          cconv.NewDefaultCustomTypeJsonConvertor[T](),
	}
}

// Configure method are configures component by passing configuration parameters.
//   - ctx context.Context
//   - config    configuration parameters to be set.
func (c *RedisCache[T]) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.connectionResolver.Configure(ctx, config)
	c.credentialResolver.Configure(ctx, config)
	c.logger.Configure(ctx, config)

	c.timeout = config.GetAsIntegerWithDefault("options.timeout", c.timeout)
	c.retries = config.GetAsIntegerWithDefault("options.retries", c.retries)
	c.dbNum = config.GetAsIntegerWithDefault("options.db_num", c.dbNum)
	if c.dbNum > 15 || c.dbNum < 0 {
		c.dbNum = 0
	}
	c.isCluster = config.GetAsBooleanWithDefault("options.cluster", c.isCluster)
}

// Sets references to dependent components.
//   - ctx context.Context
//   - references 	references to locate the component dependencies.
func (c *RedisCache[T]) SetReferences(ctx context.Context, references cref.IReferences) {
	c.connectionResolver.SetReferences(ctx, references)
	c.credentialResolver.SetReferences(ctx, references)
	c.logger.SetReferences(ctx, references)
}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *RedisCache[T]) IsOpen() bool {
	return c.client != nil
}

// Open method are opens the component.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//
// Returns: error or nil no errors occured.
func (c *RedisCache[T]) Open(ctx context.Context) error {
	var (
		connection *ccon.ConnectionParams
		credential *cauth.CredentialParams
		options    redis.Options
	)

	connection, err := c.connectionResolver.Resolve(ctx)

	if err == nil && connection == nil {
		err = cerr.NewConfigError(cctx.GetTraceId(ctx), "NO_CONNECTION", "Connection is not configured")
		return err
	}

	credential, err = c.credentialResolver.Lookup(ctx)

	if err != nil {
		return err
	}

	options.DialTimeout = time.Duration(c.timeout) * time.Millisecond
	options.DB = c.dbNum
	options.MaxRetries = c.retries

	if credential != nil {
		options.Password = credential.Password()
	}

	if connection.Uri() != "" {
		options.Addr = connection.Uri()
	} else {
		host := connection.Host()
		if host == "" {
			host = "localhost"
		}
		port := strconv.FormatInt(int64(connection.Port()), 10)
		if port == "0" {
			port = "6379"
		}
		options.Addr = host + ":" + port
	}
	if c.isCluster {
		c.client = redis.NewClusterClient(&redis.ClusterOptions{
			OnNewNode: func(client *redis.Client) {
				client = redis.NewClient(&options)
			},
			Addrs:       append([]string{}, options.Addr),
			Password:    options.Password,
			DialTimeout: options.DialTimeout,
		})
	} else {
		c.client = redis.NewClient(&options)
	}
	_, err = c.client.Ping().Result()
	return err
}

// Close method are closes component and frees used resources.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//
// Retruns: error or nil no errors occured.
func (c *RedisCache[T]) Close(ctx context.Context) error {
	if c.client != nil {
		err := c.client.Close()
		c.client = nil
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RedisCache[T]) checkOpened(traceId string) (state bool, err error) {
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
func (c *RedisCache[T]) Retrieve(ctx context.Context, key string) (value T, err error) {
	var defaultValue T

	if state, err := c.checkOpened(cctx.GetTraceId(ctx)); !state {
		return defaultValue, err
	}

	item, err := c.client.Get(key).Bytes()

	if err != nil {
		if err == redis.Nil {
			return defaultValue, nil
		}
		return defaultValue, err
	}

	if item != nil {
		val, err := c.convertor.FromJson(string(item))
		if err != nil {
			return defaultValue, err
		}
		return val, nil
	}

	return defaultValue, nil
}

// Store method are stores value in the cache with expiration time.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique value key.
//   - value             a value to store.
//   - timeout           expiration timeout in milliseconds.
//
// Retruns error or nil for success
func (c *RedisCache[T]) Store(ctx context.Context, key string, value T, timeout int64) (result T, err error) {
	state, err := c.checkOpened(cctx.GetTraceId(ctx))

	var defaultValue T

	if !state {
		return defaultValue, err
	}

	jsonVal, err := c.convertor.ToJson(value)
	if err != nil {
		return defaultValue, err
	}
	tmout := time.Duration(timeout) * time.Millisecond
	return value, c.client.Set(key, jsonVal, tmout).Err()
}

// Removes a value from the cache by its key.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique value key.
//
// Returns: error or nil for success
func (c *RedisCache[T]) Remove(ctx context.Context, key string) error {
	state, err := c.checkOpened(cctx.GetTraceId(ctx))
	if !state {
		return err
	}
	return c.client.Del(key).Err()
}

// Contains check is value stores
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique value key.
func (c *RedisCache[T]) Contains(ctx context.Context, key string) bool {
	state, err := c.checkOpened(cctx.GetTraceId(ctx))
	if !state {
		c.logger.Error(ctx, err, "Connection is not opened")
		return false
	}
	return c.client.Exists(key).Val() > 0
}
