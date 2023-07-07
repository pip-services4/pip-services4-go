package containers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	ccont "github.com/pip-services4/pip-services4-go/pip-services4-container-go/container"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	gcpctrl "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/controllers"
	gcputil "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/utils"
	httpctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ctrace "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
	rpctrace "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/trace"
	"goji.io/pattern"
)

type ICloudFunctionOverrides interface {
	crefer.IReferenceable
	// Registers all actions in this Google Function.
	//
	// Depecated: Overloading of this method has been deprecated. Use CloudFunctionService instead.
	Register()
}

// Abstract Google Function, that acts as a container to instantiate and run components
// and expose them via external entry point.
//
// When handling calls "cmd" parameter determines which what action shall be called, while
// other parameters are passed to the action itself.
//
// Container configuration for this Google Function is stored in "./config/config.yml" file.
// But this path can be overriden by CONFIG_PATH environment variable.
//
//	References
//		- *:logger:*:*:1.0							(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0						(optional) ICounters components to pass collected measurements
//		- *:controller:cloudfunc:*:1.0				(optional) ICloudFunctionController controllers to handle action requests
//		- *:controller:commandable-cloudfunc:*:1.0	(optional) ICloudFunctionController controllers to handle action requests
//
//	Example:
//		type MyCloudFunction struct {
//			*containers.CloudFunction
//		}
//
//		func NewMyCloudFunction() *MyCloudFunction {
//			c := MyCloudFunction{}
//			c.CloudFunction = containers.NewCloudFunctionWithParams("mygroup", "MyGroup Google Function")
//
//			return &c
//		}
//
//		...
//
//		cloudFunction := NewMyCloudFunction()
//		cloudFunction.Run(ctx)
//		fmt.Println("MyCloudFunction is started")
type CloudFunction struct {
	*ccont.Container

	Overrides ICloudFunctionOverrides
	// The dependency resolver.
	DependencyResolver *crefer.DependencyResolver
	// The performanc counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
	// The map of registred validation schemas.
	Schemas map[string]*cvalid.Schema
	// The map of registered actions.
	Actions map[string]http.HandlerFunc

	feedbackChan          cctx.ContextShutdownChan
	feedbackWithErrorChan cctx.ContextShutdownWithErrorChan

	// The default path to config file.
	configPath string
}

