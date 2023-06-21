package build

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
)

// MemoryMessageQueueFactory are creates MemoryMemoryMessageQueue components by their descriptors.
// Name of created message queue is taken from its descriptor.
//
// See Factory
// See MemoryMemoryMessageQueue
type MemoryMessageQueueFactory struct {
	*MessageQueueFactory
}

// NewMemoryMessageQueueFactory method are create a new instance of the factory.
func NewMemoryMessageQueueFactory() *MemoryMessageQueueFactory {
	c := MemoryMessageQueueFactory{
		MessageQueueFactory: InheritMessageQueueFactory(),
	}

	memoryQueueDescriptor := cref.NewDescriptor("pip-services", "message-queue", "memory", "*", "1.0")

	c.Register(memoryQueueDescriptor, func(locator any) any {
		name := ""
		descriptor, ok := locator.(*cref.Descriptor)
		if ok {
			name = descriptor.Name()
		}
		return c.CreateQueue(context.Background(), name)
	})

	return &c
}

// Creates a message queue component and assigns its name.
//
// Parameters:
//   - name: a name of the created message queue.
func (c *MemoryMessageQueueFactory) CreateQueue(ctx context.Context, name string) queues.IMessageQueue {
	queue := queues.NewMemoryMessageQueue(name)

	if c.Config != nil {
		queue.Configure(ctx, c.Config)
	}
	if c.References != nil {
		queue.SetReferences(ctx, c.References)
	}

	return queue
}
