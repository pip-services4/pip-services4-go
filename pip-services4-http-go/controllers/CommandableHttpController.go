package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	ccomands "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
	"goji.io/pattern"
)

// CommandableHttpController abstract service that receives remove calls via HTTP/REST protocol
// to operations automatically generated for commands defined in ICommandable components.
// Each command is exposed as POST operation that receives all parameters in body object.
//
// Commandable services require only 3 lines of code to implement a robust external
// HTTP-based remote interface.
//
//	Configuration parameters:
//		- base_route:                base route for remote URI
//		- dependencies:
//			- endpoint:              override for HTTP Endpoint dependency
//			- controller:            override for Controller dependency
//		- connection(s):
//			- discovery_key:         (optional) a key to retrieve the connection from IDiscovery
//			- protocol:              connection protocol: http or https
//			- host:                  host name or IP address
//			- port:                  port number
//			- uri:                   resource URI or connection string with all parameters in it
//
//	References:
//		- *:logger:*:*:1.0            (optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0          (optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0         (optional) IDiscovery services to resolve connection
//		- *:endpoint:http:*:1.0       (optional) HttpEndpoint reference
//
//	see clients.CommandableHttpClient
//	see RestController
//
//	Example:
//		type MyCommandableHttpController struct {
//			*CommandableHttpController
//		}
//
//		func NewMyCommandableHttpController() *MyCommandableHttpController {
//			c := MyCommandableHttpController{
//				CommandableHttpController: controllers.NewCommandableHttpController("dummies"),
//			}
//			c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
//			return &c
//		}
//
//		controller := NewMyCommandableHttpController();
//		controller.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
//			"connection.protocol", "http",
//			"connection.host", "localhost",
//			"connection.port", 8080,
//		));
//		controller.SetReferences(context.Background(), cref.NewReferencesFromTuples(
//			cref.NewDescriptor("mygroup","service","default","default","1.0"), service
//		));
//
//		opnErr := controller.Open(context.Background(), "123")
//		if opnErr == nil {
//			fmt.Println("The REST controller is running on port 8080");
//		}
type CommandableHttpController struct {
	*RestController
	commandSet  *ccomands.CommandSet
	SwaggerAuto bool
}

// NewCommandableHttpController creates a new instance of the controller.
//	Parameters:
//		- baseRoute string a controller base route.
//	Returns: *CommandableHttpController pointer on new instance CommandableHttpController
//	func NewCommandableHttpController(baseRoute string) *CommandableHttpController {
//		c := &CommandableHttpController{}
//		c.RestController = InheritRestController(c)
//		c.BaseRoute = baseRoute
//		c.SwaggerAuto = true
//		c.DependencyResolver.Put(context.Background(), "service", "none")
//		return c
//	}

// InheritCommandableHttpController creates a new instance of the controller.
//
//	Parameters:
//		- overrides references to child class that overrides virtual methods
//		- baseRoute string a service base route.
//	Returns: *CommandableHttpController pointer on new instance CommandableHttpController
func InheritCommandableHttpController(overrides IRegisterable, baseRoute string) *CommandableHttpController {
	c := &CommandableHttpController{}
	c.RestController = InheritRestController(overrides)
	c.BaseRoute = baseRoute
	c.SwaggerAuto = true
	c.DependencyResolver.Put(context.Background(), "service", "none")
	return c
}

// Configure method configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config configuration parameters to be set.
func (c *CommandableHttpController) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.RestController.Configure(ctx, config)
	c.SwaggerAuto = config.GetAsBooleanWithDefault("swagger.auto", c.SwaggerAuto)
}

// Register method are registers all service routes in HTTP endpoint.
func (c *CommandableHttpController) Register() {
	resServ, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr != nil {
		return
	}
	controller, ok := resServ.(ccomands.ICommandable)
	if !ok {
		c.Logger.Error(
			cctx.NewContextWithTraceId(context.Background(), "CommandableHttpController"),
			nil, "Can't cast Controller to ICommandable")
		return
	}
	c.commandSet = controller.GetCommandSet()

	commands := c.commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		route := command.Name()
		if route[0] != "/"[0] {
			route = "/" + route
		}

		c.RegisterRoute(http.MethodPost, route, nil, func(res http.ResponseWriter, req *http.Request) {

			// Make copy of request
			bodyBuf, bodyErr := io.ReadAll(req.Body)
			if bodyErr != nil {
				HttpResponseSender.SendError(res, req, bodyErr)
				return
			}
			_ = req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBuf))
			//-------------------------
			// TODO:: think about marshaling and error
			var params map[string]any = make(map[string]any, 0)
			json.Unmarshal(bodyBuf, &params)

			urlParams := req.URL.Query()
			for k, v := range urlParams {
				params[k] = v[0]
			}
			if reqVars, ok := req.Context().Value(pattern.AllVariables).(map[pattern.Variable]any); ok {
				for k, v := range reqVars {
					params[string(k)] = v
				}
			}

			traceId := c.GetTraceId(req)
			ctx := cctx.NewContextWithTraceId(req.Context(), traceId)
			args := cexec.NewParametersFromValue(params)
			timing := c.Instrument(ctx, c.BaseRoute+"."+command.Name())

			execRes, execErr := command.Execute(ctx, args)
			timing.EndTiming(ctx, execErr)
			c.SendResult(res, req, execRes, execErr)
		})
	}

	if c.SwaggerAuto {
		var swaggerConfig = c.config.GetSection("swagger")
		var doc = NewCommandableSwaggerDocument(c.BaseRoute, swaggerConfig, commands)
		c.RegisterOpenApiSpec(doc.ToString())
	}
}
