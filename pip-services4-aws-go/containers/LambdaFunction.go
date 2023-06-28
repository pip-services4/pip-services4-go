package containers

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	awsserv "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/controllers"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cproc "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	rpctrace "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/trace"
)

// Abstract AWS Lambda function, that acts as a container to instantiate and run components
// and expose them via external entry point.
//
// When handling calls "cmd" parameter determines which what action shall be called, while
// other parameters are passed to the action itself.
//
// Container configuration for this Lambda function is stored in "./config/config.yml" file.
// But this path can be overriden by CONFIG_PATH environment variable.
//
// Configuration parameters
//
//   - dependencies:
//   - controller:                  override for Controller dependency
//   - connections:
//   - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//   - region:                      (optional) AWS region
//   - credentials:
//   - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//   - access_id:                   AWS access/client id
//   - access_key:                  AWS access/client id
//
// References
//
//   - \*:logger:\*:\*:1.0            (optional) ILogger components to pass log messages
//   - \*:counters:\*:\*:1.0          (optional) ICounters components to pass collected measurements
//   - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connection
//   - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//
// # See LambdaClient
//
// Example:
//
//	type MyLambdaFunction struct {
//		*awscont.LambdaFunction
//		controller awstest.IMyController
//	}
//
//	func NewMyLambdaFunction() *MyLambdaFunction {
//		c := &MyLambdaFunction{}
//		c.LambdaFunction = awscont.InheriteLambdaFunction(c, "mygroup", "MyGroup lambda function")
//
//		c.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("mygroup", "controller", "*", "*", "1.0"))
//		return c
//	}
//
//	func (c *MyLambdaFunction) SetReferences(ctx context.Context, references cref.IReferences) {
//		c.LambdaFunction.SetReferences(ctx, references)
//		depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
//		if depErr == nil && depRes != nil {
//			c.controller = depRes.(awstest.IMyController)
//		}
//	}
//
//	func (c *MyLambdaFunction) getOneById(ctx context.Context, params map[string]any) (any, error) {
//		traceId, _ := params["trace_id"].(string)
//		ctx := cctx.NewContextWithTrace(ctx.Background(), traceId)
//		return c.controller.GetOneById(
//			ctx,
//			params["mydata_id"].(string),
//		)
//	}
//
//	func (c *MyLambdaFunction) Register() {
//
//		c.RegisterAction(
//			"get_mydata_by_id",
//			cvalid.NewObjectSchema().
//				WithOptionalProperty("mydata_id", cconv.String).Schema,
//			c.getOneById)
//	}
//
//
//	lambda := NewMyLambdaFunction();
//
//	lambda.Run(context.Context())
type LambdaFunction struct {
	*cproc.Container
	Overrides ILambdaFunctionOverrides

	references cref.IReferences
	// The performanc counters.
	counters *ccount.CompositeCounters
	// The tracer.
	tracer *ctrace.CompositeTracer
	// The dependency resolver.
	DependencyResolver *cref.DependencyResolver
	// The map of registred validation schemas
	schemas map[string]*cvalid.Schema
	// The map of registered actions.
	actions map[string]func(context.Context, map[string]any) (any, error)
	// The default path to config file
	configPath string
}

// Creates a new instance of this lambda function.
// - overrides	Lambda function register instance.
// - name          (optional) a container name (accessible via ContextInfo)
// - description   (optional) a container description (accessible via ContextInfo)
func InheriteLambdaFunction(overrides ILambdaFunctionOverrides, name string, description string) *LambdaFunction {
	c := &LambdaFunction{
		counters:           ccount.NewCompositeCounters(),
		tracer:             ctrace.NewCompositeTracer(),
		DependencyResolver: cref.NewDependencyResolver(),
		schemas:            make(map[string]*cvalid.Schema, 0),
		actions:            make(map[string]func(context.Context, map[string]any) (any, error), 0),
		configPath:         "./config/config.yml",
		Overrides:          overrides,
	}
	c.Container = cproc.InheritContainer(name, description, overrides)
	c.SetLogger(log.NewConsoleLogger())
	return c
}

func (c *LambdaFunction) Register() {}

