package test_queues

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/queues"
)

type natsBareMessageQueueTest struct {
	queue   *queues.NatsBareMessageQueue
	fixture *MessageQueueFixture
}

func newNatsBareMessageQueueTest() *natsBareMessageQueueTest {
	natsUri := os.Getenv("NATS_SERVICE_URI")
	natsHost := os.Getenv("NATS_SERVICE_HOST")
	if natsHost == "" {
		natsHost = "localhost"
	}

	natsPort := os.Getenv("NATS_SERVICE_PORT")
	if natsPort == "" {
		natsPort = "4222"
	}

	natsQueue := os.Getenv("NATS_QUEUE")
	if natsQueue == "" {
		natsQueue = "test"
	}

	natsToken := os.Getenv("NATS_TOKEN")
	if natsToken == "" {
		natsToken = ""
	}

	natsUser := os.Getenv("NATS_USER")
	if natsUser == "" {
		natsUser = "nats"
	}
	natsPassword := os.Getenv("NATS_PASS")
	if natsPassword == "" {
		natsPassword = "nats"
	}

	if natsUri == "" && natsHost == "" {
		return nil
	}

	queue := queues.NewNatsBareMessageQueue(natsQueue)
	queue.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connection.uri", natsUri,
		"connection.host", natsHost,
		"connection.port", natsPort,
		"credential.token", natsToken,
		"credential.username", natsUser,
		"credential.password", natsPassword,
	))

	fixture := NewMessageQueueFixture(queue)

	return &natsBareMessageQueueTest{
		queue:   queue,
		fixture: fixture,
	}
}

func (c *natsBareMessageQueueTest) setup(t *testing.T) {
	err := c.queue.Open(context.Background())
	if err != nil {
		t.Error("Failed to open queue", err)
		return
	}

	// err = c.queue.Clear("")
	// if err != nil {
	// 	t.Error("Failed to clear queue", err)
	// 	return
	// }
}

func (c *natsBareMessageQueueTest) teardown(t *testing.T) {
	err := c.queue.Close(context.Background())
	if err != nil {
		t.Error("Failed to close queue", err)
	}
}

func TestNatsBareMessageQueue(t *testing.T) {
	c := newNatsBareMessageQueueTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("Receive Send Message", c.fixture.TestReceiveSendMessage)
	c.teardown(t)

	c.setup(t)
	t.Run("On Message", c.fixture.TestOnMessage)
	c.teardown(t)
}
