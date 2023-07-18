package test_build

import (
	"testing"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	build "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/build"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/queues"
	"github.com/stretchr/testify/assert"
)

func TestNatsMessageQueueFactory(t *testing.T) {
	factory := build.NewNatsMessageQueueFactory()
	descriptor := cref.NewDescriptor("pip-services", "message-queue", "bare-nats", "test", "1.0")

	canResult := factory.CanCreate(descriptor)
	assert.NotNil(t, canResult)

	comp, err := factory.Create(descriptor)
	assert.Nil(t, err)
	assert.NotNil(t, comp)

	queue := comp.(*queues.NatsBareMessageQueue)
	assert.Equal(t, "test", queue.Name())
}
