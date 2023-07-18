package queues

import (
	"context"
	"time"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	mqcon "github.com/pip-services4/pip-services4-go/pip-services4-rabbitmq-go/connect"
	rabbitmq "github.com/streadway/amqp"
)

// Message queue that sends and receives messages via MQTT message broker.
//
// MQTT is a popular light-weight protocol to communicate IoT devices.
//
// Configuration parameters:
//
//   - topic:                         name of MQTT topic to subscribe
//
//     connection(s):
//
//   - discovery_key:               (optional) a key to retrieve the connection from  IDiscovery
//
//   - host:                        host name or IP address
//
//   - port:                        port number
//
//   - uri:                         resource URI or connection string with all parameters in it
//
//     credential(s):
//
//   - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//
//   - username:                    user name
//
//   - password:                    user password
//
// References:
//
// - *:logger:*:*:1.0             (optional) ILogger components to pass log messages
// - *:counters:*:*:1.0           (optional) ICounters components to pass collected measurements
// - *:discovery:*:*:1.0          (optional) IDiscovery services to resolve connections
// - *:credential-store:*:*:1.0   (optional) Credential stores to resolve credentials
//
// Example:
//
//	ctx := context.Background()
//	queue := queues.NewEmptyRabbitMQMessageQueue("myqueue")
//	queue.Configure(ctx, config.NewConfigParamsFromTuples(
//		"exchange", "my_exchange",
//		"queue", "my_exchange",
//		"options.auto_create", true,
//		"connection.host", "5672",
//		"connection.port", "localhost",
//		"credential.username", "user",
//		"credential.password", "password",
//	))
//	_ = queue.Open(ctx)
//	_ = queue.Send(ctx, cqueues.NewMessageEnvelope("", "mymessage", []byte("ABC")))
//	message, _ := queue.Receive(ctx, 10000*time.Millisecond)
//	if message != nil {
//		// ...
//		queue.Complete(ctx, message)
//	}
type RabbitMQMessageQueue struct {
	*cqueues.MessageQueue
	defaultCheckInterval int64
	connection           *rabbitmq.Connection
	mqChanel             *rabbitmq.Channel
	optionsResolver      *mqcon.RabbitMQConnectionResolver
	queue                string
	exchange             string
	exchangeType         string
	routingKey           string
	persistent           bool
	exclusive            bool
	autoCreate           bool
	autoDelete           bool
	noQueue              bool
	Interval             time.Duration

	cancelListen chan bool
	isListen     bool
}

//	 Creates a new instance of the message queue.
//		Parameters:
//	 	- name(optional) a queue name.
func NewEmptyRabbitMQMessageQueue(name string) *RabbitMQMessageQueue {
	c := RabbitMQMessageQueue{
		defaultCheckInterval: 1000,
		exchange:             "",
		exchangeType:         "fanout",
		routingKey:           "",
		persistent:           false,
		exclusive:            false,
		autoCreate:           false,
		autoDelete:           false,
		noQueue:              false,
		cancelListen:         make(chan bool),
	}

	c.MessageQueue = cqueues.InheritMessageQueue(
		&c, name,
		cqueues.NewMessagingCapabilities(true, true, true, true, true, false, true, false, true))
	c.Interval = time.Duration(c.defaultCheckInterval) * time.Millisecond
	c.optionsResolver = mqcon.NewRabbitMQConnectionResolver()

	return &c
}

func NewRabbitMQMessageQueueFromConfig(name string, config *cconf.ConfigParams) *RabbitMQMessageQueue {
	c := NewEmptyRabbitMQMessageQueue(name)
	if config != nil {
		c.Configure(context.Background(), config)
	}
	return c
}

func NewRabbitMQMessageQueue(name string, mqChanel *rabbitmq.Channel, queue string) *RabbitMQMessageQueue {
	c := NewEmptyRabbitMQMessageQueue(name)
	c.mqChanel = mqChanel
	c.queue = queue
	return c
}

// Configures component by passing configuration parameters.
// Parameters:
//   - ctx context.Context
//
// - config configuration parameters to be set.
func (c *RabbitMQMessageQueue) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.MessageQueue.Configure(ctx, config)

	c.Interval = time.Duration(config.GetAsLongWithDefault("interval", int64(c.defaultCheckInterval))) * time.Millisecond

	c.queue = config.GetAsStringWithDefault("queue", c.queue)
	c.exchange = config.GetAsStringWithDefault("exchange", c.exchange)

	c.exchangeType = config.GetAsStringWithDefault("options.exchange_type", c.exchangeType)
	c.routingKey = config.GetAsStringWithDefault("options.routing_key", c.routingKey)
	c.persistent = config.GetAsBooleanWithDefault("options.persistent", c.persistent)
	c.exclusive = config.GetAsBooleanWithDefault("options.exclusive", c.exclusive)
	c.autoCreate = config.GetAsBooleanWithDefault("options.auto_create", c.autoCreate)
	c.autoDelete = config.GetAsBooleanWithDefault("options.auto_delete", c.autoDelete)
	c.noQueue = config.GetAsBooleanWithDefault("options.noqueue", c.noQueue)
}

