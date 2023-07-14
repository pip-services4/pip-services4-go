package queues

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-mqtt-go/connect"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// MqttMessageQueue are message queue that sends and receives messages via MQTT message broker.
//
// Configuration parameters:
//
//   - topic:                         name of MQTT topic to subscribe
//   - connection(s):
//   - discovery_key:               (optional) a key to retrieve the connection from  IDiscovery
//   - host:                        host name or IP address
//   - port:                        port number
//   - uri:                         resource URI or connection string with all parameters in it
//   - credential(s):
//   - store_key:                   (optional) a key to retrieve the credentials from  ICredentialStore
//   - username:                    user name
//   - password:                    user password
//   - options:
//   - serialize_envelope:    (optional) true to serialize entire message as JSON, false to send only message payload (default: true)
//   - autosubscribe:        (optional) true to automatically subscribe on option (default: false)
//   - qos:                  (optional) quality of service level aka QOS (default: 0)
//   - retain:               (optional) retention flag for published messages (default: false)
//   - retry_connect:        (optional) turns on/off automated reconnect when connection is log (default: true)
//   - connect_timeout:      (optional) number of milliseconds to wait for connection (default: 30000)
//   - reconnect_timeout:    (optional) number of milliseconds to wait on each reconnection attempt (default: 1000)
//   - keepalive_timeout:    (optional) number of milliseconds to ping broker while inactive (default: 3000)
//
// References:
//   - *:logger:*:*:1.0             (optional)  ILogger components to pass log messages
//   - *:counters:*:*:1.0           (optional)  ICounters components to pass collected measurements
//   - *:discovery:*:*:1.0          (optional)  IDiscovery services to resolve connections
//   - *:credential-store:*:*:1.0   (optional) Credential stores to resolve credentials
//   - *:connection:mqtt:*:1.0      (optional) Shared connection to MQTT service
//
// See MessageQueue
// See MessagingCapabilities
//
// Example:
//
//	queue := NewMqttMessageQueue("myqueue")
//	queue.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
//		"subject", "mytopic",
//		"connection.protocol", "mqtt"
//		"connection.host", "localhost"
//		"connection.port", 1883
//	))
//
//	queue.open(context.Background(), "123")
//
//	queue.Send(context.Background(), "123", NewMessageEnvelope("", "mymessage", "ABC"))
//
//	message, err := queue.Receive(context.Background(), "123")
//	if (message != nil) {
//		...
//		queue.Complete(context.Background(), "123", message);
//	}
type MqttMessageQueue struct {
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
	//The MQTT connection component.
	Connection *connect.MqttConnection

	serializeEnvelope bool
	topic             string
	qos               byte
	retain            bool
	autoSubscribe     bool
	subscribed        bool
	messages          []cqueues.MessageEnvelope
	receiver          cqueues.IMessageReceiver
}

// Creates a new instance of the queue component.
//   - name    (optional) a queue name.
func NewMqttMessageQueue(name string) *MqttMessageQueue {
	c := MqttMessageQueue{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"topic", nil,
			"options.autosubscribe", false,
			"options.serialize_envelope", false,
			"options.retry_connect", true,
			"options.connect_timeout", 30000,
			"options.reconnect_timeout", 1000,
			"options.keepalive_timeout", 60000,
			"options.qos", 0,
			"options.retain", false,
		),
		Logger: clog.NewCompositeLogger(),
	}
	c.MessageQueue = cqueues.InheritMessageQueue(&c, name,
		cqueues.NewMessagingCapabilities(false, true, true, true, true, false, false, false, true))
	c.DependencyResolver = cref.NewDependencyResolver()
	c.DependencyResolver.Put(context.Background(), "connection", cref.NewDescriptor("pip-services", "connection", "mqtt", "*", "1.0"))
	c.DependencyResolver.Configure(context.Background(), c.defaultConfig)

	c.messages = make([]cqueues.MessageEnvelope, 0)

	return &c
}

