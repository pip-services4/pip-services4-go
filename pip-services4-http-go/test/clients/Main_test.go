package test_clients

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	test_controllers "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/controllers"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"
)

const (
	DummyRestControllerPort = iota + 4000
	DummyCommandableHttpControllerPort
)

func TestMain(m *testing.M) {

	fmt.Println("Preparing test services for clients...")

	dummyRestController := BuildTestDummyRestController()
	err := dummyRestController.Open(context.Background())
	if err != nil {
		panic(err)
	}
	defer dummyRestController.Close(context.Background())

	dummyCommandableHttpController := BuildTestDummyCommandableHttpController()
	err = dummyCommandableHttpController.Open(context.Background())
	if err != nil {
		panic(err)
	}
	defer dummyCommandableHttpController.Close(context.Background())
	time.Sleep(time.Second)
	fmt.Println("All test services started!")

	os.Exit(m.Run())
}

func BuildTestDummyRestController() *test_controllers.DummyRestController {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyRestControllerPort,
		"openapi_content", "swagger yaml or json content",
		"swagger.enable", "true",
	)

	var controller *test_controllers.DummyRestController
	srv := tsample.NewDummyService()

	controller = test_controllers.NewDummyRestController()
	controller.Configure(context.Background(), restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "rest", "default", "1.0"), controller,
	)
	controller.SetReferences(context.Background(), references)
	return controller
}

func BuildTestDummyCommandableHttpController() *test_controllers.DummyCommandableHttpController {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyCommandableHttpControllerPort,
		"swagger.enable", "true",
	)

	srv := tsample.NewDummyService()

	controller := test_controllers.NewDummyCommandableHttpController()

	controller.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "http", "default", "1.0"), controller,
	)
	controller.SetReferences(context.Background(), references)
	return controller
}
