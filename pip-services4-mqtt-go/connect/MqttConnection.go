package connect

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// MQTT connection using plain driver.
// By defining a connection and sharing it through multiple message queues
// you can reduce number of used connections.
//
// Configuration parameters
//   - client_id:               (optional) name of the client id
//   - connection(s):
//   - discovery_key:               (optional) a key to retrieve the connection from [[https://pip-services4-node.github.io/pip-services4-components-node/interfaces/connect.idiscovery.html IDiscovery]]
//   - host:                        host name or IP address
//   - port:                        port number
//   - uri:                         resource URI or connection string with all parameters in it
//   - credential(s):
//   - store_key:                   (optional) a key to retrieve the credentials from [[https://pip-services4-node.github.io/pip-services4-components-node/interfaces/auth.icredentialstore.html ICredentialStore]]
//   - username:                    user name
//   - password:                    user password
//   - options:
//   - retry_connect:        (optional) turns on/off automated reconnect when connection is log (default: true)
//   - connect_timeout:      (optional) number of milliseconds to wait for connection (default: 30000)
//   - reconnect_timeout:    (optional) number of milliseconds to wait on each reconnection attempt (default: 1000)
//   - keepalive_timeout:    (optional) number of milliseconds to ping broker while inactive (default: 3000)
//
// References
//   - \*:logger:\*:\*:1.0           (optional) ILogger components to pass log messages
//   - \*:discovery:\*:\*:1.0        (optional) IDiscovery services
//   - \*:credential-store:\*:\*:1.0 (optional) Credential stores to resolve credentials
type MqttConnection struct {
	defaultConfig *cconf.ConfigParams
	// The logger.
	Logger *clog.CompositeLogger
	// The connection resolver.
	ConnectionResolver *MqttConnectionResolver
	// The configuration options.
	Options *cconf.ConfigParams

	// The MQTT connection object.
	Connection mqtt.Client

	// Topic subscriptions
	subscriptions []*MqttSubscription
	lock          sync.Mutex

	clientId         string
	retryConnect     bool
	connectTimeout   int
	reconnectTimeout int
	keepAliveTimeout int
}

// NewMqttConnection creates a new instance of the connection component.
func NewMqttConnection() *MqttConnection {
	c := &MqttConnection{
		defaultConfig: cconf.NewConfigParamsFromTuples(
			"options.retry_connect", true,
			"options.connect_timeout", 30000,
			"options.reconnect_timeout", 1000,
			"options.keepalive_timeout", 60000,
		),

		Logger:             clog.NewCompositeLogger(),
		ConnectionResolver: NewMqttConnectionResolver(),
		Options:            cconf.NewEmptyConfigParams(),

		subscriptions: []*MqttSubscription{},

		retryConnect:     true,
		connectTimeout:   30000,
		reconnectTimeout: 60000,
		keepAliveTimeout: 1000, //!!
	}
	return c
}

// Configures component by passing configuration parameters.
// Parameters:
//   - ctx context.Context	operation context.
//   - config	configuration parameters to be set.
func (c *MqttConnection) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.ConnectionResolver.Configure(ctx, config)

	c.Options = c.Options.Override(config.GetSection("options"))

	c.clientId = config.GetAsStringWithDefault("client_id", c.clientId)
	c.retryConnect = config.GetAsBooleanWithDefault("options.retry_connect", c.retryConnect)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
	c.reconnectTimeout = config.GetAsIntegerWithDefault("options.reconnect_timeout", c.reconnectTimeout)
	c.keepAliveTimeout = config.GetAsIntegerWithDefault("options.keepalive_timeout", c.keepAliveTimeout)
}

// Sets references to dependent components.
// Parameters:
//   - ctx context.Context	operation context.
//   - references 	references to locate the component dependencies.
func (c *MqttConnection) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *MqttConnection) IsOpen() bool {
	return c.Connection != nil
}

