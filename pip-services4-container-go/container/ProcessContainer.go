package container

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cconfig "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// ProcessContainer inversion of control (IoC) container that runs as a system process.
// It processes command line arguments and handles unhandled exceptions and Ctrl-C signal
// to gracefully shutdown the container.
//
//	Command line arguments:
//		--config / -c path to JSON or YAML file with container configuration (default: "./config/config.yml")
//		--param / --params / -p value(s) to parameterize the container configuration
//		--help / -h prints the container usage help
//	see Container
//
//	Example:
//		container = NewEmptyProcessContainer()
//		container.Container.AddFactory(NewMyComponentFactory())
//		container.Run(context.Background(), os.Environ())
type ProcessContainer struct {
	*Container
	configPath            string
	feedbackChan          cctx.ContextShutdownChan
	feedbackWithErrorChan cctx.ContextShutdownWithErrorChan
}

const DefaultConfigFilePath = "./config/config.yml"

// NewEmptyProcessContainer creates a new empty instance of the container.
//
//	Returns: ProcessContainer
func NewEmptyProcessContainer() *ProcessContainer {
	c := &ProcessContainer{
		Container:             NewEmptyContainer(),
		configPath:            DefaultConfigFilePath,
		feedbackChan:          make(cctx.ContextShutdownChan),
		feedbackWithErrorChan: make(cctx.ContextShutdownWithErrorChan),
	}
	c.SetLogger(log.NewConsoleLogger())
	return c
}

// NewProcessContainer creates a new instance of the container.
//
//	Parameters:
//		- name string a container name (accessible via ContextInfo)
//		- description string a container description (accessible via ContextInfo)
//	Returns: ProcessContainer
func NewProcessContainer(name string, description string) *ProcessContainer {
	c := &ProcessContainer{
		Container:             NewContainer(name, description),
		configPath:            DefaultConfigFilePath,
		feedbackChan:          make(cctx.ContextShutdownChan),
		feedbackWithErrorChan: make(cctx.ContextShutdownWithErrorChan),
	}
	c.SetLogger(log.NewConsoleLogger())
	return c
}

// InheritProcessContainer creates a new instance of the container inherit from reference.
//
//	Parameters:
//		- name string a container name (accessible via ContextInfo)
//		- description string a container description (accessible via ContextInfo)
//		- referenceable crefer.IReferenceable
//		- referenceble object for inherit
//	Returns: *Container
func InheritProcessContainer(name string, description string,
	referenceable crefer.IReferenceable) *ProcessContainer {

	c := &ProcessContainer{
		Container:             InheritContainer(name, description, referenceable),
		configPath:            DefaultConfigFilePath,
		feedbackChan:          make(cctx.ContextShutdownChan),
		feedbackWithErrorChan: make(cctx.ContextShutdownWithErrorChan),
	}
	c.SetLogger(log.NewConsoleLogger())
	return c
}

// SetConfigPath set path for configuration file
func (c *ProcessContainer) SetConfigPath(configPath string) {
	c.configPath = configPath
}

func (c *ProcessContainer) getConfigPath(args []string) string {
	for index, arg := range args {
		nextArg := ""
		if index < len(args)-1 {
			nextArg = args[index+1]
			if strings.HasPrefix(nextArg, "-") {
				nextArg = ""
			}
		}

		if arg == "--config" || arg == "-c" {
			return nextArg
		}
	}

	return c.configPath
}

func (c *ProcessContainer) getParameters(args []string) *cconfig.ConfigParams {
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

	parameters := cconfig.NewConfigParamsFromString(line)

	for _, e := range os.Environ() {
		if env := strings.Split(e, "="); len(env) == 2 {
			parameters.SetAsObject(env[0], env[1])
		} else {
			parameters.SetAsObject(env[0], strings.Join(env[1:], "="))
		}
	}

	return parameters
}

func (c *ProcessContainer) showHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func (c *ProcessContainer) printHelp() {
	fmt.Println("Pip.Services process container - http://www.github.com/pip-services/pip-services")
	fmt.Println("run [-h] [-c <config file>] [-p <param>=<value>]*")
}

// Run the container by instantiating and running components inside the container.
// It reads the container configuration, creates, configures, references
// and opens components. On process exit it closes, unreferences and destroys
// components to gracefully shutdown.
//
//	Parameters:
//		- ctx context.Context
//		- args []string command line arguments
func (c *ProcessContainer) Run(ctx context.Context, args []string) {
	if c.showHelp(args) {
		c.printHelp()
		os.Exit(0)
		return
	}

	ctx, cancel := context.WithCancel(ctx)

	ctx, _ = cctx.AddShutdownChanToContext(ctx, c.feedbackChan)
	ctx, _ = cctx.AddErrShutdownChanToContext(ctx, c.feedbackWithErrorChan)

	ctx = cctx.NewContextWithTraceId(ctx, c.Info().Name)
	path := c.getConfigPath(args)
	parameters := c.getParameters(args)

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
