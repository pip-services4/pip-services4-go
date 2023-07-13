package test_connect

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-kafka-go/connect"
	"github.com/stretchr/testify/assert"
)

type kafkaConnectionResolverTest struct {
}

func NewKafkaConnectionResolverTest() *kafkaConnectionResolverTest {
	c := kafkaConnectionResolverTest{}
	return &c
}

func (c *kafkaConnectionResolverTest) TestSingleConnection(t *testing.T) {
	resolver := connect.NewKafkaConnectionResolver()
	resolver.Configure(context.Background(),
		cconf.NewConfigParamsFromTuples(
			"connection.protocol", "tcp",
			"connection.host", "localhost",
			"connection.port", 9092,
		),
	)

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "localhost:9092", connection.GetAsString("uri"))
	assert.Equal(t, "", connection.GetAsString("username"))
	assert.Equal(t, "", connection.GetAsString("password"))
	assert.Equal(t, "", connection.GetAsString("mechanism"))
}

func (c *kafkaConnectionResolverTest) TestClusterConnection(t *testing.T) {
	resolver := connect.NewKafkaConnectionResolver()
	resolver.Configure(context.Background(),
		cconf.NewConfigParamsFromTuples(
			"connections.0.protocol", "tcp",
			"connections.0.host", "server1",
			"connections.0.port", 9092,
			"connections.1.protocol", "tcp",
			"connections.1.host", "server2",
			"connections.1.port", 9092,
			"connections.2.protocol", "tcp",
			"connections.2.host", "server3",
			"connections.2.port", 9092,
		))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	//assert.Equal(t, "server1:9092,server2:9092,server3:9092", connection.GetAsString("uri"))
	assert.NotEqual(t, "", connection.GetAsString("uri"))
	assert.Equal(t, "", connection.GetAsString("username"))
	assert.Equal(t, "", connection.GetAsString("password"))
	assert.Equal(t, "", connection.GetAsString("mechanism"))
}

func (c *kafkaConnectionResolverTest) TestClusterConnectionWithAuth(t *testing.T) {
	resolver := connect.NewKafkaConnectionResolver()
	resolver.Configure(context.Background(),
		cconf.NewConfigParamsFromTuples(
			"connections.0.protocol", "tcp",
			"connections.0.host", "server1",
			"connections.0.port", 9092,
			"connections.1.protocol", "tcp",
			"connections.1.host", "server2",
			"connections.1.port", 9092,
			"connections.2.protocol", "tcp",
			"connections.2.host", "server3",
			"connections.2.port", 9092,
			"credential.mechanism", "plain",
			"credential.username", "test",
			"credential.password", "pass123",
		))

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	// assert.Equal(t, "server1:9092,server2:9092,server3:9092", connection.GetAsString("uri"))
	assert.NotEqual(t, "", connection.GetAsString("uri"))
	assert.Equal(t, "test", connection.GetAsString("username"))
	assert.Equal(t, "pass123", connection.GetAsString("password"))
	assert.Equal(t, "plain", connection.GetAsString("mechanism"))
}

func (c *kafkaConnectionResolverTest) TestConnectionUri(t *testing.T) {
	resolver := connect.NewKafkaConnectionResolver()
	resolver.Configure(context.Background(),
		cconf.NewConfigParamsFromTuples(
			"connection.uri", "tcp://server1:9092,server2:9092,server3:9092?param=xyz",
			"credential.mechanism", "plain",
			"credential.username", "test",
			"credential.password", "pass123",
		),
	)

	connection, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "server1:9092,server2:9092,server3:9092", connection.GetAsString("uri"))
	assert.Equal(t, "test", connection.GetAsString("username"))
	assert.Equal(t, "pass123", connection.GetAsString("password"))
	assert.Equal(t, "plain", connection.GetAsString("mechanism"))
}

func TestKafkaConnectionResolver(t *testing.T) {
	c := NewKafkaConnectionResolverTest()

	t.Run("Single Connection", c.TestSingleConnection)
	t.Run("Cluster Connection", c.TestClusterConnection)
	t.Run("Cluster Connection with Auth", c.TestClusterConnectionWithAuth)
	t.Run("Connection URI", c.TestConnectionUri)
}
