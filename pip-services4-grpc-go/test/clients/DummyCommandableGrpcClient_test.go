package test_clients

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/sample"
	testservices "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/services"
)

func TestDummyRestClient(t *testing.T) {
	ctx := context.Background()

	grpcConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3002",
	)

	var cpontroller *testservices.DummyCommandableGrpcController
	var client *DummyCommandableGrpcClient
	var fixture *DummyClientFixture

	srv := tsample.NewDummyService()

	cpontroller = testservices.NewDummyCommandableGrpcController()
	cpontroller.Configure(ctx, grpcConfig)

	references := cref.NewReferencesFromTuples(ctx,
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "grpc", "default", "1.0"), cpontroller,
	)
	cpontroller.SetReferences(ctx, references)

	cpontroller.Open(ctx)
	defer cpontroller.Close(ctx)

	client = NewDummyCommandableGrpcClient()
	fixture = NewDummyClientFixture(client)

	client.Configure(ctx, grpcConfig)
	client.SetReferences(ctx, cref.NewEmptyReferences())
	client.Open(ctx)
	defer client.Close(ctx)

	t.Run("CRUD Operations", fixture.TestCrudOperations)

}
