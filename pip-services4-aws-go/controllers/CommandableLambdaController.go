package controllers

import (
	"context"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	ccomands "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

// Abstract service that receives commands via AWS Lambda protocol
// to operations automatically generated for commands defined in ICommandable components.
// Each command is exposed as invoke method that receives command name and parameters.
//
// Commandable services require only 3 lines of code to implement a robust external
// Lambda-based remote interface.
//
// This service is intended to work inside LambdaFunction container that
// exploses registered actions externally.
//
// # Configuration parameters
//
// - dependencies:
//   - controller:            override for Controller dependency
//
// References
//
//   - *:logger:*:*:1.0               (optional) ILogger components to pass log messages
//   - *:counters:*:*:1.0             (optional) ICounters components to pass collected measurements
//
// See CommandableLambdaClient
// See LambdaController
//
// Example:
//
//	type MyCommandableLambdaController struct  {
//		*CommandableLambdaController
//	}
//
//	func NewMyCommandableLambdaController() *MyCommandableLambdaController {
//		c:= &MyCommandableLambdaController{
//			CommandableLambdaController: NewCommandableLambdaController("v1.service")
//		}
//		c.DependencyResolver.Put(context.Background(),
//			"controller",
//			cref.NewDescriptor("mygroup","controller","*","*","1.0")
//		)
//		return c
//	}
//
//	service := NewMyCommandableLambdaController();
//	service.SetReferences(context.Background(), NewReferencesFromTuples(
//	   NewDescriptor("mygroup","controller","default","default","1.0"), controller
//	))
//
//	service.Open(context.Background(),"123")
//	fmt.Println("The AWS Lambda 'v1.service' service is running")
type CommandableLambdaController struct {
	*LambdaController
	commandSet *ccomands.CommandSet
}

// Creates a new instance of the service.
// - name a service name.
func InheritCommandableLambdaController(overrides ILambdaControllerOverrides, name string) *CommandableLambdaController {
	c := &CommandableLambdaController{
		LambdaController: InheritLambdaController(overrides, name),
	}

	c.DependencyResolver.Put(context.Background(), "service", "none")
	return c
}

// Registers all actions in AWS Lambda function.
func (c *CommandableLambdaController) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		return
	}
	controller, ok := resCtrl.(ccomands.ICommandable)
	if !ok {
		c.Logger.Error(cctx.NewContextWithTraceId(context.Background(), "CommandableLambdaController"),
			nil, "Can't cast Controller to ICommandable")
		return
	}

	c.commandSet = controller.GetCommandSet()

	commands := c.commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]
		name := command.Name()
		c.RegisterAction(name, nil, func(ctx context.Context, params map[string]any) (any, error) {
			traceId, _ := params["trace_id"].(string)
			ctx = cctx.NewContextWithTraceId(ctx, traceId)
			args := cexec.NewParametersFromValue(params)
			args.Remove("trace_id")

			timing := c.Instrument(ctx, name)
			result, err := command.Execute(ctx, args)
			timing.EndTiming(ctx, err)
			return result, err

		})
	}
}
