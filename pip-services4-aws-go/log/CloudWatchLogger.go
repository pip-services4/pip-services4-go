package log

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	awsconn "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/connect"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// Logger that writes log messages to AWS Cloud Watch Log.
//
// Configuration parameters
//
//   - stream:                        (optional) Cloud Watch Log stream (default: context name)
//   - group:                         (optional) Cloud Watch Log group (default: context instance ID or hostname)
//   - connections:
//   - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//   - region:                      (optional) AWS region
//   - credentials:
//   - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//   - access_id:                   AWS access/client id
//   - access_key:                  AWS access/client id
//   - options:
//   - interval:        interval in milliseconds to save current counters measurements (default: 5 mins)
//   - reset_timeout:   timeout in milliseconds to reset the counters. 0 disables the reset (default: 0)
//
// References
//
//   - \*:context-info:\*:\*:1.0      (optional) ContextInfo to detect the context id and specify counters source
//   - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connections
//   - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//
// See Counter (in the Pip.Services components package)
// See CachedCounters (in the Pip.Services components package)
// See CompositeLogger (in the Pip.Services components package)
//
// Example:
//
//		logger := NewLogger();
//		logger.Configure(context.Background(), NewConfigParamsFromTuples(
//		    "stream", "mystream",
//		    "group", "mygroup",
//		    "connection.region", "us-east-1",
//		    "connection.access_id", "XXXXXXXXXXX",
//		    "connection.access_key", "XXXXXXXXXXX",
//		))
//		logger.SetReferences(context.Background(), NewReferencesFromTuples(
//		    NewDescriptor("pip-services", "logger", "console", "default", "1.0"),
//		    NewConsoleLogger()
//		))
//
//		err:= logger.Open(context.Background(), "123")
//		    ...
//
//	logger.SetLevel(Debug);
//
//	logger.Error(context.Background(), "123", ex, "Error occured: %s", ex.Message);
//	logger.Debug(context.Background(), "123", "Everything is OK.");
type CloudWatchLogger struct {
	*clog.CachedLogger

	timer chan bool

	connectionResolver *awsconn.AwsConnectionResolver
	client             *cloudwatchlogs.CloudWatchLogs //AmazonCloudWatchLogsClient
	connection         *awsconn.AwsConnectionParams
	connectTimeout     int

	group     string
	stream    string
	lastToken string

	logger *clog.CompositeLogger
}

// Creates a new instance of this logger.
func NewCloudWatchLogger() *CloudWatchLogger {
	c := &CloudWatchLogger{
		connectionResolver: awsconn.NewAwsConnectionResolver(),
		connectTimeout:     30000,
		group:              "undefined",
		stream:             "",
		lastToken:          "",
		logger:             clog.NewCompositeLogger(),
	}
	c.CachedLogger = clog.InheritCachedLogger(c)
	return c
}

// Configure method configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- config	configuration parameters to be set.
func (c *CloudWatchLogger) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.CachedLogger.Configure(ctx, config)
	c.connectionResolver.Configure(ctx, config)

	c.group = config.GetAsStringWithDefault("group", c.group)
	c.stream = config.GetAsStringWithDefault("stream", c.stream)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
}

// SetReferences method sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- references 	references to locate the component dependencies.
//
// See IReferences (in the Pip.Services commons package)
func (c *CloudWatchLogger) SetReferences(ctx context.Context, references cref.IReferences) {
	c.CachedLogger.SetReferences(ctx, references)
	c.logger.SetReferences(ctx, references)

	ref := references.GetOneOptional(cref.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))

	contextInfo, ok := ref.(*cctx.ContextInfo)
	if ok && c.stream == "" {
		c.stream = contextInfo.Name
	}
	if ok && c.group == "" {
		c.group = contextInfo.ContextId
	}
}

// Writes a log message to the logger destination.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- level             a log level.
//		- error             an error object associated with this message.
//		- message           a human-readable message to log.
func (c *CloudWatchLogger) Write(ctx context.Context, level clog.LevelType, ex error, message string) {
	if c.Level() < level {
		return
	}
	c.CachedLogger.Write(ctx, level, ex, message)
}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *CloudWatchLogger) IsOpen() bool {
	return c.timer != nil
}

