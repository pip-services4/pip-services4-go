package test_connect

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/connect"
	"github.com/stretchr/testify/assert"
)

type natsConnectionResolverTest struct {
	resolver *connect.NatsConnectionResolver
}

func NewNatsConnectionResolverTest() *natsConnectionResolverTest {
	c := natsConnectionResolverTest{}
	return &c
}

func (c *natsConnectionResolverTest) TestSingleConnection(t *testing.T) {
	resolver := connect.NewNatsConnectionResolver()
	resolver.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connection.protocol", "nats",
		"connection.host", "localhost",
		"connection.port", 4222,
	))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "nats://localhost:4222", connection.GetAsString("uri"))
	assert.Equal(t, "", connection.GetAsString("username"))
	assert.Equal(t, "", connection.GetAsString("password"))
	assert.Equal(t, "", connection.GetAsString("token"))
}

func (c *natsConnectionResolverTest) TestClusterConnection(t *testing.T) {
	resolver := connect.NewNatsConnectionResolver()
	resolver.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connections.0.protocol", "nats",
		"connections.0.host", "server1",
		"connections.0.port", 4222,
		"connections.1.protocol", "nats",
		"connections.1.host", "server2",
		"connections.1.port", 4222,
		"connections.2.protocol", "nats",
		"connections.2.host", "server3",
		"connections.2.port", 4222,
	))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	//assert.Equal(t, "nats://server1:4222, nats://server2:4222, nats://server3:4222", connection.GetAsString("uri"))
	assert.NotEqual(t, "", connection.GetAsString("uri"))
	assert.Equal(t, "", connection.GetAsString("username"))
	assert.Equal(t, "", connection.GetAsString("password"))
	assert.Equal(t, "", connection.GetAsString("token"))
}

func (c *natsConnectionResolverTest) TestClusterConnectionWithAuth(t *testing.T) {
	resolver := connect.NewNatsConnectionResolver()
	resolver.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connections.0.protocol", "nats",
		"connections.0.host", "server1",
		"connections.0.port", 4222,
		"connections.1.protocol", "nats",
		"connections.1.host", "server2",
		"connections.1.port", 4222,
		"connections.2.protocol", "nats",
		"connections.2.host", "server3",
		"connections.2.port", 4222,
		"credential.token", "ABC",
		"credential.username", "test",
		"credential.password", "pass123",
	))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	// assert.Equal(t, "nats://server1:4222, nats://server2:4222, nats://server3:4222", connection.GetAsString("uri"))
	assert.NotEqual(t, "", connection.GetAsString("uri"))
	assert.Equal(t, "test", connection.GetAsString("username"))
	assert.Equal(t, "pass123", connection.GetAsString("password"))
	assert.Equal(t, "ABC", connection.GetAsString("token"))
}

func TestNatsConnectionResolver(t *testing.T) {
	c := NewNatsConnectionResolverTest()

	t.Run("Single Connection", c.TestSingleConnection)
	t.Run("Cluster Connection", c.TestClusterConnection)
	t.Run("Cluster Connection with Auth", c.TestClusterConnectionWithAuth)
}
