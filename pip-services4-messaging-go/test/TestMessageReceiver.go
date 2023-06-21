package test

import (
	"context"
	"sync"

	"github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
)

type TestMessageReceiver struct {
	messages []queues.MessageEnvelope
	lock     sync.Mutex
}

func NewTestMessageReceiver() *TestMessageReceiver {
	return &TestMessageReceiver{
		messages: make([]queues.MessageEnvelope, 0),
	}
}

func (c *TestMessageReceiver) GetMessages() []queues.MessageEnvelope {
	c.lock.Lock()
	defer c.lock.Unlock()
	result := make([]queues.MessageEnvelope, len(c.messages))
	copy(result, c.messages)

	return result
}

func (c *TestMessageReceiver) GetMessageCount() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.messages)
}

func (c *TestMessageReceiver) ReceiveMessage(ctx context.Context, envelope *queues.MessageEnvelope, queue queues.IMessageQueue) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.messages = append(c.messages, *envelope)
	return nil
}

func (c *TestMessageReceiver) Clear(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.messages = make([]queues.MessageEnvelope, 0)
	return nil
}
