package count

import (
	"context"
	"os"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	clients1 "github.com/pip-services4/pip-services4-go/pip-services4-datadog-go/clients"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// Performance counters that send their metrics to DataDog service.
//
// DataDog is a popular monitoring SaaS service. It collects logs, metrics, events
// from infrastructure and applications and analyze them in a single place.
//
// ### Configuration parameters ###
//
// - connection(s):
//   - discovery_key:         (optional) a key to retrieve the connection from [[IDiscovery]]
//   - protocol:            (optional) connection protocol: http or https (default: https)
//   - host:                (optional) host name or IP address (default: api.datadoghq.com)
//   - port:                (optional) port number (default: 443)
//   - uri:                 (optional) resource URI or connection string with all parameters in it
//
// - credential:
//   - access_key:          DataDog client api key
//
// - options:
//   - retries:               number of retries (default: 3)
//   - connect_timeout:       connection timeout in milliseconds (default: 10 sec)
//   - timeout:               invocation timeout in milliseconds (default: 10 sec)
//
// ### References ###
//
// - \*:logger:\*:\*:1.0         (optional)  [[ILogger]] components to pass log messages
// - \*:counters:\*:\*:1.0         (optional) [[ICounters]] components to pass collected measurements
// - \*:discovery:\*:\*:1.0        (optional) [[IDiscovery]] services to resolve connection
//
// See [[RestService]]
// See [[CommandableHttpService]]
//
// Example:
//
//	counters := NewDataDogCounters();
//	counters.Configure(context.Background(), NewConfigParamsFromTuples(
//	    "credential.access_key", "827349874395872349875493"
//	))
//
//	err := counters.Open(context.Background(),"123")
//	    ...
//
//	counters.Increment(context.Background(), "mycomponent.mymethod.calls")
//	timing := counters.BeginTiming(context.Background(), "mycomponent.mymethod.exec_time")
//
//	  ...
//
//	timing.EndTiming(context.Background())
//
//	counters.Dump(context.Background())
type DataDogCounters struct {
	*ccount.CachedCounters
	client       *clients1.DataDogMetricsClient
	logger       *clog.CompositeLogger
	opened       bool
	source       string
	instance     string
	requestRoute string
}

// NewDataDogCounters - creates a new instance of the performance counters.
func NewDataDogCounters() *DataDogCounters {
	c := &DataDogCounters{
		client: clients1.NewDataDogMetricsClient(nil),
		logger: clog.NewCompositeLogger(),
		opened: false,
	}
	c.CachedCounters = ccount.InheritCacheCounters(c)
	c.instance, _ = os.Hostname()
	return c
}

// Configure - configures component by passing configuration parameters.
//   - config    configuration parameters to be set.
func (c *DataDogCounters) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.CachedCounters.Configure(ctx, config)
	c.client.Configure(ctx, config)

	c.source = config.GetAsStringWithDefault("source", c.source)
	c.instance = config.GetAsStringWithDefault("instance", c.instance)
}

// SetReferences - sets references to dependent components.
// - references 	references to locate the component dependencies.
func (c *DataDogCounters) SetReferences(ctx context.Context, references cref.IReferences) {
	c.logger.SetReferences(ctx, references)
	c.client.SetReferences(ctx, references)
	ref := references.GetOneOptional(cref.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))

	contextInfo, _ := ref.(*cctx.ContextInfo)

	if contextInfo != nil && c.source == "" {
		c.source = contextInfo.Name
	}
	if contextInfo != nil && c.instance == "" {
		c.instance = contextInfo.ContextId
	}
}

// IsOpen - checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *DataDogCounters) IsOpen() bool {
	return c.opened
}

// Open - opens the component.
//   - ctx context.Context execution context to trace execution through call chain.
//     Returns  error or nil no errors occured.
func (c *DataDogCounters) Open(ctx context.Context) error {
	if c.opened {
		return nil
	}

	c.opened = true
	return c.client.Open(ctx)
}

// Close - closes component and frees used resources.
//   - ctx context.Context execution context to trace execution through call chain.
//     Returns  error or null no errors occured.
func (c *DataDogCounters) Close(ctx context.Context) error {
	c.opened = false

	return c.client.Close(ctx)
}

func (c *DataDogCounters) convertCounter(counter ccount.Counter) []clients1.DataDogMetric {
	switch counter.Type {
	case ccount.Increment:
		return []clients1.DataDogMetric{{
			Metric:  counter.Name,
			Type:    clients1.Gauge,
			Host:    c.instance,
			Service: c.source,
			Points:  []clients1.DataDogMetricPoint{{Time: counter.Time, Value: (float64)(counter.Count)}},
		}}

	case ccount.LastValue:
		return []clients1.DataDogMetric{{
			Metric:  counter.Name,
			Type:    clients1.Gauge,
			Host:    c.instance,
			Service: c.source,
			Points:  []clients1.DataDogMetricPoint{{Time: counter.Time, Value: (float64)(counter.Last)}},
		}}

	case ccount.Interval:
	case ccount.Statistics:
		return []clients1.DataDogMetric{
			{
				Metric:  counter.Name + ".min",
				Type:    clients1.Gauge,
				Host:    c.instance,
				Service: c.source,
				Points:  []clients1.DataDogMetricPoint{{Time: counter.Time, Value: (float64)(counter.Min)}},
			},
			{
				Metric:  counter.Name + ".average",
				Type:    clients1.Gauge,
				Host:    c.instance,
				Service: c.source,
				Points:  []clients1.DataDogMetricPoint{{Time: counter.Time, Value: (float64)(counter.Average)}},
			},
			{
				Metric:  counter.Name + ".max",
				Type:    clients1.Gauge,
				Host:    c.instance,
				Service: c.source,
				Points:  []clients1.DataDogMetricPoint{{Time: counter.Time, Value: (float64)(counter.Max)}},
			},
		}
	}

	return nil
}

func (c *DataDogCounters) convertCounters(counters []ccount.Counter) []clients1.DataDogMetric {
	metrics := make([]clients1.DataDogMetric, 0)

	for _, counter := range counters {
		data := c.convertCounter(counter)

		if len(data) > 0 {
			metrics = append(metrics, data...)
		}
	}

	return metrics
}

// Saves the current counters measurements.
// - counters      current counters measurements to be saves.
func (c *DataDogCounters) Save(ctx context.Context, counters []ccount.Counter) error {
	metrics := c.convertCounters(counters)
	if len(metrics) == 0 {
		return nil
	}

	err := c.client.SendMetrics(cctx.NewContextWithTraceId(ctx, "datadog-counters"), metrics)
	if err != nil {
		c.logger.Error(cctx.NewContextWithTraceId(ctx, "datadog-counters"), err, "Failed to push metrics to DataDog")
	}
	return err
}
