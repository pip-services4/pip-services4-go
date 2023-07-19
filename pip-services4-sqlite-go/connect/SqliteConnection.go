package connect

import (
	"context"

	cerror "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// SqliteConnection struct help creates new connections to SQLite
// SQLite connection using plain driver.
//
// By defining a connection and sharing it through multiple persistence components
// you can reduce number of used database connections.
//
//	Configuration parameters:
//		- connection(s):
//			- discovery_key:             (optional) a key to retrieve the connection from IDiscovery
//			- host:                      host name or IP address
//			- port:                      port number (default: 8082)
//			- uri:                       resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                 (optional) a key to retrieve the credentials from ICredentialStore
//			- username:                  (optional) user name
//			- password:                  (optional) user password
//		- options:
//			- max_pool_size:             (optional) maximum connection pool size (default: 2)
//			- keep_alive:                (optional) enable connection keep alive in ms, if zero connection are keeped indefinitely (default: 0)
//			- connect_timeout:           (optional) connection timeout in milliseconds (default: 5000)
//			- socket_timeout:            (optional) socket timeout in milliseconds (default: 360000)
//			- auto_reconnect:            (optional) enable auto reconnection (default: true) (Not used)
//			- reconnect_interval:        (optional) reconnection interval in milliseconds (default: 1000) (Not used)
//			- max_page_size:             (optional) maximum page size (default: 100)
//			- replica_set:               (optional) name of replica set
//			- ssl:                       (optional) enable SSL connection (default: false) (Not release in this version)
//			- auth_source:               (optional) authentication source
//			- debug:                     (optional) enable debug output (default: false). (Not used)
//
//	References:
//		- *:logger:*:*:1.0           (optional) ILogger components to pass log messages
//		- *:discovery:*:*:1.0        (optional) IDiscovery services
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
type SqliteConnection struct {
	defaultConfig *cconf.ConfigParams
	// The logger.
	Logger *clog.CompositeLogger
	//   The connection resolver.
	ConnectionResolver *SqliteConnectionResolver
	//   The configuration options.
	Options *cconf.ConfigParams
	//   The SQLite connection object.
	Connection *sql.DB
	//   The SQLite database name.
	DatabaseName string
	//   The Sqlite database object.
	// Db *mongodrv.Database
}

// NewSqliteConnection are creates a new instance of the connection component.
// Returns *SqliteConnection with default config
func NewSqliteConnection() *SqliteConnection {
	c := SqliteConnection{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"options.max_pool_size", "2",
			"options.keep_alive", "0",
			"options.connect_timeout", "5000",
			"options.max_page_size", "100",
		),
		//The logger.
		Logger: clog.NewCompositeLogger(),
		//The connection resolver.
		ConnectionResolver: NewSqliteConnectionResolver(),
		// The configuration options.
		Options: cconf.NewEmptyConfigParams(),
	}
	return &c
}

// Configure is configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams configuration parameters to be set.
func (c *SqliteConnection) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.ConnectionResolver.Configure(ctx, config)
	c.Options = c.Options.Override(config.GetSection("options"))
}

// SetReferences are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references crefer.IReferences references to locate the component dependencies.
func (c *SqliteConnection) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// IsOpen method is checks if the component is opened.
//
//	Returns: true if the component has been opened and false otherwise.
func (c *SqliteConnection) IsOpen() bool {
	return c.Connection != nil
}

// Open method is opens the component.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil when no errors occurred.
func (c *SqliteConnection) Open(ctx context.Context) error {
	uri, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to resolve sqlite connection")
		return err
	}
	c.Logger.Debug(ctx, "Connecting to sqlite")

	con, err := sql.Open("sqlite3", uri)

	if err != nil || con == nil {
		err = cerror.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Connection to sqlite failed").WithCause(err)
		return err
	}
	c.Connection = con
	c.DatabaseName = uri
	return nil
}

// Close method is closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil when no errors occurred.
func (c *SqliteConnection) Close(ctx context.Context) error {
	if c.Connection == nil {
		return nil
	}

	err := c.Connection.Close()

	if err != nil {
		return cerror.NewConnectionError(cctx.GetTraceId(ctx), "DISCONNECT_FAILED", "Disconnect from sqlite failed: ").WithCause(err)
	}

	c.Logger.Debug(ctx, "Disconnected from sqlite database %s", c.DatabaseName)
	c.Connection = nil
	c.DatabaseName = ""

	return err
}

// GetConnection method return work connection object
//
//	Returns: *sql.DB
func (c *SqliteConnection) GetConnection() *sql.DB {
	return c.Connection
}

// GetDatabase method retrun work database object
//
//	Returns: *sql.DB
func (c *SqliteConnection) GetDatabase() *sql.DB {
	return c.Connection
}

// GetDatabaseName method reruns name of work database
//
//	Returns: string
func (c *SqliteConnection) GetDatabaseName() string {
	return c.DatabaseName
}
