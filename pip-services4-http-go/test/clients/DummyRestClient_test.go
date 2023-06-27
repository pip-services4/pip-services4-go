package test_clients

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

func TestDummyRestClient(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", DummyRestControllerPort,
		"options.trace_id_place", "headers",
	)

	client := NewDummyRestClient()
	fixture := NewDummyClientFixture(client)

	client.Configure(context.TODO(), restConfig)
	client.SetReferences(context.TODO(), cref.NewEmptyReferences())
	client.Open(context.TODO())

	t.Run("DummyRestClient.CrudOperations", fixture.TestCrudOperations)
}
