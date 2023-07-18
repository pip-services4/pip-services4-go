package lock

import (
	"context"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	ccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	clock "github.com/pip-services4/pip-services4-go/pip-services4-logic-go/lock"
)

/*
RedisLock are distributed lock that is implemented based on Redis in-memory database.

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
  - retrytimeout:         timeout in milliseconds to retry lock acquisition. (Default: 100)
  - retries:               number of retries (default: 3)
  - db_num:                database number in Redis  (default 0)

References:

- *:discovery:*:*:1.0        (optional) IDiscovery services to resolve connection
- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credential

Example:

		ctx := context.Background()

	    lock = NewRedisRedis();
	    lock.Configure(ctx, cconf.NewConfigParamsFromTuples(
	      "host", "localhost",
	      "port", 6379,
	    ));

	    err = lock.Open(ctx)
	      ...

	    result, err := lock.TryAcquireLock(ctx, "key1", 3000)
	    if result {
	    	// Processing...
	    }
	    err = lock.ReleaseLock(ctx, "key1")
	    // Continue...
*/
type RedisLock struct {
	*clock.Lock
	connectionResolver *ccon.ConnectionResolver
	credentialResolver *cauth.CredentialResolver

	lockId  string
	timeout int
	retries int
	dbNum   int

	client redis.Conn
}

// NewRedisLock method are creates a new instance of this lock.
func NewRedisLock() *RedisLock {
	c := &RedisLock{
		connectionResolver: ccon.NewEmptyConnectionResolver(),
		credentialResolver: cauth.NewEmptyCredentialResolver(),
		lockId:             keys.IdGenerator.NextLong(),
		timeout:            30000,
		retries:            3,
		dbNum:              0,
		client:             nil,
	}
	c.Lock = clock.InheritLock(c)
	return c
}

// Configure method are configures component by passing configuration parameters.
// Parameters:
//   - config    configuration parameters to be set.
func (c *RedisLock) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.connectionResolver.Configure(ctx, config)
	c.credentialResolver.Configure(ctx, config)

	c.timeout = config.GetAsIntegerWithDefault("options.timeout", c.timeout)
	c.retries = config.GetAsIntegerWithDefault("options.retries", c.retries)
	c.dbNum = config.GetAsIntegerWithDefault("options.db_num", c.dbNum)
	if c.dbNum > 15 || c.dbNum < 0 {
		c.dbNum = 0
	}
}

// SetReferences method are sets references to dependent components.
// Parameters:
//   - ctx context.Context
//   - references 	references to locate the component dependencies.
func (c *RedisLock) SetReferences(ctx context.Context, references cref.IReferences) {
	c.connectionResolver.SetReferences(ctx, references)
	c.credentialResolver.SetReferences(ctx, references)
}

// IsOpen method are checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *RedisLock) IsOpen() bool {
	return c.client != nil
}

// Open method are opens the component.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//
// Returns: error or nil no errors occured.
func (c *RedisLock) Open(ctx context.Context) error {
	var connection *ccon.ConnectionParams
	var credential *cauth.CredentialParams

	connection, err := c.connectionResolver.Resolve(ctx)

	if err == nil && connection == nil {
		err = cerr.NewConfigError(cctx.GetTraceId(ctx), "NO_CONNECTION", "Connection is not configured")
		return err
	}

	credential, err = c.credentialResolver.Lookup(ctx)
	if err != nil {
		return err
	}

	var url, host, port, password string
	var dialOpts []redis.DialOption = make([]redis.DialOption, 0)

	dialOpts = append(dialOpts, redis.DialConnectTimeout(time.Duration(c.timeout)*time.Millisecond))
	dialOpts = append(dialOpts, redis.DialDatabase(c.dbNum))

	if credential != nil {
		password = credential.Password()
		dialOpts = append(dialOpts, redis.DialPassword(password))
	}

	if connection.Uri() != "" {
		url = connection.Uri()
		c.client, err = redis.DialURL(url, dialOpts...)
	} else {
		host = connection.Host()
		if host == "" {
			host = "localhost"
		}
		port = strconv.FormatInt(int64(connection.Port()), 10)
		if port == "0" {
			port = "6379"
		}
		url = host + ":" + port
		c.client, err = redis.Dial("tcp", url, dialOpts...)
	}
	return err
}

// Close method are closes component and frees used resources.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//
// Retruns: error or nil no errors occured.
func (c *RedisLock) Close(ctx context.Context) error {
	if c.client != nil {
		err := c.client.Close()
		c.client = nil
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RedisLock) checkOpened(traceId string) (state bool, err error) {
	if !c.IsOpen() {
		err = cerr.NewInvalidStateError(traceId, "NOT_OPENED", "Connection is not opened")
		return false, err
	}

	return true, nil
}

// TryAcquireLock method are makes a single attempt to acquire a lock by its key.
// It returns immediately a positive or negative result.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique lock key to acquire.
//   - ttl               a lock timeout (time to live) in milliseconds.
//
// Returns: a lock result or error.
func (c *RedisLock) TryAcquireLock(ctx context.Context, key string, ttl int64) (result bool, err error) {
	state, err := c.checkOpened(cctx.GetTraceId(ctx))
	if !state {
		return false, err
	}

	res, err := redis.String(c.client.Do("SET", key, c.lockId, "NX", "PX", ttl))
	if err != nil && err == redis.ErrNil {
		return false, nil
	}
	return res == "OK", err
}

// ReleaseLock method are releases prevously acquired lock by its key.
//   - ctx context.Context transaction id to trace execution through call chain.
//   - key               a unique lock key to release.
//
// Returns: error or nil for success.
func (c *RedisLock) ReleaseLock(ctx context.Context, key string) error {
	state, err := c.checkOpened(cctx.GetTraceId(ctx))
	if !state {
		return err
	}

	// Start transaction on key
	_, err = c.client.Do("WATCH", key)
	if err != nil {
		return err
	}

	// Read and check if lock is the same
	keyId, err := redis.String(c.client.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		c.client.Do("UNWATCH")
		return err
	}
	// Remove the lock if it matches
	if keyId == c.lockId {
		c.client.Send("MULTI")
		c.client.Send("DEL", key)
		_, err = c.client.Do("EXEC")
	} else { // Cancel transaction if it doesn"t match
		_, err = c.client.Do("UNWATCH")
	}
	return err
}
