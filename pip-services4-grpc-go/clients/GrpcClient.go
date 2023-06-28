package clients

import (
	"context"
	"time"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	rpccon "github.com/pip-services4/pip-services4-go/pip-services4-config-go/connect"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	grpcproto "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/protos"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	rpctrace "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// GrpcClient abstract client that calls commandable HTTP service.
//
// Commandable services are generated automatically for ICommandable objects. Each command is exposed as POST operation that receives all parameters in body object.
//
//		Configuration parameters:
//
//	 	- base_route: base route for remote URI
//	 	- connection(s):
//	 		- discovery_key: (optional) a key to retrieve the connection from IDiscovery
//	 		- protocol: connection protocol: http or https
//	 		- host: host name or IP address
//	 		- port: port number
//	 		- uri: resource URI or connection string with all parameters in it
//	 	- options:
//	 		- retries: number of retries (default: 3)
//	 		- connect_timeout: connection timeout in milliseconds (default: 10 sec)
//	 		- timeout: invocation timeout in milliseconds (default: 10 sec)
//
//		References:
//
//			- *:logger:*:*:1.0 (optional) ILogger components to pass log messages
//			- *:counters:*:*:1.0 (optional) ICounters components to pass collected measurements
//			- *:discovery:*:*:1.0 (optional) IDiscovery services to resolve connection
//
// Example:
//
//	type MyGrpcClient struct{
//		*GrpcClient
//	}
//
//	func  (c *MyGrpcClient) GetData(ctx context.Context, id string) (res any, err error) {
//		req := &testproto.MyDataIdRequest{
//		    TraceId: cctx.GetTraceId(ctx),
//		    mydataId:       id,
//		}
//		reply := new(testproto.MyData)
//		err = c.Call("get_mydata_by_id", req, reply)
//		c.Instrument(ctx, traceId, "mydata.get_one_by_id")
//		if err != nil {
//		    return nil, err
//		}
//
//		result = toMyData(reply)
//		if result != nil && result.Id == "" && result.Key == "" {
//		    result = nil
//		}
//
//		return result, nil
//	}
//
//	var client = NewMyGrpcClient();
//	client.Configure(ctx, NewConfigParamsFromTuples(
//	    "connection.protocol", "http",
//	    "connection.host", "localhost",
//	    "connection.port", 8080,
//	));
//
//	result, err := client.GetData(ctx, "123", "1")
//	...
type GrpcClient struct {
	address string
	name    string

	defaultConfig *cconf.ConfigParams
	//	The GRPC client.
	Client grpcproto.CommandableClient
	// The GRPC connection
	Connection *grpc.ClientConn
	//	The connection resolver.
	ConnectionResolver *rpccon.HttpConnectionResolver
	//	The logger.
	Logger *clog.CompositeLogger
	//	The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
	//	The configuration options.
	Options *cconf.ConfigParams
	//	The connection timeout in milliseconds.
	ConnectTimeout time.Duration
	//	The invocation timeout in milliseconds.
	Timeout time.Duration
	//	The remote service uri which is calculated on open.
	Uri string
	// interceptors
	interceptors []grpc.DialOption
}

// NewGrpcClient method are creates a new instance of the client.
// Parameters:
//   - baseRoute string
//     a base route for remote service.
//
// Returns *GrpcClient
func NewGrpcClient(name string) *GrpcClient {
	c := GrpcClient{
		name: name,
	}
	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "8090",

		"options.connect_timeout", 10000,
		"options.timeout", 10000,
		"options.retries", 3,
		"options.debug", true,
	)
	c.ConnectionResolver = rpccon.NewHttpConnectionResolver()
	c.Logger = clog.NewCompositeLogger()
	c.Counters = ccount.NewCompositeCounters()
	c.Tracer = ctrace.NewCompositeTracer()
	c.Options = cconf.NewEmptyConfigParams()
	c.ConnectTimeout = 10000 * time.Millisecond
	c.Timeout = 10000 * time.Millisecond
	c.interceptors = make([]grpc.DialOption, 0)
	return &c
}

// Configure method are configures component by passing configuration parameters.
//
//		Parameters:
//			- ctx context.Context	operation context
//			- config *config.ConfigParams
//	  configuration parameters to be set.
func (c *GrpcClient) Configure(ctx context.Context, config *cconf.ConfigParams) {
	host := config.GetAsStringWithDefault("connection.host", "localhost")
	port := config.GetAsStringWithDefault("connection.port", "8090")

	c.ConnectTimeout = time.Duration(config.GetAsIntegerWithDefault("connection.connect_timeout", 10000)) * time.Millisecond
	c.Timeout = time.Duration(config.GetAsIntegerWithDefault("connection.timeout", 10000)) * time.Millisecond
	c.ConnectionResolver.Configure(ctx, config)
	c.address = host + ":" + port
}

