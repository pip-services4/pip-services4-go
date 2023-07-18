package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-rabbitmq-go/queues"
)

// Creates RabbitMQMessageQueue components by their descriptors.
// See RabbitMQMessageQueue
type DefaultRabbitMQFactory struct {
	*cbuild.Factory
}

// NewDefaultRabbitMQFactory method are create a new instance of the factory.
func NewDefaultRabbitMQFactory() *DefaultRabbitMQFactory {
	c := DefaultRabbitMQFactory{}
	c.Factory = cbuild.NewFactory()

	rabbitMQMessageQueueFactoryDescriptor := cref.NewDescriptor("pip-services", "queue-factory", "rabbitmq", "*", "1.0")
	rabbitMQMessageQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "rabbitmq", "*", "1.0")

	c.RegisterType(rabbitMQMessageQueueFactoryDescriptor, NewRabbitMQMessageQueueFactory)

	c.Register(rabbitMQMessageQueueDescriptor, func(locator interface{}) interface{} {
		name := ""
		descriptor, ok := locator.(*cref.Descriptor)
		if ok {
			name = descriptor.Name()
		}

		return queues.NewEmptyRabbitMQMessageQueue(name)
	})

	return &c
}
