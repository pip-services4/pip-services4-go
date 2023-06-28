package log

import (
	"context"
	"os"
	"time"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clients1 "github.com/pip-services4/pip-services4-go/pip-services4-datadog-go/clients"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// Logger that dumps execution logs to DataDog service.
//
// DataDog is a popular monitoring SaaS service. It collects logs, metrics, events
// from infrastructure and applications and analyze them in a single place.
//
// # Configuration parameters
//
// - level:             maximum log level to capture
// - source:            source (context) name
// - connection:
//   - discovery_key:         (optional) a key to retrieve the connection from [[IDiscovery]]
//   - protocol:              (optional) connection protocol: http or https (default: https)
//   - host:                  (optional) host name or IP address (default: http-intake.logs.datadoghq.com)
//   - port:                  (optional) port number (default: 443)
//   - uri:                   (optional) resource URI or connection string with all parameters in it
//
// - credential:
//   - access_key:      DataDog client api key
//
// - options:
//   - interval:        interval in milliseconds to save log messages (default: 10 seconds)
//   - max_cache_size:  maximum number of messages stored in this cache (default: 100)
//   - reconnect:       reconnect timeout in milliseconds (default: 60 sec)
//   - timeout:         invocation timeout in milliseconds (default: 30 sec)
//   - max_retries:     maximum number of retries (default: 3)
//
// References:
//
// - \*:context-info:\*:\*:1.0      (optional) ContextInfo to detect the context id and specify counters source
// - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connection
//
// Example:
//
//	logger := NewDataDogLogger();
//	logger.Configure(context.Background(), NewConfigParamsFromTuples(
//	    "credential.access_key", "827349874395872349875493"
//	))
//
//	err := logger.Open(context.Background(), "123")
//
//	logger.Error(context.Background(), "123", ex, "Error occured: %s", ex.message)
//	logger.Debug(context.Background(), "123", "Everything is OK.")
type DataDogLogger struct {
	*clog.CachedLogger
	client   *clients1.DataDogLogClient
	timer    chan bool
	instance string
}

// NewDataDogLogger - Creates a new instance of the logger.
func NewDataDogLogger() *DataDogLogger {
	c := DataDogLogger{
		client: clients1.NewDataDogLogClient(nil),
	}
	c.CachedLogger = log.InheritCachedLogger(&c)
	c.instance, _ = os.Hostname()
	return &c
}

// Сonfigure -  Configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *DataDogLogger) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.CachedLogger.Configure(ctx, config)
	c.client.Configure(ctx, config)
	c.instance = config.GetAsStringWithDefault("instance", c.instance)
}

// SetReferences - Sets references to dependent components.
//   - references 	references to locate the component dependencies.
func (c *DataDogLogger) SetReferences(ctx context.Context, references cref.IReferences) {
	c.CachedLogger.SetReferences(ctx, references)
	c.client.SetReferences(ctx, references)
	ref := references.GetOneOptional(cref.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))

	contextInfo, _ := ref.(*cctx.ContextInfo)

	if contextInfo != nil && c.Source() == "" {
		c.SetSource(contextInfo.Name)
	}
	if contextInfo != nil && c.instance == "" {
		c.instance = contextInfo.ContextId
	}
}

// IsOpen - Checks if the component is opened.
//
//	Returns true if the component has been opened and false otherwise.
func (c *DataDogLogger) IsOpen() bool {
	return c.timer != nil
}

// Open -  Opens the component.
//   - ctx context.Context execution context to trace execution through call chain.
//   - Returns error or nil no errors occured.
func (c *DataDogLogger) Open(ctx context.Context) error {
	if c.IsOpen() {
		return nil
	}

	err := c.client.Open(ctx)
	if err == nil {
		c.timer = c.setInterval(func() { c.Dump(ctx) }, c.Interval, true)
	}

	return err

}

// Close - Closes component and frees used resources.
//   - ctx context.Context execution context to trace execution through call chain.
//   - Returns error or nil no errors occured.
func (c *DataDogLogger) Close(ctx context.Context) error {
	err := c.Save(ctx, c.Cache)

	if err != nil {
		return err
	}

	if c.timer != nil {
		c.timer <- true
	}

	c.Cache = make([]clog.LogMessage, 0)
	c.timer = nil

	return c.client.Close(ctx)

}

func (c *DataDogLogger) convertMessage(message clog.LogMessage) clients1.DataDogLogMessage {

	tm := message.Time
	if tm.IsZero() {
		tm = time.Now().UTC()
	}
	result := clients1.DataDogLogMessage{
		Time: tm,
		Tags: map[string]string{
			"trace_id": message.TraceId,
		},
		Host:    c.instance,
		Status:  clog.LevelConverter.ToString(message.Level),
		Message: message.Message,
	}

	result.Service = message.Source
	if result.Service == "" {
		result.Service = c.Source()
	}

	if message.Error.Message != "" || message.Error.Code != "" {
		result.ErrorKind = message.Error.Type
		result.ErrorMessage = message.Error.Message
		result.ErrorStack = message.Error.StackTrace
	}

	return result
}

// Saves log messages from the cache.
// - messages  a list with log messages
// - Returns error or nil for success.
func (c *DataDogLogger) Save(ctx context.Context, messages []clog.LogMessage) error {
	if !c.IsOpen() || len(messages) == 0 {
		return nil
	}

	data := make([]clients1.DataDogLogMessage, 0)
	for _, m := range messages {
		data = append(data, c.convertMessage(m))
	}

	return c.client.SendLogs(cctx.NewContextWithTraceId(ctx, "datadog-logger"), data)
}

func (c *DataDogLogger) setInterval(someFunc func(), milliseconds int, async bool) chan bool {

	interval := time.Duration(milliseconds) * time.Millisecond
	ticker := time.NewTicker(interval)
	clear := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					go someFunc()
				} else {
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}

		}
	}()

	return clear
}
