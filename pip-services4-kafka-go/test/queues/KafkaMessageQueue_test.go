package test_queues

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-kafka-go/queues"
)

type kafkaMessageQueueTest struct {
	queue   *queues.KafkaMessageQueue
	fixture *MessageQueueFixture
}

func newKafkaMessageQueueTest() *kafkaMessageQueueTest {
	kafkaUri := os.Getenv("KAFKA_SERVICE_URI")
	kafkaHost := os.Getenv("KAFKA_SERVICE_HOST")
	if kafkaHost == "" {
		kafkaHost = "localhost"
	}

	kafkaPort := os.Getenv("KAFKA_SERVICE_PORT")
	if kafkaPort == "" {
		kafkaPort = "9092"
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "test"
	}

	kafkaUser := os.Getenv("KAFKA_USER")
	// if kafkaUser == "" {
	// 	kafkaUser = "kafka"
	// }
	kafkaPassword := os.Getenv("KAFKA_PASS")
	// if kafkaPassword == "" {
	// 	kafkaPassword = "pass123"
	// }

	if kafkaUri == "" && kafkaHost == "" {
		return nil
	}

	queue := queues.NewKafkaMessageQueue(kafkaTopic)
	queue.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connection.uri", kafkaUri,
		"connection.host", kafkaHost,
		"connection.port", kafkaPort,
		"credential.mechanism", "plain",
		"credential.username", kafkaUser,
		"credential.password", kafkaPassword,
		"options.autosubscribe", true,
		"options.num_partitions", 2,
		"options.read_partitions", "1",
		"options.write_partition", 1,
		"options.read_partitions", "1",
		"options.write_partition", "1",
		"options.listen_connection", true,
	))

	fixture := NewMessageQueueFixture(queue)

	return &kafkaMessageQueueTest{
		queue:   queue,
		fixture: fixture,
	}
}

func (c *kafkaMessageQueueTest) setup(t *testing.T) {
	err := c.queue.Open(context.Background())
	if err != nil {
		t.Error("Failed to open queue", err)
		return
	}

	err = c.queue.Clear(context.Background())
	if err != nil {
		t.Error("Failed to clear queue", err)
		return
	}
}

func (c *kafkaMessageQueueTest) teardown(t *testing.T) {
	err := c.queue.Close(context.Background())
	if err != nil {
		t.Error("Failed to close queue", err)
	}
}

func TestKafkaMessageQueue(t *testing.T) {
	c := newKafkaMessageQueueTest()
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
