package connect

import (
	"context"
	"database/sql"
	"math"
	"time"

	_ "github.com/go-sql-driver/mysql"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// MySQL connection using plain driver.
//
// By defining a connection and sharing it through multiple persistence components
// you can reduce number of used database connections.
//
//	Configuration parameters
//		- connection(s):
//			- discovery_key:        (optional) a key to retrieve the connection from IDiscovery
//			- host:                 host name or IP address
//			- port:                 port number (default: 27017)
//			- uri:                  resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:            (optional) a key to retrieve the credentials from ICredentialStore
//			- username:             user name
//			- password:             user password
//		- options:
//			- connect_timeout:      (optional) number of milliseconds to wait before timing out when connecting a new client (default: 0)
//			- idle_timeout:         (optional) number of milliseconds a client must sit idle in the pool and not be checked out (default: 10000)
//			- max_pool_size:        (optional) maximum number of clients the pool should contain (default: 10)
//
//	References
//		- *:logger:*:*:1.0           (optional) ILogger components to pass log messages
//		- *:discovery:*:*:1.0        (optional) IDiscovery services
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
type MySqlConnection struct {
	defaultConfig *cconf.ConfigParams
	// The logger.
	Logger *clog.CompositeLogger
	// The connection resolver.
	ConnectionResolver *MySqlConnectionResolver
	// The configuration options.
	Options *cconf.ConfigParams
	// The MySQL connection pool object.
	Connection *sql.DB
	// The MySQL database name.
	DatabaseName string

	retries int
}

const (
	DefaultConnectTimeout = 1000
	DefaultIdleTimeout    = 10000
	DefaultMaxPoolSize    = 3
	DefaultRetriesCount   = 3
)

// NewMySqlConnection creates a new instance of the connection component.
func NewMySqlConnection() *MySqlConnection {
	c := &MySqlConnection{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"options.connect_timeout", DefaultConnectTimeout,
			"options.idle_timeout", DefaultIdleTimeout,
			"options.max_pool_size", DefaultMaxPoolSize,
		),
		Logger:             clog.NewCompositeLogger(),
		ConnectionResolver: NewMySqlConnectionResolver(),
		Options:            cconf.NewEmptyConfigParams(),
		retries:            DefaultRetriesCount,
	}
	return c
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config configuration parameters to be set.
func (c *MySqlConnection) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.ConnectionResolver.Configure(ctx, config)
	c.Options = c.Options.Override(config.GetSection("options"))

	c.DatabaseName, _ = config.GetAsNullableString("connection.database")
}

// SetReferences references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references references to locate the component dependencies.
func (c *MySqlConnection) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// IsOpen checks if the component is opened.
//
//	Returns true if the component has been opened and false otherwise.
func (c *MySqlConnection) IsOpen() bool {
	return c.Connection != nil
}

// Open the component.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - Return 			error or nil no errors occurred.
func (c *MySqlConnection) Open(ctx context.Context) error {

	uri, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to resolve MySql connection")
		return nil
	}

	c.Logger.Debug(ctx, "Connecting to mysql")

	retries := c.retries
	for retries > 0 {
		pool, err := sql.Open("mysql", uri)
		if err != nil {
			retries--
			if retries <= 0 {
				return cerr.
					NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Connection to mysql failed").
					WithCause(err)
			}
			c.Logger.Debug(ctx, "Failed to connect to mysqls, try reconnect...")
			err = c.waitForRetry(ctx, retries)
			if err != nil {
				return err
			}
			continue
		}
		idleTimeoutMS := c.Options.GetAsIntegerWithDefault("idle_timeout", DefaultIdleTimeout)
		maxPoolSize := c.Options.GetAsIntegerWithDefault("max_pool_size", DefaultMaxPoolSize)
		connectTimeoutMS := c.Options.GetAsIntegerWithDefault("connect_timeout", DefaultConnectTimeout)

		pool.SetConnMaxIdleTime(time.Duration(idleTimeoutMS) * time.Millisecond)
		pool.SetMaxOpenConns(maxPoolSize)
		pool.SetConnMaxLifetime(time.Duration(connectTimeoutMS) * time.Millisecond)

		c.Connection = pool
		break
	}
	return nil
}

// Close component and frees used resources.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil no errors occurred
func (c *MySqlConnection) Close(ctx context.Context) error {
	if c.Connection == nil {
		return nil
	}
	c.Connection.Close()
	c.Logger.Debug(ctx, "Disconnected from mysql database %s", c.DatabaseName)
	c.Connection = nil
	c.DatabaseName = ""
	return nil
}

func (c *MySqlConnection) GetConnection() *sql.DB {
	return c.Connection
}

func (c *MySqlConnection) GetDatabaseName() string {
	return c.DatabaseName
}

func (c *MySqlConnection) waitForRetry(ctx context.Context, retries int) error {
	waitTime := DefaultConnectTimeout * int(math.Pow(float64(c.retries-retries), 2))

	select {
	case <-time.After(time.Duration(waitTime) * time.Millisecond):
		return nil
	case <-ctx.Done():
		return cerr.ApplicationErrorFactory.Create(
			&cerr.ErrorDescription{
				Type:     "Application",
				Category: "Application",
				Code:     "CONTEXT_CANCELLED",
				Message:  "request canceled by parent context",
				TraceId:  cctx.GetTraceId(ctx),
			},
		)
	}
}
