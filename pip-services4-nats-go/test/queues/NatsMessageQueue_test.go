package test_queues

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/queues"
)

type natsMessageQueueTest struct {
	queue   *queues.NatsMessageQueue
	fixture *MessageQueueFixture
}

func newNatsMessageQueueTest() *natsMessageQueueTest {
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

	queue := queues.NewNatsMessageQueue(natsQueue)
	queue.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connection.uri", natsUri,
		"connection.host", natsHost,
		"connection.port", natsPort,
		"credential.token", natsToken,
		"credential.username", natsUser,
		"credential.password", natsPassword,
		"options.autosubscribe", true,
	))

	fixture := NewMessageQueueFixture(queue)

	return &natsMessageQueueTest{
		queue:   queue,
		fixture: fixture,
	}
}

func (c *natsMessageQueueTest) setup(t *testing.T) {
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

func (c *natsMessageQueueTest) teardown(t *testing.T) {
	err := c.queue.Close(context.Background())
	if err != nil {
		t.Error("Failed to close queue", err)
	}
}

func TestNatsMessageQueue(t *testing.T) {
	c := newNatsMessageQueueTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("Send Receive Message", c.fixture.TestSendReceiveMessage)
	c.teardown(t)

	c.setup(t)
	t.Run("Receive Send Message", c.fixture.TestReceiveSendMessage)
	c.teardown(t)

	c.setup(t)
	t.Run("Send Peek Message", c.fixture.TestSendPeekMessage)
	c.teardown(t)

	c.setup(t)
	t.Run("Peek No Message", c.fixture.TestPeekNoMessage)
	c.teardown(t)

	c.setup(t)
	t.Run("On Message", c.fixture.TestOnMessage)
	c.teardown(t)
}
