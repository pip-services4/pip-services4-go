package connect

import (
	"context"
	"math"
	"time"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"

	"github.com/jackc/pgx/v4/pgxpool"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// PostgresConnection is a PostgreSQL connection using plain driver.
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
type PostgresConnection struct {
	defaultConfig *cconf.ConfigParams
	// The logger.
	Logger *clog.CompositeLogger
	// The connection resolver.
	ConnectionResolver *PostgresConnectionResolver
	// The configuration options.
	Options *cconf.ConfigParams
	// The PostgreSQL connection pool object.
	Connection *pgxpool.Pool
	// The PostgreSQL database name.
	DatabaseName string

	retries int
}

const (
	DefaultConnectTimeout = 1000
	DefaultIdleTimeout    = 10000
	DefaultMaxPoolSize    = 3
	DefaultRetriesCount   = 3
)

// NewPostgresConnection creates a new instance of the connection component.
func NewPostgresConnection() *PostgresConnection {
	c := &PostgresConnection{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"options.connect_timeout", DefaultConnectTimeout,
			"options.idle_timeout", DefaultIdleTimeout,
			"options.max_pool_size", DefaultMaxPoolSize,
		),
		Logger:             clog.NewCompositeLogger(),
		ConnectionResolver: NewPostgresConnectionResolver(),
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
func (c *PostgresConnection) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.ConnectionResolver.Configure(ctx, config)
	c.Options = c.Options.Override(config.GetSection("options"))
}

// SetReferences references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references references to locate the component dependencies.
func (c *PostgresConnection) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// IsOpen checks if the component is opened.
//
//	Returns true if the component has been opened and false otherwise.
func (c *PostgresConnection) IsOpen() bool {
	return c.Connection != nil
}

// Open the component.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - Return 			error or nil no errors occurred.
func (c *PostgresConnection) Open(ctx context.Context) error {

	uri, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to resolve Postgres connection")
		return nil
	}

	maxPoolSize := c.Options.GetAsIntegerWithDefault("max_pool_size", DefaultMaxPoolSize)
	idleTimeoutMS := c.Options.GetAsIntegerWithDefault("idle_timeout", DefaultIdleTimeout)
	connectTimeoutMS := c.Options.GetAsIntegerWithDefault("connect_timeout", DefaultConnectTimeout)

	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to parse Postgres config string")
		return nil
	}

	if connectTimeoutMS > 0 {
		config.ConnConfig.ConnectTimeout = time.Duration((int64)(connectTimeoutMS)) * time.Millisecond
	}
	if idleTimeoutMS > 0 {
		config.MaxConnIdleTime = time.Duration((int64)(idleTimeoutMS)) * time.Millisecond
	}
	if maxPoolSize > 0 {
		config.MaxConns = (int32)(maxPoolSize)
	}

	c.Logger.Debug(ctx, "Connecting to postgres")

	retries := c.retries
	for retries > 0 {
		pool, err := pgxpool.ConnectConfig(ctx, config)
		if err != nil {
			retries--
			if retries <= 0 {
				return cerr.
					NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Connection to postgres failed").
					WithCause(err)
			}
			c.Logger.Debug(ctx, "Failed to connect to postgress, try reconnect...")
			err = c.waitForRetry(ctx, retries)
			if err != nil {
				return err
			}
			continue
		}
		c.Connection = pool
		c.DatabaseName = config.ConnConfig.Database
		break
	}
	return nil
}

// Close component and frees used resources.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil no errors occurred
func (c *PostgresConnection) Close(ctx context.Context) error {
	if c.Connection == nil {
		return nil
	}
	c.Connection.Close()
	c.Logger.Debug(ctx, "Disconnected from postgres database %s", c.DatabaseName)
	c.Connection = nil
	c.DatabaseName = ""
	return nil
}

func (c *PostgresConnection) GetConnection() *pgxpool.Pool {
	return c.Connection
}

func (c *PostgresConnection) GetDatabaseName() string {
	return c.DatabaseName
}

func (c *PostgresConnection) waitForRetry(ctx context.Context, retries int) error {
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
