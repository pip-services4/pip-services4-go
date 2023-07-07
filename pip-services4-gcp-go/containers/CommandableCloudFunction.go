package containers

import (
	"context"
	"net/http"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	gcputil "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/utils"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

// Abstract Google Function function, that acts as a container to instantiate and run components
// and expose them via external entry point. All actions are automatically generated for commands
// defined in ICommandable components. Each command is exposed as an action defined by "cmd" parameter.
//
// Container configuration for this Google Function is stored in "./config/config.yml" file.
// But this path can be overridden by <code>CONFIG_PATH</code> environment variable.
//
//	References
//		- *:logger:*:*:1.0							(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0						(optional) ICounters components to pass collected measurements
//		- *:controller:cloudfunc:*:1.0       		(optional) ICloudFunctionController controllers to handle action requests
//		- *:controller:commandable-cloudfunc:*:1.0	(optional) ICloudFunctionController controllers to handle action requests
//
//	Example:
//		type MyCloudFunction struct {
//			*containers.CommandableCloudFunction
//			controller IMyController
//		}
//
//		func NewMyCloudFunction() *MyCloudFunction {
//			c := MyCloudFunction{}
//			c.CloudFunction = containers.NewCommandableCloudFunctionWithParams("mygroup", "MyGroup CloudFunction")
//
//			return &c
//		}
//
//		...
//
//		cloudFunction := NewMyCloudFunction()
//		cloudFunction.Run(ctx)
//		fmt.Println("MyCloudFunction is started")
//
// Deprecated: This component has been deprecated. Use CloudFunctionService instead.
type CommandableCloudFunction struct {
	*CloudFunction
}

// Creates a new instance of this Google Function.
func NewCommandableCloudFunction() *CommandableCloudFunction {
	c := CommandableCloudFunction{}
	c.CloudFunction = InheritCloudFunction(&c)
	return &c
}

// Creates a new instance of this Google Function.
// Parameters:
//   - name	(optional) a container name (accessible via ContextInfo)
//   - description	(optional) a container description (accessible via ContextInfo)
func NewCommandableCloudFunctionWithParams(name string, description string) *CommandableCloudFunction {
	c := CommandableCloudFunction{}
	c.CloudFunction = InheritCloudFunctionWithParams(&c, name, description)
	return &c
}

// Returns body from Google Function request.
// This method can be overloaded in child classes
// Parameters:
//   - req	Googl Function request
//
// Returns Parameters from request
func (c *CommandableCloudFunction) GetParameters(req *http.Request) *cexec.Parameters {
	return gcputil.CloudFunctionRequestHelper.GetParameters(req)
}

func (c *CommandableCloudFunction) registerCommandSet(commandSet *ccomand.CommandSet) {
	commands := commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		c.RegisterAction(command.Name(), nil, func(w http.ResponseWriter, r *http.Request) {
			traceId := c.GetTraceId(r)
			ctx := cctx.NewContextWithTraceId(r.Context(), traceId)
			args := c.GetParameters(r)

			timing := c.Instrument(ctx, command.Name())
			execRes, execErr := command.Execute(ctx, args)
			timing.EndTiming(ctx, execErr)

			httpctrl.HttpResponseSender.SendResult(w, r, execRes, execErr)
		})
	}
}

// Registers all actions in this Google Function.
//
// Deprecated: Overloading of this method has been deprecated. Use CloudFunctionService instead.
func (c *CommandableCloudFunction) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		panic(depErr)
	}

	service, ok := resCtrl.(ccomand.ICommandable)
	if !ok {
		c.Logger().Error(cctx.NewContextWithTraceId(context.Background(), "CommandableCloudController"), nil, "Can't cast Controller to ICommandable")
		return
	}

	commandSet := service.GetCommandSet()
	c.registerCommandSet(commandSet)
}
