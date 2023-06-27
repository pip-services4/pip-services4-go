package controllers

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	rpccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"

	grpcproto "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GrpcEndpoint used for creating GRPC endpoints. An endpoint is a URL, at which a given controller can be accessed by a client.
//
// Configuration parameters:
//
// Parameters to pass to the configure method for component configuration:
//
//		Configuration parameters
//
//		- connection(s) - the connection resolver"s connections:
//			- discovery_key - the key to use for connection resolving in a discovery controller;
//			- protocol - the connection"s protocol
//			- host - the target host;
//			- port - the target port;
//			- uri - the target URI.
//		- credential(s) - the HTTPS credentials:
//			- ssl_key_file - the SSL private key in PEM
//			- ssl_crt_file - the SSL certificate in PEM
//			- ssl_ca_file - the certificate authorities (root cerfiticates) in PEM
//
//		References:
//
//			- logger: "*:logger:*:*:1.0";
//			- counters: "*:counters:*:*:1.0";
//			- discovery: "*:discovery:*:*:1.0" (for the connection resolver).
//
//		Examples:
//
//	   func (c* Endpoint) MyMethod(ctx context.Context, config ConfigParams, references IReferences) {
//	       endpoint := NewGrpcEndpoint();
//	       if c.config != nil {
//	           endpoint.Configure(ctx, c._config);
//	       }
//	       if c.references != nil {
//	           endpoint.SetReferences(ctx, c.references);
//	       }
//	       ...
//
//	       err := c.endpoint.Open(ctx)
//	       if err != nil {
//	           // error ocured
//	           return err
//	       }
//	       c.Opened = true
//	       return nil
//	       ...
//	   }
type GrpcEndpoint struct {
	defaultConfig      *cconf.ConfigParams
	server             *grpc.Server
	connectionResolver *rpccon.HttpConnectionResolver
	logger             *clog.CompositeLogger
	counters           *ccount.CompositeCounters
	maintenanceEnabled bool
	fileMaxSize        int64
	uri                string
	registrations      []IRegisterable
	commandableMethods map[string]func(ctx context.Context, args *cexec.Parameters) (result any, err error)
	commandableSchemas map[string]*cvalid.Schema
	interceptors       []grpc.ServerOption
}

// NewGrpcEndpoint method are creates new instance of GrpcEndpoint
func NewGrpcEndpoint() *GrpcEndpoint {
	c := GrpcEndpoint{}
	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", 3000,

		"credential.ssl_key_file", "",
		"credential.ssl_crt_file", "",
		"credential.ssl_ca_file", "",

		"options.maintenance_enabled", false,
		"options.request_max_size", 1024*1024,
		"options.file_max_size", 200*1024*1024,
		"options.connect_timeout", 60000,
		"options.debug", true,
	)

	c.connectionResolver = rpccon.NewHttpConnectionResolver()
	c.logger = clog.NewCompositeLogger()
	c.counters = ccount.NewCompositeCounters()
	c.maintenanceEnabled = false
	c.fileMaxSize = 200 * 1024 * 1024
	c.registrations = make([]IRegisterable, 0)
	c.commandableMethods = make(map[string]func(ctx context.Context, args *cexec.Parameters) (result any, err error), 0)
	c.commandableSchemas = make(map[string]*cvalid.Schema, 0)
	c.interceptors = make([]grpc.ServerOption, 0, 0)
	return &c
}

// Configure method are configures c HttpEndpoint using the given configuration parameters.
// Configuration parameters:
//   - connection(s) - the connection resolver"s connections;
//   - "connection.discovery_key" - the key to use for connection resolving in a discovery controller;
//   - "connection.protocol" - the connection"s protocol;
//   - "connection.host" - the target host;
//   - "connection.port" - the target port;
//   - "connection.uri" - the target URI.
//   - "credential.ssl_key_file" - SSL private key in PEM
//   - "credential.ssl_crt_file" - SSL certificate in PEM
//   - "credential.ssl_ca_file" - Certificate authority (root certificate) in PEM
//     Parameters:
//   - ctx context.Context	operation context
//   - config    configuration parameters, containing a "connection(s)" section.
//
// See ConfigParams (in the PipServices "Commons" package)
func (c *GrpcEndpoint) Configure(ctx context.Context, config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.connectionResolver.Configure(ctx, config)

	c.maintenanceEnabled = config.GetAsBooleanWithDefault("options.maintenance_enabled", c.maintenanceEnabled)
	c.fileMaxSize = config.GetAsLongWithDefault("options.file_max_size", c.fileMaxSize)
}

