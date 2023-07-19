package connect

import (
	"context"
	"strings"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	ccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// SqliteConnectionResolver a helper struct  that resolves SQLite connection and credential parameters,
// validates them and generates a connection URI.
// It is able to process multiple connections to SQLite cluster nodes.
//
//	Configuration parameters
//		- connection(s):
//			- discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//			- host:                        host name or IP address
//			- port:                        port number (default: 8082)
//			- database:                    database name
//			- uri:                         resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//			- username:                    user name
//			- password:                    user password
//	References
//		- *:discovery:*:*:1.0             (optional) IDiscovery services
//		- *:credential-store:*:*:1.0      (optional) Credential stores to resolve credentials
type SqliteConnectionResolver struct {
	//The connections resolver.
	ConnectionResolver ccon.ConnectionResolver
	//The credentials resolver.
	CredentialResolver auth.CredentialResolver
}

// NewSqliteConnectionResolver creates new connection resolver
//
//	Returns: *SqliteConnectionResolver
func NewSqliteConnectionResolver() *SqliteConnectionResolver {
	sqliteCon := SqliteConnectionResolver{}
	sqliteCon.ConnectionResolver = *ccon.NewEmptyConnectionResolver()
	sqliteCon.CredentialResolver = *auth.NewEmptyCredentialResolver()
	return &sqliteCon
}

// Configure is configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams configuration parameters to be set.
func (c *SqliteConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.CredentialResolver.Configure(ctx, config)
}

// SetReferences is sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context,
//		- references crefer.IReferences references to locate the component dependencies.
func (c *SqliteConnectionResolver) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.ConnectionResolver.SetReferences(ctx, references)
	c.CredentialResolver.SetReferences(ctx, references)
}

func (c *SqliteConnectionResolver) validateConnection(traceId string, connection *ccon.ConnectionParams) error {
	uri := connection.Uri()
	if uri != "" {
		if !strings.HasPrefix(uri, "file://") {
			return cerr.NewConfigError(traceId, "WRONG_PROTOCOL", "Connection protocol must be file://")
		}
		return nil
	}

	if database, ok := connection.GetAsNullableString("database"); !ok || database == "" {
		return cerr.NewConfigError(traceId, "NO_DATABASE", "Connection database is not set")
	}
	return nil
}

func (c *SqliteConnectionResolver) validateConnections(traceId string, connections []*ccon.ConnectionParams) error {
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

func (c *SqliteConnectionResolver) composeUri(connections []*ccon.ConnectionParams, credential *auth.CredentialParams) string {
	// If there is a uri or database then return it immediately
	for _, connection := range connections {
		uri := connection.Uri()
		if uri != "" {
			// Removing file://
			return uri[7:]
		}

		database, isFetched := connection.GetAsNullableString("database")
		if isFetched && database != "" {
			return database
		}
	}

	return ""
}

// Resolve method are resolves SQLite connection URI from connection and credential parameters.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
//	Returns: uri string, err error resolved URI and error, if this occured.
func (c *SqliteConnectionResolver) Resolve(ctx context.Context) (uri string, err error) {
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
