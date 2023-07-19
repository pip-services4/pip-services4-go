package connect

import (
	"context"
	"net/url"
	"strconv"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	cconn "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// Helper class that resolves SqlServer connection and credential parameters,
// validates them and generates a connection URI.
// It is able to process multiple connections to SqlServer cluster nodes.
//
//	Configuration parameters:
//		- connection(s):
//			- discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//			- host:                        host name or IP address
//			- port:                        port number (default: 27017)
//			- database:                    database name
//			- uri:                         resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//			- username:                    user name
//			- password:                    user password
//
//	References:
//		- *:discovery:*:*:1.0             (optional) IDiscovery services
//		- *:credential-store:*:*:1.0      (optional) Credential stores to resolve credentials
type SqlServerConnectionResolver struct {
	// The connections' resolver.
	ConnectionResolver *cconn.ConnectionResolver
	// The credentials' resolver.
	CredentialResolver *cauth.CredentialResolver
}

// NewSqlServerConnectionResolver creates new connection resolver
//
//	Returns: *SqlServerConnectionResolver
func NewSqlServerConnectionResolver() *SqlServerConnectionResolver {
	mongoCon := SqlServerConnectionResolver{}
	mongoCon.ConnectionResolver = cconn.NewEmptyConnectionResolver()
	mongoCon.CredentialResolver = cauth.NewEmptyCredentialResolver()
	return &mongoCon
}

// Configure is configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *cconf.ConfigParams configuration parameters to be set.
func (c *SqlServerConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.CredentialResolver.Configure(ctx, config)
}

// SetReferences is sets references to dependent components.
// Parameters:
//   - ctx context.Context
//   - references crefer.IReferences references to locate the component dependencies.
func (c *SqlServerConnectionResolver) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.ConnectionResolver.SetReferences(ctx, references)
	c.CredentialResolver.SetReferences(ctx, references)
}

func (c *SqlServerConnectionResolver) validateConnection(traceId string, connection *cconn.ConnectionParams) error {
	uri := connection.Uri()
	if uri != "" {
		return nil
	}

	host := connection.Host()
	if host == "" {
		return cerr.NewConfigError(traceId, "NO_HOST", "Connection host is not set")
	}
	port := connection.Port()
	if port == 0 {
		return cerr.NewConfigError(traceId, "NO_PORT", "Connection port is not set")
	}
	database, ok := connection.GetAsNullableString("database")
	if !ok || database == "" {
		return cerr.NewConfigError(traceId, "NO_DATABASE", "Connection database is not set")
	}
	return nil
}

func (c *SqlServerConnectionResolver) validateConnections(traceId string, connections []*cconn.ConnectionParams) error {
	if len(connections) == 0 {
		return cerr.NewConfigError(traceId, "NO_CONNECTION", "Database connection is not set")
	}
	for _, connection := range connections {
		err := c.validateConnection(traceId, connection)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *SqlServerConnectionResolver) composeUri(connections []*cconn.ConnectionParams,
	credential *cauth.CredentialParams) string {

	// If there is an uri then return it immediately
	for _, connection := range connections {
		uri := connection.Uri()
		if uri != "" {
			return uri
		}
	}

	// Define hosts
	var hosts = ""
	for _, connection := range connections {
		host := connection.Host()
		port := connection.Port()

		if len(hosts) > 0 {
			hosts += ","
		}
		if port != 0 {
			hosts += host + ":" + strconv.Itoa(port)
		}
	}

	// Define database
	database := ""
	for _, connection := range connections {
		if database == "" {
			database, _ = connection.GetAsNullableString("database")
		}
	}
	if len(database) > 0 {
		database = "/" + database
	}

	// Define authentication part
	var auth = ""
	if credential != nil {
		var username = credential.Username()
		if len(username) > 0 {
			var password = credential.Password()
			if len(password) > 0 {
				auth = username + ":" + password + "@"
			} else {
				auth = username + "@"
			}
		}
	}
	// Define additional parameters
	consConf := cdata.NewEmptyStringValueMap()
	for _, v := range connections {
		consConf.Append(v.Value())
	}
	var options *cconf.ConfigParams
	if credential != nil {
		options = cconf.NewConfigParamsFromMaps(consConf.Value(), credential.Value())
	} else {
		options = cconf.NewConfigParamsFromValue(consConf.Value())
	}
	options.Remove("uri")
	options.Remove("host")
	options.Remove("port")
	options.Remove("database")
	options.Remove("username")
	options.Remove("password")
	params := ""
	keys := options.Keys()
	for _, key := range keys {
		if len(params) > 0 {
			params += "&"
		}
		params += key

		value := options.GetAsString(key)
		if value != "" {
			params += "=" + value
		}
	}
	if len(params) > 0 {
		params = "?" + url.PathEscape(params)
	}

	// Compose uri

	uri := "sqlserver://" + auth + hosts + database + params

	return uri
}

// Resolve method are resolves SqlServer connection URI from connection and credential parameters.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: uri string, err error resolved URI and error, if this occured.
func (c *SqlServerConnectionResolver) Resolve(ctx context.Context) (uri string, err error) {

	connections, err := c.ConnectionResolver.ResolveAll(ctx)
	// Validate connections
	if err != nil {
		return "", err
	}
	err = c.validateConnections(cctx.GetTraceId(ctx), connections)
	if err != nil {
		return "", err
	}
	credential, err := c.CredentialResolver.Lookup(ctx)
	if err != nil {
		return "", err
	}
	return c.composeUri(connections, credential), nil
}