// SetReferences method are sets references to c endpoint"s logger, counters, and connection resolver.
// References:
//   - logger: "*:logger:*:*:1.0"
//   - counters: "*:counters:*:*:1.0"
//   - discovery: "*:discovery:*:*:1.0" (for the connection resolver)
//     Parameters:
//   - ctx context.Context	operation context
//   - references    an IReferences object, containing references to a logger, counters,
//     and a connection resolver.
//
// See IReferences (in the PipServices "Commons" package)
func (c *GrpcEndpoint) SetReferences(ctx context.Context, references cref.IReferences) {
	c.logger.SetReferences(ctx, references)
	c.counters.SetReferences(ctx, references)
	c.connectionResolver.SetReferences(ctx, references)
}

// IsOpen method are return whether or not c endpoint is open with an actively listening GRPC server.
func (c *GrpcEndpoint) IsOpen() bool {
	return c.server != nil
}

// Open method are opens a connection using the parameters resolved by the referenced connection
// resolver and creates a GRPC server (controller) using the set options and parameters.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//
// Retunrns: an error if one is raised.
func (c *GrpcEndpoint) Open(ctx context.Context) (err error) {

	if c.IsOpen() {
		return nil
	}
	connection, credential, err := c.connectionResolver.Resolve(ctx)
	if err != nil {
		return err
	}
	c.uri = connection.Host() + ":" + strconv.FormatInt(int64(connection.Port()), 10)
	opts := []grpc.ServerOption{}
	if len(c.interceptors) > 0 {
		// Add interceptors
		opts = append(opts, c.interceptors...)
	}
	if connection.Protocol() == "https" {
		sslKeyFile := credential.GetAsString("ssl_key_file")
		sslCrtFile := credential.GetAsString("ssl_crt_file")
		creds, _ := credentials.NewServerTLSFromFile(sslCrtFile, sslKeyFile)
		opts = append(opts, grpc.Creds(creds))
	}
	lis, lErr := net.Listen("tcp", c.uri)
	if lErr != nil {
		return lErr
	}
	// Create instance of express application
	c.server = grpc.NewServer(opts...)
	if c.server == nil {
		return cerr.NewConnectionError(utils.ContextHelper.GetTraceId(ctx), "CAN'T_CREATE_SRV", "Opening GRPC controller failed").
			Wrap(err).WithDetails("url", c.uri)
	}

	err = c.connectionResolver.Register(ctx)
	if err != nil {
		return err
	}

	// Start operations
	c.performRegistrations()

	go func(server *grpc.Server) {
		servErr := server.Serve(lis)
		if servErr != nil {
			err := cerr.NewConnectionError(utils.ContextHelper.GetClient(ctx), "CANNOT_CONNECT", "Opening GRPC controller failed").
				Wrap(servErr).WithDetails("url", c.uri)
			panic(err)
		}
	}(c.server)

	c.logger.Debug(ctx, "Opened GRPC controller at tcp:\\\\%s", c.uri)

	return nil
}

// Close methods are closes c endpoint and the GRPC server (controller) that was opened earlier.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//
// Returns: an error if one is raised.
func (c *GrpcEndpoint) Close(ctx context.Context) (err error) {
	if c.server != nil {
		c.uri = ""

		c.commandableMethods = nil
		c.commandableSchemas = nil

		c.server.GracefulStop()
		c.logger.Debug(ctx, "Closed GRPC controller at %s", c.uri)
		c.server = nil
	}

	return nil
}

// GetServer return working gRPC server for register controllers
// Note: this server is async working in goroutione, wrap into locks if you want change this variable
// Returns *grpc.Server
func (c *GrpcEndpoint) GetServer() *grpc.Server {
	return c.server
}

// AddInterceptors method are registers a middleware for methods in GRPC endpoint.
// See https://github.com/grpc/grpc-go/tree/master/examples/features/interceptor
// Parameters:
//   - interceptors ...grpc.ServerOption
//
// interceptor functions (Stream or Unary use grpc.UnaryInterceptor() or grpc.StreamInterceptor() for inflate in grpc.ServerOption)
func (c *GrpcEndpoint) AddInterceptors(interceptors ...grpc.ServerOption) {
	c.interceptors = append(c.interceptors, interceptors...)
}

// Register method are registers a registerable object for dynamic endpoint discovery.
//   - registration      the registration to add.
//
// See IRegisterable
func (c *GrpcEndpoint) Register(registration IRegisterable) {
	c.registrations = append(c.registrations, registration)
}

