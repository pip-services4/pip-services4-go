package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	"github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/trace"
	"goji.io/pat"
)

// RestController Abstract service that receives remove calls via HTTP/REST protocol.
//
//	Configuration parameters:
//		- base_route:              base route for remote URI
//		- dependencies:
//			- endpoint:            override for HTTP Endpoint dependency
//			- controller:          override for Controller dependency
//		- connection(s):
//			- discovery_key:       (optional) a key to retrieve the connection from IDiscovery
//			- protocol:            connection protocol: http or https
//			- host:                host name or IP address
//			- port:                port number
//			- uri:                 resource URI or connection string with all parameters in it
//		- credential - the HTTPS credentials:
//			- ssl_key_file:        the SSL private key in PEM
//			- ssl_crt_file:        the SSL certificate in PEM
//			- ssl_ca_file:         the certificate authorities (root cerfiticates) in PEM
//
//	References:
//		- *:logger:*:*:1.0         (optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0       (optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0      (optional) IDiscovery services to resolve connection
//		- *:endpoint:http:*:1.0    (optional) HttpEndpoint reference
//
//	See clients.RestClient
//
//	Example:
//		type MyRestController struct {
//			*RestController
//			service IMyService
//		}
//		...
//		func NewMyRestController() *MyRestController {
//			c := MyRestController{}
//			c.RestController = services.NewRestController()
//			c.RestController.IRegisterable = &c
//			c.numberOfCalls = 0
//			c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("mygroup", "service", "*", "*", "1.0"))
//			return &c
//		}
//
//		func (c * MyRestController) SetReferences(ctx context.Context, references IReferences) {
//			c.RestController.SetReferences(ctx, references);
//			resolv := c.DependencyResolver.GetRequired("service");
//			if resolv != nil {
//				c.service, _ = resolv.(IMyService)
//			}
//		}
//
//		func (c *MyRestController) getOneById(res http.ResponseWriter, req *http.Request) {
//			params := req.URL.Query()
//			vars := mux.Vars(req)
//
//			mydataId := params.Get("mydata_id")
//			if mydataId == "" {
//				mydataId = vars["mydatay_id"]
//			}
//			result, err := c.controller.GetOneById(params.Get("tace_id"), mydataId),
//			c.SendResult(res, req, result, err)
//		}
//
//		func (c * MyRestController) Register() {
//			c.RegisterRoute(
//				"get", "get_mydata/{mydata_id}",
//				&cvalid.NewObjectSchema().
//					WithRequiredProperty("mydata_id", cconv.String).Schema,
//				c.getOneById,
//			)
//			...
//		}
//
//
//		controller := NewMyRestController();
//		controller.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
//			"connection.protocol", "http",
//			"connection.host", "localhost",
//			"connection.port", 8080,
//		));
//		controller.SetReferences(context.Background(), cref.NewReferencesFromTuples(
//			cref.NewDescriptor("mygroup","service","default","default","1.0"), service
//		));
//
//		opnRes := controller.Open(context.Background(), "123")
//		if opnErr == nil {
//			fmt.Println("The REST cntroller is running on port 8080");
//		}
type RestController struct {
	Overrides IRegisterable

	defaultConfig *cconf.ConfigParams
	config        *cconf.ConfigParams
	references    crefer.IReferences
	localEndpoint bool
	opened        bool
	//The base route.
	BaseRoute string
	//The HTTP endpoint that exposes this service.
	Endpoint *HttpEndpoint
	//The dependency resolver.
	DependencyResolver *crefer.DependencyResolver
	//The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer

	SwaggerController ISwaggerController
	SwaggerEnabled    bool
	SwaggerRoute      string
}

