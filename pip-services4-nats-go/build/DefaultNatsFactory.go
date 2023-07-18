package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/connect"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/queues"
)

// Creates NatsMessageQueue components by their descriptors.
// See NatsMessageQueue
type DefaultNatsFactory struct {
	*cbuild.Factory
}

// NewDefaultNatsFactory method are create a new instance of the factory.
func NewDefaultNatsFactory() *DefaultNatsFactory {
	c := DefaultNatsFactory{}
	c.Factory = cbuild.NewFactory()

	natsQueueFactoryDescriptor := cref.NewDescriptor("pip-services", "queue-factory", "nats", "*", "1.0")
	natsConnectionDescriptor := cref.NewDescriptor("pip-services", "connection", "nats", "*", "1.0")
	bareNatsQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "bare-nats", "*", "1.0")
	natsQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "nats", "*", "1.0")

	c.RegisterType(natsQueueFactoryDescriptor, NewNatsMessageQueueFactory)

	c.RegisterType(natsConnectionDescriptor, connect.NewNatsConnection)

	c.Register(bareNatsQueueDescriptor, func(locator interface{}) interface{} {
		name := ""
		descriptor, ok := locator.(*cref.Descriptor)
		if ok {
			name = descriptor.Name()
		}

		return queues.NewNatsBareMessageQueue(name)
	})

	c.Register(natsQueueDescriptor, func(locator interface{}) interface{} {
		name := ""
		descriptor, ok := locator.(*cref.Descriptor)
		if ok {
			name = descriptor.Name()
		}

		return queues.NewNatsMessageQueue(name)
	})

	return &c
}
