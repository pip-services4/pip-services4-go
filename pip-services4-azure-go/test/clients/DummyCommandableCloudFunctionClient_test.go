package clients_test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
)

func TestDummyCommandableAzureFunctionClient(t *testing.T) {

	appName := os.Getenv("AZURE_FUNCTION_APP_NAME")
	functionName := os.Getenv("AZURE_FUNCTION_NAME")
	protocol := os.Getenv("AZURE_FUNCTION_PROTOCOL")
	authCode := os.Getenv("AZURE_FUNCTION_AUTH_CODE")
	uri := os.Getenv("AZURE_FUNCTION_URI")

	// if uri == "" {
	// 	uri = "http://localhost:7071/api/HttpTrigger1" // uncomment for local test
	// }

	if uri == "" && (appName == "" || functionName == "" || protocol == "" || authCode == "") {
		t.Skip("No credentials set, skip TestDummyCommandableAzureFunctionClient")
	}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", uri,
		"connection.protocol", protocol,
		"connection.app_name", appName,
		"connection.function_name", functionName,
		"credential.auth_code", authCode,
	)

	client := NewDummyCommandableAzureFunctionClient()
	fixture := NewDummyClientFixture(client)

	ctx := utils.ContextHelper.NewContextWithTraceId(context.Background(), "DummyCommandableAzureFunctionClient")
	client.Configure(ctx, restConfig)
	client.SetReferences(ctx, cref.NewEmptyReferences())
	client.Open(ctx)
	defer client.Close(ctx)

	t.Run("DummyCommandableAzureFunctionClient.CrudOperations", fixture.TestCrudOperations)
}