// InheritRestController creates new instance of RestController
func InheritRestController(overrides IRegisterable) *RestController {
	rs := RestController{
		Overrides: overrides,
	}
	rs.defaultConfig = cconf.NewConfigParamsFromTuples(
		"base_route", "",
		"dependencies.endpoint", "*:endpoint:http:*:1.0",
		"dependencies.swagger", "*:swagger-controller:*:*:1.0",
	)
	rs.DependencyResolver = crefer.NewDependencyResolver()
	rs.DependencyResolver.Configure(context.TODO(), rs.defaultConfig)
	rs.Logger = clog.NewCompositeLogger()
	rs.Counters = ccount.NewCompositeCounters()
	rs.Tracer = ctrace.NewCompositeTracer()
	rs.SwaggerEnabled = false
	rs.SwaggerRoute = "swagger"
	return &rs
}

// Configure method are configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config  *cconf.ConfigParams  configuration parameters to be set.
func (c *RestController) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config
	c.DependencyResolver.Configure(ctx, config)
	c.BaseRoute = config.GetAsStringWithDefault("base_route", c.BaseRoute)
	c.SwaggerEnabled = config.GetAsBooleanWithDefault("swagger.enable", c.SwaggerEnabled)
	c.SwaggerRoute = config.GetAsStringWithDefault("swagger.route", c.SwaggerRoute)
}

// SetReferences method are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references crefer.IReferences	references to locate the component dependencies.
func (c *RestController) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.references = references

	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.Tracer.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)

	// Get endpoint
	depRes := c.DependencyResolver.GetOneOptional("endpoint")
	if depRes != nil {
		c.Endpoint = depRes.(*HttpEndpoint)
	}

	// Or create a local one
	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint(ctx)
		c.localEndpoint = true
	} else {
		c.localEndpoint = false
	}
	// Add registration callback to the endpoint
	c.Endpoint.Register(c)

	depRes = c.DependencyResolver.GetOneOptional("swagger")
	if depRes != nil {
		if _val, ok := depRes.(ISwaggerController); ok {
			c.SwaggerController = _val
		}
	}
}

// UnsetReferences method are unsets (clears) previously set references to dependent components.
func (c *RestController) UnsetReferences() {
	// Remove registration callback from endpoint
	if c.Endpoint != nil {
		c.Endpoint.Unregister(c)
		c.Endpoint = nil
	}
	c.SwaggerController = nil
}

func (c *RestController) createEndpoint(ctx context.Context) *HttpEndpoint {
	endpoint := NewHttpEndpoint()

	if c.config != nil {
		endpoint.Configure(ctx, c.config)
	}
	if c.references != nil {
		endpoint.SetReferences(ctx, c.references)
	}

	return endpoint
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- name              a method name.
//	Returns: Timing object to end the time measurement.
func (c *RestController) Instrument(ctx context.Context, name string) *trace.InstrumentTiming {
	c.Logger.Trace(ctx, "Executing %s method", name)
	c.Counters.IncrementOne(ctx, name+".exec_count")

	counterTiming := c.Counters.BeginTiming(ctx, name+".exec_time")
	traceTiming := c.Tracer.BeginTrace(ctx, name, "")

	return trace.NewInstrumentTiming(ctx, name, "exec",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// InstrumentError method are adds instrumentation to error handling.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- name    string        a method name.
//		- err     error         an occurred error
//		- result  any			(optional) an execution result
//	Returns: result any, err error (optional) an execution callback
func (c *RestController) InstrumentError(ctx context.Context, name string, errIn error,
	resIn any) (result any, err error) {

	if errIn != nil {
		c.Logger.Error(ctx, errIn, "Failed to execute %s method", name)
		c.Counters.IncrementOne(ctx, name+".exec_errors")
	}
	return resIn, errIn
}

// IsOpen method checks if the component is opened.
//
//	Returns: true if the component has been opened and false otherwise.
func (c *RestController) IsOpen() bool {
	return c.opened
}

// Open method are opens the component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error or nil no errors occured.
func (c *RestController) Open(ctx context.Context) error {
	if c.opened {
		return nil
	}

	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint(ctx)
		c.Endpoint.Register(c)
		c.localEndpoint = true
	}

	if c.localEndpoint {
		oErr := c.Endpoint.Open(ctx)
		if oErr != nil {
			c.opened = false
			return oErr
		}
	}
	c.opened = true
	return nil
}

// Close method are closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error or nil no errors occurred.
func (c *RestController) Close(ctx context.Context) error {
	if !c.opened {
		return nil
	}

	if c.Endpoint == nil {
		return cerr.NewInvalidStateError(utils.ContextHelper.GetTraceId(ctx), "NO_ENDPOINT", "HTTP endpoint is missing")
	}

	if c.localEndpoint {
		cErr := c.Endpoint.Close(ctx)
		if cErr != nil {
			c.opened = false
			return cErr
		}
	}
	c.opened = false
	return nil
}

// SendResult method are sends result as JSON object.
// That function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 200 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with approproate status code.
//
//	Parameters:
//		- req       a HTTP request object.
//		- res       a HTTP response object.
//		- result    (optional) result object to send
//		- err error (optional) error objrct to send
func (c *RestController) SendResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	HttpResponseSender.SendResult(res, req, result, err)
}

// SendCreatedResult method are sends newly created object as JSON.
// That callback function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 201 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with approproate status code.
//
//	Parameters:
//		- req       a HTTP request object.
//		- res       a HTTP response object.
//		- result    (optional) result object to send
//		- err error (optional) error objrct to send
func (c *RestController) SendCreatedResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	HttpResponseSender.SendCreatedResult(res, req, result, err)
}

