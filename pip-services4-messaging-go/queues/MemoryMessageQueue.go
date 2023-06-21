package queues

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
)

// MemoryMessageQueue Message queue that sends and receives messages within the same process by using shared memory.
// This queue is typically used for testing to mock real queues.
//
//	Configuration parameters:
//		- name: name of the message queue
//	References:
//		- *:logger:*:*:1.0           (optional)  ILogger components to pass log messages
//		- *:counters:*:*:1.0         (optional)  ICounters components to pass collected measurements
//
//	see MessageQueue
//	see MessagingCapabilities
//
//	Example:
//		queue := NewMessageQueue("myqueue");
//		queue.Send(context.Background(), "123", NewMessageEnvelop("", "mymessage", "ABC"));
//		message, err := queue.Receive(contex.Backgroudn(), "123", 10000*time.Milliseconds)
//		if (message != nil) {
//			...
//			queue.Complete(ctx, message);
//		}
type MemoryMessageQueue struct {
	MessageQueue
	messages          []*MessageEnvelope
	lockTokenSequence int
	lockedMessages    map[int]*LockedMessage
	opened            bool
	cancel            int32
}

// NewMemoryMessageQueue method are creates a new instance of the message queue.
//
//	Parameters:
//		- name  (optional) a queue name.
//	Returns: *MemoryMessageQueue
//	see MessagingCapabilities
func NewMemoryMessageQueue(name string) *MemoryMessageQueue {
	c := MemoryMessageQueue{}

	c.MessageQueue = *InheritMessageQueue(
		&c, name, NewMessagingCapabilities(true, true, true, true, true, true, true, false, true),
	)

	c.messages = make([]*MessageEnvelope, 0)
	c.lockTokenSequence = 0
	c.lockedMessages = make(map[int]*LockedMessage, 0)
	c.opened = false
	c.cancel = 0

	return &c
}

// IsOpen method are checks if the component is opened.
//
//	Returns: true if the component has been opened and false otherwise.
func (c *MemoryMessageQueue) IsOpen() bool {
	return c.opened
}

// Open method are opens the component with given connection and credential parameters.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId     	(optional) transaction id to trace execution through call chain.
//		- connection        	is a connection parameters
//		- credential        	is a credential parameters
//	Returns: error or nil no errors occured.
func (c *MemoryMessageQueue) Open(ctx context.Context) (err error) {
	c.opened = true

	c.Logger.Debug(ctx, "Opened queue %s", c.Name())

	return nil
}

// Close method are closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId 		(optional) transaction id to trace execution through call chain.
//	Returns: error or nil no errors occured.
func (c *MemoryMessageQueue) Close(ctx context.Context) (err error) {
	c.opened = false
	atomic.StoreInt32(&c.cancel, 1)

	c.Logger.Debug(ctx, "Closed queue %s", c.Name())

	return nil
}

// Clear method are clears component state.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId 		(optional) transaction id to trace execution through call chain.
//	Returns: error or nil no errors occured.
func (c *MemoryMessageQueue) Clear(ctx context.Context) (err error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.messages = make([]*MessageEnvelope, 0)
	c.lockedMessages = make(map[int]*LockedMessage, 0)
	atomic.StoreInt32(&c.cancel, 0)

	return nil
}

// ReadMessageCount method are reads the current number of messages in the queue to be delivered.
//
//	Returns: number of messages or error.
func (c *MemoryMessageQueue) ReadMessageCount() (count int64, err error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	count = (int64)(len(c.messages))
	return count, nil
}

// Send method are sends a message into the queue.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId     	(optional) transaction id to trace execution through call chain.
//		- envelope          	a message envelop to be sent.
//	Returns: error or nil for success.
func (c *MemoryMessageQueue) Send(ctx context.Context, envelope *MessageEnvelope) (err error) {
	envelope.SentTime = time.Now()

	// Add message to the queue
	c.Lock.Lock()
	c.messages = append(c.messages, envelope)
	c.Lock.Unlock()

	c.Counters.IncrementOne(ctx, "queue."+c.Name()+".sent_messages")
	ctx = utils.ContextHelper.NewContextWithTraceId(ctx, envelope.TraceId)
	c.Logger.Debug(ctx, "Sent message %s via %s", envelope.String(), c.Name())

	return nil
}