// Configures component by passing configuration parameters.
// Parameters:
//   - ctx context.Context	operation context.
//   - config    configuration parameters to be set.
func (c *MqttMessageQueue) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config

	c.DependencyResolver.Configure(ctx, config)

	c.topic = config.GetAsStringWithDefault("topic", c.topic)
	c.serializeEnvelope = config.GetAsBooleanWithDefault("options.serialize_envelope", c.serializeEnvelope)
	c.autoSubscribe = config.GetAsBooleanWithDefault("options.autosubscribe", c.autoSubscribe)
	c.qos = byte(config.GetAsIntegerWithDefault("options.qos", int(c.qos)))
	c.retain = config.GetAsBooleanWithDefault("options.retain", c.retain)
}

// Sets references to dependent components.
// Parameters:
//   - ctx context.Context	operation context.
//   - references 	references to locate the component dependencies.
func (c *MqttMessageQueue) SetReferences(ctx context.Context, references cref.IReferences) {
	c.references = references
	c.Logger.SetReferences(ctx, references)

	// Get connection
	c.DependencyResolver.SetReferences(ctx, references)
	result := c.DependencyResolver.GetOneOptional("connection")
	if dep, ok := result.(*connect.MqttConnection); ok {
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
func (c *MqttMessageQueue) UnsetReferences() {
	c.Connection = nil
}

func (c *MqttMessageQueue) createConnection() *connect.MqttConnection {
	connection := connect.NewMqttConnection()
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
func (c *MqttMessageQueue) IsOpen() bool {
	return c.opened
}

// Opens the component.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - Returns 			 error or nil no errors occured.
func (c *MqttMessageQueue) Open(ctx context.Context) (err error) {
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
		err = cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "MQTT connection is missing")
	}

	if err == nil && !c.Connection.IsOpen() {
		err = cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "MQTT connection is not opened")
	}

	if err != nil {
		return err
	}

	// Automatically subscribe if needed
	if c.autoSubscribe {
		err = c.subscribe(ctx)
		if err != nil {
			return err
		}
	}

	c.opened = true

	return err
}

// Closes component and frees used resources.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - Returns 			error or nil no errors occured.
func (c *MqttMessageQueue) Close(ctx context.Context) (err error) {
	if !c.opened {
		return nil
	}

	if c.Connection == nil {
		return cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "MQTT connection is missing")
	}

	if c.localConnection {
		err = c.Connection.Close(ctx)
	}
	if err != nil {
		return err
	}

	// Unsubscribe from topic
	if c.subscribed {
		topic := c.getTopic()
		c.Connection.Unsubscribe(ctx, topic, c)
		c.subscribed = false
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.opened = false
	c.receiver = nil
	c.messages = make([]cqueues.MessageEnvelope, 0)

	return nil
}

func (c *MqttMessageQueue) getTopic() string {
	if c.topic != "" {
		return c.topic
	}
	return c.Name()
}

func (c *MqttMessageQueue) subscribe(ctx context.Context) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// Check if already were subscribed
	if c.subscribed {
		return nil
	}

	// Subscribe to the topic
	topic := c.getTopic()
	err := c.Connection.Subscribe(ctx, topic, c.qos, c)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to subscribe to topic "+topic)
		return err
	}

	c.subscribed = true
	return nil
}

