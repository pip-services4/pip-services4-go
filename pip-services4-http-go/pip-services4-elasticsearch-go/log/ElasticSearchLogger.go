package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	crpccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

/*
ElasticSearchLogger is logger that dumps execution logs to ElasticSearch service.
ElasticSearch is a popular search index. It is often used
to store and index execution logs by itself or as a part of
ELK (ElasticSearch - Logstash - Kibana) stack.

Authentication is not supported in this version.

Configuration parameters:

- level:             maximum log level to capture
- source:            source (context) name
- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port int
  - uri:                   resource URI or connection string with all parameters in it

- options:
  - interval:        interval in milliseconds to save log messages (default: 10 seconds)
  - max_cache_size:  maximum int of messages stored in this cache (default: 100)
  - index:           ElasticSearch index name (default: "log")
  - daily:           true to create a new index every day by adding date suffix to the index  (default: false)
  - reconnect:       reconnect timeout in milliseconds (default: 60 sec)
  - timeout:         invocation timeout in milliseconds (default: 30 sec)
  - max_retries:     maximum int of retries (default: 3)
  - index_message:   true to enable indexing for message object (default: false)
  - include_type_name: Will create using a "typed" index compatible with ElasticSearch 6.x (default: false)

References:

- *:context-info:*:*:1.0      (optional)  ContextInfo to detect the context id and specify counters source
- *:discovery:*:*:1.0         (optional)  IDiscovery services to resolve connection

Example:

	logger := NewElasticSearchLogger();
	logger.Configure(contex.Background(), cconf.NewConfigParamsFromTuples(
	    "connection.protocol", "http",
	    "connection.host", "localhost",
	    "connection.port", "9200"
	));

	logger.Open(contex.Background(), "123")

	logger.Error(contex.Background(), "123", ex, "Error occured: %s", ex.message);
	logger.Debug(contex.Background(), "123", "Everything is OK.");
*/
type ElasticSearchLogger struct {
	*clog.CachedLogger
	connectionResolver *crpccon.HttpConnectionResolver

	timer        chan bool
	index        string
	dailyIndex   bool
	currentIndex string
	reconnect    int
	timeout      int
	maxRetries   int
	indexMessage bool

	includeTypeName bool

	client *esv7.Client
}

// NewElasticSearchLogger method creates a new instance of the logger.
// Retruns *ElasticSearchLogger
// pointer on new ElasticSearchLogger
func NewElasticSearchLogger() *ElasticSearchLogger {
	c := ElasticSearchLogger{}
	c.CachedLogger = clog.InheritCachedLogger(&c)
	c.connectionResolver = crpccon.NewHttpConnectionResolver()
	c.index = "log"
	c.dailyIndex = false
	c.reconnect = 60000
	c.timeout = 30000
	c.maxRetries = 3
	c.Interval = 10000
	c.indexMessage = false
	c.includeTypeName = false
	return &c
}

// Configure are configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- config  *cconf.ConfigParams   configuration parameters to be set.
func (c *ElasticSearchLogger) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.CachedLogger.Configure(ctx, config)

	c.connectionResolver.Configure(ctx, config)

	c.index = config.GetAsStringWithDefault("index", c.index)
	c.dailyIndex = config.GetAsBooleanWithDefault("daily", c.dailyIndex)
	c.reconnect = config.GetAsIntegerWithDefault("options.reconnect", c.reconnect)
	c.timeout = config.GetAsIntegerWithDefault("options.timeout", c.timeout)
	c.maxRetries = config.GetAsIntegerWithDefault("options.max_retries", c.maxRetries)
	c.indexMessage = config.GetAsBooleanWithDefault("options.index_message", c.indexMessage)
	c.includeTypeName = config.GetAsBooleanWithDefault("options.include_type_name", c.includeTypeName)
}

// SetReferences method are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- references cref.IReferences 	references to locate the component dependencies.
func (c *ElasticSearchLogger) SetReferences(ctx context.Context, references cref.IReferences) {
	c.CachedLogger.SetReferences(ctx, references)
	c.connectionResolver.SetReferences(ctx, references)
}

// IsOpen method are checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *ElasticSearchLogger) IsOpen() bool {
	return c.timer != nil
}

// Open method are ppens the component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//
// Returns error or nil, if no errors occured.
func (c *ElasticSearchLogger) Open(ctx context.Context) (err error) {
	if c.IsOpen() {
		return nil
	}

	connection, _, err := c.connectionResolver.Resolve(ctx)

	if connection == nil {
		err = cerr.NewConfigError(cctx.GetTraceId(ctx), "NO_CONNECTION", "Connection is not configured")
	}

	if err != nil {
		return err
	}

	uri := connection.Uri()

	options := esv7.Config{
		Addresses: []string{uri},
		Transport: &http.Transport{
			ResponseHeaderTimeout: (time.Duration)(c.timeout) * time.Millisecond,
			IdleConnTimeout:       (time.Duration)(c.reconnect) * time.Millisecond},
		MaxRetries: c.maxRetries,
	}

	elasticsearch, esErr := esv7.NewClient(options)
	if esErr != nil {
		return esErr
	}
	c.client = elasticsearch

	err = c.createIndexIfNeeded(ctx, true)
	if err == nil {
		c.timer = setInterval(func() { c.Dump(ctx) }, c.Interval, true)
	}

	return nil
}

