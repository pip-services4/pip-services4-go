package containers

import (
	"context"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	ccomands "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

// Abstract AWS Lambda function, that acts as a container to instantiate and run components
// and expose them via external entry point. All actions are automatically generated for commands
// defined in ICommandable components. Each command is exposed as an action defined by "cmd" parameter.
//
// Container configuration for this Lambda function is stored in "./config/config.yml" file.
// But this path can be overriden by CONFIG_PATH environment variable.
//
// Configuration parameters ###
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
//   - \*:controller:awslambda:\*:1.0       		(optional) ILambdaController controllers to handle action requests
//   - \*:controller:commandable-awslambda:\*:1.0	(optional) ILambdaController controllers to handle action requests
//
// # See LambdaClient
//
// Example:
//
//	type MyCommandableLambdaFunction struct {
//		*awscont.CommandableLambdaFunction
//	}
//
//	func NewMyCommandableLambdaFunction() *MyCommandableLambdaFunction {
//		c := &MyCommandableLambdaFunction{}
//		c.CommandableLambdaFunction = awscont.NewCommandableLambdaFunction("my_group", "My data lambda function")
//
//		c.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("my-group", "controller", "default", "*", "*"))
//		c.AddFactory(awstest.NewDummyFactory())
//		return c
//	}
//
//	lambda := NewMyCommandableLambdaFunction();
//
//	lambda.Run(context.Context())
type CommandableLambdaFunction struct {
	*LambdaFunction
}

// Creates a new instance of this lambda function.
//
//   - name          (optional) a container name (accessible via ContextInfo)
//   - description   (optional) a container description (accessible via ContextInfo)
func NewCommandableLambdaFunction(name string, description string) *CommandableLambdaFunction {
	c := &CommandableLambdaFunction{}
	c.LambdaFunction = InheriteLambdaFunction(c, name, description)
	c.DependencyResolver.Put(context.Background(), "service", "none")
	return c
}

func (c *CommandableLambdaFunction) registerCommandSet(commandSet *ccomands.CommandSet) {
	commands := commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		c.RegisterAction(command.Name(), nil, func(ctx context.Context, params map[string]any) (result any, err error) {
			traceId, _ := params["trace_id"].(string)

			args := cexec.NewParametersFromValue(params)
			timing := c.Instrument(
				cctx.NewContextWithTraceId(context.Background(), traceId),
				c.Info().Name+"."+command.Name())
			result, errRes := command.Execute(ctx, args)
			timing.EndTiming(ctx, errRes)
			return result, errRes
		})
	}
}

// Registers all actions in this lambda function.
func (c *CommandableLambdaFunction) Register() {

	ref, _ := c.DependencyResolver.GetOneRequired("service")
	controller, ok := ref.(ccomands.ICommandable)
	if ok {
		commandSet := controller.GetCommandSet()
		c.registerCommandSet(commandSet)
	}
}