func (c *MqttMessageQueue) fromMessage(message *cqueues.MessageEnvelope) ([]byte, error) {
	if message == nil {
		return nil, nil
	}

	data := message.Message
	if c.serializeEnvelope {
		message.SentTime = time.Now()
		var err error
		data, err = json.Marshal(message)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (c *MqttMessageQueue) toMessage(msg mqtt.Message) (*cqueues.MessageEnvelope, error) {
	message := cqueues.NewEmptyMessageEnvelope()

	if c.serializeEnvelope {
		err := json.Unmarshal(msg.Payload(), message)
		if err != nil {
			return nil, err
		}
	} else {
		message.MessageId = strconv.FormatUint(uint64(msg.MessageID()), 10)
		message.MessageType = msg.Topic()
		message.Message = msg.Payload()
	}
	message.SetReference(msg)

	return message, nil
}

func (c *MqttMessageQueue) OnMessage(msg mqtt.Message) {
	// Skip if it came from a wrong topic
	expectedTopic := c.getTopic()
	if !strings.Contains(expectedTopic, "*") && expectedTopic != msg.Topic() {
		return
	}

	// Deserialize message
	message, err := c.toMessage(msg)
	if message == nil || err != nil {
		c.Logger.Error(context.Background(), err, "Failed to read received message")
		return
	}

	c.Counters.IncrementOne(context.Background(), "queue."+c.Name()+".received_messages")
	c.Logger.Debug(context.Background(), message.TraceId, "Received message %s via %s", message, c.Name())

	// Send message to receiver if its set or put it into the queue
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if c.receiver != nil {
		receiver := c.receiver
		c.sendMessageToReceiver(context.Background(), receiver, message)
	} else {
		c.messages = append(c.messages, *message)
	}
}

// Clear method are clears component state.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//
// Returns error or nil no errors occured.
func (c *MqttMessageQueue) Clear(ctx context.Context) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.messages = make([]cqueues.MessageEnvelope, 0)

	return nil
}

// ReadMessageCount method are reads the current number of messages in the queue to be delivered.
// Returns number of messages or error.
func (c *MqttMessageQueue) ReadMessageCount() (int64, error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	count := (int64)(len(c.messages))
	return count, nil
}

// Peek method are peeks a single incoming message from the queue without removing it.
// If there are no messages available in the queue it returns nil.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//
// Returns: result *cqueues.MessageEnvelope, err error
// message or error.
func (c *MqttMessageQueue) Peek(ctx context.Context) (*cqueues.MessageEnvelope, error) {
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
		message = &c.messages[0]
	}
	c.Lock.Unlock()

	if message != nil {
		c.Logger.Trace(cctx.NewContextWithTraceId(ctx, message.TraceId), "Peeked message %s on %s", message, c.String())
	}

	return message, nil
}

// PeekBatch method are peeks multiple incoming messages from the queue without removing them.
// If there are no messages available in the queue it returns an empty list.
// Important: This method is not supported by MQTT.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - messageCount      a maximum number of messages to peek.
//
// Returns:          callback function that receives a list with messages or error.
func (c *MqttMessageQueue) PeekBatch(ctx context.Context, messageCount int64) ([]*cqueues.MessageEnvelope, error) {
	err := c.CheckOpen(cctx.GetTraceId(ctx))
	if err != nil {
		return nil, err
	}

	// Subscribe if needed
	err = c.subscribe(ctx)
	if err != nil {
		return nil, err
	}

	c.Lock.Lock()
	batchMessages := c.messages
	if messageCount <= (int64)(len(batchMessages)) {
		batchMessages = batchMessages[0:messageCount]
	}
	c.Lock.Unlock()

	messages := []*cqueues.MessageEnvelope{}
	for _, message := range batchMessages {
		messages = append(messages, &message)
	}

	c.Logger.Trace(ctx, "Peeked %d messages on %s", len(messages), c.Name())

	return messages, nil
}

// Receive method are receives an incoming message and removes it from the queue.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - waitTimeout  time.Duration     a timeout in milliseconds to wait for a message to come.
//
// Returns:  result *cqueues.MessageEnvelope, err error
// receives a message or error.
func (c *MqttMessageQueue) Receive(ctx context.Context, waitTimeout time.Duration) (*cqueues.MessageEnvelope, error) {
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
		message = &c.messages[0]
		c.messages = c.messages[1:]

		// Add messages to locked messages list
		messageReceived = true
		c.Lock.Unlock()
	}

	return message, nil
}

// Send method are sends a message into the queue.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - envelope *cqueues.MessageEnvelope  a message envelop to be sent.
//
// Returns: error or nil for success.
func (c *MqttMessageQueue) Send(ctx context.Context, envelop *cqueues.MessageEnvelope) error {
	err := c.CheckOpen(cctx.GetTraceId(ctx))
	if err != nil {
		return err
	}

	c.Counters.IncrementOne(ctx, "queue."+c.Name()+".sent_messages")
	c.Logger.Debug(cctx.NewContextWithTraceId(ctx, envelop.TraceId), "Sent message %s via %s", envelop.String(), c.Name())

	msg, err := c.fromMessage(envelop)
	if err != nil {
		return err
	}

	topic := c.topic
	if topic == "" {
		topic = c.Name()
	}

	err = c.Connection.Publish(ctx, topic, c.qos, c.retain, msg)
	if err != nil {
		c.Logger.Error(cctx.NewContextWithTraceId(ctx, envelop.TraceId), err, "Failed to send message via %s", c.Name())
		return err
	}

	return nil
}

