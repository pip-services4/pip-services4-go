package services

import (
	"context"
	"net/http"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	gcputil "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/utils"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

// Abstract service that receives commands via Google Function protocol
// to operations automatically generated for commands defined in ccomand.ICommandable components.
// Each command is exposed as invoke method that receives command name and parameters.
//
// Commandable services require only 3 lines of code to implement a robust external
// Google Function-based remote interface.
//
// This service is intended to work inside Google Function container that
// exploses registered actions externally.
//
//	Configuration parameters:
//		- dependencies:
//			- service:            override for Service dependency
//	References
//		- *:logger:*:*:1.0			(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0		(optional) ICounters components to pass collected measurements
//
// see CloudFunctionService
//
//	Example:
//		type MyCommandableCloudFunctionController struct {
//			*gcpsrv.CommandableCloudFunctionController
//		}
//
//		func NewMyCommandableCloudFunctionController() *MyCommandableCloudFunctionController {
//			c := MyCommandableCloudFunctionController{}
//			c.CommandableCloudFunctionController = gcpsrv.NewCommandableCloudFunctionService("mydata")
//			c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("mygroup", "service", "default", "*", "*"))
//			return &c
//		}
//
// /		...
//
//	service := NewMyCommandableCloudFunctionController()
//	service.SetReferences(crefer.NewReferencesFromTuples(
//		crefer.NewDescriptor("mygroup","controller","default","default","1.0"), controller,
//	))
//	service.Open(ctx, "123")
//	fmt.Println("The Google Function service is running")
type CommandableCloudFunctionController struct {
	*CloudFunctionController
	commandSet *ccomand.CommandSet
}

// Creates a new instance of the service.
// Parameters:
//   - name 	a service name.
func NewCommandableCloudFunctionController(name string) *CommandableCloudFunctionController {
	c := CommandableCloudFunctionController{}
	c.CloudFunctionController = InheritCloudFunctionController(&c, name)
	return &c
}

// Returns body from Google Function request.
// This method can be overloaded in child classes
// Parameters:
//   - req	Google Function request
//
// Returns Parameters from request
func (c *CommandableCloudFunctionController) GetParameters(req *http.Request) *cexec.Parameters {
	return gcputil.CloudFunctionRequestHelper.GetParameters(req)
}

// Registers all actions in Google Function.
func (c *CommandableCloudFunctionController) Register() {
	resServ, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		panic(depErr)
	}

	service, ok := resServ.(ccomand.ICommandable)
	if !ok {
		c.Logger.Error(cctx.NewContextWithTraceId(context.Background(), "CommandableCloudController"), nil, "Can't cast Service to ICommandable")
		return
	}

	c.commandSet = service.GetCommandSet()
	commands := c.commandSet.Commands()

	for index := 0; index < len(commands); index++ {
		command := commands[index]
		name := command.Name()

		c.RegisterAction(name, nil, func(w http.ResponseWriter, r *http.Request) {
			traceId := c.GetTraceId(r)
			ctx := cctx.NewContextWithTraceId(r.Context(), traceId)
			args := c.GetParameters(r)
			args.Remove("trace_id")

			timing := c.Instrument(ctx, name)
			execRes, execErr := command.Execute(ctx, args)
			timing.EndTiming(ctx, execErr)
			httpctrl.HttpResponseSender.SendResult(w, r, execRes, execErr)
		})
	}
}