// Peek meethod are peeks a single incoming message from the queue without removing it.
// If there are no messages available in the queue it returns nil.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId     	(optional) transaction id to trace execution through call chain.
//	Returns: a message or error.
func (c *MemoryMessageQueue) Peek(ctx context.Context) (result *MessageEnvelope, err error) {
	var message *MessageEnvelope

	// Pick a message
	c.Lock.Lock()
	if len(c.messages) > 0 {
		message = c.messages[0]
	}
	c.Lock.Unlock()

	if message != nil {
		ctx = utils.ContextHelper.NewContextWithTraceId(ctx, message.TraceId)
		c.Logger.Trace(ctx, "Peeked message %s on %s", message, c.String())
	}

	return message, nil
}

// PeekBatch method are peeks multiple incoming messages from the queue without removing them.
// If there are no messages available in the queue it returns an empty list.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId     	(optional) transaction id to trace execution through call chain.
//		- messageCount      	a maximum number of messages to peek.
//	Returns: a list with messages or error.
func (c *MemoryMessageQueue) PeekBatch(ctx context.Context, messageCount int64) (result []*MessageEnvelope, err error) {
	c.Lock.Lock()
	batchMessages := c.messages
	if messageCount <= (int64)(len(batchMessages)) {
		batchMessages = batchMessages[0:messageCount]
	}
	c.Lock.Unlock()

	messages := []*MessageEnvelope{}
	messages = append(messages, batchMessages...)

	c.Logger.Trace(ctx, "Peeked %d messages on %s", len(messages), c.Name())

	return messages, nil
}

// Receive method are receives an incoming message and removes it from the queue.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId     	(optional) transaction id to trace execution through call chain.
//		- waitTimeout       	a timeout in milliseconds to wait for a message to come.
//	Returns: a message or error.
func (c *MemoryMessageQueue) Receive(ctx context.Context, waitTimeout time.Duration) (*MessageEnvelope, error) {
	messageReceived := false
	var message *MessageEnvelope
	elapsedTime := time.Duration(0)

	for elapsedTime < waitTimeout && !messageReceived {
		c.Lock.Lock()
		if len(c.messages) == 0 {
			c.Lock.Unlock()
			time.Sleep(time.Duration(100) * time.Millisecond)
			elapsedTime += time.Duration(100)
			continue
		}

		// Get message from the queue
		message = c.messages[0]
		c.messages = c.messages[1:]

		// Generate and set locked token
		lockedToken := c.lockTokenSequence
		c.lockTokenSequence++
		message.SetReference(lockedToken)

		// Add messages to locked messages list
		now := time.Now().Add(waitTimeout)
		lockedMessage := &LockedMessage{
			ExpirationTime: now,
			Message:        message,
			Timeout:        waitTimeout,
		}
		c.lockedMessages[lockedToken] = lockedMessage

		messageReceived = true
		c.Lock.Unlock()
	}

	if message != nil {
		c.Counters.IncrementOne(ctx, "queue."+c.Name()+".received_messages")
		ctx = utils.ContextHelper.NewContextWithTraceId(ctx, message.TraceId)
		c.Logger.Debug(ctx, "Received message %s via %s", message, c.Name())
	}

	return message, nil
}

// RenewLock method are renews a lock on a message that makes it invisible from other receivers in the queue.
// This method is usually used to extend the message processing time.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- message       		a message to extend its lock.
//		- lockTimeout   		a locking timeout in milliseconds.
//	Returns: error or nil for success.
func (c *MemoryMessageQueue) RenewLock(ctx context.Context, message *MessageEnvelope, lockTimeout time.Duration) (err error) {
	reference := message.GetReference()
	if reference == nil {
		return nil
	}

	c.Lock.Lock()
	// Get message from locked queue
	lockedToken := reference.(int)
	if lockedMessage, ok := c.lockedMessages[lockedToken]; ok {
		// If lock is found, extend the lock
		now := time.Now()
		// Todo: Shall we skip if the message already expired?
		if lockedMessage.ExpirationTime.After(now) {
			lockedMessage.ExpirationTime = now.Add(lockedMessage.Timeout)
		}
	}
	c.Lock.Unlock()

	ctx = utils.ContextHelper.NewContextWithTraceId(ctx, message.TraceId)
	c.Logger.Trace(ctx, "Renewed lock for message %s at %s", message, c.Name())

	return nil
}

