package controllers

import (
	"context"
	"net/http"

	azureutil "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/utils"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

// Abstract service that receives commands via Azure Function protocol
// to operations automatically generated for commands defined in ccomand.ICommandable components.
// Each command is exposed as invoke method that receives command name and parameters.
//
// Commandable services require only 3 lines of code to implement a robust external
// Azure Function-based remote interface.
//
// This service is intended to work inside Azure Function container that
// exploses registered actions externally.
//
//	Configuration parameters:
//		- dependencies:
//			- service:            override for Service dependency
//	References
//		- *:logger:*:*:1.0			(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0		(optional) ICounters components to pass collected measurements
//
// see AzureFunctionService
//
//	Example:
//		type MyCommandableAzureFunctionController struct {
//			*azuresrv.CommandableAzureFunctionController
//		}
//
//		func NewMyCommandableAzureFunctionController() *MyCommandableAzureFunctionController {
//			c := MyCommandableAzureFunctionController{}
//			c.CommandableAzureFunctionController = azuresrv.NewCommandableAzureFunctionService("mydata")
//			c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("mygroup", "service", "default", "*", "*"))
//			return &c
//		}
//
//		...
//
//		service := NewMyCommandableAzureFunctionController()
//		service.SetReferences(crefer.NewReferencesFromTuples(
//			crefer.NewDescriptor("mygroup","service","default","default","1.0"), service,
//		))
//		service.Open(ctx, "123")
//		fmt.Println("The Azure Function controller is running")
type CommandableAzureFunctionController struct {
	*AzureFunctionController
	commandSet *ccomand.CommandSet
}

// Creates a new instance of the service.
// Parameters:
//   - name 	a service name.
func NewCommandableAzureFunctionController(name string) *CommandableAzureFunctionController {
	c := CommandableAzureFunctionController{}
	c.AzureFunctionController = InheritAzureFunctionController(&c, name)

	return &c
}

// Returns body from Azure Function request.
// This method can be overloaded in child classes
// Parameters:
//   - req	Azure Function request
//
// Returns Parameters from request
func (c *CommandableAzureFunctionController) GetParameters(req *http.Request) *cexec.Parameters {
	return azureutil.AzureFunctionRequestHelper.GetParameters(req)
}

// Registers all actions in Azure Function.
func (c *CommandableAzureFunctionController) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		panic(depErr)
	}

	service, ok := resCtrl.(ccomand.ICommandable)
	if !ok {
		c.Logger.Error(cctx.NewContextWithTraceId(context.Background(), "CommandableHttpService"), nil, "Can't cast Service to ICommandable")
		return
	}

	c.commandSet = service.GetCommandSet()
	commands := c.commandSet.Commands()

	for index := 0; index < len(commands); index++ {
		command := commands[index]
		name := command.Name()

		c.RegisterAction(name, nil, func(w http.ResponseWriter, r *http.Request) {
			traceId := c.GetTraceId(r)
			args := c.GetParameters(r)
			args.Remove("trace_id")

			ctx := cctx.NewContextWithTraceId(r.Context(), traceId)
			timing := c.Instrument(ctx, name)
			execRes, execErr := command.Execute(ctx, args)
			timing.EndTiming(ctx, execErr)
			httpctrl.HttpResponseSender.SendResult(w, r, execRes, execErr)
		})
	}
}
