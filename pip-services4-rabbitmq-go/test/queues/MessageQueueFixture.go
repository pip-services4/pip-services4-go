package test_queues

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	"github.com/stretchr/testify/assert"
)

type MessageQueueFixture struct {
	queue queues.IMessageQueue
}

func NewMessageQueueFixture(queue queues.IMessageQueue) *MessageQueueFixture {
	c := MessageQueueFixture{
		queue: queue,
	}
	return &c
}

func (c *MessageQueueFixture) TestSendReceiveMessage(t *testing.T) {
	envelope1 := queues.NewMessageEnvelope("123", "Test", []byte("Test message"))
	sndErr := c.queue.Send(context.Background(), envelope1)
	assert.Nil(t, sndErr)

	time.Sleep(100 * time.Millisecond)

	// if c.queue.GetCapabilities().CanMessageCount() {
	// 	count, rdErr := c.queue.MessageCount()
	// 	assert.Nil(t, rdErr)
	// 	assert.Greater(t, count, (int64)(0))
	// }

	envelope2, rcvErr := c.queue.Receive(context.Background(), 10000*time.Millisecond)
	assert.Nil(t, rcvErr)
	assert.NotNil(t, envelope2)
	assert.Equal(t, envelope1.MessageType, envelope2.MessageType)
	assert.Equal(t, envelope1.Message, envelope2.Message)
	assert.Equal(t, envelope1.TraceId, envelope2.TraceId)
}

func (c *MessageQueueFixture) TestMessageCount(t *testing.T) {
	envelope1 := queues.NewMessageEnvelope("123", "Test", []byte("Test message"))
	sendErr := c.queue.Send(context.Background(), envelope1)
	assert.Nil(t, sendErr)

	time.Sleep(500 * time.Millisecond)

	count, err := c.queue.ReadMessageCount()
	assert.Nil(t, err)
	assert.True(t, count >= 1)
}

func (c *MessageQueueFixture) TestReceiveSendMessage(t *testing.T) {
	envelope1 := queues.NewMessageEnvelope("123", "Test", []byte("Test message"))

	sendErr := c.queue.Send(context.Background(), envelope1)
	assert.Nil(t, sendErr)

	time.Sleep(500 * time.Millisecond)

	envelope2, rcvErr := c.queue.Receive(context.Background(), 10000*time.Millisecond)
	assert.Nil(t, rcvErr)
	assert.NotNil(t, envelope2)
	assert.Equal(t, envelope1.MessageType, envelope2.MessageType)
	assert.Equal(t, envelope1.Message, envelope2.Message)
	assert.Equal(t, envelope1.TraceId, envelope2.TraceId)
}

func (c *MessageQueueFixture) TestReceiveCompleteMessage(t *testing.T) {
	envelope1 := queues.NewMessageEnvelope("123", "Test", []byte("Test message"))
	sndErr := c.queue.Send(context.Background(), envelope1)
	assert.Nil(t, sndErr)

	time.Sleep(100 * time.Millisecond)

	// count, rdErr := c.queue.ReadMessageCount()
	// assert.Nil(t, rdErr)
	// assert.Greater(t, count, (int64)(0))

	envelope2, rcvErr := c.queue.Receive(context.Background(), 10000*time.Millisecond)
	assert.Nil(t, rcvErr)
	assert.NotNil(t, envelope2)
	assert.Equal(t, envelope1.MessageType, envelope2.MessageType)
	assert.Equal(t, envelope1.Message, envelope2.Message)
	assert.Equal(t, envelope1.TraceId, envelope2.TraceId)

	cplErr := c.queue.Complete(context.Background(), envelope2)
	assert.Nil(t, cplErr)
	assert.Nil(t, envelope2.GetReference())
}

