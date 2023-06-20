package test_connect

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestMemoryDiscoveryReadConnections(t *testing.T) {
	config := config.NewConfigParamsFromTuples(
		"key1.host", "10.1.1.100",
		"key1.port", "8080",
		"key2.host", "10.1.1.101",
		"key2.port", "8082",
	)

	discovery := connect.NewEmptyMemoryDiscovery()
	discovery.Configure(context.Background(), config)

	// Resolve one
	connection, err := discovery.ResolveOne(context.Background(), "key1")

	assert.Equal(t, err, nil)
	assert.Equal(t, "10.1.1.100", connection.Host())
	assert.Equal(t, 8080, connection.Port())

	connection, err = discovery.ResolveOne(context.Background(), "key2")

	assert.Equal(t, err, nil)
	assert.Equal(t, "10.1.1.101", connection.Host())
	assert.Equal(t, 8082, connection.Port())

	// Resolve all
	_, err = discovery.Register(context.Background(), "key1", connect.NewConnectionParamsFromTuples(
		"host", "10.3.3.151",
	))
	assert.Equal(t, err, nil)

	connections, err := discovery.ResolveAll(context.Background(), "key1")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(connections) > 1, true)
}
