package controllers

import (
	"context"
	"io"
	"net/http"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	pcount "github.com/pip-services4/pip-services4-go/pip-services4-prometheus-go/count"
)

// PrometheusMetricsController is controller that exposes "/metrics" route for Prometheus to scap performance metrics.
//
//	Configuration parameters:
//
//		- dependencies:
//			- endpoint:              override for HTTP Endpoint dependency
//			- prometheus-counters:   override for PrometheusCounters dependency
//		- connection(s):
//			- discovery_key:         (optional) a key to retrieve the connection from IDiscovery
//			- protocol:              connection protocol: http or https
//			- host:                  host name or IP address
//			- port:                  port number
//			- uri:                   resource URI or connection string with all parameters in it
//
//	References:
//
//		- *:logger:*:*:1.0         (optional)  ILogger components to pass log messages
//		- *:counters:*:*:1.0         (optional)  ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0        (optional)  IDiscovery services to resolve connection
//		- *:endpoint:http:*:1.0          (optional)  HttpEndpoint reference to expose REST operation
//		- *:counters:prometheus:*:1.0    PrometheusCounters reference to retrieve collected metrics
//
// See RestController
// See RestClient
//
// Example
//
//	ctx := context.Context()
//	service := NewPrometheusMetricsController();
//	service.Configure(ctx, cconf.NewConfigParamsFromTuples(
//	    "connection.protocol", "http",
//	    "connection.host", "localhost",
//	    "connection.port", "8080",
//	));
//
//	err := service.Open(ctx)
//	if  err == nil {
//	    fmt.Println("The Prometheus metrics service is accessible at http://localhost:8080/metrics");
//	    defer service.Close(ctx)
//	}
type PrometheusMetricsController struct {
	httpctrl.RestController
	cachedCounters *ccount.CachedCounters
	source         string
	instance       string
}

// NewPrometheusMetricsController are creates a new instance of c service.
// Returns *PrometheusMetricsController
// pointer on new instance
func NewPrometheusMetricsController() *PrometheusMetricsController {
	c := &PrometheusMetricsController{}
	c.RestController = *httpctrl.InheritRestController(c)
	c.DependencyResolver.Put(context.Background(), "cached-counters", cref.NewDescriptor("pip-services", "counters", "cached", "*", "1.0"))
	c.DependencyResolver.Put(context.Background(), "prometheus-counters", cref.NewDescriptor("pip-services", "counters", "prometheus", "*", "1.0"))
	return c
}

// SetReferences is sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- references cref.IReferences
//
// references to locate the component dependencies.
func (c *PrometheusMetricsController) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestController.SetReferences(ctx, references)

	resolv := c.DependencyResolver.GetOneOptional("prometheus-counters")
	c.cachedCounters = resolv.(*pcount.PrometheusCounters).CachedCounters
	if c.cachedCounters == nil {
		resolv = c.DependencyResolver.GetOneOptional("cached-counters")
		c.cachedCounters = resolv.(*ccount.CachedCounters)
	}
	ref := references.GetOneOptional(
		cref.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))
	contextInfo := ref.(*cctx.ContextInfo)

	if contextInfo != nil && c.source == "" {
		c.source = contextInfo.Name
	}
	if contextInfo != nil && c.instance == "" {
		c.instance = contextInfo.ContextId
	}
}

// Register method are registers all service routes in HTTP endpoint.
func (c *PrometheusMetricsController) Register() {
	c.RegisterRoute("get", "metrics", nil, func(res http.ResponseWriter, req *http.Request) { c.metrics(res, req) })
}

// Handles metrics requests
//
//	Parameters:
//		- req   an HTTP request
//		- res   an HTTP response
func (c *PrometheusMetricsController) metrics(res http.ResponseWriter, req *http.Request) {

	var atomicCounters []*ccount.AtomicCounter
	if c.cachedCounters != nil {
		atomicCounters = c.cachedCounters.GetAll()
	}

	counters := pcount.PrometheusCounterConverter.AtomicCountersToCounters(atomicCounters)
	body := pcount.PrometheusCounterConverter.ToString(counters, c.source, c.instance)

	res.Header().Add("content-type", "text/plain")
	res.WriteHeader(200)
	_, wrErr := io.WriteString(res, (string)(body))
	if wrErr != nil {
		c.Logger.Error(cctx.NewContextWithTraceId(req.Context(), "PrometheusMetricsController"), wrErr, "Can't write response")
	}
}
