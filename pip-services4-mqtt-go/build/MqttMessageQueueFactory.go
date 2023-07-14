package build

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/build"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	"github.com/pip-services4/pip-services4-go/pip-services4-mqtt-go/queues"
)

// MqttMessageQueueFactory are creates MqttMessageQueue components by their descriptors.
// Name of created message queue is taken from its descriptor.
//
// See Factory
// See MqttMessageQueue
type MqttMessageQueueFactory struct {
	*build.MessageQueueFactory
}

// NewMqttMessageQueueFactory method are create a new instance of the factory.
func NewMqttMessageQueueFactory() *MqttMessageQueueFactory {
	c := MqttMessageQueueFactory{
		MessageQueueFactory: build.InheritMessageQueueFactory(),
	}

	mqttQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "mqtt", "*", "1.0")

	c.Register(mqttQueueDescriptor, func(locator interface{}) interface{} {
		name := ""
		descriptor, ok := locator.(*cref.Descriptor)
		if ok {
			name = descriptor.Name()
		}
		return c.CreateQueue(name)
	})

	return &c
}

// Creates a message queue component and assigns its name.
//
// Parameters:
//   - name: a name of the created message queue.
func (c *MqttMessageQueueFactory) CreateQueue(name string) cqueues.IMessageQueue {
	queue := queues.NewMqttMessageQueue(name)

	if c.Config != nil {
		queue.Configure(context.Background(), c.Config)
	}
	if c.References != nil {
		queue.SetReferences(context.Background(), c.References)
	}

	return queue
}