// Opens the component.
// Parameters:
//   - ctx context.Context	execution context to trace execution through call chain.
//   - Returns 			 error or nil no errors occured.
func (c *CloudWatchLogger) Open(ctx context.Context) error {
	if c.IsOpen() {
		return nil
	}

	var globalErr error

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		connection, err := c.connectionResolver.Resolve(ctx)
		c.connection = connection
		if err != nil {
			globalErr = err
			return
		}

		awsCred := credentials.NewStaticCredentials(c.connection.GetAccessId(), c.connection.GetAccessKey(), "")
		sess := session.Must(session.NewSession(&aws.Config{
			MaxRetries:  aws.Int(3),
			Region:      aws.String(c.connection.GetRegion()),
			Credentials: awsCred,
		}))
		// Create new cloudwatch client.
		c.client = cloudwatchlogs.New(sess)
		c.client.APIVersion = "2014-03-28"
		c.client.Config.HTTPClient.Timeout = time.Duration((int64)(c.connectTimeout)) * time.Millisecond

		groupParam := &cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String(c.group),
		}
		_, groupErr := c.client.CreateLogGroupWithContext(ctx, groupParam)
		if _, ok := groupErr.(*cloudwatchlogs.ResourceAlreadyExistsException); !ok {
			globalErr = groupErr
			return
		}

		streamParam := &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(c.group),
			LogStreamName: aws.String(c.stream),
		}
		_, streamErr := c.client.CreateLogStreamWithContext(ctx, streamParam)

		if streamErr != nil {
			if _, ok := streamErr.(*cloudwatchlogs.ResourceAlreadyExistsException); ok {

				params := &cloudwatchlogs.DescribeLogStreamsInput{
					LogGroupName:        aws.String(c.group),
					LogStreamNamePrefix: aws.String(c.stream),
				}

				descData, describeErr := c.client.DescribeLogStreamsWithContext(ctx, params)
				if describeErr != nil {
					globalErr = describeErr
					return
				}
				if len(descData.LogStreams) > 0 {
					if descData.LogStreams[0].UploadSequenceToken != nil {
						c.lastToken = *descData.LogStreams[0].UploadSequenceToken
					}
				}
			} else {
				globalErr = streamErr
				return
			}
		} else {
			c.lastToken = ""
		}

		if c.timer == nil {
			c.timer = setInterval(func() { c.Dump(ctx) }, c.Interval, true)
		}
	}()
	wg.Wait()

	if globalErr != nil {
		return globalErr
	}
	return nil
}

// Closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- Returns       error or nil no errors occured.
func (c *CloudWatchLogger) Close(ctx context.Context) error {
	err := c.Save(ctx, c.Cache)

	if c.timer != nil {
		c.timer <- true
	}

	c.Cache = make([]clog.LogMessage, 0)
	c.timer = nil
	c.client = nil

	return err
}

func (c *CloudWatchLogger) formatMessageText(message clog.LogMessage) string {

	result := "["

	if message.Source != "" {
		result += message.Source
	} else {
		result += "---"
	}

	result += ":"

	if message.TraceId != "" {
		result += message.TraceId
	} else {
		result += "---"
	}
	result += ":" + clog.LevelConverter.ToString(message.Level) + "] " + message.Message
	if message.Error.Message != "" || message.Error.Code != "" {
		if message.Message == "" {
			result += "Error: "
		} else {
			result += ": "
		}
		result += message.Error.Message
		if message.Error.StackTrace != "" {
			result += " StackTrace: " + message.Error.StackTrace
		}
	}
	return result
}

// Saves log messages from the cache.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- messages  a list with log messages
//		- Returns   error or nil for success.
func (c *CloudWatchLogger) Save(ctx context.Context, messages []clog.LogMessage) error {
	if !c.IsOpen() || messages == nil || len(messages) == 0 {
		return nil
	}

	if c.client == nil {
		err := cerr.NewConfigError("cloudwatch_logger", "NOT_OPENED", "CloudWatchLogger is not opened")
		if err != nil {
			return err
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var events []*cloudwatchlogs.InputLogEvent
		events = make([]*cloudwatchlogs.InputLogEvent, 0)

		for _, message := range messages {
			events = append(events, &cloudwatchlogs.InputLogEvent{
				Timestamp: aws.Int64(message.Time.UnixNano() / (int64)(time.Millisecond)),
				Message:   aws.String(c.formatMessageText(message)),
			})
		}

		// get token again if saving log from another container
		describeParams := &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(c.group),
			LogStreamNamePrefix: aws.String(c.stream),
		}

		data, _ := c.client.DescribeLogStreamsWithContext(ctx, describeParams)
		if len(data.LogStreams) > 0 {
			if data.LogStreams[0].UploadSequenceToken != nil {
				c.lastToken = *data.LogStreams[0].UploadSequenceToken
			}
		}
		var token *string = nil
		if c.lastToken != "" {
			token = &c.lastToken
		}

		params := &cloudwatchlogs.PutLogEventsInput{
			LogEvents:     events,
			LogGroupName:  aws.String(c.group),
			LogStreamName: aws.String(c.stream),
			SequenceToken: token,
		}

		putRes, putErr := c.client.PutLogEvents(params)
		if putErr != nil {
			if c.logger != nil {
				c.logger.Error(cctx.NewContextWithTraceId(ctx, "cloudwatch_logger"), putErr, "putLogEvents error")
			}
		} else {
			if putRes.NextSequenceToken != nil {
				c.lastToken = *putRes.NextSequenceToken
			}
		}
	}()

	wg.Wait()

	return nil
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
