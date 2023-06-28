package count

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	awsconn "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/connect"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// Performance counters that periodically dumps counters to AWS Cloud Watch Metrics.
//
// # Configuration parameters
//
// - connections:
//   - discovery_key:         (optional) a key to retrieve the connection from IDiscovery
//   - region:                (optional) AWS region
//
// - credentials:
//   - store_key:             (optional) a key to retrieve the credentials from ICredentialStore
//   - access_id:             AWS access/client id
//   - access_key:            AWS access/client id
//
// - options:
//   - interval:              interval in milliseconds to save current counters measurements (default: 5 mins)
//   - reset_timeout:         timeout in milliseconds to reset the counters. 0 disables the reset (default: 0)
//
// References ###
//
// - \*:context-info:\*:\*:1.0      (optional) ContextInfo to detect the context id and specify counters source
// - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connections
// - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//
// See Counter (in the Pip.Services components package)
// See CachedCounters (in the Pip.Services components package)
// See CompositeLogger (in the Pip.Services components package)
//
// ### Example ###
//
//		  ctx := context.Background()
//	   counters := NewCloudWatchCounters()
//	   counters.Configure(ctx, config.NewConfigParamsFromTuples(
//	       "connection.region", "us-east-1",
//	       "connection.access_id", "XXXXXXXXXXX",
//	       "connection.access_key", "XXXXXXXXXXX"
//	   ))
//	   counters.SetReferences(ctx, NewReferencesFromTuples(
//	       NewDescriptor("pip-services", "logger", "console", "default", "1.0"),
//	       NewConsoleLogger()
//	   ))
//
//	   err := counters.Open(ctx, "123")
//	       ...
//
//	   counters.Increment(ctx, "mycomponent.mymethod.calls")
//	   timing := counters.BeginTiming(ctx, "mycomponent.mymethod.exec_time")
//
//	       ...
//
//	   timing.EndTiming(ctx, err)
//
//	   counters.Dump(ctx)
type CloudWatchCounters struct {
	*ccount.CachedCounters
	logger *cexec.CompositeLogger

	connectionResolver *awsconn.AwsConnectionResolver
	connection         *awsconn.AwsConnectionParams
	connectTimeout     int
	client             *cloudwatch.CloudWatch //AmazonCloudWatchClient
	source             string
	instance           string
	opened             bool
}

// Creates a new instance of this counters.
func NewCloudWatchCounters() *CloudWatchCounters {
	c := &CloudWatchCounters{
		logger:             cexec.NewCompositeLogger(),
		connectionResolver: awsconn.NewAwsConnectionResolver(),
		connectTimeout:     30000,
		opened:             false,
	}
	c.CachedCounters = ccount.InheritCacheCounters(c)
	return c
}

// Configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *CloudWatchCounters) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.CachedCounters.Configure(ctx, config)
	c.connectionResolver.Configure(ctx, config)

	c.source = config.GetAsStringWithDefault("source", c.source)
	c.instance = config.GetAsStringWithDefault("instance", c.instance)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
}

// Sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- references 	references to locate the component dependencies.
//
// See IReferences (in the Pip.Services commons package)
func (c *CloudWatchCounters) SetReferences(ctx context.Context, references cref.IReferences) {
	c.logger.SetReferences(ctx, references)
	c.connectionResolver.SetReferences(ctx, references)
	ref := references.GetOneOptional(
		cref.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))
	contextInfo, ok := ref.(*cctx.ContextInfo)

	if ok && c.source == "" {
		c.source = contextInfo.Name
	}

	if ok && c.instance == "" {
		c.instance = contextInfo.ContextId
	}

}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *CloudWatchCounters) IsOpen() bool {
	return c.opened
}

