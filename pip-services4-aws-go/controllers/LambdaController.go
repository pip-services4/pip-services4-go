package controllers

import (
	"context"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	rpctrace "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/trace"
)

type ILambdaControllerOverrides interface {
	Register()
}

// Abstract service that receives remove calls via AWS Lambda protocol.
//
// This service is intended to work inside LambdaFunction container that
// exploses registered actions externally.
//
// Configuration parameters
//
//   - dependencies:
//   - controller:            override for Controller dependency
//
// References:
//   - *:logger:*:*:1.0               (optional) [[ILogger]] components to pass log messages
//   - *:counters:*:*:1.0             (optional) [[ICounters]] components to pass collected measurements
//
// # See LambdaClient
//
// Example
//
//	   struct MyLambdaController struct  {
//	      *LambdaController
//	      service IMyService
//	   }
//	      ...
//		func NewMyLambdaController()* MyLambdaController {
//		   c:= &MyLambdaController{}
//		   c.LambdaController = NewLambdaController("v1.mycontroller")
//		   c.DependencyResolver.Put(
//			   context.Background(),
//		       "controller",
//		       cref.NewDescriptor("mygroup","controller","*","*","1.0")
//		   )
//		   return c
//		}
//
//		func (c * LambdaController)  SetReferences(ctx context.Context, references IReferences){
//		   c.LambdaController.SetReferences(references)
//		   ref := c.DependencyResolver.GetRequired("controller")
//		   c.service = ref.(IMyService)
//		}
//
//		func (c * LambdaController)  Register() {
//			c.RegisterAction("get_mydata", nil,  func(ctx context.Context, params map[string]any)(any, error) {
//		        traceId := params.GetAsString("trace_id")
//		        id := params.GetAsString("id")
//				return  c.controller.GetMyData(cctx.NewContextWithTraceId(ctx), id)
//		    })
//		    ...
//		}
//
//		controller := NewMyLambdaController();
//		controller.Configure(ctx context.Context, NewConfigParamsFromTuples(
//		    "connection.protocol", "http",
//		    "connection.host", "localhost",
//		    "connection.port", 8080
//		))
//		controller.SetReferences(context.Background(), cref.NewReferencesFromTuples(
//		   cref.NewDescriptor("mygroup","service","default","default","1.0"), service
//		))
//
//		controller.Open(context.Background())
//		fmt.Println("The Lambda 'v1.myservice' controller is running on port 8080");
type LambdaController struct { // ILambdaController, IOpenable, IConfigurable, IReferenceable

	name         string
	actions      []*LambdaAction
	interceptors []func(ctx context.Context, params map[string]any, next func(ctx context.Context, params map[string]any) (any, error)) (any, error)
	opened       bool

	Overrides ILambdaControllerOverrides

	// The dependency resolver.
	DependencyResolver *cref.DependencyResolver
	// The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters
	//The tracer.
	Tracer *ctrace.CompositeTracer
}

// Creates an instance of this service.
// -  name a service name to generate action cmd. LambdaController()
func InheritLambdaController(overrides ILambdaControllerOverrides, name string) *LambdaController {
	return &LambdaController{
		Overrides:          overrides,
		name:               name,
		actions:            make([]*LambdaAction, 0),
		interceptors:       make([]func(ctx context.Context, params map[string]any, next func(ctx context.Context, params map[string]any) (any, error)) (any, error), 0),
		DependencyResolver: cref.NewDependencyResolver(),
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
	}
}

// Configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		-  config    configuration parameters to be set.
func (c *LambdaController) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(ctx, config)
}

// Sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		-  references 	references to locate the component dependencies.
func (c *LambdaController) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.Tracer.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
}

// Get all actions supported by the service.
// Returns an array with supported actions.
func (c *LambdaController) GetActions() []*LambdaAction {
	return c.actions
}

// Adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
//		-  name              a method name.
//
// returns Timing object to end the time measurement.
func (c *LambdaController) Instrument(ctx context.Context, name string) *rpctrace.InstrumentTiming {
	c.Logger.Trace(ctx, "Executing %s method", name)
	c.Counters.IncrementOne(ctx, name+".exec_count")

	counterTiming := c.Counters.BeginTiming(ctx, name+".exec_time")
	traceTiming := c.Tracer.BeginTrace(ctx, name, "")
	return rpctrace.NewInstrumentTiming(ctx, name, "exec",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// Checks if the component is opened.
// Returns true if the component has been opened and false otherwise.
func (c *LambdaController) IsOpen() bool {
	return c.opened
}

// Opens the component.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
func (c *LambdaController) Open(ctx context.Context) error {
	if c.opened {
		return nil
	}

	c.Register()

	c.opened = true
	return nil
}

// Closes component and frees used resources.
//
//	Parameters:
//		- ctx context.Context	execution context to trace execution through call chain.
func (c *LambdaController) Close(ctx context.Context) error {
	if !c.opened {
		return nil
	}

	c.opened = false
	c.actions = make([]*LambdaAction, 0)
	c.interceptors = make([]func(ctx context.Context, params map[string]any, next func(ctx context.Context, params map[string]any) (any, error)) (any, error), 0)
	return nil
}

func (c *LambdaController) ApplyValidation(schema *cvalid.Schema, action func(ctx context.Context, params map[string]any) (any, error)) func(context.Context, map[string]any) (any, error) {
	// Create an action function
	actionWrapper := func(ctx context.Context, params map[string]any) (any, error) {
		// Validate object
		if schema != nil && params != nil {
			// Perform validation
			traceId, _ := params["trace_id"].(string)
			err := schema.ValidateAndReturnError(traceId, params, false)
			if err != nil {
				return nil, err
			}
		}
		return action(ctx, params)
	}

	return actionWrapper
}

func (c *LambdaController) ApplyInterceptors(action func(context.Context, map[string]any) (any, error)) func(context.Context, map[string]any) (any, error) {
	actionWrapper := action

	for index := len(c.interceptors) - 1; index >= 0; index-- {
		interceptor := c.interceptors[index]
		actionWrapper = (func(action func(context.Context, map[string]any) (any, error)) func(context.Context, map[string]any) (any, error) {
			return func(ctx context.Context, params map[string]any) (any, error) {
				return interceptor(ctx, params, action)
			}
		})(actionWrapper)
	}

	return actionWrapper
}

func (c *LambdaController) GenerateActionCmd(name string) string {
	cmd := name
	if c.name != "" {
		cmd = c.name + "." + cmd
	}
	return cmd
}

// Registers a action in AWS Lambda function.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		- name          an action name
//		- schema        a validation schema to validate received parameters.
//		- action        an action function that is called when operation is invoked.
func (c *LambdaController) RegisterAction(name string, schema *cvalid.Schema, action func(ctx context.Context, params map[string]any) (any, error)) {
	actionWrapper := c.ApplyValidation(schema, action)
	actionWrapper = c.ApplyInterceptors(actionWrapper)

	registeredAction := &LambdaAction{
		Cmd:    c.GenerateActionCmd(name),
		Schema: schema,
		Action: func(ctx context.Context, params map[string]any) (any, error) {
			return actionWrapper(ctx, params)
		},
	}
	c.actions = append(c.actions, registeredAction)
}

// Registers an action with authorization.
//
//	Parameters:
//		-  name          an action name
//		-  schema        a validation schema to validate received parameters.
//		-  authorize     an authorization interceptor
//		-  action        an action function that is called when operation is invoked.
func (c *LambdaController) RegisterActionWithAuth(name string, schema *cvalid.Schema,
	authorize func(ctx context.Context, params map[string]any, next func(context.Context, map[string]any) (any, error)) (any, error),
	action func(ctx context.Context, params map[string]any) (any, error)) {

	actionWrapper := c.ApplyValidation(schema, action)
	// Add authorization just before validation
	actionWrapper = func(ctx context.Context, params map[string]any) (any, error) {
		return authorize(ctx, params, actionWrapper)
	}
	actionWrapper = c.ApplyInterceptors(actionWrapper)

	registeredAction := &LambdaAction{
		Cmd:    c.GenerateActionCmd(name),
		Schema: schema,
		Action: func(ctx context.Context, params map[string]any) (any, error) {
			return actionWrapper(ctx, params)
		},
	}
	c.actions = append(c.actions, registeredAction)
}

// Registers a middleware for actions in AWS Lambda service.
// -  action        an action function that is called when middleware is invoked.
func (c *LambdaController) RegisterInterceptor(action func(ctx context.Context, params map[string]any, next func(ctx context.Context, params map[string]any) (any, error)) (any, error)) {
	c.interceptors = append(c.interceptors, action)
}

// Registers all service routes in HTTP endpoint.
// This method is called by the service and must be overriden
// in child classes.
func (c *LambdaController) Register() {
	c.Overrides.Register()
}

// Calls registered action in this lambda function.
// "cmd" parameter in the action parameters determin
// what action shall be called.
// This method shall only be used in testing.
//
//	Parameters:
//		- ctx context.Context	operation context.
//		-  params action parameters.
func (c *LambdaController) Act(ctx context.Context, params map[string]any) (any, error) {
	cmd, ok := params["cmd"].(string)
	traceId, _ := params["trace_id"].(string)

	if !ok || cmd == "" {
		return nil, cerr.NewBadRequestError(
			traceId,
			"NO_COMMAND",
			"Cmd parameter is missing",
		)
	}

	var action *LambdaAction
	for _, act := range c.actions {
		if act.Cmd == cmd {
			action = act
			break
		}
	}

	if action == nil {
		return nil, cerr.NewBadRequestError(
			traceId,
			"NO_ACTION",
			"Action "+cmd+" was not found",
		).
			WithDetails("command", cmd)
	}

	return action.Action(ctx, params)
}
