package connect

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// NATS connection using plain driver.
// By defining a connection and sharing it through multiple message queues
// you can reduce number of used connections.
//
//	Configuration parameters
//		- connection(s):
//			- discovery_key:             (optional) a key to retrieve the connection from IDiscovery
//			- host:                      host name or IP address
//			- port:                      port number (default: 27017)
//			- uri:                       resource URI or connection string with all parameters in it
//		- credential(s):
//			- store_key:                 (optional) a key to retrieve the credentials from ICredentialStore
//			- username:                  user name
//			- password:                  user password
//		- options:
//			- retry_connect:        (optional) turns on/off automated reconnect when connection is log (default: true)
//			- max_reconnect:        (optional) maximum reconnection attempts (default: 3)
//			- reconnect_timeout:    (optional) number of milliseconds to wait on each reconnection attempt (default: 3000)
//			- flush_timeout:        (optional) number of milliseconds to wait on flushing messages (default: 3000)
//
//	References
//		- *:logger:*:*:1.0           (optional) ILogger components to pass log messages
//		- *:discovery:*:*:1.0        (optional) IDiscovery services
//		- *:credential-store:*:*:1.0 (optional) Credential stores to resolve credentials
type NatsConnection struct {
	defaultConfig *cconf.ConfigParams
	// The logger.
	Logger *clog.CompositeLogger
	// The connection resolver.
	ConnectionResolver *NatsConnectionResolver
	// The configuration options.
	Options *cconf.ConfigParams

	// The NATS connection object.
	Connection *nats.Conn

	// Topic subscriptions
	subscriptions []*NatsSubscription

	retryConnect     bool
	maxReconnect     int
	reconnectTimeout int
	flushTimeout     int
}

// NewNatsConnection creates a new instance of the connection component.
func NewNatsConnection() *NatsConnection {
	c := &NatsConnection{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"options.retry_connect", true,
			"options.connect_timeout", 0,
			"options.reconnect_timeout", 3000,
			"options.max_reconnect", 3,
			"options.flush_timeout", 3000,
		),

		Logger:             clog.NewCompositeLogger(),
		ConnectionResolver: NewNatsConnectionResolver(),
		Options:            cconf.NewEmptyConfigParams(),

		subscriptions: []*NatsSubscription{},

		retryConnect:     true,
		maxReconnect:     3,
		reconnectTimeout: 3000,
		flushTimeout:     3000,
	}
	return c
}

// Configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *NatsConnection) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.ConnectionResolver.Configure(ctx, config)

	c.Options = c.Options.Override(config.GetSection("options"))

	c.retryConnect = config.GetAsBooleanWithDefault("options.retry_connect", c.retryConnect)
	c.maxReconnect = config.GetAsIntegerWithDefault("options.max_reconnect", c.maxReconnect)
	c.reconnectTimeout = config.GetAsIntegerWithDefault("options.reconnect_timeout", c.reconnectTimeout)
	c.flushTimeout = config.GetAsIntegerWithDefault("options.flush_timeout", c.flushTimeout)
}

// Sets references to dependent components.
//   - references 	references to locate the component dependencies.
func (c *NatsConnection) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *NatsConnection) IsOpen() bool {
	return c.Connection != nil
}

// Opens the component.
//   - ctx context.Context 	transaction id to trace execution through call chain.
//   - Return 			error or nil no errors occured.
func (c *NatsConnection) Open(ctx context.Context) error {
	options, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		return err
	}

	uri := options.GetAsString("uri")

	connectOptions := []nats.Option{
		// nats.RetryOnFailedConnect(c.retryConnect),
		nats.MaxReconnects(c.maxReconnect),
		nats.ReconnectWait(time.Millisecond * time.Duration(c.reconnectTimeout)),
	}

	username := options.GetAsString("username")
	password := options.GetAsString("password")
	if username != "" {
		connectOptions = append(connectOptions, nats.UserInfo(username, password))
	}

	token := options.GetAsString("token")
	if token != "" {
		connectOptions = append(connectOptions, nats.Token(token))
	}

	connection, err := nats.Connect(uri, connectOptions...)
	if err != nil {
		c.Logger.Error(ctx, err, "Failed to connect to NATS server at "+uri)
		return err
	}

	c.Connection = connection

	c.Logger.Debug(ctx, "Connected to NATS server at "+uri)

	return nil
}

// Closes component and frees used resources.
//   - ctx context.Context 	transaction id to trace execution through call chain.
//
// Return			 error or nil no errors occured
func (c *NatsConnection) Close(ctx context.Context) error {
	if c.Connection == nil {
		return nil
	}
	c.Connection.Close()
	c.Logger.Debug(ctx, "Disconnected to NATS server")

	c.Connection = nil
	c.subscriptions = []*NatsSubscription{}

	return nil
}

func (c *NatsConnection) GetConnection() *nats.Conn {
	return c.Connection
}

func (c *NatsConnection) ReadQueueNames() ([]string, error) {
	return []string{}, nil
}

func (c *NatsConnection) CreateQueue() error {
	return nil
}

func (c *NatsConnection) DeleteQueue() error {
	return nil
}

func (c *NatsConnection) checkOpen() error {
	if c.Connection != nil {
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
//   - subject a subject (topic) name
//   - message a message to be published
//
// Returns: error or nil for success
func (c *NatsConnection) Publish(ctx context.Context, subject string, message *nats.Msg) error {
	// Check for open connection
	err := c.checkOpen()
	if err != nil {
		return err
	}

	if subject != "" {
		message.Subject = subject
	}
	return c.Connection.PublishMsg(message)
}

// Subscribe to a topic
//
// Parameters:
//   - subject a subject (topic) name
//   - queue (optional) a queue group
//   - listener a message listener
//
// Returns: err or nil for success
func (c *NatsConnection) Subscribe(ctx context.Context, subject string, queue string, listener INatsMessageListener) error {
	// Check for open connection
	err := c.checkOpen()
	if err != nil {
		return err
	}

	// Create the subscription
	subscription := &NatsSubscription{
		Subject:    subject,
		QueueGroup: queue,
		Listener:   listener,
	}

	// Subscribe to topic
	if queue == "" {
		subscription.Handler, err = c.Connection.QueueSubscribe(subject, queue, listener.OnMessage)
	} else {
		subscription.Handler, err = c.Connection.Subscribe(subject, listener.OnMessage)
	}
	if err != nil {
		return err
	}

	// Add the subscription
	c.subscriptions = append(c.subscriptions, subscription)
	return nil
}

// Unsubscribe from a previously subscribed topic topic
//
// Parameters:
//   - subject a subject (topic) name
//   - queue (optional) a queue group
//   - listener a message listener
//
// Returns: err or nil for success
func (c *NatsConnection) Unsubscribe(ctx context.Context, subject string, queue string, listener INatsMessageListener) error {
	// Remove the subscription
	var removedSubscription *NatsSubscription
	for index, subscription := range c.subscriptions {
		if subscription.Subject == subject && subscription.QueueGroup == queue && subscription.Listener == listener {
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
		return removedSubscription.Handler.Unsubscribe()
	}

	return nil
}
