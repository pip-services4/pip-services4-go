package connect

import (
	"context"
	"strconv"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	ccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// MongoDbConnectionResolver a helper struct  that resolves MongoDB connection and credential parameters,
// validates them and generates a connection URI.
// It is able to process multiple connections to MongoDB cluster nodes.
//
//	Configuration parameters
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
//	References
//		- *:discovery:*:*:1.0             (optional) IDiscovery services
//		- *:credential-store:*:*:1.0      (optional) Credential stores to resolve credentials
type MongoDbConnectionResolver struct {
	//The connections resolver.
	ConnectionResolver ccon.ConnectionResolver
	//The credentials resolver.
	CredentialResolver auth.CredentialResolver
}

// NewMongoDbConnectionResolver creates new connection resolver
//
//	Returns: *MongoDbConnectionResolver
func NewMongoDbConnectionResolver() *MongoDbConnectionResolver {
	mongoCon := MongoDbConnectionResolver{}
	mongoCon.ConnectionResolver = *ccon.NewEmptyConnectionResolver()
	mongoCon.CredentialResolver = *auth.NewEmptyCredentialResolver()
	return &mongoCon
}

// Configure is configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams configuration parameters to be set.
func (c *MongoDbConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.CredentialResolver.Configure(ctx, config)
}

// SetReferences is sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context,
//		- references crefer.IReferences references to locate the component dependencies.
func (c *MongoDbConnectionResolver) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.ConnectionResolver.SetReferences(ctx, references)
	c.CredentialResolver.SetReferences(ctx, references)
}

func (c *MongoDbConnectionResolver) validateConnection(traceId string, connection *ccon.ConnectionParams) error {
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
	if database, ok := connection.GetAsNullableString("database"); !ok || database == "" {
		return cerr.NewConfigError(traceId, "NO_DATABASE", "Connection database is not set")
	}
	return nil
}

func (c *MongoDbConnectionResolver) validateConnections(traceId string, connections []*ccon.ConnectionParams) error {
	if connections == nil || len(connections) == 0 {
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

func (c *MongoDbConnectionResolver) composeUri(connections []*ccon.ConnectionParams, credential *auth.CredentialParams) string {
	// Define hosts
	hosts := ""
	// Define database
	database := ""
	// Define additional parameters
	consConf := cdata.NewEmptyStringValueMap()

	for _, connection := range connections {
		uri := connection.Uri()
		// If there is a uri then return it immediately
		if uri != "" {
			return uri
		}

		host := connection.Host()
		port := connection.Port()

		if len(hosts) > 0 {
			hosts += ","
		}
		if port != 0 {
			hosts += host + ":" + strconv.Itoa(port)
		} else {
			hosts += host
		}

		// Take database name from the first connection that has it
		if database == "" {
			database = connection.GetAsString("database")
		}

		consConf.Append(connection.Value())
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
		params = "?" + params
	}

	// Compose uri
	uri := "mongodb://" + auth + hosts + database + params

	return uri
}

// Resolve method are resolves MongoDB connection URI from connection and credential parameters.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: uri string, err error resolved URI and error, if this occured.
func (c *MongoDbConnectionResolver) Resolve(ctx context.Context) (uri string, err error) {
	connections, err := c.ConnectionResolver.ResolveAll(ctx)
	if err != nil {
		return "", err
	}
	//Validate connections
	err = c.validateConnections(cctx.GetTraceId(ctx), connections)
	if err != nil {
		return "", err
	}
	credential, _ := c.CredentialResolver.Lookup(ctx)
	return c.composeUri(connections, credential), nil
}
