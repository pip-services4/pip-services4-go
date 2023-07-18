package queues

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
)

//NatsMessageQueue are message queue that sends and receives messages via NATS message broker.
//
//	Configuration parameters:
//
//		- subject:                       name of NATS topic (subject) to subscribe
//		- queue_group:                   name of NATS queue group
//		- connection(s):
//			- discovery_key:               (optional) a key to retrieve the connection from  IDiscovery
//			- host:                        host name or IP address
//			- port:                        port number
//			- uri:                         resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                   (optional) a key to retrieve the credentials from  ICredentialStore
//			- username:                    user name
//			- password:                    user password
//		- options:
//			- serialize_message:    (optional) true to serialize entire message as JSON, false to send only message payload (default: true)
//			- autosubscribe:        (optional) true to automatically subscribe on option (default: false)
//			- retry_connect:        (optional) turns on/off automated reconnect when connection is log (default: true)
//			- max_reconnect:        (optional) maximum reconnection attempts (default: 3)
//			- reconnect_timeout:    (optional) number of milliseconds to wait on each reconnection attempt (default: 3000)
//			- flush_timeout:        (optional) number of milliseconds to wait on flushing messages (default: 3000)
//
//
// References:
//
//		- *:logger:*:*:1.0             (optional)  ILogger components to pass log messages
//		- *:counters:*:*:1.0           (optional)  ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0          (optional)  IDiscovery services to resolve connections
//		- *:credential-store:*:*:1.0   (optional) Credential stores to resolve credentials
//		- *:connection:nats:*:1.0      (optional) Shared connection to NATS service
//
// See MessageQueue
// See MessagingCapabilities
//
//	Example:
//		ctx := context.Background()
//		queue := NewNatsMessageQueue("myqueue")
//		queue.Configure(ctx, cconf.NewConfigParamsFromTuples(
//			"subject", "mytopic",
//			"queue_group", "mygroup",
//			"connection.protocol", "nats"
//			"connection.host", "localhost"
//			"connection.port", 1883
//		))
//
//		_ = queue.Open(ctx)
//
//		_ = queue.Send(ctx NewMessageEnvelope("", "mymessage", "ABC"))
//
//		message, err := queue.Receive(ctx 10000*time.Milliseconds)
//		if (message != nil) {
//			...
//			queue.Complete(ctx, message);
//		}

type NatsMessageQueue struct {
	*NatsAbstractMessageQueue

	autoSubscribe bool
	subscribed    bool
	messages      []*cqueues.MessageEnvelope
	receiver      cqueues.IMessageReceiver
}

// NewNatsMessageQueue are creates a new instance of the message queue.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- name  string (optional) a queue name.
func NewNatsMessageQueue(name string) *NatsMessageQueue {
	c := NatsMessageQueue{}

	c.NatsAbstractMessageQueue = InheritNatsAbstractMessageQueue(&c, name,
		cqueues.NewMessagingCapabilities(false, true, true, true, true, false, false, false, true))

	c.messages = make([]*cqueues.MessageEnvelope, 0)

	return &c
}

// Configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- config    configuration parameters to be set.
func (c *NatsMessageQueue) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.NatsAbstractMessageQueue.Configure(ctx, config)

	c.autoSubscribe = config.GetAsBooleanWithDefault("options.autosubscribe", c.autoSubscribe)
}

// Opens the component with given connection and credential parameters.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//
// Returns error or nil no errors occured.
func (c *NatsMessageQueue) Open(ctx context.Context) error {
	if c.IsOpen() {
		return nil
	}

	err := c.NatsAbstractMessageQueue.Open(ctx)
	if err != nil {
		return err
	}

	// Subscribe right away
	if c.autoSubscribe {
		err = c.subscribe(ctx)
		if err != nil {
			c.Close(ctx)
			return err
		}
	}

	return nil
}

// Close method are Closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//
// Returns error or nil no errors occured.
func (c *NatsMessageQueue) Close(ctx context.Context) error {
	if !c.IsOpen() {
		return nil
	}

	err := c.NatsAbstractMessageQueue.Close(ctx)

	// Unsubscribe from topic
	if c.subscribed {
		subject := c.SubscriptionSubject()
		c.Connection.Unsubscribe(ctx, subject, c.QueueGroup, c)
		c.subscribed = false
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.receiver = nil
	c.messages = make([]*cqueues.MessageEnvelope, 0)

	return err
}

func (c *NatsMessageQueue) subscribe(ctx context.Context) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// Check if already were subscribed
	if c.subscribed {
		return nil
	}

	// Subscribe to the topic
	subject := c.SubscriptionSubject()
	err := c.Connection.Subscribe(ctx, subject, c.QueueGroup, c)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to subscribe to subject "+subject)
		return err
	}

	c.subscribed = true
	return nil
}

// Clear method are clears component state.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//
// Returns error or nil no errors occured.
func (c *NatsMessageQueue) Clear(ctx context.Context) (err error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.messages = make([]*cqueues.MessageEnvelope, 0)
	c.receiver = nil

	return nil
}

// ReadMessageCount method are reads the current number of messages in the queue to be delivered.
// Returns number of messages or error.
func (c *NatsMessageQueue) ReadMessageCount() (count int64, err error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	count = (int64)(len(c.messages))
	return count, nil
}

