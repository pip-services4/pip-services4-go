package connect

import (
	"context"
	"os"
	"strings"
	"time"

	kafka "github.com/Shopify/sarama"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// Kafka connection using plain driver.
// By defining a connection and sharing it through multiple message queues
// you can reduce number of used connections.
//
// ### Configuration parameters ###
//   - client_id:               (optional) name of the client id
//   - connection(s):
//   - discovery_key:             (optional) a key to retrieve the connection from IDiscovery
//   - host:                      host name or IP address
//   - port:                      port number (default: 27017)
//   - uri:                       resource URI or connection string with all parameters in it
//   - credential(s):
//   - store_key:                 (optional) a key to retrieve the credentials from ICredentialStore
//   - username:                  user name
//   - password:                  user password
//   - options:
//   - acks                  (optional) control the number of required acks: -1 - all, 0 - none, 1 - only leader (default: -1)
//   - num_partitions:       (optional) number of partitions of the created topic (default: 1)
//   - replication_factor:   (optional) kafka replication factor of the topic (default: 1)
//   - log_level:            (optional) log level 0 - None, 1 - Error, 2 - Warn, 3 - Info, 4 - Debug (default: 1)
//   - connect_timeout:      (optional) number of milliseconds to connect to broker (default: 1000)
//   - max_retries:          (optional) maximum retry attempts (default: 5)
//   - retry_timeout:        (optional) number of milliseconds to wait on each reconnection attempt (default: 30000)
//   - request_timeout:      (optional) number of milliseconds to wait on flushing messages (default: 30000)
//
// ### References ###
//   - *:logger:*:*:1.0           (optional) ILogger components to pass log messages
//   - *:discovery:*:*:1.0        (optional) IDiscovery services
//   - *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
type KafkaConnection struct {
	defaultConfig *cconf.ConfigParams
	// The logger.
	Logger *clog.CompositeLogger
	// The connection resolver.
	ConnectionResolver *KafkaConnectionResolver
	// The configuration options.
	Options *cconf.ConfigParams

	// The Kafka connection object.
	connection kafka.SyncProducer

	// Kafka admin client object
	adminClient kafka.ClusterAdmin

	// Kafka client
	client kafka.Client

	// Topic subscriptions
	subscriptions []*KafkaSubscription

	clientId          string
	logLevel          int
	connectTimeout    int
	maxRetries        int
	retryTimeout      int
	requestTimeout    int
	numPartitions     int
	replicationFactor int

	acks int
}

// NewKafkaConnection creates a new instance of the connection component.
func NewKafkaConnection() *KafkaConnection {
	c := &KafkaConnection{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			// "client_id", nil,
			"options.log_level", 1,
			"options.connect_timeout", 1000,
			"options.retry_timeout", 30000,
			"options.max_retries", 5,
			"options.request_timeout", 30000,
		),

		Logger:             clog.NewCompositeLogger(),
		ConnectionResolver: NewKafkaConnectionResolver(),
		Options:            cconf.NewEmptyConfigParams(),

		subscriptions: []*KafkaSubscription{},

		logLevel:          1,
		connectTimeout:    100,
		maxRetries:        3,
		retryTimeout:      30000,
		requestTimeout:    30000,
		numPartitions:     1,
		replicationFactor: 1,
		acks:              -1,
	}

	c.clientId, _ = os.Hostname()

	return c
}

//		Configures component by passing configuration parameters.
//		Parameters:
//			- ctx context.Context	operation context
//	  	- config    configuration parameters to be set.
func (c *KafkaConnection) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.ConnectionResolver.Configure(ctx, config)

	c.Options = c.Options.Override(config.GetSection("options"))

	c.clientId = config.GetAsStringWithDefault("client_id", c.clientId)
	c.logLevel = config.GetAsIntegerWithDefault("options.log_level", c.logLevel)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
	c.maxRetries = config.GetAsIntegerWithDefault("options.max_retries", c.maxRetries)
	c.retryTimeout = config.GetAsIntegerWithDefault("options.retry_timeout", c.retryTimeout)
	c.requestTimeout = config.GetAsIntegerWithDefault("options.request_timeout", c.requestTimeout)

	c.numPartitions = config.GetAsIntegerWithDefault("options.num_partitions", c.numPartitions)
	c.replicationFactor = config.GetAsIntegerWithDefault("options.replication_factor",
		c.replicationFactor)

	c.acks = config.GetAsIntegerWithDefault("options.acks",
		c.acks)
}

