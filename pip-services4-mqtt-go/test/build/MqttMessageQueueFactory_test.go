package test_build

import (
	"testing"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-mqtt-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-mqtt-go/queues"
	"github.com/stretchr/testify/assert"
)

func TestMqttMessageQueueFactory(t *testing.T) {
	factory := build.NewMqttMessageQueueFactory()
	descriptor := cref.NewDescriptor("pip-services", "message-queue", "mqtt", "test", "1.0")

	canResult := factory.CanCreate(descriptor)
	assert.NotNil(t, canResult)

	comp, err := factory.Create(descriptor)
	assert.Nil(t, err)
	assert.NotNil(t, comp)

	queue := comp.(*queues.MqttMessageQueue)
	assert.Equal(t, "test", queue.Name())
}
