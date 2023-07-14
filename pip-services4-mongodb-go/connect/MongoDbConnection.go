package connect

import (
	"context"
	"time"

	cerror "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	mongodrv "go.mongodb.org/mongo-driver/mongo"
	mongoclopt "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// MongoDbConnection struct help creates new connections to MongoDB
// MongoDB connection using plain driver.
//
// By defining a connection and sharing it through multiple persistence components
// you can reduce number of used database connections.
//
//	Configuration parameters:
//		- connection(s):
//			- discovery_key:             (optional) a key to retrieve the connection from IDiscovery
//			- host:                      host name or IP address
//			- port:                      port number (default: 27017)
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
type MongoDbConnection struct {
	defaultConfig *cconf.ConfigParams
	// The logger.
	Logger *clog.CompositeLogger
	//   The connection resolver.
	ConnectionResolver *MongoDbConnectionResolver
	//   The configuration options.
	Options *cconf.ConfigParams
	//   The MongoDB connection object.
	Connection *mongodrv.Client
	//   The MongoDB database name.
	DatabaseName string
	//   The MongoDb database object.
	Db *mongodrv.Database
}

// NewMongoDbConnection are creates a new instance of the connection component.
// Returns *MongoDbConnection with default config
func NewMongoDbConnection() *MongoDbConnection {
	c := MongoDbConnection{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"options.max_pool_size", "2",
			"options.keep_alive", "0",
			"options.connect_timeout", "5000",
			"options.max_page_size", "100",
		),
		//The logger.
		Logger: clog.NewCompositeLogger(),
		//The connection resolver.
		ConnectionResolver: NewMongoDbConnectionResolver(),
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
func (c *MongoDbConnection) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.ConnectionResolver.Configure(ctx, config)
	c.Options = c.Options.Override(config.GetSection("options"))
}

// SetReferences are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references crefer.IReferences references to locate the component dependencies.
func (c *MongoDbConnection) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// IsOpen method is checks if the component is opened.
//
//	Returns: true if the component has been opened and false otherwise.
func (c *MongoDbConnection) IsOpen() bool {
	return c.Connection != nil
}

func (c *MongoDbConnection) composeSettings(settings *mongoclopt.ClientOptions) {
	maxPoolSize := (uint64)(c.Options.GetAsInteger("max_pool_size"))
	keepAlive := c.Options.GetAsInteger("keep_alive")
	MaxConnIdleTime := (time.Duration)(keepAlive) * time.Millisecond
	connectTimeoutMS := c.Options.GetAsInteger("connect_timeout")
	ConnectTimeout := (time.Duration)(connectTimeoutMS) * time.Millisecond
	socketTimeoutMS := c.Options.GetAsInteger("socket_timeout")
	SocketTimeout := (time.Duration)(socketTimeoutMS) * time.Millisecond

	replicaSet := c.Options.GetAsString("replica_set")
	authSource := c.Options.GetAsString("auth_source")
	authUser := c.Options.GetAsString("auth_user")
	authPassword := c.Options.GetAsString("auth_password")

	settings.SetMaxPoolSize(maxPoolSize)
	settings.SetMaxConnIdleTime(MaxConnIdleTime)
	settings.SetConnectTimeout(ConnectTimeout)
	settings.SetSocketTimeout(SocketTimeout)

	if replicaSet != "" {
		settings.SetReplicaSet(replicaSet)
	}

	// TODO: Relase configure TLS(SSL) connection to MongoDB
	// ssl := c.Options.GetAsNullableBoolean("ssl")
	// if ssl != nil {
	// 	settings.ssl = ssl
	// }

	// Auth params
	if authSource != "" && authUser != "" && authPassword != "" {
		authParams := mongoclopt.Credential{
			AuthSource: authSource,
			Username:   authUser,
			Password:   authPassword,
		}
		settings.SetAuth(authParams)
	}
}

// Open method is opens the component.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil when no errors occurred.
func (c *MongoDbConnection) Open(ctx context.Context) error {
	uri, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to resolve MongoDb connection")
		return err
	}
	c.Logger.Debug(ctx, "Connecting to mongodb")

	settings := mongoclopt.Client()
	settings.ApplyURI(uri)
	c.composeSettings(settings)

	//settings.useNewUrlParser = true;
	//settings.useUnifiedTopology = true;

	client, err := mongodrv.NewClient(settings)

	if err != nil {
		err = cerror.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Create client for mongodb failed").WithCause(err)
		return err
	}
	cs, _ := connstring.Parse(uri)
	c.DatabaseName = cs.Database
	err = client.Connect(ctx)
	if err != nil {
		err = cerror.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Connection to mongodb failed").WithCause(err)
		return err
	}
	c.Connection = client
	c.Db = client.Database(c.DatabaseName)
	return nil
}

// Close method is closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: error or nil when no errors occurred.
func (c *MongoDbConnection) Close(ctx context.Context) error {
	if c.Connection == nil {
		return nil
	}

	err := c.Connection.Disconnect(ctx)

	if err != nil {
		return cerror.NewConnectionError(cctx.GetTraceId(ctx), "DISCONNECT_FAILED", "Disconnect from mongodb failed: ").WithCause(err)
	}

	c.Logger.Debug(ctx, "Disconnected from mongodb database %s", c.DatabaseName)
	c.Connection = nil
	c.Db = nil
	c.DatabaseName = ""

	return err
}

// GetConnection method return work connection object
//
//	Returns: *mongodrv.Client
func (c *MongoDbConnection) GetConnection() *mongodrv.Client {
	return c.Connection
}

// GetDatabase method retrun work database object
//
//	Returns: *mongodrv.Database
func (c *MongoDbConnection) GetDatabase() *mongodrv.Database {
	return c.Db
}

// GetDatabaseName method reruns name of work database
//
//	Returns: string
func (c *MongoDbConnection) GetDatabaseName() string {
	return c.DatabaseName
}
