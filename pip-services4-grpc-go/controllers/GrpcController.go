package controllers

import (
	"context"
	"encoding/json"
	"strings"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	crun "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	rpctrace "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/trace"
	"google.golang.org/grpc"
)

type IGrpcControllerOverrides interface {
	Register()
}

// GrpcController abstract controller that receives remove calls via GRPC protocol.
//
//	Configuration parameters:
//
//	- dependencies:
//		- endpoint:              override for GRPC Endpoint dependency
//		- service:            override for Controller dependency
//	- connection(s):
//		- discovery_key:         (optional) a key to retrieve the connection from connect.idiscovery.html IDiscovery
//		- protocol:              connection protocol: http or https
//		- host:                  host name or IP address
//		- port:                  port number
//		- uri:                   resource URI or connection string with all parameters in it
//	- credential - the HTTPS credentials:
//		- ssl_key_file:         the SSL private key in PEM
//		- ssl_crt_file:         the SSL certificate in PEM
//		- ssl_ca_file:          the certificate authorities (root cerfiticates) in PEM
//
//	References:
//
//		- *:logger:*:*:1.0               (optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0             (optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0            (optional) IDiscovery controllers to resolve connection
//		- *:endpoint:grpc:*:1.0           (optional) GrpcEndpoint reference
//
// # See GrpcClient
//
// Example:
//
//	type MyGrpcController struct{
//	   *GrpcController
//	   service IMyService;
//	}
//	...
//
//	func NewMyGrpcController() *MyGrpcController {
//		c := MyGrpcController{}
//		c.GrpcController = grpcservices.NewGrpcService("Mydata.Mydatas")
//		c.GrpcController.IRegisterable = &c
//		c.numberOfCalls = 0
//		c.DependencyResolver.Put(context.Context(), "service", cref.NewDescriptor("mygroup", "service", "*", "*", "*"))
//		return &c
//	}
//
//	func (c*MyGrpcService) SetReferences(ctx context.Context, references: IReferences) {
//		c.service.SetReferences(references);
//		resolv, err := c.DependencyResolver.GetOneRequired("service")
//		if err == nil && resolv != nil {
//		    c.service = resolv.(grpctest.IMyService)
//		    return
//		}
//		panic("Can't resolve 'service' reference")
//	}
//
//	func (c*MyGrpcService) Register() {
//		protos.RegisterMyDataServer(c.Endpoint.GetServer(), c)
//		...
//	}
//
//	service := NewMyGrpcService();
//	service.Configure(ctx, cconf.NewConfigParamsFromTuples(
//	    "connection.protocol", "http",
//	    "connection.host", "localhost",
//	    "connection.port", 8080,
//	));
//	service.SetReferences(ctx, cref.NewReferencesFromTuples(
//	   cref.NewDescriptor("mygroup","controller","default","default","1.0"), controller
//	));
//
//	err := service.Open(ctx)
//	if  err == nil {
//	   fmt.Println("The GRPC service is running on port 8080");
//	}
type GrpcController struct {
	Overrides IGrpcControllerOverrides

	defaultConfig *cconf.ConfigParams
	serviceName   string
	config        *cconf.ConfigParams
	references    cref.IReferences
	localEndpoint bool
	opened        bool
	//  The GRPC endpoint that exposes c service.
	Endpoint *GrpcEndpoint
	//  The dependency resolver.
	DependencyResolver *cref.DependencyResolver
	//  The logger.
	Logger *clog.CompositeLogger
	//  The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
}

// InheritGrpcService methods are creates new instance NewGrpcService
// Parameters:
//   - overrides a reference to child class that overrides virtual methods
//   - serviceName string
//
// service name from XYZ.pb.go, set "" for use default gRPC commandable protobuf
// Return *GrpcService
func InheritGrpcService(overrides IGrpcControllerOverrides, serviceName string) *GrpcController {
	c := &GrpcController{
		Overrides: overrides,
	}
	c.serviceName = serviceName
	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"dependencies.endpoint", "*:endpoint:grpc:*:1.0",
	)
	c.DependencyResolver = cref.NewDependencyResolverWithParams(context.Background(), c.defaultConfig, cref.NewEmptyReferences())
	c.Logger = clog.NewCompositeLogger()
	c.Counters = ccount.NewCompositeCounters()
	c.Tracer = ctrace.NewCompositeTracer()
	return c
}

// Configure method are configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- config   configuration parameters to be set.
func (c *GrpcController) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.config = config
	c.DependencyResolver.Configure(ctx, config)
}

// SetReferences method are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- references 	references to locate the component dependencies.
func (c *GrpcController) SetReferences(ctx context.Context, references cref.IReferences) {
	c.references = references
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.Tracer.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
	// Get endpoint
	res := c.DependencyResolver.GetOneOptional("endpoint")
	c.Endpoint, _ = res.(*GrpcEndpoint)
	// Or create a local one
	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint()
		c.localEndpoint = true
	} else {
		c.localEndpoint = false
	}
	// Add registration callback to the endpoint
	c.Endpoint.Register(c)
}

