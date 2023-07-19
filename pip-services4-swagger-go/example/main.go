package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	crun "github.com/pip-services4/pip-services4-go/pip-services4-components-go/run"
	cservices "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	ccontrollers "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/controllers"
	controllers "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/example/controllers"
	logic "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/example/logic"
)

func main() {
	ctx := context.Background()

	// Create components
	logger := clog.NewConsoleLogger()
	counter := ccount.NewLogCounters()
	service := logic.NewDummyService()
	httpEndpoint := cservices.NewHttpEndpoint()
	restController := controllers.NewDummyRestController()
	httpController := controllers.NewDummyCommandableHttpController()
	statusController := cservices.NewStatusRestController()
	heartbeatController := cservices.NewHeartbeatRestController()
	swaggerController := ccontrollers.NewSwaggerController()

	components := []any{
		logger,
		counter,
		service,
		httpEndpoint,
		restController,
		httpController,
		statusController,
		heartbeatController,
		swaggerController,
	}

	// Configure components
	logger.Configure(ctx, cconf.NewConfigParamsFromTuples(
		"level", "trace",
	))

	httpEndpoint.Configure(ctx, cconf.NewConfigParamsFromTuples(
		"connection.prototol", "http",
		"connection.host", "localhost",
		"connection.port", 8080,
	))

	restController.Configure(ctx, cconf.NewConfigParamsFromTuples(
		"swagger.enable", true,
	))

	httpController.Configure(ctx, cconf.NewConfigParamsFromTuples(
		"base_route", "dummies2",
		"swagger.enable", true,
	))

	// Set references
	references := cref.NewReferencesFromTuples(ctx,
		cref.NewDescriptor("pip-services", "logger", "console", "default", "1.0"), logger,
		cref.NewDescriptor("pip-services", "counter", "log", "default", "1.0"), counter,
		cref.NewDescriptor("pip-services", "endpoint", "http", "default", "1.0"), httpEndpoint,
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), service,
		cref.NewDescriptor("pip-services-dummies", "controller", "rest", "default", "1.0"), restController,
		cref.NewDescriptor("pip-services-dummies", "controller", "commandable-http", "default", "1.0"), httpController,
		cref.NewDescriptor("pip-services", "status-controller", "rest", "default", "1.0"), statusController,
		cref.NewDescriptor("pip-services", "heartbeat-controller", "rest", "default", "1.0"), heartbeatController,
		cref.NewDescriptor("pip-services", "swagger-controller", "http", "default", "1.0"), swaggerController,
	)

	cref.Referencer.SetReferences(ctx, references, components)

	// Open components
	err := crun.Opener.Open(ctx, components)
	if err != nil {
		logger.Error(ctx, err, "Failed to open components")
		return
	}

	// Wait until user presses ENTER
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press ENTER to stop the microservice...")
	reader.ReadString('\n')

	// Close components
	err = crun.Closer.Close(ctx, components)
	if err != nil {
		logger.Error(ctx, err, "Failed to close components")
	}
}
