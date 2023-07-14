package test_queues

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-mqtt-go/queues"
)

type mqttMessageQueueTest struct {
	queue   *queues.MqttMessageQueue
	fixture *MessageQueueFixture
}

func newMqttMessageQueueTest() *mqttMessageQueueTest {
	mqttUri := os.Getenv("MQTT_SERVICE_URI")
	mqttHost := os.Getenv("MQTT_SERVICE_HOST")
	if mqttHost == "" {
		mqttHost = "localhost"
	}

	mqttPort := os.Getenv("MQTT_SERVICE_PORT")
	if mqttPort == "" {
		mqttPort = "1883"
	}

	mqttTopic := os.Getenv("MQTT_TOPIC")
	if mqttTopic == "" {
		mqttTopic = "test"
	}

	mqttUser := os.Getenv("MQTT_USER")
	if mqttUser == "" {
		mqttUser = "mqtt"
	}
	mqttPassword := os.Getenv("MQTT_PASS")
	if mqttPassword == "" {
		mqttPassword = "mqtt"
	}

	if mqttUri == "" && mqttHost == "" {
		return nil
	}

	queue := queues.NewMqttMessageQueue(mqttTopic)
	queue.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"connection.uri", mqttUri,
		"connection.host", mqttHost,
		"connection.port", mqttPort,
		"credential.username", mqttUser,
		"credential.password", mqttPassword,
		"options.autosubscribe", true,
		"options.serialize_envelope", true,
	))

	fixture := NewMessageQueueFixture(queue)

	return &mqttMessageQueueTest{
		queue:   queue,
		fixture: fixture,
	}
}

func (c *mqttMessageQueueTest) setup(t *testing.T) {
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

func (c *mqttMessageQueueTest) teardown(t *testing.T) {
	err := c.queue.Close(context.Background())
	if err != nil {
		t.Error("Failed to close queue", err)
	}
}

func TestMqttMessageQueue(t *testing.T) {
	c := newMqttMessageQueueTest()
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
