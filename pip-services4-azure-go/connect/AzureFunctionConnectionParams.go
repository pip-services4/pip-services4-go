package connect

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	cconn "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// Contains connection parameters to authenticate against Azure
// and connect to specific Azure Functions.
//
// The class is able to compose and parse Azure Function connection parameters.
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
// In addition to standard parameters CredentialParams may contain any number of custom parameters.
//
// see AzureConnectionResolver
//
//	Example:
//		connection := connect.NewAzureFunctionConnectionParamsFromTuples(
//			"connection.uri", "http://myapp.azurewebsites.net/api/myfunction",
//			"connection.protocol", "http",
//			"connection.app_name", "myapp",
//			"connection.function_name", "myfunction",
//			"connection.auth_code", "code",
//		)
//		uri, _ := connection.FunctionUri()          	// Result: "http://myapp.azurewebsites.net/api/myfunction"
//		protocol, _ := connection.Protocol()        	// Result: "http"
//		appName, _ := connection.AppName()    			// Result: "myapp"
//		functionName, _ := connection.FunctionName()	// Result: "myfunction"
//		authCode, _ := connection.AuthCode()  			// Result: "code"
type AzureFunctionConnectionParams struct {
	cconf.ConfigParams
}

// Creates an new instance of the connection parameters.
func NewEmptyAzureFunctionConnectionParams() *AzureFunctionConnectionParams {
	return &AzureFunctionConnectionParams{ConfigParams: *cconf.NewEmptyConfigParams()}
}

// Creates an new instance of the connection parameters.
// Parameters:
//   - values	(optional) an object to be converted into key-value pairs to initialize this connection.
func NewAzureFunctionConnectionParams(values map[string]string) *AzureFunctionConnectionParams {
	return &AzureFunctionConnectionParams{ConfigParams: *cconf.NewConfigParamsFromMaps(values)}
}

// Creates a new AzureFunctionConnectionParams object filled with key-value pairs serialized as a string.
// Parameters:
//   - line	a string with serialized key-value pairs as "key1=value1;key2=value2;..."
//     Example: "Key1=123;Key2=ABC;Key3=2016-09-16T00:00:00.00Z"
//
// Returns a new AzureFunctionConnectionParams object.
func NewAzureFunctionConnectionParamsFromString(line string) *AzureFunctionConnectionParams {
	maps := cdata.NewStringValueMapFromString(line)
	return NewAzureFunctionConnectionParamsFromMaps(maps.Value())
}

// Retrieves AzureFunctionConnectionParams from configuration parameters.
// The values are retrieves from "connection" and "credential" sections.
// Parameters:
//   - config	configuration parameters
//
// Returns the generated AzureFunctionConnectionParams object.
func NewAzureFunctionConnectionParamsFromConfig(config *cconf.ConfigParams) *AzureFunctionConnectionParams {
	result := NewEmptyAzureFunctionConnectionParams()

	credentials := cauth.NewManyCredentialParamsFromConfig(config)
	for _, credential := range credentials {
		result.Append(credential.Value())
	}

	connections := cconn.NewManyConnectionParamsFromConfig(config)
	for _, credential := range connections {
		result.Append(credential.Value())
	}

	return result
}

// Creates a new ConfigParams object filled with provided key-value pairs called tuples.
// Tuples parameters contain a sequence of key1, value1, key2, value2, ... pairs.
// Parameters:
//   - tuples	the tuples to fill a new ConfigParams object.
//
// Returns a new ConfigParams object.
func NewAzureFunctionConnectionParamsFromTuples(tuples ...any) *AzureFunctionConnectionParams {
	config := cconf.NewConfigParamsFromTuples(tuples...)
	return NewAzureFunctionConnectionParamsFromConfig(config)
}

// Retrieves AzureFunctionConnectionParams from multiple configuration parameters.
// The values are retrieves from "connection" and "credential" sections.
// Parameters:
//   - configs	a list with configuration parameters
//
// Returns the generated AzureFunctionConnectionParams object.
func NewAzureFunctionConnectionParamsFromMaps(configs ...map[string]string) *AzureFunctionConnectionParams {
	config := cconf.NewConfigParamsFromMaps(configs...)
	return NewAzureFunctionConnectionParams(config.Value())
}

// Gets the Azure Platform service connection protocol.
// Returns the Azure service connection protocol.
func (c *AzureFunctionConnectionParams) Protocol() (string, bool) {
	return c.GetAsNullableString("protocol")
}

// Sets the Azure Platform service connection protocol.
// Parameters:
//   - value	a new Azure service connection protocol.
func (c *AzureFunctionConnectionParams) SetProtocol(value string) {
	c.SetAsObject("protocol", value)
}

// Gets the Azure Platform service uri.
// Returns the Azure sevice uri.
func (c *AzureFunctionConnectionParams) FunctionUri() (string, bool) {
	return c.GetAsNullableString("uri")
}

// Sets the Azure Platform service uri.
// Parameters:
//   - value	a new Azure service uri.
func (c *AzureFunctionConnectionParams) SetFunctionUri(value string) {
	c.SetAsObject("uri", value)
}

// Gets the Azure app name.
// Returns the Azure app name.
func (c *AzureFunctionConnectionParams) AppName() (string, bool) {
	return c.GetAsNullableString("app_name")
}

// Sets the Azure app name.
// Parameters:
//   - value	the Azure app name.
func (c *AzureFunctionConnectionParams) SetAppName(value string) {
	c.SetAsObject("app_name", value)
}

// Gets the Azure function name.
// Returns the Azure function name.
func (c *AzureFunctionConnectionParams) FunctionName() (string, bool) {
	return c.GetAsNullableString("app_name")
}

// Sets the Azure function name.
// Parameters:
//   - value	the Azure function name.
func (c *AzureFunctionConnectionParams) SetFunctionName(value string) {
	c.SetAsObject("app_name", value)
}

// Gets the Azure auth code.
// Returns the Azure auth code.
func (c *AzureFunctionConnectionParams) AuthCode() (string, bool) {
	return c.GetAsNullableString("app_name")
}

// Sets the Azure auth code.
// Parameters:
//   - value	the Azure auth code.
func (c *AzureFunctionConnectionParams) SetAuthCode(value string) {
	c.SetAsObject("app_name", value)
}

// Validates this connection parameters
// Parameters:
//   - ctx context.Context	execution context to trace execution through call chain.
func (c *AzureFunctionConnectionParams) Validate(ctx context.Context) error {
	_, uriOk := c.FunctionUri()
	protocol, protocolOk := c.Protocol()
	_, appNameOk := c.AppName()
	_, functionNameOk := c.FunctionName()

	if !uriOk && (!appNameOk || !functionNameOk || !protocolOk) {
		return cerr.NewConfigError(
			utils.ContextHelper.GetTraceId(ctx),
			"NO_CONNECTION_URI",
			"No uri, app_name and function_name is not configured in Auzre function uri",
		)
	}

	if protocolOk && protocol != "http" && protocol != "https" {
		return cerr.NewConfigError(
			utils.ContextHelper.GetTraceId(ctx),
			"WRONG_PROTOCOL",
			"Protocol is not supported by REST connection",
		).WithDetails("protocol", protocol)
	}

	return nil
}