func (c *LambdaFunction) getConfigPath() string {
	res := os.Getenv("CONFIG_PATH")
	if res == "" {
		return c.configPath
	}
	return res
}

func (c *LambdaFunction) getParameters() *cconf.ConfigParams {
	parameters := cconf.NewConfigParamsFromValue(os.Environ())
	return parameters
}

func (c *LambdaFunction) captureErrors(ctx context.Context) {
	if r := recover(); r != nil {
		err, _ := r.(error)
		c.Logger().Fatal(ctx, err, "Process is terminated")
		os.Exit(1)
	}
}

func (c *LambdaFunction) captureExit(ctx context.Context) {
	c.Logger().Info(ctx, "Press Control-C to stop the microservice...")

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	go func() {
		select {
		case <-ch:
			c.Close(ctx)
			c.Logger().Info(ctx, "Googbye!")
			os.Exit(0)

		case <-ctx.Done():
			c.Close(ctx)

			if ctx.Err() != nil {
				c.Logger().Error(ctx, ctx.Err(), "Application crashed.")
				os.Exit(1)
			} else {
				c.Logger().Info(ctx, "Googbye!")
				os.Exit(0)
			}

		}

	}()
}

// Sets references to dependent components.
// Parameters:
//   - ctx context.Context	operation context.
//   - references 	references to locate the component dependencies.
func (c *LambdaFunction) SetReferences(ctx context.Context, references cref.IReferences) {
	c.references = references
	c.counters.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
	c.Overrides.Register()
}

// Adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
// Parameters:
//   - ctx context.Context	execution context to trace execution through call chain.
//   - name	a method name.
//
// Returns Timing object to end the time measurement.
func (c *LambdaFunction) Instrument(ctx context.Context, name string) *rpctrace.InstrumentTiming {
	c.Logger().Trace(ctx, "Executing %s method", name)
	c.counters.IncrementOne(ctx, name+".exec_count")

	counterTiming := c.counters.BeginTiming(ctx, name+".exec_time")
	traceTiming := c.tracer.BeginTrace(ctx, name, "")
	return rpctrace.NewInstrumentTiming(ctx, name, "exec",
		c.Logger(), c.counters, counterTiming, traceTiming)
}

// InstrumentError method are adds instrumentation to error handling.
// Parameters:
//   - ctx context.Context	execution context to trace execution through call chain.
//   - name    string         a method name.
//   - err     error          an occured error
//   - result  any    (optional) an execution result
//
// Returns:  result any, err error
// (optional) an execution callback
func (c *LambdaFunction) InstrumentError(ctx context.Context, name string, errIn error,
	resIn any) (result any, err error) {
	if errIn != nil {
		c.Logger().Error(ctx, errIn, "Failed to execute %s method", name)
		c.counters.IncrementOne(ctx, name+".exec_errors")
	}
	return resIn, errIn
}

// Runs this lambda function, loads container configuration,
// instantiate components and manage their lifecycle,
// makes this function ready to access action calls.
// Parameters:
//   - ctx context.Context	operation context.
func (c *LambdaFunction) Run(ctx context.Context) error {
	traceId := c.Info().Name
	ctx = cctx.NewContextWithTraceId(ctx, traceId)

	path := c.getConfigPath()
	parameters := c.getParameters()
	c.ReadConfigFromFile(ctx, path, parameters)

	c.captureErrors(ctx)
	c.captureExit(ctx)
	return c.Open(ctx)
}

//	Opens the component.
//
// Parameters:
//   - ctx context.Context	execution context to trace execution through call chain.
func (c *LambdaFunction) Open(ctx context.Context) error {
	if c.IsOpen() {
		return nil
	}

	err := c.Container.Open(ctx)
	if err != nil {
		return err
	}

	c.RegisterControllers()
	return nil
}

