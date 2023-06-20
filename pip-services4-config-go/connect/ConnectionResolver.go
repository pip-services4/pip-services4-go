package connect

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// ConnectionResolver helper class to retrieve component connections.
// If connections are configured to be retrieved from IDiscovery, it automatically locates
// IDiscovery in component references and retrieve connections from there using discovery_key parameter.
//
//	Configuration parameters
//		- connection:
//			- discovery_key: (optional) a key to retrieve the connection from IDiscovery
//			- ... other connection parameters
//		- connections: alternative to connection
//			- [connection params 1]: first connection parameters
//				- ... connection parameters for key 1
//			- [connection params N]: Nth connection parameters
//				- ... connection parameters for key N
//	References:
//		- *:discovery:*:*:1.0 (optional) IDiscovery services to resolve connections
//
//	see ConnectionParams
//	see IDiscovery
//
//	Example:
//		config = NewConfigParamsFromTuples(
//			"connection.host", "10.1.1.100",
//			"connection.port", 8080
//		);
//		connectionResolver := NewConnectionResolver();
//		connectionResolver.Configure(context.Background(), config);
//		connectionResolver.SetReferences(context.Background(), references);
//		res, err := connectionResolver.Resolve("123");
type ConnectionResolver struct {
	connections []*ConnectionParams
	references  refer.IReferences
}

// NewEmptyConnectionResolver creates a new instance of connection resolver.
//
//	Returns: *ConnectionResolver
func NewEmptyConnectionResolver() *ConnectionResolver {
	return &ConnectionResolver{
		connections: []*ConnectionParams{},
		references:  nil,
	}
}

// NewConnectionResolver Creates a new instance of connection resolver.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams component configuration parameters
//		- references refer.IReferences component references
//	Returns: *ConnectionResolver
func NewConnectionResolver(ctx context.Context, config *config.ConfigParams, references refer.IReferences) *ConnectionResolver {
	c := &ConnectionResolver{
		connections: []*ConnectionParams{},
		references:  references,
	}

	if config != nil {
		c.Configure(ctx, config)
	}

	return c
}

