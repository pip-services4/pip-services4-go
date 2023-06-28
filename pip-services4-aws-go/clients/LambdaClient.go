package clients

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	awscon "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/connect"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	ckeys "github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	rpctrace "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/trace"
)

// Abstract client that calls AWS Lambda Functions.
//
// When making calls "cmd" parameter determines which what action shall be called, while
// other parameters are passed to the action itself.
//
// Configuration parameters:
//
//   - connections:
//   - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//   - region:                      (optional) AWS region
//   - credentials:
//   - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//   - access_id:                   AWS access/client id
//   - access_key:                  AWS access/client id
//   - options:
//   - connect_timeout:             (optional) connection timeout in milliseconds (default: 10 sec)
//
// References:
//
//   - \*:logger:\*:\*:1.0            (optional) ILogger components to pass log messages
//
//   - \*:counters:\*:\*:1.0          (optional) ICounters components to pass collected measurements
//
//   - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connection
//
//   - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//
//     See LambdaFunction
//     See CommandableLambdaClient
//
// Example:
//
//	type MyLambdaClient struct  {
//		*LambdaClient
//		...
//	}
//
//	func (c* MyLambdaClient) getData(ctx context.Context, id string)(result MyData, err error){
//		timing := c.Instrument(ctx, "myclient.get_data");
//		callRes, callErr := c.Call(ctx ,"get_data", map[string]interface{ "id": id })
//		if callErr != nil {
//			return callErr
//		}
//		defer timing.EndTiming(ctx, nil)
//		return awsclient.HandleLambdaResponse[*cdata.DataPage[MyData]](calValue)
//	}
//	...
//
//
//	client = NewMyLambdaClient();
//	client.Configure(context.Background(), NewConfigParamsFromTuples(
//	    "connection.region", "us-east-1",
//	    "connection.access_id", "XXXXXXXXXXX",
//	    "connection.access_key", "XXXXXXXXXXX",
//	    "connection.arn", "YYYYYYYYYYYYY"
//	))
//
//	data, err := client.GetData(context.Background(), "1")
//	...
type LambdaClient struct {
	// The reference to AWS Lambda Function.
	Lambda *lambda.Lambda
	// The opened flag.
	Opened bool
	// The AWS connection parameters
	Connection     *awscon.AwsConnectionParams
	connectTimeout int
	// The dependencies resolver.
	DependencyResolver *cref.DependencyResolver
	// The connection resolver.
	ConnectionResolver *awscon.AwsConnectionResolver
	// The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
}

func NewLambdaClient() *LambdaClient {
	c := &LambdaClient{
		Opened:             false,
		connectTimeout:     10000,
		DependencyResolver: cref.NewDependencyResolver(),
		ConnectionResolver: awscon.NewAwsConnectionResolver(),
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
	}
	return c
}

// Configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- config    configuration parameters to be set.
func (c *LambdaClient) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.DependencyResolver.Configure(ctx, config)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
}

// Sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- references	references to locate the component dependencies.
func (c *LambdaClient) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a services.InstrumentTiming object that is used to end the time measurement.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- name string a method name.
//	Returns: services.InstrumentTiming object to end the time measurement.
func (c *LambdaClient) Instrument(ctx context.Context, name string) *rpctrace.InstrumentTiming {
	c.Logger.Trace(ctx, "Calling %s method", name)
	c.Counters.IncrementOne(ctx, name+".call_count")
	counterTiming := c.Counters.BeginTiming(ctx, name+".call_time")
	traceTiming := c.Tracer.BeginTrace(ctx, name, "")
	return rpctrace.NewInstrumentTiming(ctx, name, "call",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *LambdaClient) IsOpen() bool {
	return c.Opened
}

// Opens the component.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- Return 			 error or nil no errors occured.
func (c *LambdaClient) Open(ctx context.Context) error {
	if c.IsOpen() {
		return nil
	}

	wg := sync.WaitGroup{}
	var errGlobal error

	wg.Add(1)
	go func() {
		defer wg.Done()
		connection, err := c.ConnectionResolver.Resolve(ctx)
		c.Connection = connection
		errGlobal = err

		awsCred := credentials.NewStaticCredentials(c.Connection.GetAccessId(), c.Connection.GetAccessKey(), "")
		sess := session.Must(session.NewSession(&aws.Config{
			MaxRetries:  aws.Int(3),
			Region:      aws.String(c.Connection.GetRegion()),
			Credentials: awsCred,
		}))
		// Create new cloudwatch client.
		c.Lambda = lambda.New(sess)
		c.Lambda.Config.HTTPClient.Timeout = time.Duration((int64)(c.connectTimeout)) * time.Millisecond
		c.Logger.Debug(ctx, "Lambda client connected to %s", c.Connection.GetArn())

	}()
	wg.Wait()
	if errGlobal != nil {
		c.Opened = false
		return errGlobal
	}
	return nil
}

// Closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- Returns 			 error or null no errors occured.
func (c *LambdaClient) Close(ctx context.Context) error {
	// Todo: close listening?
	c.Opened = false
	c.Lambda = nil
	return nil
}

// Performs AWS Lambda Function invocation.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- invocationType    an invocation type: "RequestResponse" or "Event"
//		- cmd               an action name to be called.
//		- args              action arguments
//
// Returns           result or error.
func (c *LambdaClient) Invoke(ctx context.Context, invocationType string, cmd string, args map[string]any) (result *lambda.InvokeOutput, err error) {
	traceId := cctx.GetTraceId(ctx)
	if cmd == "" {
		err = cerr.NewUnknownError("", "NO_COMMAND", "Missing cmd")
		c.Logger.Error(ctx, err, "Failed to call %s", cmd)
		return nil, err
	}

	//args = _.clone(args)

	args["cmd"] = cmd
	if traceId != "" {
		args["trace_id"] = traceId
	} else {
		args["trace_id"] = ckeys.IdGenerator.NextLong()
	}

	payloads, jsonErr := json.Marshal(args)

	if jsonErr != nil {
		c.Logger.Error(ctx, jsonErr, "Failed to call %s", cmd)
		return nil, jsonErr
	}

	params := &lambda.InvokeInput{
		FunctionName:   aws.String(c.Connection.GetArn()),
		InvocationType: aws.String(invocationType),
		LogType:        aws.String("None"),
		Payload:        payloads,
	}

	data, lambdaErr := c.Lambda.InvokeWithContext(ctx, params)

	if lambdaErr != nil {
		err = cerr.NewInvocationError(
			traceId,
			"CALL_FAILED",
			"Failed to invoke lambda function").WithCause(err)
		return nil, err
	}

	return data, nil
}

// Calls a AWS Lambda Function action.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- cmd               an action name to be called.
//		- params            (optional) action parameters.
//		- Returns           result and error.
func (c *LambdaClient) Call(ctx context.Context, cmd string, params map[string]any) (result *lambda.InvokeOutput, err error) {
	return c.Invoke(ctx, "RequestResponse", cmd, params)
}

// Calls a AWS Lambda Function action asynchronously without waiting for response.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		- cmd               an action name to be called.
//		- params            (optional) action parameters.
//		- Returns           error or null for success.
func (c *LambdaClient) CallOneWay(ctx context.Context, cmd string, params map[string]any) error {
	_, err := c.Invoke(ctx, "Event", cmd, params)
	return err
}