func (c *RabbitMQMessageQueue) checkOpened(traceId string) error {
	if c.mqChanel == nil {
		return cerr.NewInvalidStateError(traceId, "NOT_OPENED", "The queue is not opened")
	}
	return nil
}

// Checks if the component is opened.
// Retruns : true if the component has been opened and false otherwise.
func (c *RabbitMQMessageQueue) IsOpen() bool {
	return c.connection != nil && c.mqChanel != nil
}

// Opens the component with given connection and credential parameters.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - connections connection parameters
//   - credential credential parameters
func (c *RabbitMQMessageQueue) Open(ctx context.Context) error {

	connection, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		return err
	}

	credential, err := c.CredentialResolver.Lookup(ctx)
	if err != nil {
		return err
	}

	options, err := c.optionsResolver.Compose(ctx, connection, credential)
	if err != nil {
		return err
	}

	if c.queue == "" && c.exchange == "" {
		return cerr.NewConfigError(cctx.GetTraceId(ctx),
			"NO_QUEUE",
			"Queue or exchange are not defined in connection parameters")
	}

	conn, err := rabbitmq.Dial(options.GetAsString("uri"))
	if err != nil {
		return err
	}
	c.connection = conn
	c.mqChanel, err = conn.Channel()
	if err != nil {
		return err
	}

	// Automatically create queue, exchange and binding
	if c.autoCreate {
		if c.exchange != "" {
			c.mqChanel.ExchangeDeclare(
				c.exchange,
				c.exchangeType,
				c.persistent,
				c.autoDelete,
				false,
				false,
				nil,
			)
		}

		if !c.noQueue {

			if c.queue == "" {
				res, err := c.mqChanel.QueueDeclare(
					"",
					c.persistent,
					true,
					true,
					false,
					nil,
				)
				if err != nil {
					return err
				}
				c.queue = res.Name
			} else {
				c.mqChanel.QueueDeclare(
					c.queue,
					c.persistent,
					c.exclusive,
					c.autoDelete,
					false,
					nil,
				)
			}

			c.mqChanel.QueueBind(
				c.queue,
				c.routingKey,
				c.exchange,
				false,
				nil,
			)

		}
	}
	return nil
}

