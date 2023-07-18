package queues

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	cconn "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-nats-go/connect"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// Abstract NATS message queue with ability to connect to NATS server.
type NatsAbstractMessageQueue struct {
	*cqueues.MessageQueue

	defaultConfig   *cconf.ConfigParams
	config          *cconf.ConfigParams
	references      cref.IReferences
	opened          bool
	localConnection bool

	//The dependency resolver.
	DependencyResolver *cref.DependencyResolver
	//The logger.
	Logger *clog.CompositeLogger
	//The NATS connection component.
	Connection *connect.NatsConnection
	//The NATS connection object.
	Client *nats.Conn

	// SerializeEnvelop bool
	Subject    string
	QueueGroup string
}

// Creates a new instance of the queue component.
//   - overrides a queue overrides
//   - name    (optional) a queue name.
func InheritNatsAbstractMessageQueue(overrides cqueues.IMessageQueueOverrides, name string, capabilities *cqueues.MessagingCapabilities) *NatsAbstractMessageQueue {
	c := &NatsAbstractMessageQueue{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"subject", nil,
			"queue_group", nil,
			// "options.serialize_envelop", true,
			"options.retry_connect", true,
			"options.connect_timeout", 0,
			"options.reconnect_timeout", 3000,
			"options.max_reconnect", 3,
			"options.flush_timeout", 3000,
		),
		Logger: clog.NewCompositeLogger(),
	}
	c.MessageQueue = cqueues.InheritMessageQueue(overrides, name, capabilities)
	c.DependencyResolver = cref.NewDependencyResolver()
	c.DependencyResolver.Configure(context.Background(), c.defaultConfig)
	return c
}

// Configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- config    configuration parameters to be set.
func (c *NatsAbstractMessageQueue) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config

	c.DependencyResolver.Configure(ctx, config)

	// c.SerializeEnvelop = config.GetAsBooleanWithDefault("options.serialize_envelop", c.SerializeEnvelop)
	c.Subject = config.GetAsStringWithDefault("topic", c.Subject)
	c.Subject = config.GetAsStringWithDefault("subject", c.Subject)
	c.QueueGroup = config.GetAsStringWithDefault("group", c.QueueGroup)
	c.QueueGroup = config.GetAsStringWithDefault("queue_group", c.QueueGroup)
}

// Sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- references 	references to locate the component dependencies.
func (c *NatsAbstractMessageQueue) SetReferences(ctx context.Context, references cref.IReferences) {
	c.references = references
	c.Logger.SetReferences(ctx, references)

	// Get connection
	c.DependencyResolver.SetReferences(ctx, references)
	result := c.DependencyResolver.GetOneOptional("connection")
	if dep, ok := result.(*connect.NatsConnection); ok {
		c.Connection = dep
	}
	// Or create a local one
	if c.Connection == nil {
		c.Connection = c.createConnection()
		c.localConnection = true
	} else {
		c.localConnection = false
	}
}

// Unsets (clears) previously set references to dependent components.
func (c *NatsAbstractMessageQueue) UnsetReferences() {
	c.Connection = nil
}

func (c *NatsAbstractMessageQueue) createConnection() *connect.NatsConnection {
	connection := connect.NewNatsConnection()
	if c.config != nil {
		connection.Configure(context.Background(), c.config)
	}
	if c.references != nil {
		connection.SetReferences(context.Background(), c.references)
	}
	return connection
}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *NatsAbstractMessageQueue) IsOpen() bool {
	return c.opened
}

// Opens the component.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- Returns 			 error or nil no errors occured.
func (c *NatsAbstractMessageQueue) Open(ctx context.Context) (err error) {
	if c.opened {
		return nil
	}

	if c.Connection == nil {
		c.Connection = c.createConnection()
		c.localConnection = true
	}

	if c.localConnection {
		err = c.Connection.Open(ctx)
	}

	if err == nil && c.Connection == nil {
		err = cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "NATS connection is missing")
	}

	if err == nil && !c.Connection.IsOpen() {
		err = cerr.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "NATS connection is not opened")
	}

	c.opened = true

	if err != nil {
		return err
	}
	c.Client = c.Connection.GetConnection()

	return err
}

// OpenWithParams method are opens the component with given connection and credential parameters.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- connections        connection parameters
//		- credential        credential parameters
//
// Returns error or nil no errors occured.
func (c *NatsAbstractMessageQueue) OpenWithParams(ctx context.Context, connections []*cconn.ConnectionParams,
	credential *cauth.CredentialParams) error {
	panic("Not supported")
}

// Closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- Returns 			error or nil no errors occured.
func (c *NatsAbstractMessageQueue) Close(ctx context.Context) (err error) {
	if !c.opened {
		return nil
	}

	if c.Connection == nil {
		return cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "NATS connection is missing")
	}

	if c.localConnection {
		err = c.Connection.Close(ctx)
	}
	if err != nil {
		return err
	}

	// Todo: Flush messages?
	c.opened = false
	c.Client = nil

	return nil
}

func (c *NatsAbstractMessageQueue) CheckOpen(traceId string) error {
	if !c.IsOpen() {
		err := cerr.NewInvalidStateError(
			traceId,
			"NOT_OPENED",
			"The queue is not opened",
		)
		return err
	}
	return nil
}

func (c *NatsAbstractMessageQueue) SubscriptionSubject() string {
	if c.Subject != "" {
		return c.Subject
	}
	return c.Name()
}