// SendDeletedResult method are sends deleted object as JSON.
// That callback function call be called directly or passed
// as a parameter to business logic components.
// If object is not nil it returns 200 status code.
// For nil results it returns 204 status code.
// If error occur it sends ErrorDescription with appropriate status code.
//
//	Parameters:
//		- req       a HTTP request object.
//		- res       a HTTP response object.
//		- result    (optional) result object to send
//		- err error (optional) error objrct to send
func (c *RestController) SendDeletedResult(res http.ResponseWriter, req *http.Request, result any, err error) {
	HttpResponseSender.SendDeletedResult(res, req, result, err)
}

// SendError method are sends error serialized as ErrorDescription object
// and appropriate HTTP status code.
// If status code is not defined, it uses 500 status code.
//
//	Parameters:
//		- req       a HTTP request object.
//		- res       a HTTP response object.
//		- error     an error object to be sent.
func (c *RestController) SendError(res http.ResponseWriter, req *http.Request, err error) {
	HttpResponseSender.SendError(res, req, err)
}

func (c *RestController) appendBaseRoute(route string) string {

	if route == "" {
		route = "/"
	}

	if c.BaseRoute != "" && len(c.BaseRoute) > 0 {
		baseRoute := c.BaseRoute
		if len(route) == 0 {
			route = "/"
		}
		if route[0] != '/' {
			route = "/" + route
		}
		if baseRoute[0] != '/' {
			baseRoute = "/" + baseRoute
		}
		route = baseRoute + route
	}
	return route

}

// RegisterRoute method are registers a route in HTTP endpoint.
//
//	Parameters:
//		- method        HTTP method: "get", "head", "post", "put", "delete"
//		- route         a command route. Base route will be added to this route
//		- schema        a validation schema to validate received parameters.
//		- action        an action function that is called when operation is invoked.
func (c *RestController) RegisterRoute(method string, route string, schema *cvalid.Schema,
	action func(res http.ResponseWriter, req *http.Request)) {

	if c.Endpoint == nil {
		return
	}
	route = c.appendBaseRoute(route)
	c.Endpoint.RegisterRoute(method, route, schema, action)
}

