package connect

import (
	"context"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
)

type ICompositeConnectionResolverOverrides interface {
	ValidateConnection(ctx context.Context, connection *ConnectionParams) error

	ValidateCredential(ctx context.Context, credential *auth.CredentialParams) error

	ComposeOptions(connections []*ConnectionParams, credential *auth.CredentialParams, parameters *config.ConfigParams) *config.ConfigParams

	MergeConnection(options *config.ConfigParams, connection *ConnectionParams) *config.ConfigParams

	MergeCredential(options *config.ConfigParams, credential *auth.CredentialParams) *config.ConfigParams

	MergeOptional(options *config.ConfigParams, parameters *config.ConfigParams) *config.ConfigParams

	FinalizeOptions(options *config.ConfigParams) *config.ConfigParams
}

// CompositeConnectionResolver helper class that resolves connection and credential parameters,
// validates them and generates connection options.
//
//	Configuration parameters
//		- connection(s):
//			- discovery_key:               (optional) a key to retrieve the connection from [IDiscovery]]
//			- protocol:                    communication protocol
//			- host:                        host name or IP address
//			- port:                        port number
//			- uri:                         resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                   (optional) a key to retrieve the credentials from [ICredentialStore]]
//			- username:                    user name
//			- password:                    user password
//
//	References
//		- *:discovery:*:*:1.0          (optional) [IDiscovery]] services to resolve connections
//		- *:credential-store:*:*:1.0   (optional) Credential stores to resolve credentials
type CompositeConnectionResolver struct {
	Overrides ICompositeConnectionResolverOverrides

	// The connection options
	Options *config.ConfigParams

	// The connections resolver.
	ConnectionResolver *ConnectionResolver

	// The credentials resolver.
	CredentialResolver *auth.CredentialResolver

	// The cluster support (multiple connections)
	ClusterSupported bool

	// The default protocol
	DefaultProtocol string

	// The default port number
	DefaultPort int

	// The list of supported protocols
	SupportedProtocols []string
}

// InheritCompositeConnectionResolver creates new CompositeConnectionResolver
//
//	Parameters:
//		- overrides a child reference with overrides for virtual methods
//	return *CompositeConnectionResolver
func InheritCompositeConnectionResolver(overrides ICompositeConnectionResolverOverrides) *CompositeConnectionResolver {
	return &CompositeConnectionResolver{
		Overrides:          overrides,
		ConnectionResolver: NewEmptyConnectionResolver(),
		CredentialResolver: auth.NewEmptyCredentialResolver(),
		ClusterSupported:   true,
		DefaultPort:        0,
	}
}

