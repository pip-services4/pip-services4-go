package clients_test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

func TestDummyCommandableCloudFunctionClient(t *testing.T) {

	functionName := os.Getenv("GCP_FUNCTION_NAME")
	protocol := os.Getenv("GCP_FUNCTION_PROTOCOL")
	region := os.Getenv("GCP_FUNCTION_REGION")
	projectId := os.Getenv("GCP_PROJECT_ID")
	uri := os.Getenv("GCP_FUNCTION_URI")

	// if uri == "" {
	// 	uri = "http://localhost:3000" // uncomment for local test
	// }

	if uri == "" && (region == "" || functionName == "" || protocol == "" || projectId == "") {
		t.Skip("No credentials set, skip TestDummyCommandableCloudFunctionClient")
	}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", uri,
		"connection.protocol", protocol,
		"connection.region", region,
		"connection.function", functionName,
		"connection.project_id", projectId,
	)

	client := NewDummyCommandableCloudFunctionClient()
	fixture := NewDummyClientFixture(client)

	ctx := cctx.NewContextWithTraceId(context.Background(), "DummyCommandableCloudFunctionClient")
	client.Configure(ctx, restConfig)
	client.SetReferences(ctx, cref.NewEmptyReferences())
	client.Open(ctx)
	defer client.Close(ctx)

	t.Run("DummyCommandableCloudFunctionClient.CrudOperations", fixture.TestCrudOperations)
}
