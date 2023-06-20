package connect

import "context"

// IDiscovery interface for discovery services which are used to store and resolve
// connection parameters to connect to external services.
type IDiscovery interface {

	// Register connection parameters into the discovery service.
	Register(ctx context.Context, key string,
		connection *ConnectionParams) (result *ConnectionParams, err error)

	// ResolveOne a single connection parameters by its key.
	ResolveOne(ctx context.Context, key string) (result *ConnectionParams, err error)

	// ResolveAll all connection parameters by their key.
	ResolveAll(ctx context.Context, key string) (result []*ConnectionParams, err error)
}