// Configure configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config configuration parameters to be set.
func (c *CompositeConnectionResolver) Configure(ctx context.Context, config *config.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.CredentialResolver.Configure(ctx, config)
	c.Options = config.GetSection("options")
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references references to locate the component dependencies.
func (c *CompositeConnectionResolver) SetReferences(ctx context.Context, references refer.IReferences) {
	c.ConnectionResolver.SetReferences(ctx, references)
	c.CredentialResolver.SetReferences(ctx, references)
}

// Resolve connection options from connection and credential parameters.
//   - ctx context.Context execution context to trace execution through call chain.
//   - return 			 resolved options or error.
func (c *CompositeConnectionResolver) Resolve(ctx context.Context) (options *config.ConfigParams, err error) {
	var connections []*ConnectionParams
	var credential *auth.CredentialParams

	connections, err = c.ConnectionResolver.ResolveAll(ctx)

	// Validate if cluster (multiple connections) is supported
	if err == nil && len(connections) > 0 && !c.ClusterSupported {
		err = cerr.NewConfigError(
			cctx.GetTraceId(ctx),
			"MULTIPLE_CONNECTIONS_NOT_SUPPORTED",
			"Multiple (cluster) connections are not supported",
		)
	}

	// Validate connections
	if err == nil {
		for _, connection := range connections {
			err = c.ValidateConnection(ctx, connection)
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		return nil, err
	}

	credential, err = c.CredentialResolver.Lookup(context.Background())
	if credential == nil {
		credential = auth.NewEmptyCredentialParams()
	}
	// Validate credential
	if err == nil {
		err = c.ValidateCredential(ctx, credential)
	}

	if err != nil {
		return nil, err
	}

	return c.ComposeOptions(connections, credential, c.Options), nil
}

// Compose Composite connection options from connection and credential parameters.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- connections connection parameters
//		- credential parameters
//		- parameters optional parameters
//	Returns: resolved options or error.
func (c *CompositeConnectionResolver) Compose(ctx context.Context, connections []*ConnectionParams, credential *auth.CredentialParams,
	parameters *config.ConfigParams) (options *config.ConfigParams, err error) {

	// Validate connection parameters
	for _, connection := range connections {
		err = c.Overrides.ValidateConnection(ctx, connection)
		if err != nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	// Validate credential parameters
	err = c.Overrides.ValidateCredential(ctx, credential)

	if err != nil {
		return nil, err
	}

	return c.Overrides.ComposeOptions(connections, credential, parameters), nil
}

// ValidateConnection validates connection parameters.
// This method can be override in child classes.
//
//	Parameters:
//		- ctx context.Context a context to trace execution through call chain.
//		- connection    parameters to be validated
//	Returns: error or nil if validation was successful
func (c *CompositeConnectionResolver) ValidateConnection(ctx context.Context, connection *ConnectionParams) error {
	traceId := cctx.GetTraceId(ctx)

	if connection == nil {
		return cerr.NewConfigError(traceId, "NO_CONNECTION", "Connection parameters are not set is not set")
	}

	// URI usually contains all information
	uri := connection.Uri()
	if uri != "" {
		return nil
	}

	protocol := connection.ProtocolWithDefault(c.DefaultProtocol)
	if protocol == "" {
		return cerr.NewConfigError(traceId, "NO_PROTOCOL", "Connection protocol is not set")
	}
	if c.SupportedProtocols != nil && indexOf(c.SupportedProtocols, protocol) < 0 {
		return cerr.NewConfigError(traceId, "UNSUPPORTED_PROTOCOL", "The protocol "+protocol+" is not supported")
	}

	var host = connection.Host()
	if host == "" {
		return cerr.NewConfigError(traceId, "NO_HOST", "Connection host is not set")
	}

	var port = connection.PortWithDefault(c.DefaultPort)
	if port == 0 {
		return cerr.NewConfigError(traceId, "NO_PORT", "Connection port is not set")
	}

	return nil
}

// ValidateCredential parameters.
// This method can be override in child classes.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- credential    parameters to be validated
//	Returns: error or nil if validation was successful
func (c *CompositeConnectionResolver) ValidateCredential(ctx context.Context, credential *auth.CredentialParams) error {
	return nil
}

// ComposeOptions composes connection and credential parameters into connection options.
// This method can be override in child classes.
//
//	Parameters:
//		- connections a list of connection parameters
//		- credential parameters
//		- parameters optional parameters
//	Returns: a composed connection options.
func (c *CompositeConnectionResolver) ComposeOptions(connections []*ConnectionParams, credential *auth.CredentialParams, parameters *config.ConfigParams) *config.ConfigParams {
	// Connection options
	options := config.NewEmptyConfigParams()

	// Merge connection parameters
	for _, connection := range connections {
		options = c.Overrides.MergeConnection(options, connection)
	}

	// Merge credential parameters
	options = c.Overrides.MergeCredential(options, credential)

	// Merge optional parameters
	options = c.Overrides.MergeOptional(options, parameters)

	// Perform final processing
	options = c.Overrides.FinalizeOptions(options)

	return options
}

// MergeConnection merges connection options with connection parameters
// This method can be override in child classes.
//
//	Parameters:
//		-  options connection options
//		-  connection parameters to be merged
//	Returns: merged connection options.
func (c *CompositeConnectionResolver) MergeConnection(options *config.ConfigParams, connection *ConnectionParams) *config.ConfigParams {
	var mergedOptions = options.SetDefaults(connection.ConfigParams)
	return mergedOptions
}

// MergeCredential merges connection options with credential parameters
// This method can be overriden in child classes.
//
//	Parameters:
//		- options connection options
//		- credential parameters to be merged
//	Returns: merged connection options.
func (c *CompositeConnectionResolver) MergeCredential(options *config.ConfigParams, credential *auth.CredentialParams) *config.ConfigParams {
	var mergedOptions = options.Override(credential.ConfigParams)
	return mergedOptions
}

// MergeOptional merges connection options with optional parameters
// This method can be overriden in child classes.
//
//	Parameters:
//		- options connection options
//		- parameters optional parameters to be merged
//	Returns merged connection options.
func (c *CompositeConnectionResolver) MergeOptional(options *config.ConfigParams, parameters *config.ConfigParams) *config.ConfigParams {
	var mergedOptions = options.SetDefaults(parameters)
	return mergedOptions
}

// FinalizeOptions finalize merged options
// This method can be overriden in child classes.
//
//	Parameters:
//		- options connection options
//	Returns: finalized connection options
func (c *CompositeConnectionResolver) FinalizeOptions(options *config.ConfigParams) *config.ConfigParams {
	return options
}

func indexOf(a []string, e string) int {
	for i := range a {
		if e == a[i] {
			return i
		}
	}
	return -1
}
