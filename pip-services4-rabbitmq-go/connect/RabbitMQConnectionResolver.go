package connect

import (
	"context"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	ccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// RabbitMQConnectionResolver helper class that resolves RabbitMQ connection and credential parameters,
// validates them and generates connection options.
//
//	Configuration parameters:
//
// - connection(s):
//   - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//   - host:                        host name or IP address
//   - port:                        port number
//   - uri:                         resource URI or connection string with all parameters in it
//
// - credential(s):
//
//   - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//
//   - username:                    user name
//
//   - password:                    user password
//
//     References:
//
// - *:discovery:*:*:1.0          (optional) IDiscovery services to resolve connections
// - *:credential-store:*:*:1.0   (optional) Credential stores to resolve credentials
type RabbitMQConnectionResolver struct {
	// The connections resolver.
	ConnectionResolver *ccon.ConnectionResolver
	//The credentials resolver.
	CredentialResolver *cauth.CredentialResolver
}

func NewRabbitMQConnectionResolver() *RabbitMQConnectionResolver {
	c := RabbitMQConnectionResolver{}
	c.ConnectionResolver = ccon.NewEmptyConnectionResolver()
	c.CredentialResolver = cauth.NewEmptyCredentialResolver()
	return &c
}

// Configure are configures component by passing configuration parameters.
// Parameters:
//   - ctx context.Context
//   - config   *cconf.ConfigParams
//
// configuration parameters to be set.
func (c *RabbitMQConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.CredentialResolver.Configure(ctx, config)
}

// SetReferences are sets references to dependent components.
// Parameters:
//   - ctx context.Context
//   - references  cref.IReferences
//
// references to locate the component dependencies.
func (c *RabbitMQConnectionResolver) SetReferences(ctx context.Context, references cref.IReferences) {
	c.ConnectionResolver.SetReferences(ctx, references)
	c.CredentialResolver.SetReferences(ctx, references)
}

func (c *RabbitMQConnectionResolver) validateConnection(traceId string, connection *ccon.ConnectionParams) error {
	if connection == nil {
		return cerr.NewConfigError(traceId, "NO_CONNECTION", "RabbitMQ connection is not set")
	}

	uri := connection.Uri()
	if uri != "" {
		return nil
	}

	protocol := connection.GetAsString("protocol")
	if protocol == "" {
		//return cerr.NewConfigError(traceId, "NO_PROTOCOL", "Connection protocol is not set")
		connection.SetAsObject("protocol", "amqp")
	}

	host := connection.Host()
	if host == "" {
		return cerr.NewConfigError(traceId, "NO_HOST", "Connection host is not set")
	}

	port := connection.Port()
	if port == 0 {
		return cerr.NewConfigError(traceId, "NO_PORT", "Connection port is not set")
	}

	return nil
}

func (c *RabbitMQConnectionResolver) composeOptions(connection *ccon.ConnectionParams, credential *cauth.CredentialParams) *cconf.ConfigParams {

	// Define additional parameters parameters
	if credential == nil {
		credential = cauth.NewEmptyCredentialParams()
	}
	options := connection.Override(credential.ConfigParams)

	// Compose uri
	if _, ok := options.Get("uri"); !ok {
		credential := ""
		if username, ok := options.Get("username"); ok {
			credential = username.(string)
		}
		if password, ok := options.Get("password"); ok {
			credential += ":" + password.(string)
		}
		uri := ""
		if credential == "" {
			uri = options.GetAsString("protocol") + "://" + options.GetAsString("host")
		} else {
			uri = options.GetAsString("protocol") + "://" + credential + "@" + options.GetAsString("host")
		}
		if _, ok := options.Get("port"); ok {
			uri = uri + ":" + options.GetAsString("port")
		}
		options.SetAsObject("uri", uri)
	}
	return options
}

// Resolves RabbitMQ connection options from connection and credential parameters.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//
// Retruns options *cconf.ConfigParams, err error
// receives resolved options or error.
func (c *RabbitMQConnectionResolver) Resolve(ctx context.Context) (options *cconf.ConfigParams, err error) {
	var connection *ccon.ConnectionParams
	var credential *cauth.CredentialParams
	var errCred, errConn error

	connection, errConn = c.ConnectionResolver.Resolve(ctx)
	// Validate connections
	if errConn == nil {
		errConn = c.validateConnection(cctx.GetTraceId(ctx), connection)
	}

	credential, errCred = c.CredentialResolver.Lookup(ctx)
	// Credentials are not validated right now

	if errConn != nil {
		return nil, errConn
	}
	if errCred != nil {
		return nil, errCred
	}
	options = c.composeOptions(connection, credential)
	return options, nil
}

// Compose method are composes RabbitMQ connection options from connection and credential parameters.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - connection  *ccon.ConnectionParams    connection parameters
//   - credential  *cauth.CredentialParams   credential parameters
//
// Returns: options *cconf.ConfigParams, err error
// resolved options or error.
func (c *RabbitMQConnectionResolver) Compose(ctx context.Context, connection *ccon.ConnectionParams, credential *cauth.CredentialParams) (options *cconf.ConfigParams, err error) {
	// Validate connections
	err = c.validateConnection(cctx.GetTraceId(ctx), connection)
	if err != nil {
		return nil, err
	} else {
		options := c.composeOptions(connection, credential)
		return options, nil
	}
}