// Complete method are permanently removes a message from the queue.
// This method is usually used to remove the message after successful processing.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- message  				a message to remove.
//	Returns: error or nil for success.
func (c *MemoryMessageQueue) Complete(ctx context.Context, message *MessageEnvelope) (err error) {
	reference := message.GetReference()
	if reference == nil {
		return nil
	}

	c.Lock.Lock()
	lockedToken := reference.(int)
	delete(c.lockedMessages, lockedToken)
	message.SetReference(nil)
	c.Lock.Unlock()

	ctx = utils.ContextHelper.NewContextWithTraceId(ctx, message.TraceId)
	c.Logger.Trace(ctx, "Completed message %s at %s", message, c.Name())

	return nil
}

// Abandon method are returns message into the queue and makes it available for all subscribers to receive it again.
// This method is usually used to return a message which could not be processed at the moment
// to repeat the attempt. Messages that cause unrecoverable errors shall be removed permanently
// or/and send to dead letter queue.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- message   			a message to return.
//	Returns: error or nil for success.
func (c *MemoryMessageQueue) Abandon(ctx context.Context, message *MessageEnvelope) (err error) {
	reference := message.GetReference()
	if reference == nil {
		return nil
	}

	c.Lock.Lock()
	// Get message from locked queue
	lockedToken := reference.(int)
	if lockedMessage, ok := c.lockedMessages[lockedToken]; ok {
		// Remove from locked messages
		delete(c.lockedMessages, lockedToken)
		message.SetReference(nil)

		// Skip if it is already expired
		if lockedMessage.ExpirationTime.Before(time.Now()) {
			c.Lock.Unlock()
			return nil
		}
	} else { // Skip if it absent
		c.Lock.Unlock()
		return nil
	}
	c.Lock.Unlock()

	c.Logger.Trace(ctx, message.TraceId, "Abandoned message %s at %s", message, c.Name())

	// Add back to message queue
	ctx = utils.ContextHelper.NewContextWithTraceId(ctx, message.TraceId)
	return c.Send(ctx, message)
}

// MoveToDeadLetter method are permanently removes a message from the queue and sends it to dead letter queue.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- message   			a message to be removed.
//	Returns: error or nil for success.
func (c *MemoryMessageQueue) MoveToDeadLetter(ctx context.Context, message *MessageEnvelope) (err error) {
	reference := message.GetReference()
	if reference == nil {
		return nil
	}

	c.Lock.Lock()
	if lockedToken, ok := reference.(int); ok {
		delete(c.lockedMessages, lockedToken)
		message.SetReference(nil)
	}
	c.Lock.Unlock()

	c.Counters.IncrementOne(ctx, "queue."+c.Name()+".dead_messages")
	ctx = utils.ContextHelper.NewContextWithTraceId(ctx, message.TraceId)
	c.Logger.Trace(ctx, "Moved to dead message %s at %s", message, c.Name())

	return nil
}

// Listen method are listens for incoming messages and blocks the current thread until queue is closed.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId     	(optional) transaction id to trace execution through call chain.
//		- receiver          	a receiver to receive incoming messages.
//	see IMessageReceiver
//	see Receive
func (c *MemoryMessageQueue) Listen(ctx context.Context, receiver IMessageReceiver) error {
	c.Logger.Trace(ctx, "", "Started listening messages at %s", c.String())

	// Unset cancellation token
	atomic.StoreInt32(&c.cancel, 0)

	for atomic.LoadInt32(&c.cancel) == 0 {
		message, err := c.Receive(ctx, time.Duration(1000)*time.Millisecond)
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to receive the message")
		}

		if message != nil && atomic.LoadInt32(&c.cancel) == 0 {
			// Todo: shall we recover after panic here??
			func(message *MessageEnvelope) {
				defer func() {
					if r := recover(); r != nil {
						err := fmt.Sprintf("%v", r)
						c.Logger.Error(ctx, nil, "Failed to process the message - "+err)
					}
				}()

				err = receiver.ReceiveMessage(ctx, message, c)
				if err != nil {
					c.Logger.Error(ctx, err, "Failed to process the message")
				}
			}(message)
		}
	}

	return nil
}

// EndListen method are ends listening for incoming messages.
// When c method is call listen unblocks the thread and execution continues.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- correlationId     	(optional) transaction id to trace execution through call chain.
func (c *MemoryMessageQueue) EndListen(ctx context.Context) {
	atomic.StoreInt32(&c.cancel, 1)
}
