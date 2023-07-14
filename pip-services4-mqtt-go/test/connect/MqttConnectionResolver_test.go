package test_connect

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-mqtt-go/connect"
	"github.com/stretchr/testify/assert"
)

type mqttConnectionResolverTest struct {
	resolver *connect.MqttConnectionResolver
}

func NewMqttConnectionResolverTest() *mqttConnectionResolverTest {
	c := mqttConnectionResolverTest{}
	return &c
}

func (c *mqttConnectionResolverTest) TestSingleConnection(t *testing.T) {
	resolver := connect.NewMqttConnectionResolver()
	resolver.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connection.protocol", "tcp",
		"connection.host", "localhost",
		"connection.port", 1883,
	))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "tcp://localhost:1883", connection.GetAsString("uri"))
	assert.Equal(t, "", connection.GetAsString("username"))
	assert.Equal(t, "", connection.GetAsString("password"))
}

func (c *mqttConnectionResolverTest) TestSingleConnectionWithAuth(t *testing.T) {
	resolver := connect.NewMqttConnectionResolver()
	resolver.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connection.protocol", "tcp",
		"connection.host", "localhost",
		"connection.port", 1883,
		"credential.username", "test",
		"credential.password", "pass123",
	))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "tcp://localhost:1883", connection.GetAsString("uri"))
	assert.Equal(t, "test", connection.GetAsString("username"))
	assert.Equal(t, "pass123", connection.GetAsString("password"))
}

func (c *mqttConnectionResolverTest) TestClusterConnection(t *testing.T) {
	resolver := connect.NewMqttConnectionResolver()
	resolver.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connections.0.protocol", "tcp",
		"connections.0.host", "server1",
		"connections.0.port", 1883,
		"connections.1.protocol", "tcp",
		"connections.1.host", "server2",
		"connections.1.port", 1883,
		"connections.2.protocol", "tcp",
		"connections.2.host", "server3",
		"connections.2.port", 1883,
	))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	//assert.Equal(t, "tcp://server1:1883,tcp://server2:1883,tcp://server3:1883", connection.GetAsString("uri"))
	assert.NotEqual(t, "", connection.GetAsString("uri"))
	assert.Equal(t, "", connection.GetAsString("username"))
	assert.Equal(t, "", connection.GetAsString("password"))
}

func (c *mqttConnectionResolverTest) TestClusterConnectionWithAuth(t *testing.T) {
	resolver := connect.NewMqttConnectionResolver()
	resolver.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connections.0.protocol", "tcp",
		"connections.0.host", "server1",
		"connections.0.port", 1883,
		"connections.1.protocol", "tcp",
		"connections.1.host", "server2",
		"connections.1.port", 1883,
		"connections.2.protocol", "tcp",
		"connections.2.host", "server3",
		"connections.2.port", 1883,
		"credential.username", "test",
		"credential.password", "pass123",
	))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	// assert.Equal(t, "tcp://server1:1883,tcp://server2:1883,tcp://server3:1883", connection.GetAsString("uri"))
	assert.NotEqual(t, "", connection.GetAsString("uri"))
	assert.Equal(t, "test", connection.GetAsString("username"))
	assert.Equal(t, "pass123", connection.GetAsString("password"))
}

func TestMqttConnectionResolver(t *testing.T) {
	c := NewMqttConnectionResolverTest()

	t.Run("Single Connection", c.TestSingleConnection)
	t.Run("Single Connection with Auth", c.TestSingleConnectionWithAuth)
	t.Run("Cluster Connection", c.TestClusterConnection)
	t.Run("Cluster Connection with Auth", c.TestClusterConnectionWithAuth)
}