// Registers all lambda controllers in the container.
func (c *LambdaFunction) RegisterControllers() {
	// Extract regular and commandable Lambda controllers from references
	controllersRefs := c.references.GetOptional(
		cref.NewDescriptor("*", "controller", "awslambda", "*", "*"),
	)

	controllers := make([]awsserv.ILambdaController, 0)

	for _, ctrl := range controllersRefs {
		if s, ok := ctrl.(awsserv.ILambdaController); ok {
			controllers = append(controllers, s)
		}
	}

	cmdControllersRefs := c.references.GetOptional(
		cref.NewDescriptor("*", "controller", "commandable-awslambda", "*", "*"),
	)

	for _, ctrl := range cmdControllersRefs {
		if s, ok := ctrl.(awsserv.ILambdaController); ok {
			controllers = append(controllers, s)
		}
	}

	// Register actions defined in those controllers
	for _, ctrl := range controllers {
		actions := ctrl.GetActions()
		for _, action := range actions {
			c.Logger().Debug(context.Background(), "RegisterControllers", "Register commmand: %v", action.Cmd)
			c.RegisterAction(action.Cmd, action.Schema, action.Action)
		}
	}
}

// Registers an action in this lambda function.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- cmd           a action/command name.
//		- schema        a validation schema to validate received parameters.
//		- action        an action function that is called when action is invoked.
func (c *LambdaFunction) RegisterAction(cmd string, schema *cvalid.Schema,
	action func(ctx context.Context, params map[string]any) (result any, err error)) error {

	if cmd == "" {
		return cerr.NewUnknownError("", "NO_COMMAND", "Missing command")
	}

	if action == nil {
		return cerr.NewUnknownError("", "NO_ACTION", "Missing action")
	}

	// Hack!!! Wrapping action to preserve prototyping context
	actionCurl := func(ctx context.Context, params map[string]any) (any, error) {
		// Perform validation
		if schema != nil {
			traceId, _ := params["trace_id"].(string)
			err := schema.ValidateAndReturnError(traceId, params, false)
			if err != nil {
				return nil, err
			}
		}

		return action(ctx, params)
	}

	c.actions[cmd] = actionCurl
	return nil
}

func (c *LambdaFunction) execute(ctx context.Context, body map[string]any) (string, error) {
	cmd, ok := body["cmd"].(string)
	traceId, _ := body["trace_id"].(string)

	if !ok || cmd == "" {
		err := cerr.NewBadRequestError(
			traceId,
			"NO_COMMAND",
			"Cmd parameter is missing")
		ctx.Done()
		return "ERROR", err
	}

	action := c.actions[cmd]
	if action == nil {
		err := cerr.NewBadRequestError(
			traceId,
			"NO_ACTION",
			"Action "+cmd+" was not found").
			WithDetails("command", cmd)

		ctx.Done()
		return "ERROR", err
	}

	res, err := action(ctx, body)
	ctx.Done()
	resStr := "ERROR"
	if res != nil {
		convRes, convErr := json.Marshal(res)
		if convRes == nil || convErr != nil {
			err = convErr
		} else {
			resStr = (string)(convRes)
		}
	}
	return resStr, err
}

func (c *LambdaFunction) Handler(ctx context.Context, event map[string]any) (string, error) { //handler(event: any, context: any) {
	// If already started then execute
	if c.IsOpen() {
		if event != nil {
			return c.execute(ctx, event)
		}
	} else { // Start before execute
		err := c.Run(ctx)
		if err != nil {
			ctx.Done()
			return "", err
		}
		if event != nil {
			return c.execute(ctx, event)
		}
	}
	err := cerr.NewBadRequestError(
		"Lambda",
		"NO_EVENT",
		"Event is empty")
	return "ERROR", err
}

// Gets entry point into this lambda function.
//   - event     an incoming event object with invocation parameters.
//   - context   a context object with local references.
func (c *LambdaFunction) GetHandler() func(ctx context.Context, event map[string]any) (string, error) {

	// Return plugin function
	return func(ctx context.Context, event map[string]any) (string, error) {
		// Calling run with changed context
		return c.Handler(ctx, event)
	}
}

// Calls registered action in this lambda function.
// "cmd" parameter in the action parameters determin
// what action shall be called.
//
// This method shall only be used in testing.
//   - params action parameters.
//   - callback callback function that receives action result or error.
func (c *LambdaFunction) Act(params map[string]any) (string, error) {
	ctx := context.TODO()
	return c.GetHandler()(ctx, params)
}