// Opens the component.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- Returns          error or null no errors occured.
func (c *CloudWatchCounters) Open(ctx context.Context) error {
	if c.opened {
		return nil
	}

	c.opened = true

	wg := sync.WaitGroup{}
	var errGlobal error

	wg.Add(1)
	go func() {
		defer wg.Done()
		connection, err := c.connectionResolver.Resolve(ctx)
		c.connection = connection
		errGlobal = err

		awsCred := credentials.NewStaticCredentials(c.connection.GetAccessId(), c.connection.GetAccessKey(), "")
		sess := session.Must(session.NewSession(&aws.Config{
			MaxRetries:  aws.Int(3),
			Region:      aws.String(c.connection.GetRegion()),
			Credentials: awsCred,
		}))
		// Create new cloudwatch client.
		c.client = cloudwatch.New(sess)
		c.client.APIVersion = "2010-08-01"
		c.client.Config.HTTPClient.Timeout = time.Duration((int64)(c.connectTimeout)) * time.Millisecond

	}()
	wg.Wait()
	if errGlobal != nil {
		c.opened = false
		return errGlobal
	}
	return nil
}

// Closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- Return         error or nil no errors occured.
func (c *CloudWatchCounters) Close(ctx context.Context) error {
	c.opened = false
	c.client = nil
	return nil
}

func (c *CloudWatchCounters) getCounterData(counter ccount.Counter, now time.Time, dimensions []*cloudwatch.Dimension) *cloudwatch.MetricDatum {

	value := &cloudwatch.MetricDatum{
		MetricName: aws.String(counter.Name),
		Unit:       aws.String(None),
		Dimensions: dimensions,
	}
	tm := counter.Time
	if tm.IsZero() {
		tm = time.Now().UTC()
	}
	value.SetTimestamp(tm)

	switch counter.Type {
	case ccount.Increment:
		value.Value = aws.Float64((float64)(counter.Count))
		value.Unit = aws.String(Count)
		break
	case ccount.Interval:
		value.Unit = aws.String(Milliseconds)
		//value.Value = counter.average;
		value.StatisticValues = &cloudwatch.StatisticSet{
			SampleCount: aws.Float64((float64)(counter.Count)),
			Maximum:     aws.Float64((float64)(counter.Max)),
			Minimum:     aws.Float64((float64)(counter.Min)),
			Sum:         aws.Float64((float64)(counter.Count) * (float64)(counter.Average)),
		}
		break
	case ccount.Statistics:
		//value.Value = counter.average;
		value.StatisticValues = &cloudwatch.StatisticSet{
			SampleCount: aws.Float64((float64)(counter.Count)),
			Maximum:     aws.Float64((float64)(counter.Max)),
			Minimum:     aws.Float64((float64)(counter.Min)),
			Sum:         aws.Float64((float64)(counter.Count) * (float64)(counter.Average)),
		}
		break
	case ccount.LastValue:
		value.Value = aws.Float64((float64)(counter.Last))
		break
	case ccount.Timestamp:
		value.Value = aws.Float64((float64)(counter.Time.UnixNano()) / (float64)(time.Millisecond)) // Convert to milliseconds UnixTimeStamp
		break
	}

	return value
}

// Saves the current counters measurements.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- counters      current counters measurements to be saves.
func (c *CloudWatchCounters) Save(ctx context.Context, counters []ccount.Counter) error {
	if c.client == nil {
		return nil
	}

	var dimensions []*cloudwatch.Dimension
	dimensions = make([]*cloudwatch.Dimension, 0)
	dimensions = append(dimensions, &cloudwatch.Dimension{
		Name:  aws.String("InstanceID"),
		Value: aws.String(c.instance),
	})

	now := time.Now()

	var data []*cloudwatch.MetricDatum
	data = make([]*cloudwatch.MetricDatum, 0)

	params := &cloudwatch.PutMetricDataInput{
		MetricData: data,
		Namespace:  aws.String(c.source),
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for _, counter := range counters {
			data = append(data, c.getCounterData(counter, now, dimensions))

			if len(data) >= 20 {
				params.MetricData = data
				_, err := c.client.PutMetricDataWithContext(ctx, params)
				if err != nil && c.logger != nil {
					c.logger.Error(cctx.NewContextWithTraceId(ctx, "cloudwatch_counters"), err, "putMetricData error")
				}
				data = make([]*cloudwatch.MetricDatum, 0)
			}
		}

	}()

	wg.Wait()

	params.MetricData = data

	if len(data) > 0 {
		_, err := c.client.PutMetricDataWithContext(ctx, params)
		if err != nil && c.logger != nil {
			c.logger.Error(cctx.NewContextWithTraceId(ctx, "cloudwatch_counters"), err, "putMetricData error")
		}
	}
	return nil
}
