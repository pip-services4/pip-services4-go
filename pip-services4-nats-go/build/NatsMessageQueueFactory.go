package build

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/build"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	"github.com/pip-services4/pip-services4-go/pip-services4-nats-go/queues"
)

// NatsMessageQueueFactory are creates NatsMessageQueue components by their descriptors.
// Name of created message queue is taken from its descriptor.
//
// See Factory
// See NatsMessageQueue
type NatsMessageQueueFactory struct {
	*build.MessageQueueFactory
}

// NewNatsMessageQueueFactory method are create a new instance of the factory.
func NewNatsMessageQueueFactory() *NatsMessageQueueFactory {
	c := NatsMessageQueueFactory{
		MessageQueueFactory: build.InheritMessageQueueFactory(),
	}

	bareNatsQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "bare-nats", "*", "1.0")
	natsQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "nats", "*", "1.0")

	c.Register(bareNatsQueueDescriptor, func(locator interface{}) interface{} {
		name := ""
		descriptor, ok := locator.(*cref.Descriptor)
		if ok {
			name = descriptor.Name()
		}
		return c.CreateBareQueue(name)
	})

	c.Register(natsQueueDescriptor, func(locator interface{}) interface{} {
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
func (c *NatsMessageQueueFactory) CreateQueue(name string) cqueues.IMessageQueue {
	queue := queues.NewNatsMessageQueue(name)

	if c.Config != nil {
		queue.Configure(context.Background(), c.Config)
	}
	if c.References != nil {
		queue.SetReferences(context.Background(), c.References)
	}

	return queue
}

// Creates a message queue component and assigns its name.
//
// Parameters:
//   - name: a name of the created message queue.
func (c *NatsMessageQueueFactory) CreateBareQueue(name string) cqueues.IMessageQueue {
	queue := queues.NewNatsBareMessageQueue(name)

	if c.Config != nil {
		queue.Configure(context.Background(), c.Config)
	}
	if c.References != nil {
		queue.SetReferences(context.Background(), c.References)
	}

	return queue
}
