package connect

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	cconn "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// Contains connection parameters to authenticate against Google
// and connect to specific Google Cloud Platform.
//
// The class is able to compose and parse Google Platform connection parameters.
//
//	Configuration parameters
//		- connections:
//		     - uri:           full connection uri with specific app and function name
//		     - protocol:      connection protocol
//		     - project_id:    is your Google Cloud Platform project ID
//		     - region:        is the region where your function is deployed
//		     - function:      is the name of the HTTP function you deployed
//		     - org_id:        organization name
//
//		- credentials:
//		    - account: the service account name
//		    - auth_token:    Google-generated ID token or null if using custom auth (IAM)
//
// In addition to standard parameters CredentialParams may contain any number of custom parameters.
//
// see GcpConnectionResolver
//
//	Example:
//		connection := connect.NewGcpConnectionParamsFromTuples(
//			"connection.uri", "http://east-my_test_project.cloudfunctions.net/myfunction",
//			"connection.protocol", "http",
//			"connection.region", "east",
//			"connection.function", "myfunction",
//			"connection.project_id", "my_test_project",
//			"credential.auth_token", "1234",
//		)
//		uri, _ := connection.Uri()               // Result: "http://east-my_test_project.cloudfunctions.net/myfunction"
//		region, _ := connection.Region()         // Result: "east"
//		protocol, _ := connection.Protocol()     // Result: "http"
//		functionName, _ := connection.Function() // Result: "myfunction"
//		projectId, _ := connection.ProjectId()   // Result: "my_test_project"
//		authToken, _ := connection.AuthToken()   // Result: "1234"
type GcpConnectionParams struct {
	*cconf.ConfigParams
}

// Creates an new instance of the connection parameters.
func NewEmptyGcpConnectionParams() *GcpConnectionParams {
	return &GcpConnectionParams{ConfigParams: cconf.NewEmptyConfigParams()}
}

// Creates an new instance of the connection parameters.
// Parameters:
//   - values	(optional) an object to be converted into key-value pairs to initialize this connection.
func NewGcpConnectionParams(values map[string]string) *GcpConnectionParams {
	return &GcpConnectionParams{ConfigParams: cconf.NewConfigParamsFromMaps(values)}
}

// Creates a new GcpConnectionParams object filled with key-value pairs serialized as a string.
// Parameters:
//   - line	a string with serialized key-value pairs as "key1=value1;key2=value2;..."
//     Example: "Key1=123;Key2=ABC;Key3=2016-09-16T00:00:00.00Z"
//
// Returns a new GcpConnectionParams object.
func NewGcpConnectionParamsFromString(line string) *GcpConnectionParams {
	maps := cdata.NewStringValueMapFromString(line)
	return NewGcpConnectionParamsFromMaps(maps.Value())
}