// RenewLock method are renews a lock on a message that makes it invisible from other receivers in the queue.
// This method is usually used to extend the message processing time.
// Important: This method is not supported by MQTT.
// Parameters:
//   - ctx context.Context	operation context.
//   - message   *cqueues.MessageEnvelope    a message to extend its lock.
//   - lockTimeout  time.Duration  a locking timeout in milliseconds.
//
// Returns: error
// receives an error or nil for success.
func (c *MqttMessageQueue) RenewLock(ctx context.Context, message *cqueues.MessageEnvelope, lockTimeout time.Duration) (err error) {
	// Not supported
	return nil
}

// Complete method are permanently removes a message from the queue.
// This method is usually used to remove the message after successful processing.
// Important: This method is not supported by MQTT.
// Parameters:
//   - ctx context.Context	operation context.
//   - message  *cqueues.MessageEnvelope a message to remove.
//
// Returns: error
// error or nil for success.
func (c *MqttMessageQueue) Complete(ctx context.Context, message *cqueues.MessageEnvelope) error {
	// Not supported
	return nil
}

// Abandon method are returnes message into the queue and makes it available for all subscribers to receive it again.
// This method is usually used to return a message which could not be processed at the moment
// to repeat the attempt. Messages that cause unrecoverable errors shall be removed permanently
// or/and send to dead letter queue.
// Important: This method is not supported by MQTT.
// Parameters:
//   - ctx context.Context	operation context.
//   - message *cqueues.MessageEnvelope  a message to return.
//
// Returns: error
//
//	error or nil for success.
func (c *MqttMessageQueue) Abandon(ctx context.Context, message *cqueues.MessageEnvelope) error {
	// Not supported
	return nil
}

// Permanently removes a message from the queue and sends it to dead letter queue.
// Important: This method is not supported by MQTT.
// Parameters:
//   - ctx context.Context	operation context.
//   - message  *cqueues.MessageEnvelope a message to be removed.
//
// Returns: error
//
//	error or nil for success.
func (c *MqttMessageQueue) MoveToDeadLetter(ctx context.Context, message *cqueues.MessageEnvelope) error {
	// Not supported
	return nil
}

func (c *MqttMessageQueue) sendMessageToReceiver(ctx context.Context, receiver cqueues.IMessageReceiver, message *cqueues.MessageEnvelope) {
	traceId := message.TraceId

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("%v", r)
			c.Logger.Error(cctx.NewContextWithTraceId(ctx, traceId), nil, "Failed to process the message - "+err)
		}
	}()

	err := receiver.ReceiveMessage(ctx, message, c)
	if err != nil {
		c.Logger.Error(cctx.NewContextWithTraceId(ctx, traceId), err, "Failed to process the message")
	}
}

// Listens for incoming messages and blocks the current thread until queue is closed.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - receiver    cqueues.IMessageReceiver      a receiver to receive incoming messages.
//
// See IMessageReceiver
// See receive
func (c *MqttMessageQueue) Listen(ctx context.Context, receiver cqueues.IMessageReceiver) error {
	err := c.CheckOpen(cctx.GetTraceId(ctx))
	if err != nil {
		return err
	}

	// Subscribe if needed
	err = c.subscribe(ctx)
	if err != nil {
		return err
	}

	c.Logger.Trace(ctx, "Started listening messages at %s", c.Name())

	// Get all collected messages
	c.Lock.Lock()
	batchMessages := c.messages
	c.messages = []cqueues.MessageEnvelope{}
	c.Lock.Unlock()

	// Resend collected messages to receiver
	for _, message := range batchMessages {
		receiver.ReceiveMessage(ctx, &message, c)
	}

	// Set the receiver
	c.Lock.Lock()
	c.receiver = receiver
	c.Lock.Unlock()

	return nil
}

// EndListen method are ends listening for incoming messages.
// When this method is call listen unblocks the thread and execution continues.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
func (c *MqttMessageQueue) EndListen(ctx context.Context) {
	c.Lock.Lock()
	c.receiver = nil
	c.Lock.Unlock()
}
