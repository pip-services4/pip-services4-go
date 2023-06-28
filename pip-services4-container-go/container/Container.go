package container

import (
	"context"
	"errors"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cconfig "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-container-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-container-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-container-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// Container Inversion of control (IoC) container that creates components and manages their lifecycle.
// The container is driven by configuration, that usually stored in JSON or YAML file.
// The configuration contains a list of components identified by type or locator,
// followed by component configuration.
// On container start it performs the following actions:
//
// Creates components using their types or calls registered factories to create components using their locators
// Configures components that implement IConfigurable interface and passes them their configuration parameters
// Sets references to components that implement IReferenceable interface and passes
// them references of all components in the container
// Opens components that implement IOpenable interface
// On container stop actions are performed in reversed order:
//
// Closes components that implement ICloseable interface
// Unsets references in components that implement IUnreferenceable interface
// Destroys components in the container.
// The component configuration can be parameterized by dynamic values.
// That allows specialized containers to inject parameters from command line or from environment variables.
//
// The container automatically creates a ContextInfo component that carries detail information about
// the container and makes it available for other components.
//
//	see IConfigurable (in the PipServices "Commons" package)
//	see IReferenceable (in the PipServices "Commons" package)
//	see IOpenable (in the PipServices "Commons" package)
//
//	Configuration parameters:
//		name: the context (container or process) name
//		description: human-readable description of the context
//		properties: entire section of additional descriptive properties
//		- ...
//	Example:
//		======= config.yml ========
//		- descriptor: mygroup:mycomponent1:default:default:1.0
//			param1: 123
//			param2: ABC
//		- type: mycomponent2,mypackage
//			param1: 321
//			param2: XYZ
//		============================
//
//		container := NewEmptyContainer();
//		container.AddFactory(newMyComponentFactory());
//
//		parameters := NewConfigParamsFromValue(os.Environ());
//		container.ReadConfigFromFile(context.Background(), "123", "./config/config.yml", parameters);
//
//		err := container.Open(context.Background(), "123")
//		fmt.Println("Container is opened"))
//		...
//		err = container.Close(context.Background(), "123")
//		fmt.Println("Container is closed")
type Container struct {
	logger          log.ILogger
	factories       *cbuild.CompositeFactory
	info            *cctx.ContextInfo
	config          config.ContainerConfig
	References      *refer.ContainerReferences
	referenceable   crefer.IReferenceable
	unreferenceable crefer.IUnreferenceable
}

// NewEmptyContainer creates a new empty instance of the container.
//
//	Returns *Container
func NewEmptyContainer() *Container {
	return &Container{
		logger:    log.NewNullLogger(),
		factories: build.NewDefaultContainerFactory(),
		info:      cctx.NewContextInfo(),
	}
}

// NewContainer creates a new instance of the container.
//
//	Parameters:
//		- name string a container name (accessible via ContextInfo)
//		- description string a container description (accessible via ContextInfo)
//	Returns: *Container
func NewContainer(name string, description string) *Container {
	c := NewEmptyContainer()

	c.info.Name = name
	c.info.Description = description

	return c
}

// InheritContainer creates a new instance of the container inherit from reference.
//
//	Parameters:
//		- name string a container name (accessible via ContextInfo)
//		- description string a container description (accessible via ContextInfo)
//		- referenceable crefer.IReferenceable
//		- referenceble object for inherit
//	Returns: *Container
func InheritContainer(name string, description string,
	referenceable crefer.IReferenceable) *Container {

	c := NewEmptyContainer()

	c.info.Name = name
	c.info.Description = description
	c.referenceable = referenceable
	c.unreferenceable, _ = referenceable.(crefer.IUnreferenceable)

	return c
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *cconfig.ConfigParams configuration parameters to be set.
func (c *Container) Configure(ctx context.Context, conf *cconfig.ConfigParams) {
	c.config, _ = config.ReadContainerConfigFromConfig(conf)
}

// ReadConfigFromFile container configuration from JSON or YAML file and parameterizes it with given values.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- path string a path to configuration file
//		- parameters *cconfig.ConfigParams values to parameters the configuration or null to skip parameterization.
func (c *Container) ReadConfigFromFile(ctx context.Context,
	path string, parameters *cconfig.ConfigParams) error {

	var err error
	c.config, err = config.ContainerConfigReader.ReadFromFile(ctx, path, parameters)
	// c.logger.Trace(ctx, config.String())
	return err
}

func (c *Container) initReferences(ctx context.Context, references crefer.IReferences) {
	contextInfoRef := references.GetOneOptional(
		crefer.NewDescriptor("pip-services", "context-info",
			"*", "*", "1.0",
		),
	)
	if existingInfo, ok := contextInfoRef.(*cctx.ContextInfo); !ok {
		references.Put(
			ctx,
			crefer.NewDescriptor(
				"pip-services",
				"context-info",
				"default", "default", "1.0",
			),
			c.info,
		)
	} else {
		c.info = existingInfo
	}

	references.Put(
		ctx,
		crefer.NewDescriptor(
			"pip-services",
			"factory",
			"container", "default", "1.0",
		),
		c.factories,
	)
}

func (c *Container) Logger() log.ILogger {
	return c.logger
}

func (c *Container) SetLogger(logger log.ILogger) {
	c.logger = logger
}

func (c *Container) Info() *cctx.ContextInfo {
	return c.info
}

// AddFactory a factory to the container.
// The factory is used to create components added to the container by their locators (descriptors).
//
//	Parameters:
//		- factory IFactory a component factory to be added.
func (c *Container) AddFactory(factory cbuild.IFactory) {
	c.factories.Add(factory)
}

// IsOpen checks if the component is opened.
//
//	Returns bool true if the component has been opened and false otherwise.
func (c *Container) IsOpen() bool {
	return c.References != nil
}

// Open the component.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error
func (c *Container) Open(ctx context.Context) error {
	var err error

	if c.References != nil {
		return cerr.NewInvalidStateError(
			cctx.GetTraceId(ctx), "ALREADY_OPENED", "Container was already opened",
		)
	}

	defer func() {
		if r := recover(); r != nil {
			recoverErr, ok := r.(error)
			if !ok {
				msg := cconv.StringConverter.ToString(r)
				recoverErr = errors.New(msg)
			}
			err = recoverErr
			c.logger.Error(ctx, recoverErr, "Failed to start container")
			panic(err)
		}
	}()

	c.logger.Trace(ctx, "Starting container.")

	// Create references with configured components
	c.References = refer.NewContainerReferences()
	c.initReferences(ctx, c.References)
	err = c.References.PutFromConfig(ctx, c.config)
	if err != nil {
		return err
	}

	if c.referenceable != nil {
		c.referenceable.SetReferences(ctx, c.References)
	}

	// Get custom description if available
	infoDescriptor := crefer.NewDescriptor("*", "context-info", "*", "*", "*")
	info, ok := c.References.GetOneOptional(infoDescriptor).(*cctx.ContextInfo)
	if ok {
		c.info = info
	}

	// Get reference to logger
	c.logger = log.NewCompositeLoggerFromReferences(ctx, c.References)

	// Open references
	err = c.References.Open(ctx)
	if err == nil {
		c.logger.Info(ctx, "Container %s started", c.info.Name)
	} else {
		c.logger.Fatal(ctx, err, "Failed to start container")
		_ = c.Close(ctx)
	}

	return err
}

// Close component and frees used resources.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//	Returns: error
func (c *Container) Close(ctx context.Context) error {
	// Skip if container wasn't opened
	if c.References == nil {
		return nil
	}

	var err error

	defer func() {
		if r := recover(); r != nil {
			recoverErr, ok := r.(error)
			if !ok {
				msg := cconv.StringConverter.ToString(r)
				recoverErr = errors.New(msg)
			}
			err = recoverErr
			c.logger.Error(ctx, recoverErr, "Failed to stop container")
			panic(err)
		}
	}()

	c.logger.Trace(ctx, "Stopping %s container", c.info.Name)

	// Unset references for child container
	if c.unreferenceable != nil {
		c.unreferenceable.UnsetReferences(ctx)
	}

	// Close and dereference components
	err = c.References.Close(ctx)

	c.References = nil

	if err == nil {
		c.logger.Info(ctx, "Container %s stopped", c.info.Name)
	} else {
		c.logger.Error(ctx, err, "Failed to stop container")
	}

	return err
}
