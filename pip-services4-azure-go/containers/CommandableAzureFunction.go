package containers

import (
	"context"
	"net/http"

	azureutil "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/utils"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

// Abstract Azure Function function, that acts as a container to instantiate and run components
// and expose them via external entry point. All actions are automatically generated for commands
// defined in ICommandable components. Each command is exposed as an action defined by "cmd" parameter.
//
// Container configuration for this Azure Function is stored in "./config/config.yml" file.
// But this path can be overridden by <code>CONFIG_PATH</code> environment variable.
//
//	References
//		- *:logger:*:*:1.0							(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0						(optional) ICounters components to pass collected measurements
//		- *:controller:azurefunc:*:1.0       		(optional) IAzureFunctionController controllers to handle action requests
//		- *:controller:commandable-azurefunc:*:1.0	(optional) IAzureFunctionController controllers to handle action requests
//
//	Example:
//		type MyAzureFunction struct {
//			*containers.CommandableAzureFunction
//			controller IMyController
//		}
//
//		func NewMyAzureFunction() *MyAzureFunction {
//			c := MyAzureFunction{}
//			c.AzureFunction = containers.NewCommandableAzureFunctionWithParams("mygroup", "MyGroup AzureFunction")
//
//			return &c
//		}
//
//		...
//
//		AzureFunction := NewMyAzureFunction()
//		AzureFunction.Run(ctx)
//		fmt.Println("MyAzureFunction is started")
//
// Deprecated: This component has been deprecated. Use AzureFunctionService instead.
type CommandableAzureFunction struct {
	*AzureFunction
}

// Creates a new instance of this Azure Function.
func NewCommandableAzureFunction() *CommandableAzureFunction {
	c := CommandableAzureFunction{}
	c.AzureFunction = InheritAzureFunction(&c)
	return &c
}

// Creates a new instance of this Azure Function.
// Parameters:
//   - name	(optional) a container name (accessible via ContextInfo)
//   - description	(optional) a container description (accessible via ContextInfo)
func NewCommandableAzureFunctionWithParams(name string, description string) *CommandableAzureFunction {
	c := CommandableAzureFunction{}
	c.AzureFunction = InheritAzureFunctionWithParams(&c, name, description)
	return &c
}

// Returns body from Azure Function request.
// This method can be overloaded in child classes
// Parameters:
//   - req	Googl Function request
//
// Returns Parameters from request
func (c *CommandableAzureFunction) GetParameters(req *http.Request) *cexec.Parameters {
	return azureutil.AzureFunctionRequestHelper.GetParameters(req)
}

func (c *CommandableAzureFunction) registerCommandSet(commandSet *ccomand.CommandSet) {
	commands := commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		c.RegisterAction(command.Name(), nil, func(w http.ResponseWriter, r *http.Request) {
			traceId := c.GetTraceId(r)
			ctx := utils.ContextHelper.NewContextWithTraceId(r.Context(), traceId)
			args := c.GetParameters(r)

			timing := c.Instrument(ctx, command.Name())
			execRes, execErr := command.Execute(ctx, args)
			timing.EndTiming(ctx, execErr)

			httpctrl.HttpResponseSender.SendResult(w, r, execRes, execErr)
		})
	}
}

// Registers all actions in this Azure Function.
//
// Deprecated: Overloading of this method has been deprecated. Use AzureFunctionService instead.
func (c *CommandableAzureFunction) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		panic(depErr)
	}

	controller, ok := resCtrl.(ccomand.ICommandable)
	if !ok {
		c.Logger().Error(utils.ContextHelper.NewContextWithTraceId(context.Background(), "CommandableAzureController"),
			nil, "Can't cast Service to ICommandable")
		return
	}

	commandSet := controller.GetCommandSet()
	c.registerCommandSet(commandSet)
}