// RegisterRouteWithAuth method are registers a route with authorization in HTTP endpoint.
//
//	Parameters:
//		- method        HTTP method: "get", "head", "post", "put", "delete"
//		- route         a command route. Base route will be added to this route
//		- schema        a validation schema to validate received parameters.
//		- authorize     an authorization interceptor
//		- action        an action function that is called when operation is invoked.
func (c *RestController) RegisterRouteWithAuth(method string, route string, schema *cvalid.Schema,
	authorize func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc),
	action func(res http.ResponseWriter, req *http.Request)) {

	if c.Endpoint == nil {
		return
	}
	route = c.appendBaseRoute(route)
	c.Endpoint.RegisterRouteWithAuth(
		method, route, schema,
		func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
			if authorize != nil {
				authorize(res, req, next)
			} else {
				next.ServeHTTP(res, req)
			}
		}, action)
}

// RegisterInterceptor method are registers a middleware for a given route in HTTP endpoint.
//
//	Parameters:
//		- route         a command route. Base route will be added to this route
//		- action        an action function that is called when middleware is invoked.
func (c *RestController) RegisterInterceptor(route string,
	action func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc)) {

	if c.Endpoint == nil {
		return
	}
	route = c.appendBaseRoute(route)
	c.Endpoint.RegisterInterceptor(route, action)
}

// GetParam methods helps get all params from query
//
//	Parameters:
//		- req  incoming request
//		- name parameter name
//	Returns value or empty string if param not exists
func (c *RestController) GetParam(req *http.Request, name string) string {
	param := req.URL.Query().Get(name)
	if param == "" {
		param = pat.Param(req, name)
	}
	return param
}

// DecodeBody methods helps decode body
//
//		Parameters:
//	  - req   	- incoming request
//	  - target  	- pointer on target variable for decode
//
// Returns error
func (c *RestController) DecodeBody(req *http.Request, target any) error {
	bodyBytes, err := io.ReadAll(req.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, target)

	if err != nil {
		return err
	}

	_ = req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

// GetPagingParams methods helps decode paging params
//
//	Parameters:
//		- req  incoming request
//	Returns: paging params
func (c *RestController) GetPagingParams(req *http.Request) *cquery.PagingParams {

	pagingParams := make(map[string]string, 0)
	pagingParams["skip"] = c.GetParam(req, "skip")
	pagingParams["take"] = c.GetParam(req, "take")
	pagingParams["total"] = c.GetParam(req, "total")

	return cquery.NewPagingParamsFromValue(pagingParams)
}

// GetFilterParams methods helps decode filter params
//
//	Parameters:
//		- req incoming request
//	Returns: filter params
func (c *RestController) GetFilterParams(req *http.Request) *cquery.FilterParams {

	params := req.URL.Query()

	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")
	delete(params, "trace_id")

	return cquery.NewFilterParamsFromValue(params)
}

// GetTraceId method returns TraceId from request
//
//	Parameters:
//		- req *http.Request  request
//	Returns: string trace_id or empty string
func (c *RestController) GetTraceId(req *http.Request) string {
	traceId := req.URL.Query().Get("trace_id")
	if traceId == "" {
		traceId = req.Header.Get("trace_id")
	}
	return traceId
}

func (c *RestController) RegisterOpenApiSpecFromFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		c.Logger.Error(
			utils.ContextHelper.NewContextWithTraceId(context.Background(), "RestController"),
			err,
			"Can't read swagger file by path %s",
			path,
		)
		return
	}
	c.RegisterOpenApiSpec((string)(content))
}

func (c *RestController) RegisterOpenApiSpec(content string) {
	if c.SwaggerEnabled {
		c.RegisterRoute(http.MethodGet,
			c.SwaggerRoute, nil, func(res http.ResponseWriter, req *http.Request) {
				res.Header().Add("Content-Length", cconv.StringConverter.ToString(len(content)))
				res.Header().Add("Content-Type", "application/x-yaml")
				res.WriteHeader(200)
				_, _ = io.WriteString(res, content)
			})

		if c.SwaggerController != nil {
			c.SwaggerController.RegisterOpenApiSpec(c.BaseRoute, c.SwaggerRoute)
		}
	}
}

// Register method are registers all service routes in HTTP endpoint.
func (c *RestController) Register() {
	// Override in child classes
	c.Overrides.Register()
}
