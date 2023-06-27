package test_clients

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

func TestDummyCommandableHttpClient(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyCommandableHttpControllerPort,
	)

	client := NewDummyCommandableHttpClient()
	fixture := NewDummyClientFixture(client)

	client.Configure(context.Background(), restConfig)
	client.SetReferences(context.Background(), cref.NewEmptyReferences())
	_ = client.Open(context.Background())
	t.Run("CRUD Operations", fixture.TestCrudOperations)
}