// Configure Configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration parameters to be set.
func (c *ConnectionResolver) Configure(ctx context.Context, config *config.ConfigParams) {
	connections := NewManyConnectionParamsFromConfig(config)
	c.connections = append(c.connections, connections...)
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references refer.IReferences references to locate the component dependencies.
func (c *ConnectionResolver) SetReferences(ctx context.Context, references refer.IReferences) {
	c.references = references
}

// GetAll gets all connections configured in component configuration.
// Redirect to Discovery services is not done at this point.
// If you need fully fleshed connection use resolve method instead.
//
//	Returns: []*ConnectionParams a list with connection parameters
func (c *ConnectionResolver) GetAll() []*ConnectionParams {
	return c.connections
}

// Add a new connection to component connections
//
//	Parameters:
//		- connection *ConnectionParams new connection parameters to be added
func (c *ConnectionResolver) Add(connection *ConnectionParams) {
	c.connections = append(c.connections, connection)
}

func (c *ConnectionResolver) resolveInDiscovery(ctx context.Context,
	connection *ConnectionParams) (result *ConnectionParams, err error) {

	if !connection.UseDiscovery() {
		return connection, nil
	}

	key := connection.DiscoveryKey()
	if c.references == nil {
		return nil, nil
	}

	discoveryDescriptor := refer.NewDescriptor("*", "discovery", "*", "*", "*")
	components := c.references.GetOptional(discoveryDescriptor)
	if len(components) == 0 {
		err := refer.NewReferenceError(ctx, discoveryDescriptor)
		return nil, err
	}

	for _, component := range components {
		if discovery, ok := component.(IDiscovery); ok && discovery != nil {
			connection, err = discovery.ResolveOne(ctx, key)
			if connection != nil || err != nil {
				return connection, err
			}
		}
	}

	return nil, nil
}

// Resolve a single component connection. If connections are configured to be retrieved
// from Discovery service it finds a IDiscovery and resolves the connection there.
//
//	see IDiscovery
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: *ConnectionParams, error resolved connection or error.
func (c *ConnectionResolver) Resolve(ctx context.Context) (*ConnectionParams, error) {
	if len(c.connections) == 0 {
		return nil, nil
	}

	resolveConnections := make([]*ConnectionParams, 0)

	for _, connection := range c.connections {
		if !connection.UseDiscovery() {
			return connection, nil
		}

		resolveConnections = append(resolveConnections, connection)
	}

	for _, connection := range resolveConnections {
		c, err := c.resolveInDiscovery(ctx, connection)
		if c != nil || err != nil {
			return c, err
		}
	}

	return nil, nil
}

func (c *ConnectionResolver) resolveAllInDiscovery(ctx context.Context,
	connection *ConnectionParams) (result []*ConnectionParams, err error) {

	if !connection.UseDiscovery() {
		return []*ConnectionParams{connection}, nil
	}

	key := connection.DiscoveryKey()
	if c.references == nil {
		return nil, nil
	}

	discoveryDescriptor := refer.NewDescriptor("*", "discovery", "*", "*", "*")
	components := c.references.GetOptional(discoveryDescriptor)
	if len(components) == 0 {
		err := refer.NewReferenceError(ctx, discoveryDescriptor)
		return nil, err
	}

	resolvedConnections := make([]*ConnectionParams, 0)

	for _, component := range components {
		if discovery, ok := component.(IDiscovery); ok && discovery != nil {
			connections, err := discovery.ResolveAll(ctx, key)
			if err != nil {
				return nil, err
			}
			if connections != nil {
				for _, c := range connections {
					resolvedConnections = append(resolvedConnections, c)
				}
			}
		}
	}

	return resolvedConnections, nil
}

// ResolveAll resolves all component connection. If connections are configured to be
// retrieved from Discovery service it finds a IDiscovery and resolves the connection there.
//
//	see IDiscovery
//	Parameters:
//		-  transaction id to trace execution through call chain.
//	Returns: []*ConnectionParams, error resolved connections or error.
func (c *ConnectionResolver) ResolveAll(ctx context.Context) ([]*ConnectionParams, error) {
	resolvedConnections := make([]*ConnectionParams, 0)
	resolveConnections := make([]*ConnectionParams, 0)

	for _, connection := range c.connections {
		if !connection.UseDiscovery() {
			resolvedConnections = append(resolvedConnections, connection)
		} else {
			resolveConnections = append(resolveConnections, connection)
		}
	}

	for _, connection := range resolveConnections {
		connections, err := c.resolveAllInDiscovery(ctx, connection)
		if err != nil {
			return nil, err
		}
		resolvedConnections = append(resolvedConnections, connections...)
	}

	return resolvedConnections, nil
}

func (c *ConnectionResolver) registerInDiscovery(ctx context.Context,
	connection *ConnectionParams) (result bool, err error) {

	if !connection.UseDiscovery() {
		return false, nil
	}

	key := connection.DiscoveryKey()
	if c.references == nil {
		return false, nil
	}

	discoveryDescriptor := refer.NewDescriptor("*", "discovery", "*", "*", "*")
	components := c.references.GetOptional(discoveryDescriptor)
	if len(components) == 0 {
		err := refer.NewReferenceError(ctx, discoveryDescriptor)
		return false, err
	}

	registered := false

	for _, component := range components {
		if discovery, ok := component.(IDiscovery); ok && discovery != nil {
			_, err = discovery.Register(ctx, key, connection)
			if err != nil {
				return false, err
			}
			registered = true
		}
	}

	return registered, nil
}

// Register the given connection in all referenced discovery services.
// This method can be used for dynamic service discovery.
//
//	see IDiscovery
//	Parameters:
//		-  transaction id to trace execution through call chain.
//		- connection *ConnectionParams a connection to register.
//	Returns: error
func (c *ConnectionResolver) Register(ctx context.Context, connection *ConnectionParams) error {
	registered, err := c.registerInDiscovery(ctx, connection)
	if registered {
		c.connections = append(c.connections, connection)
	}
	return err
}