// Close method are closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//
// Returns error or nil, if no errors occured.
func (c *ElasticSearchLogger) Close(ctx context.Context) (err error) {
	svErr := c.Save(ctx, c.Cache)
	if svErr == nil {
		return svErr
	}

	if c.timer != nil {
		c.timer <- true
	}

	c.Cache = make([]clog.LogMessage, 0)

	close(c.timer)
	c.timer = nil
	c.client = nil
	return nil
}

func (c *ElasticSearchLogger) getCurrentIndex() string {
	if !c.dailyIndex {
		return c.index
	}
	now := time.Now()
	return c.index + "-" + now.UTC().Format("20060102")
}

func (c *ElasticSearchLogger) createIndexIfNeeded(ctx context.Context, force bool) (err error) {
	newIndex := c.getCurrentIndex()
	if !force && c.currentIndex == newIndex {
		return nil
	}

	c.currentIndex = newIndex
	exists, err := c.client.Indices.Exists([]string{c.currentIndex}, c.client.Indices.Exists.WithContext(ctx))
	if err != nil || exists.StatusCode == 404 {
		return err
	}

	indBody := `{
		"settings": {
			"number_of_shards": "1"
		},
		"mappings": {
			` + c.getIndexSchema() + `
		}
	}`

	resp, err := c.client.Indices.Create(c.currentIndex,
		c.client.Indices.Create.WithContext(ctx),
		c.client.Indices.Create.WithBody(strings.NewReader(indBody)),
	)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return err
		}
		// Skip already exist errors
		if strings.Contains(e["error"].(map[string]interface{})["type"].(string), "resource_already_exists") {
			return nil
		}
		err = cerr.NewError(e["error"].(map[string]interface{})["type"].(string)).WithCauseString(e["error"].(map[string]interface{})["reason"].(string))
		return err
	}
	return nil
}

func (c *ElasticSearchLogger) getIndexSchema() string {
	schema := `"properties": {
					"time": { "type": "date", "index": true },
					"source": { "type": "keyword", "index": true },
					"level": { "type": "keyword", "index": true },
					"trace_id": { "type": "text", "index": true },
					"error": {
						"type": "object",
						"properties": {
							"type": { "type": "keyword", "index": true },
							"category": { "type": "keyword", "index": true },
							"status": { "type": "integer", "index": false },
							"code": { "type": "keyword", "index": true },
							"message": { "type": "text", "index": false },
							"details": { "type": "object" },
							"trace_id": { "type": "text", "index": false },
							"cause": { "type": "text", "index": false },
							"stack_trace": { "type": "text", "index": false }
						}
					},
					"message": { "type": "text", "index":` + strconv.FormatBool(c.indexMessage) + ` }
				}`

	if c.includeTypeName {
		return fmt.Sprintf(`"log_message": {%s}`, schema)
	} else {
		return schema
	}
}

// Save method are saves log messages from the cache.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- messages []clog.LogMessage a list with log messages
//
// Retruns error or nil for success.
func (c *ElasticSearchLogger) Save(ctx context.Context, messages []clog.LogMessage) (err error) {

	if !c.IsOpen() || len(messages) == 0 {
		return nil
	}

	err = c.createIndexIfNeeded(cctx.NewContextWithTraceId(ctx, "elasticsearch_logger"), false)

	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	for _, message := range messages {
		meta := []byte(fmt.Sprintf(`{ "index": %s}%s`, c.getLogItem(), "\n"))
		data, err := json.Marshal(message)

		if err != nil {
			c.Logger.Error(ctx, err, "Cannot encode message "+err.Error())
		}
		data = append(data, "\n"...)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	resp, err := c.client.Bulk(bytes.NewReader(buf.Bytes()), c.client.Bulk.WithContext(ctx), c.client.Bulk.WithIndex(c.currentIndex))
	if err != nil {
		c.Logger.Error(ctx, err, "Failure indexing batch %s", err.Error())
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	buf.Reset()

	if resp.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return err
		}
		err = cerr.NewError(e["error"].(map[string]interface{})["type"].(string)).WithCauseString(e["error"].(map[string]interface{})["reason"].(string))
	}
	return err
}

func (c *ElasticSearchLogger) getLogItem() string {
	if c.includeTypeName {
		return fmt.Sprintf(`{ "_index":"%s", "_type":"log_message", "_id":"%s"}`,
			c.currentIndex,
			keys.IdGenerator.NextLong(),
		) // ElasticSearch 6.x
	} else {
		return fmt.Sprintf(`{ "_index":"%s", "_id": "%s"}`,
			c.currentIndex,
			keys.IdGenerator.NextLong(),
		) // ElasticSearch 7.x
	}
}

func setInterval(someFunc func(), milliseconds int, async bool) chan bool {

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