// Peek method are peeks a single incoming message from the queue without removing it.
// If there are no messages available in the queue it returns nil.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//
// Returns: result *cqueues.MessageEnvelope, err error
// message or error.
func (c *NatsMessageQueue) Peek(ctx context.Context) (*cqueues.MessageEnvelope, error) {
	err := c.CheckOpen(cctx.GetTraceId(ctx))
	if err != nil {
		return nil, err
	}

	// Subscribe if needed
	err = c.subscribe(ctx)
	if err != nil {
		return nil, err
	}

	var message *cqueues.MessageEnvelope

	// Pick a message
	c.Lock.Lock()
	if len(c.messages) > 0 {
		message = c.messages[0]
	}
	c.Lock.Unlock()

	if message != nil {
		c.Logger.Trace(cctx.NewContextWithTraceId(ctx, message.TraceId), "Peeked message %s on %s", message, c.String())
	}

	return message, nil
}

// PeekBatch method are peeks multiple incoming messages from the queue without removing them.
// If there are no messages available in the queue it returns an empty list.
// Important: This method is not supported by NATS.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- messageCount      a maximum number of messages to peek.
//
// Returns:          receives a list with messages or error.
func (c *NatsMessageQueue) PeekBatch(ctx context.Context, messageCount int64) ([]*cqueues.MessageEnvelope, error) {
	err := c.CheckOpen(cctx.GetTraceId(ctx))
	if err != nil {
		return nil, err
	}

	// Subscribe if needed
	err = c.subscribe(ctx)
	if err != nil {
		return nil, err
	}

	messages := make([]*cqueues.MessageEnvelope, messageCount)

	c.Lock.Lock()
	if messageCount <= (int64)(len(messages)) {
		copy(messages, c.messages[0:messageCount])
	}
	c.Lock.Unlock()

	c.Logger.Trace(ctx, "Peeked %d messages on %s", len(messages), c.Name())

	return messages, nil
}

// Receive method are receives an incoming message and removes it from the queue.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- waitTimeout  time.Duration     a timeout in milliseconds to wait for a message to come.
//
// Returns:  result *cqueues.MessageEnvelope, err error
// receives a message or error.
func (c *NatsMessageQueue) Receive(ctx context.Context, waitTimeout time.Duration) (*cqueues.MessageEnvelope, error) {
	err := c.CheckOpen(cctx.GetTraceId(ctx))
	if err != nil {
		return nil, err
	}

	// Subscribe if needed
	err = c.subscribe(ctx)
	if err != nil {
		return nil, err
	}

	messageReceived := false
	var message *cqueues.MessageEnvelope
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

		// Add messages to locked messages list
		messageReceived = true
		c.Lock.Unlock()
	}

	return message, nil
}

// Function thath process incoming messages
//
//	Parameters:
//		- msg *nats.Msg	message from the NATs
func (c *NatsMessageQueue) OnMessage(msg *nats.Msg) {
	// Deserialize message
	message, err := c.ToMessage(msg)
	if err != nil {
		c.Logger.Error(cctx.NewContextWithTraceId(context.Background(), message.TraceId), err, "Failed to read received message")
	}

	c.Counters.IncrementOne(context.Background(), "queue."+c.Name()+".received_messages")
	c.Logger.Debug(cctx.NewContextWithTraceId(context.Background(), message.TraceId), "Received message %s via %s", msg, c.Name())

	// Send message to receiver if its set or put it into the queue
	c.Lock.Lock()
	if c.receiver != nil {
		receiver := c.receiver
		c.Lock.Unlock()
		c.sendMessageToReceiver(receiver, message)
	} else {
		c.messages = append(c.messages, message)
		c.Lock.Unlock()
	}
}

func (c *NatsMessageQueue) sendMessageToReceiver(receiver cqueues.IMessageReceiver, message *cqueues.MessageEnvelope) {
	ctx := cctx.NewContextWithTraceId(context.Background(), message.TraceId)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%v", r)
			c.Logger.Error(ctx, nil, "Failed to process the message - "+err)
		}
	}()

	err := receiver.ReceiveMessage(context.Background(), message, c)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to process the message")
	}
}

// Listens for incoming messages and blocks the current thread until queue is closed.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- receiver    cqueues.IMessageReceiver      a receiver to receive incoming messages.
//
// See IMessageReceiver
// See receive
func (c *NatsMessageQueue) Listen(ctx context.Context, receiver cqueues.IMessageReceiver) error {
	err := c.CheckOpen(cctx.GetTraceId(ctx))
	if err != nil {
		return err
	}

	// Subscribe if needed
	err = c.subscribe(ctx)
	if err != nil {
		return err
	}

	c.Logger.Trace(ctx, "", "Started listening messages at %s", c.Name())

	// Get all collected messages
	c.Lock.Lock()
	batchMessages := c.messages
	c.messages = []*cqueues.MessageEnvelope{}
	c.Lock.Unlock()

	// Resend collected messages to receiver
	for _, message := range batchMessages {
		receiver.ReceiveMessage(ctx, message, c)
	}

	// Set the receiver
	c.Lock.Lock()
	c.receiver = receiver
	c.Lock.Unlock()

	return nil
}

// EndListen method are ends listening for incoming messages.
// When this method is call listen unblocks the thread and execution continues.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
func (c *NatsMessageQueue) EndListen(ctx context.Context) {
	c.Lock.Lock()
	c.receiver = nil
	c.Lock.Unlock()
}