// Unregister mwthod are unregisters a registerable object, so that it is no longer used in dynamic
// endpoint discovery.
//   - registration      the registration to remove.
//
// See IRegisterable
func (c *GrpcEndpoint) Unregister(registration IRegisterable) {
	for i := 0; i < len(c.registrations); {
		if c.registrations[i] == registration {
			if i == len(c.registrations)-1 {
				c.registrations = c.registrations[:i]
			} else {
				c.registrations = append(c.registrations[:i], c.registrations[i+1:]...)
			}
		} else {
			i++
		}
	}
}

func (c *GrpcEndpoint) performRegistrations() {
	for _, registration := range c.registrations {
		registration.Register()
	}
	c.registerCommandableController()
}

func (c *GrpcEndpoint) registerCommandableController() {
	if len(c.commandableMethods) == 0 {
		return
	}
	invokeMediator := InvokeComandMediator{InvokeFunc: c.invoke}
	grpcproto.RegisterCommandableServer(c.server, &invokeMediator)
}

// RegisterController method are registers a controller with related implementation
//   - implementation the service implementation method Invoke.
func (c *GrpcEndpoint) RegisterController(sd *grpc.ServiceDesc, implementation any) {
	if c.server != nil {
		c.server.RegisterService(sd, implementation)
	}
}

// Invoke method for implements interface grpcproto.CommandableServer
//
//	Parameters:
//		- ctx context.Context	operation context
//		- request *grpcproto.InvokeRequest request struct
//
// Returns response *grpcproto.InvokeReply and error invocation
func (c *GrpcEndpoint) invoke(ctx context.Context, request *grpcproto.InvokeRequest) (response *grpcproto.InvokeReply, err error) {

	method := request.Method
	var action func(ctx context.Context, args *cexec.Parameters) (result any, err error)
	if len(c.commandableMethods) > 0 {
		action = c.commandableMethods[method]
	}
	traceId := request.TraceId
	// Handle method not found
	if action == nil {
		appErr := cerr.NewInvocationError(traceId, "METHOD_NOT_FOUND", "Method "+method+" was not found").
			WithDetails("method", method)

		var errDesc grpcproto.ErrorDescription
		errDescJson, _ := json.Marshal(appErr)
		json.Unmarshal(errDescJson, &errDesc)
		response = &grpcproto.InvokeReply{
			Error:       &errDesc,
			ResultEmpty: true,
			ResultJson:  "",
		}
		return response, nil
	}
	// Convert arguments
	argsEmpty := request.ArgsEmpty
	argsJson := request.ArgsJson
	var args *cexec.Parameters = cexec.NewEmptyParameters()

	if !argsEmpty && argsJson != "" {
		var buf map[string]any
		err := json.Unmarshal([]byte(argsJson), &buf)
		if err == nil {
			args.Append(buf)
		}
	}
	// Call command action
	ctx = utils.ContextHelper.NewContextWithTraceId(ctx, traceId)
	result, err := action(ctx, args)
	// Process result and generate response
	if err != nil {
		appErr := cerr.ErrorDescriptionFactory.Create(err)
		var errDesc grpcproto.ErrorDescription
		errDescJson, _ := json.Marshal(appErr)
		json.Unmarshal(errDescJson, &errDesc)
		response = &grpcproto.InvokeReply{
			Error:       &errDesc,
			ResultEmpty: true,
			ResultJson:  "",
		}
	} else {
		resJson, _ := json.Marshal(result)
		response = &grpcproto.InvokeReply{
			Error:       nil,
			ResultEmpty: result == nil || string(resJson) == "null",
			ResultJson:  string(resJson),
		}
	}
	return response, err
}

// RegisterCommandableMethod method are registers a commandable method in c objects GRPC server (controller) by the given name.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//		- method        the GRPC method name.
//		- schema        the schema to use for parameter validation.
//		- action        the action to perform at the given route.
func (c *GrpcEndpoint) RegisterCommandableMethod(method string, schema *cvalid.Schema,
	action func(ctx context.Context, args *cexec.Parameters) (result any, err error)) {

	if c.commandableMethods == nil {
		c.commandableMethods = make(map[string]func(ctx context.Context, args *cexec.Parameters) (result any, err error))
	}
	c.commandableMethods[method] = action
	if c.commandableSchemas == nil {
		c.commandableSchemas = make(map[string]*cvalid.Schema)
	}
	c.commandableSchemas[method] = schema
}
