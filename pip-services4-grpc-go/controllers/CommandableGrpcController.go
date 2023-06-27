package controllers

import (
	"context"

	crun "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	ccomands "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

// CommandableGrpcController abstract service that receives commands via GRPC protocol
// to operations automatically generated for commands defined in ICommandable components.
// Each command is exposed as invoke method that receives command name and parameters.
//
// Commandable services require only 3 lines of code to implement a robust external
// GRPC-based remote interface.
//
//	Configuration parameters:
//
//		- dependencies:
//			- endpoint:              override for HTTP Endpoint dependency
//			- controller:            override for Controller dependency
//		- connection(s):
//			- discovery_key:         (optional) a key to retrieve the connection from  IDiscovery
//			- protocol:              connection protocol: http or https
//			- host:                  host name or IP address
//			- port:                  port number
//			- uri:                   resource URI or connection string with all parameters in it
//
//	References:
//
//		- *:logger:*:*:1.0               (optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0             (optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0            (optional) IDiscovery services to resolve connection
//		- *:endpoint:grpc:*:1.0          (optional) GrpcEndpoint reference
//
// See CommandableGrpcClient
// See GrpcService
//
// Example:

//	type MyCommandableGrpcController struct {
//		*CommandableGrpcController
//	}
//
//	func NewCommandableGrpcController() *CommandableGrpcController {
//		c := DumMyCommandableGrpcController{}
//		c.CommandableGrpcController = grpcservices.NewCommandableGrpcController("mycontroller")
//		c.DependencyResolver.Put("service", cref.NewDescriptor("mygroup", "service", "default", "*", "*"))
//		return &c
//	}
//
// controller := NewMyCommandableGrpcController();
// controller.Configure(ctx, cconf.NewConfigParamsFromTuples(
//
//	"connection.protocol", "http",
//	"connection.host", "localhost",
//	"connection.port", "8080",
//
// ));
// controller.SetReferences(ctx, cref.NewReferencesFromTuples(
//
//	cref.NewDescriptor("mygroup","controller","default","default","1.0"), service
//
// ));
//
// opnErr := controller.Open(ctx)
//
//	if opnErr == nil {
//		fmt.Println("The GRPC controller is running on port 8080");
//	}
type CommandableGrpcController struct {
	*GrpcController
	name       string
	commandSet *ccomands.CommandSet
}

// InheritCommandableGrpcController method are creates a new instance of the service.
//   - name a service name.
func InheritCommandableGrpcController(overrides IGrpcControllerOverrides, name string) *CommandableGrpcController {
	c := &CommandableGrpcController{}
	c.GrpcController = InheritGrpcService(overrides, "")
	c.name = name
	c.DependencyResolver.Put(context.Background(), "service", "none")
	return c
}

// Register method are registers all service command in gRPC endpoint.
func (c *CommandableGrpcController) Register() {

	resCtrl, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		return
	}
	controller, ok := resCtrl.(ccomands.ICommandable)
	if !ok {
		c.Logger.Error(utils.ContextHelper.NewContextWithTraceId(context.Background(), "CommandableHttpController"),
			nil, "Can't cast Controller to ICommandable")
		return
	}
	c.commandSet = controller.GetCommandSet()

	commands := c.commandSet.Commands()
	var index = 0
	for index = 0; index < len(commands); index++ {
		command := commands[index]

		method := c.name + "." + command.Name()

		c.RegisterCommandableMethod(method, nil,
			func(ctx context.Context, args *crun.Parameters) (result any, err error) {
				timing := c.Instrument(ctx, method)
				res, err := command.Execute(ctx, args)
				timing.EndTiming(ctx, err)
				return res, err
			})
	}
}
