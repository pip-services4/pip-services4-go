package test_connect

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/connect"
	"github.com/stretchr/testify/assert"
)

type natsConnectionTest struct {
	connection *connect.NatsConnection
}

func newNatsConnectionTest() *natsConnectionTest {
	natsUri := os.Getenv("NATS_SERVICE_URI")
	natsHost := os.Getenv("NATS_SERVICE_HOST")
	if natsHost == "" {
		natsHost = "localhost"
	}

	natsPort := os.Getenv("NATS_SERVICE_PORT")
	if natsPort == "" {
		natsPort = "4222"
	}

	natsToken := os.Getenv("NATS_TOKEN")
	// if natsToken == "" {
	// 	natsToken = ""
	// }
	natsUser := os.Getenv("NATS_USER")
	// if natsUser == "" {
	// 	natsUser = ""
	// }
	natsPassword := os.Getenv("NATS_PASS")
	// if natsPassword == "" {
	// 	natsPassword = ""
	// }

	if natsUri == "" && natsHost == "" {
		return nil
	}

	connection := connect.NewNatsConnection()
	connection.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connection.uri", natsUri,
		"connection.host", natsHost,
		"connection.port", natsPort,
		"credential.token", natsToken,
		"credential.username", natsUser,
		"credential.password", natsPassword,
	))

	return &natsConnectionTest{
		connection: connection,
	}
}

func (c *natsConnectionTest) TestOpenClose(t *testing.T) {
	ctx := context.Background()

	err := c.connection.Open(ctx)
	assert.Nil(t, err)
	assert.True(t, c.connection.IsOpen())
	assert.NotNil(t, c.connection.GetConnection())

	err = c.connection.Close(ctx)
	assert.Nil(t, err)
	assert.False(t, c.connection.IsOpen())
	assert.Nil(t, c.connection.GetConnection())
}

func TestNatsConnection(t *testing.T) {
	c := newNatsConnectionTest()
	if c == nil {
		return
	}

	t.Run("Open and Close", c.TestOpenClose)
}
