package test_clients

import (
	"context"
	"testing"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/test/sample"
)

func TestDummyDirectClient(t *testing.T) {

	client := NewDummyDirectClient()
	references := cref.NewReferencesFromTuples(
		context.Background(),
		cref.NewDescriptor(
			"pip-services-dummies", "controller", "default",
			"default", "1.0",
		), tsample.NewDummyService(),
	)
	client.SetReferences(context.Background(), references)
	client.Open(context.Background())
	defer client.Close(context.Background())

	fixture := NewDummyClientFixture(client)
	t.Run("CRUD Operations", fixture.TestCrudOperations)
}