// Opens the component.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//   - Return 			error or nil no errors occured.
func (c *MqttConnection) Open(ctx context.Context) error {
	options, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	opts := mqtt.NewClientOptions()

	uri := options.GetAsString("uri")
	uris := strings.Split(uri, ",")
	for _, uri = range uris {
		opts.AddBroker(uri)
	}

	user := options.GetAsString("username")
	if user != "" {
		opts.SetUsername(user)
	}
	passwd := options.GetAsString("password")
	if passwd != "" {
		opts.SetPassword(passwd)
	}

	opts.SetClientID(c.clientId)
	opts.SetAutoReconnect(c.retryConnect)
	opts.SetConnectTimeout(time.Millisecond * time.Duration(c.connectTimeout))
	opts.SetConnectRetryInterval(time.Millisecond * time.Duration(c.reconnectTimeout))
	opts.SetKeepAlive(time.Millisecond * time.Duration(c.keepAliveTimeout))
	//opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		c.Logger.Error(ctx, err, "Failed to connect to MQTT broker at "+uri)
		return err
	}

	c.Connection = client

	c.Logger.Debug(ctx, "Connected to MQTT broker at "+uri)

	return nil
}

// Closes component and frees used resources.
// Parameters:
//   - ctx context.Context	transaction id to trace execution through call chain.
//
// Return			 error or nil no errors occured
func (c *MqttConnection) Close(ctx context.Context) error {
	if c.Connection == nil {
		return nil
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.Connection.Disconnect(250)
	c.Connection = nil
	c.subscriptions = []*MqttSubscription{}

	c.Logger.Debug(ctx, "Disconnected from MQTT broker")

	return nil
}

func (c *MqttConnection) GetConnection() mqtt.Client {
	return c.Connection
}

func (c *MqttConnection) ReadQueueNames() ([]string, error) {
	return []string{}, nil
}

func (c *MqttConnection) CreateQueue() error {
	return nil
}

func (c *MqttConnection) DeleteQueue() error {
	return nil
}

func (c *MqttConnection) checkOpen() error {
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
//   - ctx context.Context	operation context.
//   - topic a topic name
//   - qos quality of service (QOS) for the message
//   - retained retained flag for the message
//   - data a message to be published
//
// Returns: error or nil for success
func (c *MqttConnection) Publish(ctx context.Context, topic string, qos byte, retained bool, data []byte) error {
	// Check for open connection
	err := c.checkOpen()
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	token := c.Connection.Publish(topic, qos, retained, data)
	token.Wait()
	return token.Error()
}

// Subscribe to a topic
//
// Parameters:
// Parameters:
//   - ctx context.Context	operation context.
//   - topic a topic name
//   - qos quality of service (QOS) for the subscription
//   - listener a message listener
//
// Returns: err or nil for success
func (c *MqttConnection) Subscribe(ctx context.Context, topic string, qos byte, listener IMqttMessageListener) error {
	// Check for open connection
	err := c.checkOpen()
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// Create the subscription
	subscription := &MqttSubscription{
		Topic:    topic,
		Qos:      qos,
		Listener: listener,
	}

	// Subscribe to topic
	token := c.Connection.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		// The listener can be removed to keep other subscriptions to the same topic alive
		if atomic.LoadInt32(&subscription.Skip) == 0 {
			subscription.Listener.OnMessage(msg)
		}
	})
	token.Wait()
	err = token.Error()
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
//   - ctx context.Context	operation context.
//   - topic a topic name
//   - qos quality of service (QOS) for the subscription
//   - listener a message listener
//
// Returns: err or nil for success
func (c *MqttConnection) Unsubscribe(ctx context.Context, topic string, listener IMqttMessageListener) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Remove the subscription
	var removedSubscription *MqttSubscription
	for index, subscription := range c.subscriptions {
		if subscription.Topic == topic && subscription.Listener == listener {
			removedSubscription = subscription
			c.subscriptions = append(c.subscriptions[:index], c.subscriptions[index+1:]...)
			break
		}
	}

	// If nothing to remove then skip
	if removedSubscription == nil {
		return nil
	}

	// Unset listener to avoid receiving subscriptions
	atomic.StoreInt32(&removedSubscription.Skip, 1)

	// Check if there are more subscriptions to the same topic
	hasMoreSubscriptions := false
	for _, subscription := range c.subscriptions {
		if subscription.Topic == topic {
			hasMoreSubscriptions = true
			break
		}
	}

	// Unsubscribe from the topic if nobody else listens
	if !hasMoreSubscriptions {
		token := c.Connection.Unsubscribe(topic)
		token.Wait()
		return token.Error()
	}

	return nil
}
