package test_controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/keys"
	services "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
	tlogic "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"
)

const (
	StatusRestControllerPort = iota + 3000
	HeartbeatRestControllerPort
	HttpEndpointControllertPort
	DummyRestControllertPort
	DummyOpenAPIFileRestControllerPort
	DummyCommandableHttpControllerPort
	DummyCommandableSwaggerHttpControllerPort
)

func TestMain(m *testing.M) {

	fmt.Println("Preparing test services...")

	statusRestController := BuildTestStatusRestController()
	err := statusRestController.Open(context.Background())
	if err != nil {
		panic(err)
	}
	defer statusRestController.Close(context.Background())

	heartbeatRestController := BuildTestHeartbeatRestController()
	err = heartbeatRestController.Open(context.Background())
	if err != nil {
		panic(err)
	}
	defer heartbeatRestController.Close(context.Background())

	httpEndpointController, endpoint := BuildTestHttpEndpointController()
	err = endpoint.Open(context.Background())
	if err != nil {
		panic(err)
	} else {
		err = httpEndpointController.Open(context.Background())
		if err != nil {
			panic(err)
		} else {
			defer endpoint.Close(context.Background())
			defer httpEndpointController.Close(context.Background())
		}
	}

	// Prepare shutdown context and channel
	shutdownCtx, restControllerCancel := context.WithTimeout(context.Background(), time.Minute*3)
	shutdownChan := make(cctx.ContextShutdownWithErrorChan)
	shutdownCtx, _ = cctx.AddErrShutdownChanToContext(shutdownCtx, shutdownChan)

	dummyRestController := BuildTestDummyRestController()
	err = dummyRestController.Open(shutdownCtx)
	if err != nil {
		panic(err)
	}
	defer dummyRestController.Close(context.Background())

	dummyOpenAPIFileRestController, filename := BuildTestDummyOpenAPIFileRestController()
	err = dummyOpenAPIFileRestController.Open(context.Background())
	if err != nil {
		panic(err)
	}
	defer dummyOpenAPIFileRestController.Close(context.Background())
	//defer os.Remove(filename)
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}()

	dummyCommandableHttpController := BuildTestDummyCommandableHttpController()
	err = dummyCommandableHttpController.Open(context.Background())
	if err != nil {
		panic(err)
	}
	defer dummyCommandableHttpController.Close(context.Background())

	dummyCommandableSwaggerHttpController := BuildTestDummyCommandableSwaggerHttpController()
	err = dummyCommandableSwaggerHttpController.Open(context.Background())
	if err != nil {
		panic(err)
	}
	defer dummyCommandableSwaggerHttpController.Close(context.Background())
	time.Sleep(time.Second)
	fmt.Println("All test services started!")

	code := m.Run()

	noc := dummyRestController.GetNumberOfCalls()
	fmt.Println("Number of calls:", noc, "from 4")
	if noc != 4 {
		panic("Number of calls test failed!")
	}

	go func() {
		getResponse, getErr := http.Get(
			fmt.Sprintf(
				"http://localhost:%d/dummies/check/graceful_shutdown",
				DummyRestControllertPort,
			),
		)
		fmt.Println(getResponse)
		fmt.Println(getErr)
	}()

	for {
		select {
		case err := <-shutdownChan:
			restControllerCancel()
			if err == nil {
				panic("invalid shutdown error")
			}
			if err.Error() != "called from DummyController.CheckGracefulShutdownContext" {
				panic("invalid shutdown error")
			}
			fmt.Println("rest service shutdown successful")
			os.Exit(code)
		case <-shutdownCtx.Done():
			fmt.Println("rest service shutdown by timeout")
			os.Exit(1)
		}
	}
}

func BuildTestStatusRestController() *services.StatusRestController {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", StatusRestControllerPort,
		"cors_headers", "trace_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	controller := services.NewStatusRestController()
	controller.Configure(context.Background(), restConfig)

	contextInfo := cctx.NewContextInfo()
	contextInfo.Name = "Test"
	contextInfo.Description = "This is a test container"

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services", "context-info", "default", "default", "1.0"), contextInfo,
		cref.NewDescriptor("pip-services", "status-controller", "http", "default", "1.0"), controller,
	)
	controller.SetReferences(context.Background(), references)
	return controller
}

func BuildTestHttpEndpointController() (*DummyRestController, *services.HttpEndpoint) {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", HttpEndpointControllertPort,
		"cors_headers", "trace_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	srv := tlogic.NewDummyService()
	controller := NewDummyRestController()
	controller.Configure(context.Background(), cconf.NewConfigParamsFromTuples(
		"base_route",
		"/api/v1",
	))

	endpoint := services.NewHttpEndpoint()
	endpoint.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "rest", "default", "1.0"), controller,
		cref.NewDescriptor("pip-services", "endpoint", "http", "default", "1.0"), endpoint,
	)
	controller.SetReferences(context.Background(), references)
	return controller, endpoint
}

func BuildTestDummyRestController() *DummyRestController {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyRestControllertPort,
		"openapi_content", "swagger yaml or json content",
		"swagger.enable", "true",
		"cors_headers", "trace_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	var controller *DummyRestController
	srv := tlogic.NewDummyService()

	controller = NewDummyRestController()
	controller.Configure(context.Background(), restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "rest", "default", "1.0"), controller,
	)
	controller.SetReferences(context.Background(), references)
	return controller
}

func BuildTestDummyOpenAPIFileRestController() (*DummyRestController, string) {

	openApiContent := "swagger yaml content from file"
	filename := path.Join(".", "dummy_"+keys.IdGenerator.NextLong()+".tmp")

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(([]byte)(openApiContent))
	if err != nil {
		panic(err)
	}
	//err = file.Close()
	//if err != nil {
	//	panic(err)
	//}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyOpenAPIFileRestControllerPort,
		"openapi_file", filename, // for test only
		"swagger.enable", "true",
		"cors_headers", "trace_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	var controller *DummyRestController
	srv := tlogic.NewDummyService()

	controller = NewDummyRestController()
	controller.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "rest", "default", "1.0"), controller,
	)
	controller.SetReferences(context.Background(), references)
	return controller, filename
}

func BuildTestDummyCommandableHttpController() *DummyCommandableHttpController {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyCommandableHttpControllerPort,
		"swagger.enable", "true",
		"cors_headers", "trace_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	srv := tlogic.NewDummyService()

	controller := NewDummyCommandableHttpController()

	controller.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "http", "default", "1.0"), controller,
	)
	controller.SetReferences(context.Background(), references)
	return controller
}

func BuildTestDummyCommandableSwaggerHttpController() *DummyCommandableHttpController {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyCommandableSwaggerHttpControllerPort,
		"swagger.enable", "true",
		"swagger.auto", false,
		"cors_headers", "trace_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	srv := tlogic.NewDummySchema()

	controller := NewDummyCommandableHttpController()

	controller.Configure(context.Background(), restConfig)

	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "http", "default", "1.0"), controller,
	)
	controller.SetReferences(context.Background(), references)
	return controller
}

func BuildTestHeartbeatRestController() *services.HeartbeatRestController {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", HeartbeatRestControllerPort,
		"cors_headers", "trace_id, access_token, Accept, Content-Type, Content-Length, X-CSRF-Token",
		"cors_origins", "*",
	)

	controller := services.NewHeartbeatRestController()
	controller.Configure(context.Background(), restConfig)
	return controller
}