// Creates a new instance of this Google Function function.
func NewCloudFunction() *CloudFunction {
	c := CloudFunction{
		DependencyResolver: crefer.NewDependencyResolver(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
		Schemas:            make(map[string]*cvalid.Schema),
		Actions:            make(map[string]http.HandlerFunc),
		configPath:         "./config/config.yml",
	}

	c.Container = ccont.InheritContainer("", "", &c)
	c.Overrides = &c
	c.SetLogger(clog.NewConsoleLogger())

	return &c
}

// Creates a new instance of this Google Function function.
// Parameters:
//   - name		(optional) a container name (accessible via ContextInfo)
//   - description		(optional) a container description (accessible via ContextInfo)
func NewCloudFunctionWithParams(name string, description string) *CloudFunction {
	c := CloudFunction{
		DependencyResolver: crefer.NewDependencyResolver(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
		Schemas:            make(map[string]*cvalid.Schema),
		Actions:            make(map[string]http.HandlerFunc),
		configPath:         "./config/config.yml",
	}

	c.Container = ccont.InheritContainer(name, description, &c)
	c.Overrides = &c
	c.SetLogger(clog.NewConsoleLogger())

	return &c
}

// InheritCloudFunction creates new instance of CloudFunction
func InheritCloudFunction(overrides ICloudFunctionOverrides) *CloudFunction {
	c := CloudFunction{
		Overrides:          overrides,
		DependencyResolver: crefer.NewDependencyResolver(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
		Schemas:            make(map[string]*cvalid.Schema),
		Actions:            make(map[string]http.HandlerFunc),
		configPath:         "./config/config.yml",
	}

	c.Container = ccont.InheritContainer("", "", overrides)
	c.SetLogger(clog.NewConsoleLogger())

	return &c
}

// InheritCloudFunction creates new instance of CloudFunction
// Parameters:
//   - name		(optional) a container name (accessible via ContextInfo)
//   - description		(optional) a container description (accessible via ContextInfo)
func InheritCloudFunctionWithParams(overrides ICloudFunctionOverrides, name string, description string) *CloudFunction {
	c := CloudFunction{
		Overrides:          overrides,
		DependencyResolver: crefer.NewDependencyResolver(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
		Schemas:            make(map[string]*cvalid.Schema),
		Actions:            make(map[string]http.HandlerFunc),
		configPath:         "./config/config.yml",
	}

	c.Container = ccont.InheritContainer("", "", overrides)
	c.SetLogger(clog.NewConsoleLogger())

	return &c
}

func (c *CloudFunction) Register() {}

// SetConfigPath set path for configuration file
// Parameters:
//   - configPath	path to config file
func (c *CloudFunction) SetConfigPath(configPath string) {
	c.configPath = configPath
}

func (c *CloudFunction) getConfigPath() string {
	env := os.Getenv("CONFIG_PATH")
	if env != "" {
		return env
	}

	return c.configPath
}

func (c *CloudFunction) getConfigParameters() *cconf.ConfigParams {
	args := os.Environ()

	line := ""

	for index := 0; index < len(args); index++ {
		arg := args[index]
		nextArg := ""
		if index < len(args)-1 {
			nextArg = args[index+1]
			if strings.HasPrefix(nextArg, "-") {
				nextArg = ""
			}
		}

		if nextArg != "" {
			if arg == "--param" || arg == "--params" || arg == "-p" {
				if line != "" {
					line = line + ";"
				}
				line = line + nextArg
				index++
			}
		}
	}

	parameters := cconf.NewConfigParamsFromString(line)

	for _, e := range os.Environ() {
		if env := strings.Split(e, "="); len(env) == 2 {
			parameters.SetAsObject(env[0], env[1])
		} else {
			parameters.SetAsObject(env[0], strings.Join(env[1:], "="))
		}
	}

	return parameters
}

// SetReferences sets references to dependent components.
//
//	see IReferences
//	Parameters:
//		- ctx context.Context
//		- references IReferences references to locate the component dependencies.
func (c *CloudFunction) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Counters.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)

	c.Overrides.Register()
}

// Open opens the component.
//
//	Parameters:
//		- ctx context.Context a context to trace execution through call chain.
//	Return: error
func (c *CloudFunction) Open(ctx context.Context) error {
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

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a Timing object that is used to end the time measurement.
//
//	Parameters:
//		- ctx context.Context a context to trace execution through call chain.
//		- name              a method name.
//	Returns: Timing object to end the time measurement.
func (c *CloudFunction) Instrument(ctx context.Context, name string) *rpctrace.InstrumentTiming {
	c.Logger().Trace(ctx, "Executing %s method", name)
	c.Counters.IncrementOne(ctx, name+".exec_count")

	counterTiming := c.Counters.BeginTiming(ctx, name+".exec_time")
	traceTiming := c.Tracer.BeginTrace(ctx, name, "")

	return rpctrace.NewInstrumentTiming(ctx, name, "exec",
		c.Logger(), c.Counters, counterTiming, traceTiming)
}

// Runs this Google Function, loads container configuration,
// instantiate components and manage their lifecycle,
// makes this function ready to access action calls.
//
//	Parameters:
//		- ctx context.Context
func (c *CloudFunction) Run(ctx context.Context) {
	traceId := c.Info().Name

	ctx, cancel := context.WithCancel(ctx)

	ctx, _ = cctx.AddShutdownChanToContext(ctx, c.feedbackChan)
	ctx, _ = cctx.AddErrShutdownChanToContext(ctx, c.feedbackWithErrorChan)
	ctx = cctx.NewContextWithTraceId(ctx, traceId)

	path := c.getConfigPath()
	parameters := c.getConfigParameters()

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				msg := cconv.StringConverter.ToString(r)
				err = errors.New(msg)
			}
			_ = c.Close(ctx)
			cancel()
			c.Logger().Fatal(ctx, err, "Process is terminated")
			os.Exit(1)
		}
	}()

	err := c.ReadConfigFromFile(ctx, path, parameters)
	if err != nil {
		c.Logger().Fatal(ctx, err, "Process is terminated")
		os.Exit(1)
		return
	}

	c.Logger().Info(ctx, "Press Control-C to stop the microservice...")

	err = c.Open(ctx)
	if err != nil {
		_ = c.Close(ctx)
		cancel()
		c.Logger().Fatal(ctx, err, "Process is terminated")
		os.Exit(1)
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGABRT)

	select {
	case err := <-c.feedbackWithErrorChan:
		msg := cconv.StringConverter.ToString(err)
		err = errors.New(msg)
		_ = c.Close(ctx)
		cancel()
		c.Logger().Fatal(ctx, err, "Process is terminated")
		os.Exit(1)
		break
	case <-c.feedbackChan:
		_ = c.Close(ctx)
		cancel()
		c.Logger().Info(ctx, "Goodbye!")
		os.Exit(0)
		break
	case <-ch:
		_ = c.Close(ctx)
		cancel()
		c.Logger().Info(ctx, "Goodbye!")
		os.Exit(0)
		break
	}
}

