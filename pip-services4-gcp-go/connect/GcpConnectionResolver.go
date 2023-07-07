package connect

import (
	"context"
	"net/url"
	"strings"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	refer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	cconn "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
)

// Helper class to retrieve Google connection and credential parameters,
// validate them and compose a GcpConnectionParams value.
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
//	References
//		- *:credential-store:*:*:1.0	(optional) Credential stores to resolve credentials
//
// see ConnectionParams (in the Pip.Services components package)
//
//	Example:
//		config := config.NewConfigParamsFromTuples(
//			"connection.uri", "http://east-my_test_project.cloudfunctions.net/myfunction",
//			"connection.protocol", "http",
//			"connection.region", "east",
//			"connection.function", "myfunction",
//			"connection.project_id", "my_test_project",
//			"credential.auth_token", "1234",
//		)
//		ctx := context.Background()
//		connectionResolver := connect.NewGcpConnectionResolver()
//		connectionResolver.Configure(ctx, config)
//		connectionResolver.SetReferences(ctx, references)
//		connectionParams, _ := connectionResolver.Resolve("123")
type GcpConnectionResolver struct {
	// The connection resolver.
	connectionResolver *cconn.ConnectionResolver
	// The credential resolver.
	credentialResolver *cauth.CredentialResolver
}

// Creates new instance of GcpConnectionResolver
func NewGcpConnectionResolver() *GcpConnectionResolver {
	return &GcpConnectionResolver{
		connectionResolver: cconn.NewEmptyConnectionResolver(),
		credentialResolver: cauth.NewEmptyCredentialResolver(),
	}
}

// Configure Configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration parameters to be set.
func (c *GcpConnectionResolver) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.connectionResolver.Configure(ctx, config)
	c.credentialResolver.Configure(ctx, config)
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references refer.IReferences references to locate the component dependencies.
func (c *GcpConnectionResolver) SetReferences(ctx context.Context, references refer.IReferences) {
	c.connectionResolver.SetReferences(ctx, references)
	c.credentialResolver.SetReferences(ctx, references)
}

// Resolves connection and credential parameters and generates a single
// GcpConnectionParams value.
// Parameters:
//   - ctx context.Context a context to trace execution through call chain.
//
// Returns GcpConnectionParams value or error.
//
// see IDiscovery (in the Pip.Services components package)
func (c *GcpConnectionResolver) Resolve(ctx context.Context) (*GcpConnectionParams, error) {
	connection := NewEmptyGcpConnectionParams()

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

func (c *GcpConnectionResolver) composeConnection(connection *GcpConnectionParams) *GcpConnectionParams {
	connection = NewGcpConnectionParamsFromMaps(connection.Value())

	uri, uriOk := connection.Uri()

	if !uriOk || uri == "" {
		protocol, _ := connection.Protocol()
		functionName, functionNameOk := connection.Function()
		projectId, _ := connection.ProjectId()
		region, _ := connection.Region()

		// https://YOUR_REGION-YOUR_PROJECT_ID.cloudfunctions.net/FUNCTION_NAME
		uri = protocol + "://" + region + "-" + projectId + ".cloudfunctions.net"
		if functionNameOk {
			uri += "/" + functionName
		}

		connection.SetUri(uri)
	} else {
		parsed, err := url.Parse(uri)
		if err != nil {
			panic(err)
		}

		var functionName string

		protocol := parsed.Scheme
		region := parseRegion(uri)
		projectId := parseProjectId(uri)

		if parsed.Path != "" {
			functionName = parsed.Path[1:]
		}

		connection.SetRegion(region)
		connection.SetProjectId(projectId)
		connection.SetFunction(functionName)
		connection.SetProtocol(protocol)
	}

	return connection
}

func parseProjectId(uri string) string {
	var startInd int
	var endInd int

	if !strings.Contains(uri, "-") {
		return ""
	}

	startInd = strings.Index(uri, "-") + 1
	endInd = strings.Index(uri, ".")

	return uri[startInd:endInd]
}

func parseRegion(uri string) string {
	if !strings.Contains(uri, "-") {
		return ""
	}

	endInd := strings.Index(uri, "-")
	startInd := strings.Index(uri, "//") + 2

	return uri[startInd:endInd]
}
