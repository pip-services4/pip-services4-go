package controllers

import (
	"context"
	"net/http"
	"time"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// StatusRestController is a service that returns microservice status information via HTTP/REST protocol.
//
//	The service responds on /status route (can be changed) with a JSON object:
//		{
//			"id":            unique container id (usually hostname)
//			"name":          container name (from ContextInfo)
//			"description":   container description (from ContextInfo)
//			"start_time":    time when container was started
//			"current_time":  current time in UTC
//			"uptime":        duration since container start time in milliseconds
//			"properties":    additional container properties (from ContextInfo)
//			"components":    descriptors of components registered in the container
//		}
//
//	Configuration parameters:
//		- baseroute:          base route for remote URI
//		- route:              status route (default: "status")
//		- dependencies:
//			- endpoint:       override for HTTP Endpoint dependency
//			- controller:     override for Controller dependency
//		- connection(s):
//			- discovery_key:  (optional) a key to retrieve the connection from IDiscovery
//			- protocol:       connection protocol: http or https
//			- host:           host name or IP address
//			- port:           port number
//			- uri:            resource URI or connection string with all parameters in it
//
//	References:
//		- *:logger:*:*:1.0       (optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0     (optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0    (optional) IDiscovery services to resolve connection
//		- *:endpoint:http:*:1.0  (optional) HttpEndpoint reference
//
//	see: RestController
//	see: clients.RestClient
//
//	Example:
//		service = NewStatusController();
//		service.Configure(context.Background(), cref.NewConfigParamsFromTuples(
//			"connection.protocol", "http",
//			"connection.host", "localhost",
//			"connection.port", 8080,
//		));
//
//		opnErr:= service.Open(context.Background(), "123")
//		if opnErr == nil {
//			fmt.Println("The Status service is accessible at http://localhost:8080/status");
//		}
type StatusRestController struct {
	*RestController
	startTime   time.Time
	references2 crefer.IReferences
	contextInfo *cctx.ContextInfo
	route       string
}

// NewStatusRestController method are creates a new instance of this service.
func NewStatusRestController() *StatusRestController {
	c := &StatusRestController{}
	c.RestController = InheritRestController(c)
	c.startTime = time.Now()
	c.route = "status"
	c.DependencyResolver.Put(
		context.Background(),
		"context-info",
		crefer.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"),
	)
	return c
}

// Configure method are configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams  configuration parameters to be set.
func (c *StatusRestController) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.RestController.Configure(ctx, config)
	c.route = config.GetAsStringWithDefault("route", c.route)
}

// SetReferences method are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references crefer.IReferences	references to locate the component dependencies.
func (c *StatusRestController) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.references2 = references
	c.RestController.SetReferences(ctx, references)

	depRes := c.DependencyResolver.GetOneOptional("context-info")
	if depRes != nil {
		if _val, ok := depRes.(*cctx.ContextInfo); ok {
			c.contextInfo = _val
		}
	}
}

// Register method are registers all service routes in HTTP endpoint.
func (c *StatusRestController) Register() {
	c.RegisterRoute(http.MethodGet, c.route, nil, c.status)
}

// Handles status requests
//
//	Parameters:
//		- req  *http.Request an HTTP request
//		- res  http.ResponseWriter  an HTTP response
func (c *StatusRestController) status(res http.ResponseWriter, req *http.Request) {

	id := ""
	if c.contextInfo != nil {
		id = c.contextInfo.ContextId
	}

	name := "Unknown"
	if c.contextInfo != nil {
		name = c.contextInfo.Name
	}

	description := ""
	if c.contextInfo != nil {
		description = c.contextInfo.Description
	}

	uptime := time.Since(c.startTime)

	properties := make(map[string]string, 0)
	if c.contextInfo != nil {
		properties = c.contextInfo.Properties
	}

	var components []string
	if c.references2 != nil {
		for _, locator := range c.references2.GetAllLocators() {
			components = append(components, cconv.StringConverter.ToString(locator))
		}
	}

	status := make(map[string]any)

	status["id"] = id
	status["name"] = name
	status["description"] = description
	status["start_time"] = cconv.StringConverter.ToString(c.startTime)
	status["current_time"] = cconv.StringConverter.ToString(time.Now())
	status["uptime"] = uptime
	status["properties"] = properties
	status["components"] = components

	c.SendResult(res, req, status, nil)
}