// UnsetReferences method are unsets (clears) previously set references to dependent components.
func (c *GrpcController) UnsetReferences() {
	// Remove registration callback from endpoint
	if c.Endpoint != nil {
		c.Endpoint.Unregister(c)
		c.Endpoint = nil
	}
}

func (c *GrpcController) createEndpoint() *GrpcEndpoint {
	endpoint := NewGrpcEndpoint()
	if c.config != nil {
		endpoint.Configure(context.Background(), c.config)
	}
	if c.references != nil {
		endpoint.SetReferences(context.Background(), c.references)
	}
	return endpoint
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//		- name              a method name.
//
// Return Timing object to end the time measurement.
func (c *GrpcController) Instrument(ctx context.Context, name string) *rpctrace.InstrumentTiming {
	c.Logger.Trace(ctx, "Executing %s method", name)
	c.Counters.IncrementOne(ctx, name+".exec_count")

	counterTiming := c.Counters.BeginTiming(ctx, name+".exec_time")
	traceTiming := c.Tracer.BeginTrace(ctx, name, "")
	return rpctrace.NewInstrumentTiming(ctx, name, "exec",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// InstrumentError method are adds instrumentation to error handling.
//	Parameters:
//		- ctx context.Context	operation context
//		- name              a method name.
//		- errIn               an occured error
//		- resIn            (optional) an execution result
// Returns: result any, err error
// input result and error
// func (c *GrpcService) InstrumentError(ctx context.Context, name string, errIn error,
// 	resIn any) (result any, err error) {
// 	if errIn != nil {
// 		c.Logger.Error(ctx, errIn, "Failed to execute %s method", name)
// 		c.Counters.IncrementOne(name + ".exec_errors")
// 	}
// 	return resIn, errIn
// }

// IsOpen method are checks if the component is opened.
// Return true if the component has been opened and false otherwise.
func (c *GrpcController) IsOpen() bool {
	return c.opened
}

// Open method are opens the component.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//
// Returns: error or nil no errors occured.
func (c *GrpcController) Open(ctx context.Context) (err error) {
	if c.opened {
		return nil
	}

	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint()
		c.Endpoint.Register(c)
		c.localEndpoint = true
	}

	if c.localEndpoint {
		opnErr := c.Endpoint.Open(ctx)
		if opnErr != nil {
			c.opened = false
			return opnErr
		}
	}
	c.opened = true
	return nil
}

// Close method are closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//
// Returns: error or nil no errors occured.
func (c *GrpcController) Close(ctx context.Context) (err error) {
	if !c.opened {
		return nil
	}

	if c.Endpoint == nil {
		return cerr.NewInvalidStateError(utils.ContextHelper.GetTraceId(ctx), "NO_Endpoint", "HTTP endpoint is missing")
	}

	if c.localEndpoint {
		clsErr := c.Endpoint.Close(ctx)
		if clsErr != nil {
			c.opened = false
			return clsErr
		}
	}
	c.opened = false
	return nil
}

// RegisterCommandableMethod method are registers a commandable method in c objects GRPC server (service) by the given name.,
//
//	Parameters:
//		- ctx context.Context	operation context
//		- method        the GRPC method name.
//		- schema        the schema to use for parameter validation.
//		- action        the action to perform at the given route.
func (c *GrpcController) RegisterCommandableMethod(method string, schema *cvalid.Schema,
	action func(ctx context.Context, data *crun.Parameters) (result any, err error)) {
	c.Endpoint.RegisterCommandableMethod(method, schema, action)
}

// Registers a middleware for methods in GRPC endpoint.
//
//	Parameters:
//		- action        an action function that is called when middleware is invoked.
func (c *GrpcController) RegisterUnaryInterceptor(action func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)) {
	if c.Endpoint == nil {
		return
	}

	c.Endpoint.AddInterceptors(grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.HasPrefix(info.FullMethod, "/"+c.serviceName+"/") {
			return action(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}))
}

// Register method are registers all service routes in HTTP endpoint.
func (c *GrpcController) Register() {
	// Override in child classes
	c.Overrides.Register()
}

func (c *GrpcController) ValidateRequest(request any, schema *cvalid.Schema) error {

	buf, err := json.Marshal(request)
	if err != nil {
		return err
	}

	validateObj := make(map[string]any)
	err = json.Unmarshal(buf, &validateObj)
	if err != nil {
		return err
	}

	validateErr := schema.ValidateAndReturnError("", validateObj, false)
	if validateErr != nil {
		return validateErr
	}
	return nil
}
