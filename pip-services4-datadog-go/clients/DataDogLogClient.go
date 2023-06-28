package clients

import (
	"context"
	"time"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cauth "github.com/pip-services4/pip-services4-go/pip-services4-config-go/auth"
	httpclient "github.com/pip-services4/pip-services4-go/pip-services4-http-go/clients"
)

type DataDogLogClient struct {
	*httpclient.RestClient
	defaultConfig      *cconf.ConfigParams
	credentialResolver *cauth.CredentialResolver
}

func NewDataDogLogClient(config *cconf.ConfigParams) *DataDogLogClient {

	c := &DataDogLogClient{
		RestClient:         httpclient.NewRestClient(),
		credentialResolver: cauth.NewEmptyCredentialResolver(),
	}
	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "https",
		"connection.host", "http-intake.logs.datadoghq.com",
		"connection.port", 443,
		"credential.internal_network", "true",
	)

	if config != nil {
		c.Configure(context.Background(), config)
	}
	c.BaseRoute = "v1"
	return c
}

func (c *DataDogLogClient) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = c.defaultConfig.Override(config)
	c.RestClient.Configure(ctx, config)
	c.credentialResolver.Configure(ctx, config)
}

func (c *DataDogLogClient) SetReferences(ctx context.Context, refs cref.IReferences) {
	c.RestClient.SetReferences(ctx, refs)
	c.credentialResolver.SetReferences(ctx, refs)
}

func (c *DataDogLogClient) Open(ctx context.Context) error {
	credential, err := c.credentialResolver.Lookup(ctx)
	if err != nil {
		return err
	}

	if credential == nil || credential.AccessKey() == "" {
		err = cerr.NewConfigError(cctx.GetTraceId(ctx), "NO_ACCESS_KEY", "Missing access key in credentials")
		return err
	}
	if c.Headers.Value() == nil {
		c.Headers = cdata.NewEmptyStringValueMap()
	}
	c.Headers.SetAsObject("DD-API-KEY", credential.AccessKey())
	return c.RestClient.Open(ctx)
}

func (c *DataDogLogClient) convertTags(tags map[string]string) string {
	if tags == nil {
		return ""
	}

	builder := ""

	for key, val := range tags {
		if builder != "" {
			builder += ","
		}
		builder += key + ":" + val
	}
	return builder
}

func (c *DataDogLogClient) convertMessage(message DataDogLogMessage) interface{} {

	timestamp := message.Time
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}
	result := map[string]interface{}{
		"timestamp": cconv.StringConverter.ToString(timestamp),
		"service":   message.Service,
		"message":   message.Message,
	}

	if message.Status != "" {
		result["status"] = message.Status
	} else {
		result["status"] = "INFO"
	}

	if message.Source != "" {
		result["ddsource"] = message.Source
	} else {
		result["ddsource"] = "pip-services"
	}

	if message.Tags != nil {
		result["ddtags"] = c.convertTags(message.Tags)
	}
	if message.Host != "" {
		result["host"] = message.Host
	}
	if message.LoggerName != "" {
		result["logger.name"] = message.LoggerName
	}
	if message.ThreadName != "" {
		result["logger.thread_name"] = message.ThreadName
	}
	if message.ErrorMessage != "" {
		result["error.message"] = message.ErrorMessage
	}
	if message.ErrorKind != "" {
		result["error.kind"] = message.ErrorKind
	}
	if message.ErrorStack != "" {
		result["error.stack"] = message.ErrorStack
	}

	return result
}

func (c *DataDogLogClient) convertMessages(messages []DataDogLogMessage) []interface{} {
	result := make([]interface{}, 0)

	for _, msg := range messages {
		result = append(result, c.convertMessage(msg))
	}
	return result
}

func (c *DataDogLogClient) SendLogs(ctx context.Context, messages []DataDogLogMessage) error {
	data := c.convertMessages(messages)

	// Commented instrumentation because otherwise it will never stop sending logs...
	//let timing = c.instrument(ctx, "datadog.send_logs");
	result, err := c.Call(ctx, "post", "input", nil, data)
	//timing.endTiming();
	_, err = c.InstrumentError(ctx, "datadog.send_logs", err, result)
	return err

}
