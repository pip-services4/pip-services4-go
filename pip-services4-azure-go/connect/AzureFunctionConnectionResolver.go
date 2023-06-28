package connect

import (
	"context"
	"net/url"
	"strings"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	cconn "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// Helper class to retrieve Google connection and credential parameters,
// validate them and compose a AzureFunctionConnectionParams value.
//
//	Configuration parameters
//		- connections:
//		     - uri:           full connection uri with specific app and function name
//		     - protocol:      connection protocol
//		     - app_name:      alternative app name
//		     - function_name: application function name
//		- credentials:
//		     - auth_code:     authorization code or null if using custom auth
//
//	References
//		- *:discovery:*:*:1.0			(optional) IDiscovery services to resolve connections
//		- *:credential-store:*:*:1.0	(optional) Credential stores to resolve credentials
//
// see ConnectionParams (in the Pip.Services components package)
//
//	Example:
//		config := config.NewConfigParamsFromTuples(
//			"connection.uri", "http://myapp.azurewebsites.net/api/myfunction",
//			"connection.app_name", "myapp",
//			"connection.function_name", "myfunction",
//			"credential.auth_code", "XXXXXXXXXX",
//		)
//		ctx := context.Background()
//		connectionResolver := connect.NewAzureConnectionResolver()
//		connectionResolver.Configure(ctx, config)
//		connectionResolver.SetReferences(ctx, references)
//
//		connectionParams, _ := connectionResolver.Resolve("123")
type AzureFunctionConnectionResolver struct {
	// The connection resolver.
	connectionResolver *cconn.ConnectionResolver
	// The credential resolver.
	credentialResolver *cauth.CredentialResolver
}

// Creates new instance of AzureFunctionConnectionResolver
func NewAzureFunctionConnectionResolver() *AzureFunctionConnectionResolver {
	return &AzureFunctionConnectionResolver{
		connectionResolver: cconn.NewEmptyConnectionResolver(),
		credentialResolver: cauth.NewEmptyCredentialResolver(),
	}
}

// Configure Configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration parameters to be set.
func (c *AzureFunctionConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.connectionResolver.Configure(ctx, config)
	c.credentialResolver.Configure(ctx, config)
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references refer.IReferences references to locate the component dependencies.
func (c *AzureFunctionConnectionResolver) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.connectionResolver.SetReferences(ctx, references)
	c.credentialResolver.SetReferences(ctx, references)
}

// Resolves connection and credential parameters and generates a single
// AzureConnectionParams value.
// Parameters:
//   - ctx context.Context	execution context to trace execution through call chain.
//
// Returns AzureConnectionParams value and error
func (c *AzureFunctionConnectionResolver) Resolve(ctx context.Context) (*AzureFunctionConnectionParams, error) {
	connection := NewEmptyAzureFunctionConnectionParams()

	connectionParams, err := c.connectionResolver.Resolve(ctx)
	if err != nil {
		return nil, err
	}
	connection.Append(connectionParams.Value())

	credentialParams, err := c.credentialResolver.Lookup(ctx)
	if err != nil {
		return nil, err
	}

	if credentialParams != nil {
		connection.Append(credentialParams.Value())
	}

	// Perform validation
	err = connection.Validate(ctx)
	if err != nil {
		return nil, err
	}

	connection = c.composeConnection(connection)

	return connection, nil
}

func (c *AzureFunctionConnectionResolver) composeConnection(connection *AzureFunctionConnectionParams) *AzureFunctionConnectionParams {
	connection = NewAzureFunctionConnectionParamsFromMaps(connection.Value())

	uri, ok := connection.FunctionUri()

	if !ok {
		protocol, _ := connection.Protocol()
		appName, _ := connection.AppName()
		functionName, _ := connection.FunctionName()

		// http://myapp.azurewebsites.net/api/myfunction
		uri = protocol + "://" + appName + ".azurewebsites.net/api/" + functionName

		connection.SetFunctionUri(uri)
	} else {
		parsed, err := url.Parse(uri)
		if err != nil {
			panic(err)
		}

		protocol := parsed.Scheme
		appName := strings.Replace(parsed.Hostname(), ".azurewebsites.net", "", 1)

		if strings.Contains(appName, ":") { // remove port number
			res := strings.Split(appName, ":")
			appName = res[0]
		}

		functionName := strings.Replace(parsed.Path, "/api/", "", 1)

		connection.SetProtocol(protocol)
		connection.SetAppName(appName)
		connection.SetFunctionName(functionName)
	}

	return connection
}
