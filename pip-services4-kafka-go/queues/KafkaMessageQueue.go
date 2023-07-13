package queues

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	kafka "github.com/Shopify/sarama"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	connect "github.com/pip-services4/pip-services4-go/pip-services4-kafka-go/connect"
	cqueues "github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// KafkaMessageQueue are message queue that sends and receives messages via Kafka message broker.
//
// Configuration parameters:
//
//   - topic:                         name of Kafka topic to subscribe
//   - group_id:                      (optional) consumer group id (default: default)
//   - from_beginning:                (optional) restarts receiving messages from the beginning (default: false)
//   - read_partitions:               (optional) number of partitions to be consumed concurrently (default: 1)
//   - autocommit:                    (optional) turns on/off autocommit (default: true)
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
//   - read_partitions:      	(optional) list of partition indexes to be read (default: all, set for example: "1;5;7")
//   - write_partition:		(optional) list of partition indexes to be read (default: auto (-1))
//   - autosubscribe:        	(optional) true to automatically subscribe on option (default: false)
//   - log_level:            	(optional) log level 0 - None, 1 - Error, 2 - Warn, 3 - Info, 4 - Debug (default: 1)
//   - connect_timeout:      	(optional) number of milliseconds to connect to broker (default: 1000)
//   - max_retries:          	(optional) maximum retry attempts (default: 5)
//   - retry_timeout:        	(optional) number of milliseconds to wait on each reconnection attempt (default: 30000)
//   - request_timeout:      	(optional) number of milliseconds to wait on flushing messages (default: 30000)
//
// References:
//
//   - *:logger:*:*:1.0             (optional)  ILogger components to pass log messages
//   - *:counters:*:*:1.0           (optional)  ICounters components to pass collected measurements
//   - *:discovery:*:*:1.0          (optional)  IDiscovery services to resolve connections
//   - *:credential-store:*:*:1.0   (optional) Credential stores to resolve credentials
//   - *:connection:kafka:*:1.0      (optional) Shared connection to Kafka service
//
// See MessageQueue
// See MessagingCapabilities
//
// Example:
//
//		ctx := context.Context()
//	    queue := NewKafkaMessageQueue("myqueue")
//	    queue.Configure(ctx, cconf.NewConfigParamsFromTuples(
//	      "subject", "mytopic",
//	      "connection.protocol", "kafka",
//	      "connection.host", "localhost",
//	      "connection.port", 1883,
//	    ))
//
//	    _ = queue.Open(ctx)
//
//	    _ = queue.Send(ctx, NewMessageEnvelope("", "mymessage", "ABC"))
//
//	    message, err := queue.Receive(ctx, 10000*time.Milliseconds)
//		if (message != nil) {
//			...
//			queue.Complete(ctx, message);
//		}
type KafkaMessageQueue struct {
	*cqueues.MessageQueue

	defaultConfig   *cconf.ConfigParams
	config          *cconf.ConfigParams
	references      cref.IReferences
	opened          bool
	localConnection bool

	// The dependency resolver.
	DependencyResolver *cref.DependencyResolver
	// The logger.
	Logger *clog.CompositeLogger
	// The Kafka connection component.
	Connection *connect.KafkaConnection

	topic         string
	groupId       string
	fromBeginning bool
	autoCommit    bool
	autoSubscribe bool
	subscribed    bool
	messages      []*cqueues.MessageEnvelope
	receiver      cqueues.IMessageReceiver

	writePartition     int
	readablePartitions []int32

	ready chan bool
}

// Creates a new instance of the queue component.
// Parameters:
//   - name    (optional) a queue name.
func NewKafkaMessageQueue(name string) *KafkaMessageQueue {
	c := KafkaMessageQueue{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"topic", nil,
			"group_id", "default",
			"from_beginning", false,
			"read_partitions", 1,
			"autocommit", true,
			"options.autosubscribe", false,
			"options.log_level", 1,
			"options.connect_timeout", 1000,
			"options.retry_timeout", 30000,
			"options.max_retries", 5,
			"options.request_timeout", 30000,
		),
		Logger: clog.NewCompositeLogger(),

		writePartition:     -1,
		readablePartitions: make([]int32, 0),

		ready: make(chan bool),
	}
	c.MessageQueue = cqueues.InheritMessageQueue(&c, name,
		cqueues.NewMessagingCapabilities(false, true, true, true, true, false, false, false, true))
	c.DependencyResolver = cref.NewDependencyResolver()
	c.DependencyResolver.Configure(context.Background(), c.defaultConfig)

	c.messages = make([]*cqueues.MessageEnvelope, 0)

	return &c
}

