package test_connect

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestConnectionResolverConfigure(t *testing.T) {
	restConfig := config.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", 3000,
	)
	connectionResolver := connect.NewConnectionResolver(context.Background(), restConfig, nil)
	connections := connectionResolver.GetAll()
	assert.Len(t, connections, 1)

	connection := connections[0]
	assert.Equal(t, "http", connection.Protocol())
	assert.Equal(t, "localhost", connection.Host())
	assert.Equal(t, 3000, connection.Port())
}

func TestConnectionResolverRegister(t *testing.T) {
	connection := connect.NewEmptyConnectionParams()
	connectionResolver := connect.NewEmptyConnectionResolver()

	err := connectionResolver.Register(context.Background(), connection)
	assert.Nil(t, err)

	connections := connectionResolver.GetAll()
	assert.Len(t, connections, 0)

	err = connectionResolver.Register(context.Background(), connection)
	assert.Nil(t, err)

	connections = connectionResolver.GetAll()
	assert.Len(t, connections, 0)

	connection.SetDiscoveryKey("Discovery key")
	err = connectionResolver.Register(context.Background(), connection)
	assert.Nil(t, err)

	connections = connectionResolver.GetAll()
	assert.Len(t, connections, 0)
}

func TestConnectionResolverResolve(t *testing.T) {
	restConfig := config.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", 3000,
	)
	connectionResolver := connect.NewConnectionResolver(context.Background(), restConfig, nil)

	connection, err := connectionResolver.Resolve(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, connection)

	restConfigDiscovery := config.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", 3000,
		"connection.discovery_key", "Discovery key",
	)
	references := refer.NewEmptyReferences()
	connectionResolver = connect.NewConnectionResolver(context.Background(), restConfigDiscovery, references)

	connection, err = connectionResolver.Resolve(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, connection)
}