func (c *MessageQueueFixture) TestReceiveAbandonMessage(t *testing.T) {
	envelope1 := queues.NewMessageEnvelope("123", "Test", []byte("Test message"))
	sndErr := c.queue.Send(context.Background(), envelope1)
	assert.Nil(t, sndErr)

	envelope2, rcvErr := c.queue.Receive(context.Background(), 10000*time.Millisecond)
	assert.Nil(t, rcvErr)
	assert.NotNil(t, envelope2)
	assert.Equal(t, envelope1.MessageType, envelope2.MessageType)
	assert.Equal(t, envelope1.Message, envelope2.Message)
	assert.Equal(t, envelope1.TraceId, envelope2.TraceId)

	abdErr := c.queue.Abandon(context.Background(), envelope2)
	assert.Nil(t, abdErr)

	envelope2, rcvErr = c.queue.Receive(context.Background(), 10000*time.Millisecond)
	assert.Nil(t, rcvErr)
	assert.NotNil(t, envelope2)
	assert.Equal(t, envelope1.MessageType, envelope2.MessageType)
	assert.Equal(t, envelope1.Message, envelope2.Message)
	assert.Equal(t, envelope1.TraceId, envelope2.TraceId)
}

func (c *MessageQueueFixture) TestSendPeekMessage(t *testing.T) {
	envelope1 := queues.NewMessageEnvelope("123", "Test", []byte("Test message"))
	sndErr := c.queue.Send(context.Background(), envelope1)
	assert.Nil(t, sndErr)

	time.Sleep(100 * time.Millisecond)

	envelope2, pkErr := c.queue.Peek(context.Background())
	assert.Nil(t, pkErr)
	assert.NotNil(t, envelope2)
	assert.Equal(t, envelope1.MessageType, envelope2.MessageType)
	assert.Equal(t, envelope1.Message, envelope2.Message)
	assert.Equal(t, envelope1.TraceId, envelope2.TraceId)

	// pop message from queue for next test
	_, rcvErr := c.queue.Receive(context.Background(), 10000*time.Millisecond)
	assert.Nil(t, rcvErr)
}

func (c *MessageQueueFixture) TestPeekNoMessage(t *testing.T) {
	envelope, pkErr := c.queue.Peek(context.Background())
	assert.Nil(t, pkErr)
	assert.Nil(t, envelope)
}

func (c *MessageQueueFixture) TestMoveToDeadMessage(t *testing.T) {
	envelope1 := queues.NewMessageEnvelope("123", "Test", []byte("Test message"))
	sndErr := c.queue.Send(context.Background(), envelope1)
	assert.Nil(t, sndErr)

	envelope2, rcvErr := c.queue.Receive(context.Background(), 10000*time.Millisecond)
	assert.Nil(t, rcvErr)
	assert.NotNil(t, envelope2)
	assert.Equal(t, envelope1.MessageType, envelope2.MessageType)
	assert.Equal(t, envelope1.Message, envelope2.Message)
	assert.Equal(t, envelope1.TraceId, envelope2.TraceId)

	mvErr := c.queue.MoveToDeadLetter(context.Background(), envelope2)
	assert.Nil(t, mvErr)
}

func (c *MessageQueueFixture) TestOnMessage(t *testing.T) {
	envelope1 := queues.NewMessageEnvelope("123", "Test", []byte("Test message"))
	receiver := &TestMsgReceiver{}
	c.queue.BeginListen(context.Background(), receiver)

	time.Sleep(1000 * time.Millisecond)

	sndErr := c.queue.Send(context.Background(), envelope1)
	assert.Nil(t, sndErr)

	time.Sleep(1000 * time.Millisecond)

	envelope2 := receiver.GetEnvelope()
	assert.NotNil(t, envelope2)
	assert.Equal(t, envelope1.MessageType, envelope2.MessageType)
	assert.Equal(t, envelope1.Message, envelope2.Message)
	assert.Equal(t, envelope1.TraceId, envelope2.TraceId)

	c.queue.EndListen(context.Background())
}

type TestMsgReceiver struct {
	lock      sync.Mutex
	_envelope *queues.MessageEnvelope
}

func (c *TestMsgReceiver) GetEnvelope() queues.MessageEnvelope {
	c.lock.Lock()
	defer c.lock.Unlock()
	return *c._envelope
}

func (c *TestMsgReceiver) ReceiveMessage(ctx context.Context, envelope *queues.MessageEnvelope, queue queues.IMessageQueue) (err error) {
	c.lock.Lock()
	c._envelope = envelope
	c.lock.Unlock()
	return nil
}