//		Sets references to dependent components.
//		Parameters:
//			- ctx context.Context	operation context
//	 	- references 	references to locate the component dependencies.
func (c *KafkaConnection) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// Checks if the component is opened.
// Returns: true if the component has been opened and false otherwise.
func (c *KafkaConnection) IsOpen() bool {
	return c.connection != nil
}

func (c *KafkaConnection) createConfig() ([]string, *kafka.Config, error) {
	options, err := c.ConnectionResolver.Resolve(context.Background())
	if err != nil {
		return nil, nil, err
	}

	uri := options.GetAsString("uri")
	brokers := strings.Split(uri, ",")

	config := kafka.NewConfig()
	config.ClientID = c.clientId
	config.Admin.Retry.Max = c.maxRetries
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = c.maxRetries
	config.Consumer.Return.Errors = true
	config.Producer.Partitioner = kafka.NewManualPartitioner
	config.Producer.RequiredAcks = kafka.RequiredAcks(c.acks)

	config.Net.DialTimeout = time.Millisecond * time.Duration(c.connectTimeout)

	username := options.GetAsString("username")
	password := options.GetAsString("password")
	if username != "" {
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.Enable = true

		mechanism := options.GetAsString("mechanism")
		switch mechanism {
		case "scram-sha-256":
			config.Net.SASL.Mechanism = kafka.SASLTypeSCRAMSHA256
		case "scram-sha-512":
			config.Net.SASL.Mechanism = kafka.SASLTypeSCRAMSHA512
		default:
			config.Net.SASL.Mechanism = kafka.SASLTypePlaintext
		}
	}

	return brokers, config, nil
}

//		Opens the component.
//		Parameters:
//			- ctx context.Context	transaction id to trace execution through call chain.
//	  	- Return 			error or nil no errors occured.
func (c *KafkaConnection) Open(ctx context.Context) error {
	brokers, config, err := c.createConfig()
	if err != nil {
		return err
	}

	uri := strings.Join(brokers, ",")
	connection, err := kafka.NewSyncProducer(brokers, config)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to connect to Kafka broker at "+uri)
		return err
	}

	c.connection = connection

	c.Logger.Debug(ctx, "Connected to Kafka broker at "+uri)

	return nil
}

//	Closes component and frees used resources.
//	Parameters:
//		- ctx context.Context	transaction id to trace execution through call chain.
//
// Return			 error or nil no errors occured
func (c *KafkaConnection) Close(ctx context.Context) error {
	if c.connection == nil {
		return nil
	}

	// Close admin client
	if c.adminClient != nil {
		c.adminClient.Close()
		c.adminClient = nil
	}

	// Close producer
	c.connection.Close()
	c.Logger.Debug(ctx, "Disconnected to Kafka broker")

	// Close all consumers
	for _, subscription := range c.subscriptions {
		if subscription.Handler != nil {
			(*subscription.Handler).Close()
		}
	}

	c.connection = nil
	c.subscriptions = []*KafkaSubscription{}

	return nil
}

// Returns connection object
func (c *KafkaConnection) GetConnection() kafka.SyncProducer {
	return c.connection
}

func (c *KafkaConnection) connectToAdmin() error {
	err := c.checkOpen()
	if err != nil {
		return err
	}

	brokers, config, err := c.createConfig()
	if err != nil {
		return err
	}

	uri := strings.Join(brokers, ",")
	client, err := kafka.NewClient(brokers, config)
	if err != nil {
		c.Logger.Error(context.Background(), err, "Failed to connect to Kafka broker at "+uri)
		return err
	}

	admin, err := kafka.NewClusterAdminFromClient(client)
	if err != nil {
		c.Logger.Error(context.Background(), err, "Failed Kafka admin broker creation at "+uri)
		return err
	}

	c.client = client
	c.adminClient = admin
	return nil
}

// Reads a list of registered queue names.
// If connection doesn't support this function returnes an empty list.
// Returns queue names.
func (c *KafkaConnection) ReadQueueNames() ([]string, error) {
	err := c.connectToAdmin()
	if err != nil {
		return nil, err
	}

	return c.client.Topics()
}

