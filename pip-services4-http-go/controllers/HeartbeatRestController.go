package controllers

import (
	"context"
	"net/http"
	"time"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

// HeartbeatRestController service returns heartbeat via HTTP/REST protocol.
// The service responds on /heartbeat route (can be changed)
// with a string with the current time in UTC.
// This service route can be used to health checks by load-balancers and
// container orchestrators.
//
//	Configuration parameters:
//		- baseroute:           base route for remote URI (default: "")
//		- route:               route to heartbeat operation (default: "heartbeat")
//		- dependencies:
//			- endpoint:        override for HTTP Endpoint dependency
//		- connection(s):
//			- discovery_key:   (optional) a key to retrieve the connection from IDiscovery
//			- protocol:        connection protocol: http or https
//			- host:            host name or IP address
//			- port:            port number
//			- uri:             resource URI or connection string with all parameters in it
//
//	References:
//		- *:logger:*:*:1.0       (optional)  ILogger components to pass log messages
//		- *:counters:*:*:1.0     (optional)  ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0    (optional)  IDiscovery services to resolve connection
//		- *:endpoint:http:*:1.0  (optional) HttpEndpoint reference
//
//	see RestController
//	see clients.RestClient
//
//	Example:
//		controller := NewHeartbeatController();
//		controller.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
//			"route", "ping",
//			"connection.protocol", "http",
//			"connection.host", "localhost",
//			"connection.port", 8080,
//		));
//
//		opnErr := controller.Open(context.Background())
//		if opnErr == nil {
//			fmt.Println("The Heartbeat controller is accessible at http://+:8080/ping");
//		}
type HeartbeatRestController struct {
	*RestController
	route string
}

// NewHeartbeatRestController creates a new instance of c service.
func NewHeartbeatRestController() *HeartbeatRestController {
	c := &HeartbeatRestController{}
	c.RestController = InheritRestController(c)
	c.route = "heartbeat"
	return c
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config configuration parameters to be set.
func (c *HeartbeatRestController) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.RestController.Configure(ctx, config)
	c.route = config.GetAsStringWithDefault("route", c.route)
}

// Register all service routes in HTTP endpoint.
func (c *HeartbeatRestController) Register() {
	c.RegisterRoute(http.MethodGet, c.route, nil, func(res http.ResponseWriter, req *http.Request) { c.heartbeat(req, res) })
}

// Handles heartbeat requests
//
//	Parameters:
//		- req   an HTTP request
//		- res   an HTTP response
func (c *HeartbeatRestController) heartbeat(req *http.Request, res http.ResponseWriter) {
	c.SendResult(res, req, time.Now(), nil)
}