// Retrieves GcpConnectionParams from configuration parameters.
// The values are retrieves from "connection" and "credential" sections.
// Parameters:
//   - config	configuration parameters
//
// Returns the generated GcpConnectionParams object.
func NewGcpConnectionParamsFromConfig(config *cconf.ConfigParams) *GcpConnectionParams {
	result := NewEmptyGcpConnectionParams()

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
func NewGcpConnectionParamsFromTuples(tuples ...any) *GcpConnectionParams {
	config := cconf.NewConfigParamsFromTuples(tuples...)
	return NewGcpConnectionParamsFromConfig(config)
}

// Retrieves GcpConnectionParams from multiple configuration parameters.
// The values are retrieves from "connection" and "credential" sections.
// Parameters:
//   - configs	a list with configuration parameters
//
// Returns the generated GcpConnectionParams object.
func NewGcpConnectionParamsFromMaps(configs ...map[string]string) *GcpConnectionParams {
	config := cconf.NewConfigParamsFromMaps(configs...)
	return NewGcpConnectionParams(config.Value())
}

// Gets the Google Platform service connection protocol.
// Returns the Google service connection protocol.
func (c *GcpConnectionParams) Protocol() (string, bool) {
	return c.GetAsNullableString("protocol")
}

// Sets the Google Platform service connection protocol.
// Parameters:
//   - value	a new Google service connection protocol.
func (c *GcpConnectionParams) SetProtocol(value string) {
	c.SetAsObject("protocol", value)
}

// Gets the Google Platform service uri.
// Returns the Google sevice uri.
func (c *GcpConnectionParams) Uri() (string, bool) {
	return c.GetAsNullableString("uri")
}

// Sets the Google Platform service uri.
// Parameters:
//   - value	a new Google service uri.
func (c *GcpConnectionParams) SetUri(value string) {
	c.SetAsObject("uri", value)
}

// Gets the Google function name.
// Returns the Google function name.
func (c *GcpConnectionParams) Function() (string, bool) {
	return c.GetAsNullableString("function")
}

// Sets the Google function name.
// Parameters:
//   - value	a new Google function name.
func (c *GcpConnectionParams) SetFunction(value string) {
	c.SetAsObject("function", value)
}

// Gets the region where your function is deployed.
// Returns the region of deployed function.
func (c *GcpConnectionParams) Region() (string, bool) {
	return c.GetAsNullableString("region")
}

// Sets the region where your function is deployed.
// Parameters:
//   - value	the region of deployed function.
func (c *GcpConnectionParams) SetRegion(value string) {
	c.SetAsObject("region", value)
}

// Gets the Google Cloud Platform project ID.
// Returns the project ID.
func (c *GcpConnectionParams) ProjectId() (string, bool) {
	return c.GetAsNullableString("project_id")
}

// Sets the Google Cloud Platform project ID.
// Parameters:
//   - value	a new project ID.
func (c *GcpConnectionParams) SetProjectId(value string) {
	c.SetAsObject("project_id", value)
}

// Gets an ID token with the request to authenticate themselves
// Returns the ID token.
func (c *GcpConnectionParams) AuthToken() (string, bool) {
	return c.GetAsNullableString("auth_token")
}

// Sets an ID token with the request to authenticate themselves
// Parameters:
//   - value	a new ID token.
func (c *GcpConnectionParams) SetAuthToken(value string) {
	c.SetAsObject("auth_token", value)
}

// Gets the service account name
// Returns the account name.
func (c *GcpConnectionParams) Account() (string, bool) {
	return c.GetAsNullableString("account")
}

// Sets the service account name
// Parameters:
//   - value	a new account name.
func (c *GcpConnectionParams) SetAccount(value string) {
	c.SetAsObject("account", value)
}

// Get organization name
// Returns the organization name.
func (c *GcpConnectionParams) OrgId() (string, bool) {
	return c.GetAsNullableString("org_id")
}

// Sets organization name
// Parameters:
//   - value	a new organization name.
func (c *GcpConnectionParams) SetOrgId(value string) {
	c.SetAsObject("org_id", value)
}

// Validates this connection parameters
// Parameters:
//   - ctx context.Context a context to trace execution through call chain.
func (c *GcpConnectionParams) Validate(ctx context.Context) error {
	_, uriOk := c.Uri()
	protocol, protocolOk := c.Protocol()
	_, functionNameOk := c.Function()
	_, regionOk := c.Region()
	_, projectIdOk := c.ProjectId()

	if !uriOk && (!projectIdOk || !regionOk || !functionNameOk || !protocolOk) {
		return cerr.NewConfigError(
			cctx.GetTraceId(ctx),
			"NO_CONNECTION_URI",
			"No uri, project_id, region and function is configured in Google function uri",
		)
	}

	if projectIdOk && protocol != "http" && protocol != "https" {
		return cerr.NewConfigError(
			cctx.GetTraceId(ctx),
			"WRONG_PROTOCOL", "Protocol is not supported by REST connection",
		).WithDetails("protocol", protocol)
	}

	return nil
}