// Registers all Google Function controllers in the container.
func (c *CloudFunction) RegisterControllers() {
	// Extract regular and commandable Google Function controllers from references
	controllers := c.References.GetOptional(crefer.NewDescriptor("*", "controller", "cloudfunc", "*", "*"))
	cmdControllers := c.References.GetOptional(crefer.NewDescriptor("*", "controller", "commandable-cloudfunc", "*", "*"))
	controllers = append(controllers, cmdControllers...)

	// Register actions defined in those controllers
	for _, service := range controllers {
		// Check if the service implements required interface
		if _val, ok := service.(gcpctrl.ICloudFunctionController); ok {
			actions := _val.GetActions()
			for _, action := range actions {
				c.RegisterAction(action.Cmd, action.Schema, action.Action)
			}
		}
	}
}

// Registers an action in this Google Function.
//
//	Parameters:
//		- cmd		a action/command name.
//		- schema	a validation schema to validate received parameters.
//		- action	an action function that is called when action is invoked.
//
// Deprecated: This method has been deprecated. Use CloudFunctionService instead.
func (c *CloudFunction) RegisterAction(cmd string, schema *cvalid.Schema, action http.HandlerFunc) {
	if cmd == "" {
		panic("NO_COMMAND: cmd parameter is missing")
	}

	if action == nil {
		panic("NO_ACTION: Missing action")
	}

	if c.Actions[cmd] != nil {
		panic("DUPLICATED_ACTION: Action already exists")
	}

	// Hack!!! Wrapping action to preserve prototyping request
	actionCurl := func(w http.ResponseWriter, r *http.Request) {
		// Perform validation
		if schema != nil {
			var params = make(map[string]any, 0)
			for k, v := range r.URL.Query() {
				params[k] = v[0]
			}

			if reqVars, ok := r.Context().Value(pattern.AllVariables).(map[pattern.Variable]any); ok {
				for k, v := range reqVars {
					params[string(k)] = v
				}
			}

			// Make copy of request
			bodyBuf, bodyErr := io.ReadAll(r.Body)
			if bodyErr != nil {
				httpctrl.HttpResponseSender.SendError(w, r, bodyErr)
				return
			}
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBuf))
			//-------------------------
			var body any
			_ = json.Unmarshal(bodyBuf, &body)
			params["body"] = body

			traceId := c.GetTraceId(r)
			err := schema.ValidateAndReturnError(traceId, params, false)
			if err != nil {
				httpctrl.HttpResponseSender.SendError(w, r, err)
				return
			}
		}

		action(w, r)
	}

	c.Actions[cmd] = actionCurl
}

// Executes this Google Function and returns the result.
// This method can be overloaded in child classes
// if they need to change the default behavior
//
//	Parameters:
//		- res the function response
//		- req the function request
func (c *CloudFunction) Execute(res http.ResponseWriter, req *http.Request) {
	traceId := c.GetTraceId(req)
	cmd, err := c.GetCommand(req)

	if err != nil {
		err := cerr.NewBadRequestError(
			traceId,
			"INVALID_JSON",
			"Invalid json format",
		)

		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	if cmd == "" {
		err = cerr.NewBadRequestError(
			traceId,
			"NO_COMMAND",
			"Cmd parameter is missing",
		)

		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	action := c.Actions[cmd]
	if action == nil {
		err = cerr.NewBadRequestError(
			traceId,
			"NO_ACTION",
			"Action "+cmd+" was not found",
		)

		httpctrl.HttpResponseSender.SendError(res, req, err)
		return
	}

	action(res, req)
}

func (c *CloudFunction) handler(res http.ResponseWriter, req *http.Request) {
	// If already started then execute
	if c.IsOpen() {
		c.Execute(res, req)
		return
	}
	// Start before execute
	c.Run(req.Context())
	c.Execute(res, req)
}

// Gets entry point into this Google Function.
//
//	Parameters:
//		- res the function response
//		- req the function request
func (c *CloudFunction) GetHandler() http.HandlerFunc {
	return c.handler
}

// Returns traceId from Googel Function request.
// This method can be overloaded in child classes
//
//	Parameters:
//		- req	Googel Function request
//
// Returns traceId from request
func (c *CloudFunction) GetTraceId(r *http.Request) string {
	return gcputil.CloudFunctionRequestHelper.GetTraceId(r)
}

// Returns command from Google Function request.
// This method can be overloaded in child classes
//
//	Parameters:
//		- req	Googel Function request
//
// Returns command from request
func (c *CloudFunction) GetCommand(r *http.Request) (string, error) {
	return gcputil.CloudFunctionRequestHelper.GetCommand(r)
}