// SetReferences method are sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- references  cref.IReferences
//
// references to locate the component dependencies.
func (c *GrpcClient) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.Tracer.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a rpctrace.InstrumentTiming object that is used to end the time measurement.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//		- name string a method name.
//	Returns: rpctrace.InstrumentTiming object to end the time measurement.
func (c *GrpcClient) Instrument(ctx context.Context, name string) *rpctrace.InstrumentTiming {
	c.Logger.Trace(ctx, "Calling %s method", name)
	c.Counters.IncrementOne(ctx, name+".call_count")
	counterTiming := c.Counters.BeginTiming(ctx, name+".call_time")
	traceTiming := c.Tracer.BeginTrace(ctx, name, "")
	return rpctrace.NewInstrumentTiming(ctx, name, "call",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// // InstrumentError mrthod are adds instrumentation to error handling.
// //   - ctx context.Context     a context to trace execution through call chain.
// //   - name              a method name.
// //   - err               an occured error
// //   - result            (optional) an execution result
// // Retruns: result any, err error
// // input result and error.
// func (c *GrpcClient) InstrumentError(ctx context.Context, name string, inErr error, inRes any) (result any, err error) {
// 	if inErr != nil {
// 		c.Logger.Error(ctx, inErr, "Failed to call %s method", name)
// 		c.Counters.IncrementOne(ctx, name+".call_errors")
// 	}

// 	return inRes, inErr
// }

// IsOpen method are checks if the component is opened.
// Returns bool
// true if the component has been opened and false otherwise.
func (c *GrpcClient) IsOpen() bool {
	return c.Connection != nil
}

// AddInterceptors method are registers a middleware for methods in gRPC client.
// See https://github.com/grpc/grpc-go/tree/master/examples/features/interceptor
// Parameters:
//   - interceptors ...grpc.DialOption
//
// interceptor functions (Stream or Unary use grpc.WithUnaryInterceptor() or grpc.WithStreamInterceptor() for inflate in grpc.ServerOption)
func (c *GrpcClient) AddInterceptors(interceptors ...grpc.DialOption) {
	c.interceptors = append(c.interceptors, interceptors...)
}

// Open method are opens the component.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//
// transaction id to trace execution through call chain.
// Returns error
// error or nil
func (c *GrpcClient) Open(ctx context.Context) error {

	if c.IsOpen() {
		return nil
	}
	connection, credential, err := c.ConnectionResolver.Resolve(ctx)
	if err != nil {
		return err
	}
	c.Uri = connection.Uri()

	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(ctx, c.ConnectTimeout)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{Timeout: c.Timeout}),
	}

	if len(c.interceptors) > 0 {
		// Add interceptors
		opts = append(opts, c.interceptors...)
	}

	if connection.Protocol() == "https" {
		//sslKeyFile := credential.GetAsString("ssl_key_file")
		sslCrtFile := credential.GetAsString("ssl_crt_file")
		transport, err := credentials.NewClientTLSFromFile(sslCrtFile, c.name)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.WithTransportCredentials(transport))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.DialContext(ctx, c.address, opts...)
	if err != nil {
		return err
	}
	c.Connection = conn
	c.Client = grpcproto.NewCommandableClient(conn)
	return nil
}

// Close method are closes component and frees used resources.
//
//		Parameters:
//			- ctx context.Context	a context to trace execution through call chain.
//	  transaction id to trace execution through call chain.
//
// Returns error
func (c *GrpcClient) Close(ctx context.Context) error {
	if c.Connection != nil {
		c.Connection.Close()
		c.Connection = nil
	}
	return nil
}

// Call method are calls a remote method via gRPC protocol.
//
//	Parameters:
//		- method string gRPC method name
//		- request any request query parameters.
//		- response any
//		- response body object.
//
// Returns error
func (c *GrpcClient) Call(method string, request any, response any) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	method = "/" + c.name + "/" + method
	err := c.Connection.Invoke(ctx, method, request, response)
	return err
}

// CallWithContext method are calls a remote method via gRPC protocol.
//
//	Parameters:
//		- ctx context.Context	a context to trace execution through call chain.
//		- method string   gRPC method name
//		- request any request query parameters.
//		- response any
//		- response body object.
//
// Returns error
func (c *GrpcClient) CallWithContext(ctx context.Context, method string, request any, response any) error {
	method = "/" + c.name + "/" + method
	err := c.Connection.Invoke(ctx, method, request, response)
	return err
}

// AddFilterParams method are adds filter parameters (with the same name as they defined)
// to invocation parameter map.
//
//	Parameters:
//		- params        invocation parameters.
//		- filter        (optional) filter parameters
//
// Return invocation parameters with added filter parameters.
func (c *GrpcClient) AddFilterParams(params *cdata.StringValueMap, filter *cquery.FilterParams) *cdata.StringValueMap {

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	if filter != nil {
		for k, v := range filter.Value() {
			params.Put(k, v)
		}
	}
	return params
}

// AddPagingParams method are adds paging parameters (skip, take, total) to invocation parameter map.
//
//	Parameters:
//		- ctx context.Context	operation context
//		- params        invocation parameters.
//		- paging        (optional) paging parameters
//
// Return invocation parameters with added paging parameters.
func (c *GrpcClient) AddPagingParams(params *cdata.StringValueMap, paging *cquery.PagingParams) *cdata.StringValueMap {
	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}

	if paging != nil {
		params.Put("total", paging.Total)
		params.Put("skip", paging.Skip)
		params.Put("take", paging.Take)
	}
	return params
}
