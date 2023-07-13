package test_build

import (
	"testing"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	build "github.com/pip-services4/pip-services4-go/pip-services4-kafka-go/build"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-kafka-go/queues"
	"github.com/stretchr/testify/assert"
)

func TestKafkaMessageQueueFactory(t *testing.T) {
	factory := build.NewKafkaMessageQueueFactory()
	descriptor := cref.NewDescriptor("pip-services", "message-queue", "kafka", "test", "1.0")

	canResult := factory.CanCreate(descriptor)
	assert.NotNil(t, canResult)

	comp, err := factory.Create(descriptor)
	assert.Nil(t, err)
	assert.NotNil(t, comp)

	queue := comp.(*queues.KafkaMessageQueue)
	assert.Equal(t, "test", queue.Name())
}
