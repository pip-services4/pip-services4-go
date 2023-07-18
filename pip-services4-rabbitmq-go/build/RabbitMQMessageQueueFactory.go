package build

import (
	"context"

	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	queues "github.com/pip-services4/pip-services4-go/pip-services4-rabbitmq-go/queues"
)

type RabbitMQMessageQueueFactory struct {
	*cbuild.Factory
	config     *cconf.ConfigParams
	references cref.IReferences
}

func NewRabbitMQMessageQueueFactory() *RabbitMQMessageQueueFactory {
	c := RabbitMQMessageQueueFactory{}
	c.Factory = cbuild.NewFactory()

	memoryQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "rabbitmq", "*", "*")

	c.Register(memoryQueueDescriptor, func(locator interface{}) interface{} {
		name := ""
		descriptor, ok := locator.(*cref.Descriptor)
		if ok {
			name = descriptor.Name()
		}

		return queues.NewEmptyRabbitMQMessageQueue(name)
	})
	return &c
}

func (c *RabbitMQMessageQueueFactory) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.config = config
}

func (c *RabbitMQMessageQueueFactory) SetReferences(ctx context.Context, references cref.IReferences) {
	c.references = references
}

// Creates a message queue component and assigns its name.
//
// Parameters:
//   - name: a name of the created message queue.
func (c *RabbitMQMessageQueueFactory) CreateQueue(name string) cqueues.IMessageQueue {
	queue := queues.NewEmptyRabbitMQMessageQueue(name)

	if c.config != nil {
		queue.Configure(context.Background(), c.config)
	}
	if c.references != nil {
		queue.SetReferences(context.Background(), c.references)
	}

	return queue
}