// Converts MessageEnvelope to NATs message structure
//
//	Parameters:
//		- message *cqueues.MessageEnvelope message object
//
// Returns: NATs message structure
func (c *NatsAbstractMessageQueue) FromMessage(message *cqueues.MessageEnvelope) (*nats.Msg, error) {
	if message == nil {
		return nil, nil
	}

	// if c.SerializeEnvelop {
	// 	data, err := json.Marshal(message)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	msg := nats.NewMsg(c.Name())
	// 	msg.Data = data
	// 	return msg, nil
	// } else {
	// 	msg := nats.NewMsg(c.Name())
	// 	msg.Data = message.Message
	// 	return msg, nil
	// }

	msg := nats.NewMsg(c.Name())
	msg.Data = message.Message
	msg.Header.Add("message_id", message.MessageId)
	msg.Header.Add("trace_id", message.TraceId)
	msg.Header.Add("message_type", message.MessageType)
	msg.Header.Add("sent_time", cconv.StringConverter.ToString(message.SentTime))
	return msg, nil
}

// Converts NATs structure to MessageEnvelope
//
//	Parameters:
//		- msg *nats.Msg message object
//
// Returns: MessageEnvelope structure
func (c *NatsAbstractMessageQueue) ToMessage(msg *nats.Msg) (*cqueues.MessageEnvelope, error) {
	if msg == nil {
		return nil, nil
	}

	// if c.SerializeEnvelop {
	// 	envelop := cqueues.MessageEnvelope{}
	// 	err := json.Unmarshal(msg.Data, &envelop)
	// 	return &envelop, err
	// } else {
	// 	envelop := cqueues.NewMessageEnvelope("", "", msg.Data)
	// 	return envelop, nil
	// }

	message := cqueues.NewEmptyMessageEnvelope()
	message.MessageId = msg.Header.Get("message_id")
	message.TraceId = msg.Header.Get("trace_id")
	message.MessageType = msg.Header.Get("message_type")
	message.SentTime = cconv.DateTimeConverter.ToDateTime(msg.Header.Get("sent_time"))
	message.Message = msg.Data
	message.SetReference(msg)

	return message, nil
}

// Clear method are clears component state.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//
// Returns error or nil no errors occured.
func (c *NatsAbstractMessageQueue) Clear(ctx context.Context) error {
	// Not supported
	return nil
}

// ReadMessageCount method are reads the current number of messages in the queue to be delivered.
// Returns number of messages or error.
func (c *NatsAbstractMessageQueue) ReadMessageCount() (int64, error) {
	// Not supported
	return 0, nil
}

// Send method are sends a message into the queue.
//
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//		- envelope *cqueues.MessageEnvelope  a message envelop to be sent.
//
// Returns: error or nil for success.
func (c *NatsAbstractMessageQueue) Send(ctx context.Context, envelop *cqueues.MessageEnvelope) error {
	err := c.CheckOpen(cctx.GetTraceId(ctx))
	if err != nil {
		return err
	}

	msg, err := c.FromMessage(envelop)
	if err != nil {
		return err
	}

	subject := c.Name()
	if subject == "" {
		subject = c.Subject
	}

	err = c.Connection.Publish(ctx, subject, msg)
	if err != nil {
		c.Logger.Error(cctx.NewContextWithTraceId(ctx, envelop.TraceId), err, "Failed to send message via %s", c.Name())
		return err
	}

	c.Counters.IncrementOne(ctx, "queue."+c.Name()+".sent_messages")
	c.Logger.Debug(cctx.NewContextWithTraceId(ctx, envelop.TraceId), "Sent message %s via %s", envelop.String(), c.Name())

	return nil
}

// RenewLock method are renews a lock on a message that makes it invisible from other receivers in the queue.
// This method is usually used to extend the message processing time.
// Important: This method is not supported by NATS.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- message   *cqueues.MessageEnvelope    a message to extend its lock.
//		- lockTimeout  time.Duration  a locking timeout in milliseconds.
//
// Returns: error
// receives an error or nil for success.
func (c *NatsAbstractMessageQueue) RenewLock(ctx context.Context, message *cqueues.MessageEnvelope, lockTimeout time.Duration) (err error) {
	// Not supported
	return nil
}

// Complete method are permanently removes a message from the queue.
// This method is usually used to remove the message after successful processing.
// Important: This method is not supported by NATS.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- message  *cqueues.MessageEnvelope a message to remove.
//
// Returns: error
// error or nil for success.
func (c *NatsAbstractMessageQueue) Complete(ctx context.Context, message *cqueues.MessageEnvelope) (err error) {
	// Not supported
	return nil
}

// Abandon method are returnes message into the queue and makes it available for all subscribers to receive it again.
// This method is usually used to return a message which could not be processed at the moment
// to repeat the attempt. Messages that cause unrecoverable errors shall be removed permanently
// or/and send to dead letter queue.
// Important: This method is not supported by NATS.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- message *cqueues.MessageEnvelope  a message to return.
//
// Returns: error
//
//	error or nil for success.
func (c *NatsAbstractMessageQueue) Abandon(ctx context.Context, message *cqueues.MessageEnvelope) (err error) {
	// Not supported
	return nil
}

// Permanently removes a message from the queue and sends it to dead letter queue.
// Important: This method is not supported by NATS.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- message  *cqueues.MessageEnvelope a message to be removed.
//
// Returns: error
//
//	error or nil for success.
func (c *NatsAbstractMessageQueue) MoveToDeadLetter(ctx context.Context, message *cqueues.MessageEnvelope) (err error) {
	// Not supported
	return nil
}
