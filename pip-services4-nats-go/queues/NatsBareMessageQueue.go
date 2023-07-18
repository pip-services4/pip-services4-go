package queues

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
)

// NatsBareMessageQueue are message queue that sends and receives messages via NATS message broker.
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
//			- retry_connect:        (optional) turns on/off automated reconnect when connection is log (default: true)
//			- max_reconnect:        (optional) maximum reconnection attempts (default: 3)
//			- reconnect_timeout:    (optional) number of milliseconds to wait on each reconnection attempt (default: 3000)
//			- flush_timeout:        (optional) number of milliseconds to wait on flushing messages (default: 3000)
//
//	References:
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
//		Example:
//			ctx := context.Background()
//			queue := NewNatsBareMessageQueue("myqueue")
//			queue.Configure(ctx, cconf.NewConfigParamsFromTuples(
//				"subject", "mytopic",
//				"queue_group", "mygroup",
//				"connection.protocol", "nats"
//				"connection.host", "localhost"
//				"connection.port", 1883
//			))
//
//			_ = queue.Open(ctx)
//
//	   	_ = queue.Send(ctx,  NewMessageEnvelope("", "mymessage", "ABC"))
//
//	   	message, err := queue.Receive(ctx, 10000*time.Milliseconds)
//		  	if (message != nil) {
//		  		...
//		  		queue.Complete(ctx, message);
//		  	}
type NatsBareMessageQueue struct {
	*NatsAbstractMessageQueue
	subscription *nats.Subscription
}

// NewNatsBareMessageQueue are creates a new instance of the message queue.
// Parameters:
//   - name  string (optional) a queue name.
func NewNatsBareMessageQueue(name string) *NatsBareMessageQueue {
	c := NatsBareMessageQueue{}
	c.NatsAbstractMessageQueue = InheritNatsAbstractMessageQueue(&c, name,
		cqueues.NewMessagingCapabilities(false, true, true, false, false, false, false, false, false))
	return &c
}

// Peek method are peeks a single incoming message from the queue without removing it.
// If there are no messages available in the queue it returns nil.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//
// Returns: result *cqueues.MessageEnvelope, err error
// message or error.
func (c *NatsBareMessageQueue) Peek(ctx context.Context) (*cqueues.MessageEnvelope, error) {
	// Not supported
	return nil, nil
}

// PeekBatch method are peeks multiple incoming messages from the queue without removing them.
// If there are no messages available in the queue it returns an empty list.
// Important: This method is not supported by NATS.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- messageCount      a maximum number of messages to peek.
//
// Returns:          callback function that receives a list with messages or error.
func (c *NatsBareMessageQueue) PeekBatch(ctx context.Context, messageCount int64) ([]*cqueues.MessageEnvelope, error) {
	// Not supported
	return []*cqueues.MessageEnvelope{}, nil
}

// Receive method are receives an incoming message and removes it from the queue.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- waitTimeout  time.Duration     a timeout in milliseconds to wait for a message to come.
//
// Returns:  result *cqueues.MessageEnvelope, err error
// receives a message or error.
func (c *NatsBareMessageQueue) Receive(ctx context.Context, waitTimeout time.Duration) (*cqueues.MessageEnvelope, error) {
	err := c.CheckOpen("")
	if err != nil {
		return nil, err
	}

	// Create a temporary subscription
	var subscription *nats.Subscription

	c.Lock.Lock()
	defer c.Lock.Unlock()

	if c.QueueGroup != "" {
		subscription, err = c.Client.QueueSubscribeSync(c.SubscriptionSubject(), c.QueueGroup)
	} else {
		subscription, err = c.Client.SubscribeSync(c.SubscriptionSubject())
	}
	if err != nil {
		return nil, err
	}

	defer subscription.Unsubscribe()

	// Wait for a message
	msg, err := subscription.NextMsg(waitTimeout)
	if err != nil {
		return nil, err
	}

	if msg != nil {
		message, err := c.ToMessage(msg)
		if err != nil {
			return nil, err
		}

		c.Counters.IncrementOne(ctx, "queue."+c.Name()+".received_messages")
		c.Logger.Debug(cctx.NewContextWithTraceId(ctx, message.TraceId), "Received message %s via %s", msg, c.Name())

		// Convert the message and return
		return message, nil
	}

	return nil, nil
}

func (c *NatsBareMessageQueue) receiveMessage(ctx context.Context, receiver cqueues.IMessageReceiver) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		// Deserialize message
		message, err := c.ToMessage(msg)
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to read received message")
		}

		c.Counters.IncrementOne(ctx, "queue."+c.Name()+".received_messages")
		c.Logger.Debug(cctx.NewContextWithTraceId(ctx, message.TraceId), "Received message %s via %s", msg, c.Name())

		// Pass the message to receiver and recover after panic
		func(message *cqueues.MessageEnvelope) {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Sprintf("%v", r)
					c.Logger.Error(cctx.NewContextWithTraceId(ctx, message.TraceId), nil, "Failed to process the message - "+err)
				}
			}()

			err = receiver.ReceiveMessage(ctx, message, c)
			if err != nil {
				c.Logger.Error(cctx.NewContextWithTraceId(ctx, message.TraceId), err, "Failed to process the message")
			}
		}(message)
	}
}

// Listens for incoming messages and blocks the current thread until queue is closed.
//
//	Parameters:
//		- ctx context.Context	 transaction id to trace execution through call chain.
//		- receiver    cqueues.IMessageReceiver      a receiver to receive incoming messages.
//
// See IMessageReceiver
// See receive
func (c *NatsBareMessageQueue) Listen(ctx context.Context, receiver cqueues.IMessageReceiver) error {
	err := c.CheckOpen("")
	if err != nil {
		return err
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()

	if c.QueueGroup != "" {
		c.subscription, err = c.Client.QueueSubscribe(c.SubscriptionSubject(), c.QueueGroup, c.receiveMessage(ctx, receiver))
	} else {
		c.subscription, err = c.Client.Subscribe(c.SubscriptionSubject(), c.receiveMessage(ctx, receiver))
	}

	return err
}

// EndListen method are ends listening for incoming messages.
// When this method is call listen unblocks the thread and execution continues.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
func (c *NatsBareMessageQueue) EndListen(ctx context.Context) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	if c.subscription != nil {
		c.subscription.Unsubscribe()
		c.subscription = nil
	}
}