// Configures component by passing configuration parameters.
// Parameters:
//   - ctx context.Context	operation context
//   - config    configuration parameters to be set.
func (c *KafkaMessageQueue) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config

	c.DependencyResolver.Configure(ctx, config)

	c.topic = config.GetAsStringWithDefault("topic", c.topic)
	c.groupId = config.GetAsStringWithDefault("group_id", c.groupId)
	c.fromBeginning = config.GetAsBooleanWithDefault("from_beginning", c.fromBeginning)
	c.autoCommit = config.GetAsBooleanWithDefault("autocommit", c.autoCommit)
	c.autoSubscribe = config.GetAsBooleanWithDefault("options.autosubscribe", c.autoSubscribe)

	c.writePartition = config.GetAsIntegerWithDefault("options.write_partition", c.writePartition)

	if partitions, ok := config.GetAsNullableString("options.read_partitions"); ok {
		for _, strVal := range strings.Split(partitions, ";") {
			val, err := strconv.Atoi(strVal)
			if err != nil {
				continue
			}
			c.readablePartitions = append(c.readablePartitions, int32(val))
		}
	}
}

// Sets references to dependent components.
// Parameters:
//   - ctx context.Context
//   - references 	references to locate the component dependencies.
func (c *KafkaMessageQueue) SetReferences(ctx context.Context, references cref.IReferences) {
	c.references = references
	c.Logger.SetReferences(ctx, references)

	// Get connection
	c.DependencyResolver.SetReferences(ctx, references)
	result := c.DependencyResolver.GetOneOptional("connection")
	if dep, ok := result.(*connect.KafkaConnection); ok {
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
// Parameters:
//   - ctx context.Context	operation context
func (c *KafkaMessageQueue) UnsetReferences(ctx context.Context) {
	c.Connection = nil
}

func (c *KafkaMessageQueue) createConnection() *connect.KafkaConnection {
	connection := connect.NewKafkaConnection()

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
func (c *KafkaMessageQueue) IsOpen() bool {
	return c.opened
}

// Opens the component.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - Returns 			 error or nil no errors occured.
func (c *KafkaMessageQueue) Open(ctx context.Context) (err error) {
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
		err = cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "Kafka connection is missing")
	}

	if err == nil && !c.Connection.IsOpen() {
		err = cerr.NewConnectionError(cctx.GetTraceId(ctx), "CONNECT_FAILED", "Kafka connection is not opened")
	}

	if err != nil {
		return err
	}

	// Create topic if it does not exist
	topics, err := c.Connection.ReadQueueNames()
	if err != nil {
		return err
	}

	found := false
	for _, v := range topics {
		if v == c.getTopic() {
			found = true
			break
		}
	}

	if !found {
		err := c.Connection.CreateQueue(c.getTopic())
		if err != nil {
			return err
		}
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
//   - ctx context.Context transaction id to trace execution through call chain.
//   - Returns 			error or nil no errors occured.
func (c *KafkaMessageQueue) Close(ctx context.Context) (err error) {
	if !c.opened {
		return nil
	}

	if c.Connection == nil {
		return cerr.NewInvalidStateError(cctx.GetTraceId(ctx), "NO_CONNECTION", "Kafka connection is missing")
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
		c.Connection.Unsubscribe(ctx, topic, c.groupId, c)
		c.subscribed = false
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.opened = false
	c.receiver = nil
	c.messages = make([]*cqueues.MessageEnvelope, 0)

	return nil
}

func (c *KafkaMessageQueue) getTopic() string {
	if c.topic != "" {
		return c.topic
	}
	return c.Name()
}

func (c *KafkaMessageQueue) subscribe(ctx context.Context) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	// Check if already were subscribed
	if c.subscribed {
		return nil
	}

	// Subscribe to the topic
	topic := c.getTopic()
	config := kafka.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = c.autoCommit
	// config.Consumer.Offsets.Initial = kafka.OffsetOldest

	err := c.Connection.Subscribe(ctx, topic, c.groupId, config, c)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to subscribe to topic "+topic)
		return err
	}

	c.subscribed = true
	return nil
}

// Set bool channel with ready flag for consumer
//
//	Parameters:
//		- chFlag	bool channel
func (c *KafkaMessageQueue) SetReady(chFlag chan bool) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.ready = chFlag
}

// Returns: channel with bool flag ready
func (c *KafkaMessageQueue) Ready() chan bool {
	return c.ready
}

func (c *KafkaMessageQueue) fromMessage(message *cqueues.MessageEnvelope) (*kafka.ProducerMessage, error) {
	if message == nil {
		return nil, nil
	}

	headers := []kafka.RecordHeader{
		{
			Key:   []byte("trace_id"),
			Value: []byte(message.TraceId),
		},
		{
			Key:   []byte("message_type"),
			Value: []byte(message.MessageType),
		},
	}

	msg := &kafka.ProducerMessage{}
	msg.Topic = c.getTopic()
	msg.Key = kafka.StringEncoder(message.MessageId)
	msg.Value = kafka.ByteEncoder(message.Message)
	msg.Headers = headers

	msg.Timestamp = time.Now()

	return msg, nil
}

func (c *KafkaMessageQueue) toMessage(msg *connect.KafkaMessage) (*cqueues.MessageEnvelope, error) {
	messageType := c.getHeaderByKey(msg.Message.Headers, "message_type")
	traceId := c.getHeaderByKey(msg.Message.Headers, "trace_id")

	message := cqueues.NewMessageEnvelope(traceId, messageType, nil)
	message.MessageId = string(msg.Message.Key)
	message.SentTime = msg.Message.Timestamp
	message.Message = msg.Message.Value
	message.SetReference(msg)

	return message, nil
}

func (c *KafkaMessageQueue) getHeaderByKey(headers []*kafka.RecordHeader, key string) string {
	for _, header := range headers {
		if key == string(header.Key) {
			return string(header.Value)
		}
	}
	return ""
}

// Setup is run at the beginning of a new session, before ConsumeClaim
// Send ready flag into channel
// Returns: error
func (c *KafkaMessageQueue) Setup(kafka.ConsumerGroupSession) error {
	// Mark the consumer as ready
	c.ready <- true
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *KafkaMessageQueue) Cleanup(kafka.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *KafkaMessageQueue) ConsumeClaim(session kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29

	for {
		// if len(c.readablePartitions) == 0 || slices.Contains(c.readablePartitions, claim.Partition()) {
		select {
		case msg := <-claim.Messages():
			if msg != nil {
				message := &connect.KafkaMessage{
					Message: msg,
					Session: session,
				}

				c.OnMessage(session.Context(), message)

				if c.autoCommit {
					session.MarkMessage(msg, "")
					session.Commit()
				}
			}
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
		// }
	}
}

// Callback for processing messages from kafka
// Parameters:
//   - ctx context.Context	operation context
//   - msg *connect.KafkaMessage	consumer message
func (c *KafkaMessageQueue) OnMessage(ctx context.Context, msg *connect.KafkaMessage) {
	// // Skip if it came from a wrong topic
	// expectedTopic := c.getTopic()
	// if !strings.Contains(expectedTopic, "*") && expectedTopic != msg.Topic {
	// 	return
	// }

	// Deserialize message
	message, err := c.toMessage(msg)

	if message == nil || err != nil {
		c.Logger.Error(ctx, err, "Failed to read received message")
		return
	}

	c.Counters.IncrementOne(ctx, "queue."+c.Name()+".received_messages")
	c.Logger.Debug(cctx.NewContextWithTraceId(ctx, message.TraceId), "Received message %s via %s", message, c.Name())

	// Send message to receiver if its set or put it into the queue
	c.Lock.Lock()
	if c.receiver != nil {
		receiver := c.receiver
		c.Lock.Unlock()
		c.sendMessageToReceiver(ctx, receiver, message)
	} else {
		c.messages = append(c.messages, message)
		c.Lock.Unlock()
	}
}

// Clear method are clears component state.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//
// Returns error or nil no errors occured.
func (c *KafkaMessageQueue) Clear(ctx context.Context) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.messages = make([]*cqueues.MessageEnvelope, 0)

	return nil
}

// ReadMessageCount method are reads the current number of messages in the queue to be delivered.
// Returns number of messages or error.
func (c *KafkaMessageQueue) ReadMessageCount() (int64, error) {
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
func (c *KafkaMessageQueue) Peek(ctx context.Context) (*cqueues.MessageEnvelope, error) {
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
// Important: This method is not supported by Kafka.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - messageCount      a maximum number of messages to peek.
//
// Returns:          callback function that receives a list with messages or error.
func (c *KafkaMessageQueue) PeekBatch(ctx context.Context, messageCount int64) ([]*cqueues.MessageEnvelope, error) {
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
	messages = append(messages, batchMessages...)

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
func (c *KafkaMessageQueue) Receive(ctx context.Context, waitTimeout time.Duration) (*cqueues.MessageEnvelope, error) {
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

// Send method are sends a message into the queue.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - envelope *cqueues.MessageEnvelope  a message envelop to be sent.
//
// Returns: error or nil for success.
func (c *KafkaMessageQueue) Send(ctx context.Context, envelop *cqueues.MessageEnvelope) error {
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

	topic := c.Name()
	if topic == "" {
		topic = c.topic
	}

	if c.writePartition != -1 {
		msg.Partition = int32(c.writePartition)
	}

	err = c.Connection.Publish(ctx, topic, []*kafka.ProducerMessage{msg})
	if err != nil {
		c.Logger.Error(cctx.NewContextWithTraceId(ctx, envelop.TraceId), err, "Failed to send message via %s", c.Name())
		return err
	}

	return nil
}

// RenewLock method are renews a lock on a message that makes it invisible from other receivers in the queue.
// This method is usually used to extend the message processing time.
// Important: This method is not supported by Kafka.
// Parameters:
//   - ctx context.Context	operation context
//   - message   *cqueues.MessageEnvelope    a message to extend its lock.
//   - lockTimeout  time.Duration  a locking timeout in milliseconds.
//
// Returns: error
// receives an error or nil for success.
func (c *KafkaMessageQueue) RenewLock(ctx context.Context, message *cqueues.MessageEnvelope, lockTimeout time.Duration) (err error) {
	// Not supported
	return nil
}

// Complete method are permanently removes a message from the queue.
// This method is usually used to remove the message after successful processing.
// Parameters:
//   - ctx context.Context	operation context
//   - message  *cqueues.MessageEnvelope a message to remove.
//
// Returns: error
// error or nil for success.
func (c *KafkaMessageQueue) Complete(ctx context.Context, message *cqueues.MessageEnvelope) error {
	err := c.CheckOpen("")
	if err != nil {
		return err
	}

	msg := message.GetReference().(*connect.KafkaMessage)

	// Skip on autocommit
	if c.autoCommit || msg == nil {
		return nil
	}

	// Commit the message offset so it won't come back
	msg.Session.MarkOffset(msg.Message.Topic, msg.Message.Partition, msg.Message.Offset, "")
	msg.Session.Commit()
	message.SetReference(nil)

	return nil
}

//		Abandon method are returnes message into the queue and makes it available for all subscribers to receive it again.
//		This method is usually used to return a message which could not be processed at the moment
//		to repeat the attempt. Messages that cause unrecoverable errors shall be removed permanently
//		or/and send to dead letter queue.
//		Parameters:
//			- ctx context.Context	operation context
//			- message *cqueues.MessageEnvelope  a message to return.
//		Returns: error
//	 error or nil for success.
func (c *KafkaMessageQueue) Abandon(ctx context.Context, message *cqueues.MessageEnvelope) error {
	err := c.CheckOpen("")
	if err != nil {
		return err
	}

	msg := message.GetReference().(*connect.KafkaMessage)

	// Skip on autocommit
	if c.autoCommit || msg == nil {
		return nil
	}

	// Seek to the message offset so it will come back
	msg.Session.ResetOffset(msg.Message.Topic, msg.Message.Partition, msg.Message.Offset, "")
	msg.Session.Commit()
	message.SetReference(nil)

	return nil
}

// Permanently removes a message from the queue and sends it to dead letter queue.
// Important: This method is not supported by Kafka.
// Parameters:
//   - ctx context.Context	operation context
//   - message  *cqueues.MessageEnvelope a message to be removed.
//
// Returns: error
// error or nil for success.
func (c *KafkaMessageQueue) MoveToDeadLetter(ctx context.Context, message *cqueues.MessageEnvelope) error {
	// Not supported
	return nil
}

func (c *KafkaMessageQueue) sendMessageToReceiver(ctx context.Context, receiver cqueues.IMessageReceiver, message *cqueues.MessageEnvelope) {
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
func (c *KafkaMessageQueue) Listen(ctx context.Context, receiver cqueues.IMessageReceiver) error {
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
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
func (c *KafkaMessageQueue) EndListen(ctx context.Context) {
	c.Lock.Lock()
	c.receiver = nil
	c.Lock.Unlock()
}
