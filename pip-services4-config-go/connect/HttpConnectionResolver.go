package connect

import (
	"context"
	"net/url"
	"strconv"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
)

// HttpConnectionResolver helper class to retrieve connections for HTTP-based services abd clients.
//
// In addition to regular functions of ConnectionResolver is able to parse http:// URIs
// and validate connection parameters before returning them.
//
//	Configuration parameters:
//		- connection:
//		- discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//		- ...                          other connection parameters
//
//		- connections:                   alternative to connection
//			- [connection params 1]:       first connection parameters
//			-  ...
//			- [connection params N]:       Nth connection parameters
//			-  ...
//
//	References:
//		- *:discovery:*:*:1.0            (optional) IDiscovery services
//
//	see: ConnectionParams
//	see: ConnectionResolver
//
//	Example:
//		config := cconf.NewConfigParamsFromTuples(
//			"connection.host", "10.1.1.100",
//			"connection.port", 8080,
//		);
//
//		connectionResolver = NewHttpConnectionResolver();
//		connectionResolver.Configure(context.Background(), config);
//		connectionResolver.SetReferences(context.Background(), references);
//
//		connection, err := connectionResolver.Resolve(ctx)
//		// Now use connection...
type HttpConnectionResolver struct {
	//The base connection resolver.
	ConnectionResolver *ConnectionResolver
	//The base credential resolver.
	CredentialResolver *cauth.CredentialResolver
}

// NewHttpConnectionResolver creates new instance NewHttpConnectionResolver
//
//	Returns: pointer on NewHttpConnectionResolver
func NewHttpConnectionResolver() *HttpConnectionResolver {
	return &HttpConnectionResolver{NewEmptyConnectionResolver(), cauth.NewEmptyCredentialResolver()}
}

// Configure method are configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams  configuration parameters to be set.
func (c *HttpConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.CredentialResolver.Configure(ctx, config)
}

// SetReferences method are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references crefer.IReferences	references to locate the component dependencies.
func (c *HttpConnectionResolver) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.ConnectionResolver.SetReferences(ctx, references)
	c.CredentialResolver.SetReferences(ctx, references)
}

func (c *HttpConnectionResolver) validateConnection(ctx context.Context, connection *ConnectionParams, credential *cauth.CredentialParams) error {
	traceId := utils.ContextHelper.GetTraceId(ctx)
	if connection == nil {
		return cerr.NewConfigError(traceId, "NO_CONNECTION", "HTTP connection is not set")
	}
	uri := connection.Uri()
	if uri != "" {
		return nil
	}

	protocol := connection.Protocol() // "http"
	if protocol != "http" && protocol != "https" {
		return cerr.NewConfigError(traceId, "WRONG_PROTOCOL", "Protocol is not supported by REST connection").WithDetails("protocol", protocol)
	}
	host := connection.Host()
	if host == "" {
		return cerr.NewConfigError(traceId, "NO_HOST", "Connection host is not set")
	}
	port := connection.Port()
	if port == 0 {
		return cerr.NewConfigError(traceId, "NO_PORT", "Connection port is not set")
	}
	// Check HTTPS credentials
	if protocol == "https" {
		// Check for credential
		if credential == nil {
			return cerr.NewConfigError(traceId, "NO_CREDENTIAL", "SSL certificates are not configured for HTTPS protocol")
		} else {
			if _, ok := credential.GetAsNullableString("internal_network"); !ok {
				if _, ok := credential.GetAsNullableString("ssl_key_file"); !ok {
					return cerr.NewConfigError(
						traceId, "NO_SSL_KEY_FILE", "SSL key file is not configured in credentials")
				} else if _, ok := credential.GetAsNullableString("ssl_crt_file"); !ok {
					return cerr.NewConfigError(
						traceId, "NO_SSL_CRT_FILE", "SSL crt file is not configured in credentials")
				}
			}
		}
	}

	return nil
}

func (c *HttpConnectionResolver) updateConnection(connection *ConnectionParams) {
	if connection == nil {
		return
	}

	uri := connection.Uri()

	if uri == "" {
		protocol := connection.Protocol() // "http"
		host := connection.Host()
		port := connection.Port()

		uri := protocol + "://" + host
		if port != 0 {
			uri += ":" + strconv.Itoa(port)
		}
		connection.SetUri(uri)
	} else {
		address, _ := url.Parse(uri)
		//protocol := ("" + address.protocol).replace(":", "")
		protocol := address.Scheme

		connection.SetProtocol(protocol)
		connection.SetHost(address.Hostname())
		port, _ := strconv.Atoi(address.Port())
		connection.SetPort(port)
	}
}

// Resolve method are resolves a single component connection. If connections are configured to be retrieved
// from Discovery service it finds a IDiscovery and resolves the connection there.
//
//	Parameters:
//		- correlationId  string     (optional) transaction id to trace execution through call chain.
//	Returns: connection *ConnectionParams, credential *cauth.CredentialParams, err error
//		resolved connection and credential or error.
func (c *HttpConnectionResolver) Resolve(ctx context.Context) (connection *ConnectionParams, credential *cauth.CredentialParams, err error) {

	connection, err = c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		return nil, nil, err
	}

	credential, err = c.CredentialResolver.Lookup(ctx)
	if err != nil {
		return nil, nil, err
	}
	err = c.validateConnection(ctx, connection, credential)
	if err == nil && connection != nil {
		c.updateConnection(connection)
	}

	return connection, credential, err
}

// ResolveAll method are resolves all component connection. If connections are configured to be retrieved
// from Discovery service it finds a IDiscovery and resolves the connection there.
//
//	Parameters:
//		- correlationId  string   (optional) transaction id to trace execution through call chain.
//	Returns:  connections []*ConnectionParams, credential *cauth.CredentialParams, err error
//		resolved connections and credential or error.
func (c *HttpConnectionResolver) ResolveAll(ctx context.Context) (connections []*ConnectionParams, credential *cauth.CredentialParams, err error) {

	connections, err = c.ConnectionResolver.ResolveAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	credential, err = c.CredentialResolver.Lookup(ctx)
	if connections == nil {
		connections = make([]*ConnectionParams, 0)
	}

	for _, connection := range connections {
		if err == nil {
			err = c.validateConnection(ctx, connection, credential)
		}
		if err == nil && connection != nil {
			c.updateConnection(connection)
		}
	}
	return connections, credential, err
}

// Register method are registers the given connection in all referenced discovery services.
// c method can be used for dynamic service discovery.
//
//	Parameters:
//		- correlationId  string   (optional) transaction id to trace execution through call chain.
//	Returns: error nil if registered connection or error.
func (c *HttpConnectionResolver) Register(ctx context.Context) error {

	connection, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		return err
	}

	credential, err := c.CredentialResolver.Lookup(ctx)
	// Validate connection
	if err == nil {
		err = c.validateConnection(ctx, connection, credential)
	}
	if err == nil {
		return c.ConnectionResolver.Register(ctx, connection)
	} else {
		return err
	}
}