// Close mwthod are closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context transaction id to trace execution through call chain.
func (c *RabbitMQMessageQueue) Close(ctx context.Context) (err error) {

	c.Lock.Lock()
	if c.isListen {
		c.isListen = false
		c.Lock.Unlock()
		close(c.cancelListen)
	} else {
		c.Lock.Unlock()
	}

	if c.mqChanel != nil {
		err = c.mqChanel.Close()
		if err != nil {
			return err
		}
	}

	if c.connection != nil {
		err = c.connection.Close()
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.connection = nil
	c.mqChanel = nil

	c.Logger.Trace(ctx, "Closed queue %s", c.queue)
	return err
}

// ReadMessageCount method are reads the current number of messages in the queue to be delivered.
// Returns count int64, err error
// number of messages or error.
func (c *RabbitMQMessageQueue) ReadMessageCount() (count int64, err error) {

	err = c.checkOpened("")
	if err != nil {
		c.Logger.Error(context.Background(), err, "RabbitMQMessageQueue:MessageCount: "+err.Error())
		return 0, err
	}

	if c.queue == "" {
		return 0, nil
	}
	queueInfo, err := c.mqChanel.QueueInspect(c.queue)
	if err != nil {
		c.Logger.Error(context.Background(), err, "RabbitMQMessageQueue:MessageCount: "+err.Error())
		return 0, err
	}
	return int64(queueInfo.Messages), nil

}

func (c *RabbitMQMessageQueue) toMessage(envelope *rabbitmq.Delivery) *cqueues.MessageEnvelope {
	if envelope == nil {
		return nil
	}

	message := cqueues.MessageEnvelope{
		MessageId:   envelope.MessageId,
		MessageType: envelope.Type,
		TraceId:     envelope.CorrelationId,
		Message:     envelope.Body,
		SentTime:    time.Now(),
	}
	message.SetReference(envelope)

	return &message
}

//	 Send method are sends a message into the queue.
//		Parameters:
//			- ctx context.Context transaction id to trace execution through call chain.
//			- message a message envelop to be sent.
func (c *RabbitMQMessageQueue) Send(ctx context.Context, message *cqueues.MessageEnvelope) (err error) {
	err = c.checkOpened(cctx.GetTraceId(ctx))
	if err != nil {
		return err
	}

	messageBuffer := rabbitmq.Publishing{
		ContentType: "text/plain",
	}

	if message.TraceId != "" {
		messageBuffer.CorrelationId = message.TraceId
	}
	if message.MessageId != "" {
		messageBuffer.MessageId = message.MessageId
	}

	if message.MessageType != "" {
		messageBuffer.Type = message.MessageType
	}

	messageBuffer.Body = []byte(message.Message)

	err = c.mqChanel.Publish(c.exchange, c.routingKey, false, false, messageBuffer)

	c.Counters.IncrementOne(ctx, "queue."+c.Name()+".sent_messages")
	c.Logger.Debug(cctx.NewContextWithTraceId(ctx, message.TraceId), "Sent message %s via %s", message, c.Name())
	return err
}

// Peeks a single incoming message from the queue without removing it.
// If there are no messages available in the queue it returns nil.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//
// Returns: a message
func (c *RabbitMQMessageQueue) Peek(ctx context.Context) (result *cqueues.MessageEnvelope, err error) {
	err = c.checkOpened(cctx.GetTraceId(ctx))
	if err != nil {
		return nil, err
	}

	envelope, ok, err := c.mqChanel.Get(c.queue, false)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	message := c.toMessage(&envelope)
	if message != nil {
		c.Logger.Trace(cctx.NewContextWithTraceId(ctx, message.TraceId), "Peeked message %s on %s", message, c.Name())
	}

	return message, nil
}

// PeekBatch method are peeks multiple incoming messages from the queue without removing them.
// If there are no messages available in the queue it returns an empty list.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
//   - messageCount a maximum number of messages to peek.
//
// Returns: a list with messages
func (c *RabbitMQMessageQueue) PeekBatch(ctx context.Context, messageCount int64) (result []*cqueues.MessageEnvelope, err error) {
	err = c.checkOpened(cctx.GetTraceId(ctx))
	if err != nil {
		return nil, err
	}
	err = nil
	messages := make([]*cqueues.MessageEnvelope, 0)
	for messageCount > 0 {
		envelope, ok, getErr := c.mqChanel.Get(c.queue, false)
		if getErr != nil || !ok {
			err = getErr
			break
		}
		message := c.toMessage(&envelope)
		messages = append(messages, message)
		messageCount--
	}
	c.Logger.Trace(ctx, "Peeked %s messages on %s", len(messages), c.Name())
	return messages, err
}

//	 Receive method are receives an incoming message and removes it from the queue.
//		Parameters:
//			- ctx context.Context  transaction id to trace execution through call chain.
//			- waitTimeout a timeout in milliseconds to wait for a message to come.
//		Returns: a message
func (c *RabbitMQMessageQueue) Receive(ctx context.Context, waitTimeout time.Duration) (result *cqueues.MessageEnvelope, err error) {

	err = c.checkOpened(cctx.GetTraceId(ctx))
	if err != nil {
		return nil, err
	}
	err = nil

	var message *cqueues.MessageEnvelope
	timeout := waitTimeout

	for {
		if timeout <= 0 {
			break
		}
		// Read the message and exit if received
		env, ok, getErr := c.mqChanel.Get(c.queue, false) // true
		if ok && getErr == nil {
			c.Lock.Lock()
			defer c.Lock.Unlock()
			message = c.toMessage(&env)
			break
		}
		timeout = timeout - c.Interval
	}

	if message != nil {
		c.Counters.IncrementOne(ctx, "queue."+c.Name()+".received_messages")
		c.Logger.Debug(cctx.NewContextWithTraceId(ctx, message.TraceId), "Received message %s via %s", message, c.Name())
	}

	return message, nil
}

// Renews a lock on a message that makes it invisible from other receivers in the queue.
// This method is usually used to extend the message processing time.
// Important: This method is not supported by RabbitMQ.
// Parameters:
//   - ctx context.Context
//   - message a message to extend its lock.
//   - lockTimeout a locking timeout in milliseconds.
func (c *RabbitMQMessageQueue) RenewLock(ctx context.Context, message *cqueues.MessageEnvelope, lockTimeout time.Duration) (err error) {
	// Operation is not supported
	return nil
}

// Returnes message into the queue and makes it available for all subscribers to receive it again.
// This method is usually used to return a message which could not be processed at the moment
// to repeat the attempt.Messages that cause unrecoverable errors shall be removed permanently
// or/and send to dead letter queue.
// Important: This method is not supported by RabbitMQ.
// Parameters:
//   - ctx context.Context
//   - message a message to return.
func (c *RabbitMQMessageQueue) Abandon(ctx context.Context, message *cqueues.MessageEnvelope) (err error) {
	err = c.checkOpened("")
	if err != nil {
		return err
	}
	err = nil

	// Make the message immediately visible
	envelope, ok := message.GetReference().(*rabbitmq.Delivery)
	if ok {
		err = c.mqChanel.Nack(envelope.DeliveryTag, false, true)
		if err != nil {
			return err
		}
		message.SetReference(nil)
		c.Logger.Trace(cctx.NewContextWithTraceId(ctx, message.TraceId), "Abandoned message %s at %c", message, c.Name())
	}
	return nil
}

// Permanently removes a message from the queue.
// This method is usually used to remove the message after successful processing.
// Important: This method is not supported by RabbitMQ.
// Parameters:
//   - ctx context.Context
//   - message a message to remove.
func (c *RabbitMQMessageQueue) Complete(ctx context.Context, message *cqueues.MessageEnvelope) (err error) {
	err = c.checkOpened("")
	if err != nil {
		return err
	}
	err = nil
	envelope, ok := message.GetReference().(*rabbitmq.Delivery)
	if ok {
		c.mqChanel.Ack(envelope.DeliveryTag, false)
		message.SetReference(nil)
		c.Logger.Trace(cctx.NewContextWithTraceId(ctx, message.TraceId), "Completed message %s at %s", message, c.Name())
	}
	return nil
}

// Permanently removes a message from the queue and sends it to dead letter queue.
// Important: This method is not supported by RabbitMQ.
// Parameters:
//   - ctx context.Context
//   - message a message to be removed.
//
// Returns: error
func (c *RabbitMQMessageQueue) MoveToDeadLetter(ctx context.Context, message *cqueues.MessageEnvelope) (err error) {
	// Operation is not supported
	return nil
}

//		Listens for incoming messages and blocks the current thread until queue is closed.
//		Parameters:
//			- ctx context.Context transaction id to trace execution through call chain.
//			- receiver IMessageReceiver callback interface to receive incoming messages.
//	 Returns: listen error
func (c *RabbitMQMessageQueue) Listen(ctx context.Context, receiver cqueues.IMessageReceiver) error {
	err := c.checkOpened("")
	if err != nil {
		c.Logger.Error(ctx, err, "RabbitMQMessageQueue:Listen: Can't start listen "+err.Error())
		return nil
	}

	c.Logger.Debug(ctx, "Started listening messages at %s", c.Name())

	messageChannel, err := c.mqChanel.Consume(
		c.queue,
		c.exchange,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		c.Logger.Error(ctx, err, "RabbitMQMessageQueue:Listen: Can't consume to queue"+err.Error())
		return nil
	}

	c.EndListen(ctx)
	c.cancelListen = make(chan bool)

	go func() {
		c.Lock.Lock()
		c.isListen = true
		c.Lock.Unlock()

		for {
			select {
			case msg := <-messageChannel:
				{
					message := c.toMessage(&msg)
					c.Counters.IncrementOne(ctx, "queue."+c.Name()+".received_messages")
					c.Logger.Debug(cctx.NewContextWithTraceId(ctx, message.TraceId), "Received message %s via %s", message, c.Name())
					recvErr := receiver.ReceiveMessage(ctx, message, c)
					if recvErr != nil {
						c.Logger.Error(cctx.NewContextWithTraceId(ctx, message.TraceId), recvErr, "Processing received message %s error in queue %s", message, c.Name())
					}
					c.mqChanel.Ack(msg.DeliveryTag, false)
				}
			case <-c.cancelListen:
				{
					return
				}
			case <-ctx.Done():
				{
					return
				}
			}
		}
	}()

	return nil
}

// Ends listening for incoming messages.
// When this method is call listen unblocks the thread and execution continues.
// Parameters:
//   - ctx context.Context transaction id to trace execution through call chain.
func (c *RabbitMQMessageQueue) EndListen(ctx context.Context) {
	c.Lock.Lock()
	if c.isListen {
		c.isListen = false
		c.Lock.Unlock()
		close(c.cancelListen)
	} else {
		c.Lock.Unlock()
	}
}

//		Clear method are clears component state.
//		Parameters:
//			- ctx context.Context transaction id to trace execution through call chain.
//	 Returns: error
func (c *RabbitMQMessageQueue) Clear(ctx context.Context) (err error) {
	err = c.checkOpened("")
	if err != nil {
		return err
	}

	count := 0
	if c.queue != "" {
		count, err = c.mqChanel.QueuePurge(c.queue, false)
	}

	if err != nil {
		return err
	}

	c.Logger.Trace(ctx, "Cleared  %s messages in queue %s", count, c.Name())

	return nil
}
