package connect

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// MemoryDiscovery discovery service that keeps connections in memory.
//
//	Configuration parameters
//		[connection key 1]:
//		... connection parameters for key 1
//		[connection key 2]:
//		... connection parameters for key N
//	see IDiscovery
//	see ConnectionParams
//	Example
//		config := NewConfigParamsFromTuples(
//			"key1.host", "10.1.1.100",
//			"key1.port", "8080",
//			"key2.host", "10.1.1.100",
//			"key2.port", "8082"
//		);
//		discovery := NewMemoryDiscovery();
//		discovery.ReadConnections(config);
//		conn, err := discovery.ResolveOne("123", "key1");
//
// Result: host=10.1.1.100;port=8080
type MemoryDiscovery struct {
	items map[string][]*ConnectionParams
}

// NewEmptyMemoryDiscovery creates a new instance of discovery service.
//
//	Returns: *MemoryDiscovery
func NewEmptyMemoryDiscovery() *MemoryDiscovery {
	return &MemoryDiscovery{
		items: map[string][]*ConnectionParams{},
	}
}

// NewMemoryDiscovery creates a new instance of discovery service.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration with connection parameters.
//	Returns: *MemoryDiscovery
func NewMemoryDiscovery(ctx context.Context, config *config.ConfigParams) *MemoryDiscovery {
	c := &MemoryDiscovery{
		items: map[string][]*ConnectionParams{},
	}

	if config != nil {
		c.Configure(ctx, config)
	}

	return c
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration parameters to be set.
func (c *MemoryDiscovery) Configure(ctx context.Context, config *config.ConfigParams) {
	c.ReadConnections(config)
}

// ReadConnections from configuration parameters. Each section represents an individual Connectionparams
//
//	Parameters:
//		- ctx context.Context
//		- config *configure.ConfigParams configuration parameters to be read
func (c *MemoryDiscovery) ReadConnections(config *config.ConfigParams) {
	c.items = make(map[string][]*ConnectionParams)

	if config.Len() > 0 {
		connectionSections := config.GetSectionNames()
		for _, key := range connectionSections {
			connection := config.GetSection(key)
			c.items[key] = []*ConnectionParams{NewConnectionParamsFromValue(connection)}
		}
	}
}

// Register connection parameters into the discovery service.
//
//	Parameters:
//		-  transaction id to trace execution through call chain.
//		- key string a key to uniquely identify the connection parameters.
//		- connection *ConnectionParams
//	Returns: *ConnectionParams, error registered connection or error.
func (c *MemoryDiscovery) Register(ctx context.Context, key string,
	connection *ConnectionParams) (result *ConnectionParams, err error) {

	if connection != nil {
		connections, ok := c.items[key]
		if ok && connections != nil {
			connections = append(connections, connection)
		} else {
			connections = []*ConnectionParams{connection}
		}
		c.items[key] = connections
	}

	return connection, nil
}

// ResolveOne a single connection parameters by its key.
//
//	Parameters:
//		- correlationId: string transaction id to trace execution through call chain.
//		- key: string a key to uniquely identify the connection.
//	Returns: *ConnectionParams, error receives found connection or error.
func (c *MemoryDiscovery) ResolveOne(ctx context.Context,
	key string) (result *ConnectionParams, err error) {

	connections, _ := c.ResolveAll(ctx, key)
	if len(connections) > 0 {
		return connections[0], nil
	}

	return nil, nil
}

// ResolveAll connection parameters by its key.
//
//	Parameters:
//		- correlationId: string transaction id to trace execution through call chain.
//		- key: string a key to uniquely identify the connection.
//	Returns: *ConnectionParams, error receives found connection or error.
func (c *MemoryDiscovery) ResolveAll(ctx context.Context,
	key string) (result []*ConnectionParams, err error) {
	connections := c.items[key]

	if connections == nil {
		connections = []*ConnectionParams{}
	}

	return connections, nil
}
