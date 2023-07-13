package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-kafka-go/connect"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-kafka-go/queues"
)

// Creates KafkaMessageQueue components by their descriptors.
// See KafkaMessageQueue
type DefaultKafkaFactory struct {
	*cbuild.Factory
}

// NewDefaultKafkaFactory method are create a new instance of the factory.
func NewDefaultKafkaFactory() *DefaultKafkaFactory {
	c := DefaultKafkaFactory{}
	c.Factory = cbuild.NewFactory()

	kafkaQueueFactoryDescriptor := cref.NewDescriptor("pip-services", "queue-factory", "kafka", "*", "1.0")
	kafkaConnectionDescriptor := cref.NewDescriptor("pip-services", "connection", "kafka", "*", "1.0")
	kafkaQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "kafka", "*", "1.0")

	c.RegisterType(kafkaQueueFactoryDescriptor, NewKafkaMessageQueueFactory)

	c.RegisterType(kafkaConnectionDescriptor, connect.NewKafkaConnection)

	c.Register(kafkaQueueDescriptor, func(locator interface{}) interface{} {
		name := ""
		descriptor, ok := locator.(*cref.Descriptor)
		if ok {
			name = descriptor.Name()
		}

		return queues.NewKafkaMessageQueue(name)
	})

	return &c
}