// Creates a message queue.
// If connection doesn't support this function it exists without error.
// Parameters:
//   - name string	the name of the queue to be created.
//
// Returns: the name of the queue to be created.
func (c *KafkaConnection) CreateQueue(name string) error {
	err := c.checkOpen()
	if err != nil {
		return err
	}

	err = c.connectToAdmin()
	if err != nil {
		return err
	}

	err = c.adminClient.CreateTopic(name, &kafka.TopicDetail{
		NumPartitions:     int32(c.numPartitions),
		ReplicationFactor: int16(c.replicationFactor),
	}, false)

	return err
}

// Deletes a message queue.
// If connection doesn't support this function it exists without error.
// Parameters:
//   - name string	the name of the queue to be deleted.
func (c *KafkaConnection) DeleteQueue(name string) error {
	err := c.checkOpen()
	if err != nil {
		return err
	}

	err = c.connectToAdmin()
	if err != nil {
		return err
	}

	return c.adminClient.DeleteTopic(name)
}

func (c *KafkaConnection) checkOpen() error {
	if c.connection != nil {
		return nil
	}

	return cerr.NewInvalidStateError(
		"",
		"NOT_OPEN",
		"Connection was not opened",
	)
}

// Publish a message to a specified topic
//
// Parameters:
//   - ctx context.Context	operation context
//   - topic a topic name
//   - messages messages to be published
//
// Returns: error or nil for success
func (c *KafkaConnection) Publish(ctx context.Context, topic string, messages []*kafka.ProducerMessage) error {
	// Check for open connection
	err := c.checkOpen()
	if err != nil {
		return err
	}

	// Assign topic to messages
	for _, message := range messages {
		message.Topic = topic
	}

	return c.connection.SendMessages(messages)
}

// Subscribe to a topic
//
// Parameters:
//   - ctx context.Context	operation context
//   - topic a subject (topic) name
//   - groupId (optional) a consumer group id
//   - config consumer configuration parameters
//   - listener a message listener
//
// Returns: err or nil for success
func (c *KafkaConnection) Subscribe(ctx context.Context, topic string, groupId string, config *kafka.Config, listener IKafkaMessageListener) error {
	// Check for open connection
	err := c.checkOpen()
	if err != nil {
		return err
	}

	brokers, consumerConfig, err := c.createConfig()
	if err != nil {
		return err
	}

	uri := strings.Join(brokers, ",")

	consumerConfig.Consumer.Offsets.AutoCommit.Enable = config.Consumer.Offsets.AutoCommit.Enable
	consumerConfig.Consumer.Return.Errors = true

	consumer, err := kafka.NewConsumerGroup(brokers, groupId, consumerConfig)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to connect Kafka consumer at "+uri)
		return err
	}

	// Consume messages in a separate thread
	go func(ctx context.Context, consumer kafka.ConsumerGroup, topic string, listener IKafkaMessageListener) {
		err = consumer.Consume(ctx, []string{topic}, listener)
		if err != nil {
			c.Logger.Error(ctx, err, "Failed to consume messages "+uri)
		}

		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}

		listener.SetReady(make(chan bool))
	}(ctx, consumer, topic, listener)

	<-listener.Ready() // wait cunsumer start

	// Create the subscription
	subscription := &KafkaSubscription{
		Topic:    topic,
		GroupId:  groupId,
		Listener: listener,
		Handler:  &consumer,
	}

	// Add the subscription
	c.subscriptions = append(c.subscriptions, subscription)
	return nil
}

// Unsubscribe from a previously subscribed topic topic
//
// Parameters:
//   - ctx context.Context	operation context
//   - topic a topic name
//   - groupId (optional) a consumer group id
//   - listener a message listener
//
// Returns: err or nil for success
func (c *KafkaConnection) Unsubscribe(ctx context.Context, topic string, groupId string, listener IKafkaMessageListener) error {
	// Remove the subscription
	var removedSubscription *KafkaSubscription
	for index, subscription := range c.subscriptions {
		if subscription.Topic == topic && subscription.GroupId == groupId && subscription.Listener == listener {
			removedSubscription = subscription
			c.subscriptions = append(c.subscriptions[:index], c.subscriptions[index+1:]...)
			break
		}
	}

	// If nothing to remove then skip
	if removedSubscription == nil {
		return nil
	}

	// Unsubscribe from the topic
	if removedSubscription.Handler != nil {
		return (*removedSubscription.Handler).Close()
	}

	return nil
}
