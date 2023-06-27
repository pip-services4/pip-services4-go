package test_services

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/protos"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/sample"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func TestDummyGrpcController(t *testing.T) {
	ctx := context.Background()

	grpcConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3004",
	)

	var Dummy1 tsample.Dummy
	var Dummy2 tsample.Dummy

	var controller *DummyGrpcController

	var client protos.DummiesClient
	srv := tsample.NewDummyService()

	controller = NewDummyGrpcController()
	controller.Configure(ctx, grpcConfig)

	references := cref.NewReferencesFromTuples(ctx,
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "grpc", "default", "1.0"), controller,
	)
	controller.SetReferences(ctx, references)

	controller.Open(ctx)
	defer controller.Close(ctx)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("localhost:3004", opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client = protos.NewDummiesClient(conn)

	Dummy1 = tsample.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	Dummy2 = tsample.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	// Test CRUD Operations
	// Create first dummy
	protoDummy := protos.Dummy{}
	protoDummy.Id = Dummy1.Id
	protoDummy.Key = Dummy1.Key
	protoDummy.Content = Dummy1.Content
	request := protos.DummyObjectRequest{Dummy: &protoDummy}
	dummy, err := client.CreateDummy(context.TODO(), &request)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, protoDummy.Content, dummy.Content)
	assert.Equal(t, protoDummy.Key, dummy.Key)

	dummy1 := dummy

	// Create another dummy
	protoDummy.Id = Dummy2.Id
	protoDummy.Key = Dummy2.Key
	protoDummy.Content = Dummy2.Content
	request = protos.DummyObjectRequest{Dummy: &protoDummy}
	dummy, err = client.CreateDummy(context.TODO(), &request)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, protoDummy.Content, dummy.Content)
	assert.Equal(t, protoDummy.Key, dummy.Key)

	// Get all dummies
	requestPage := protos.DummiesPageRequest{}
	dummies, err := client.GetDummies(context.TODO(), &requestPage)
	assert.Nil(t, err)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	protoDummy.Id = dummy1.Id
	protoDummy.Key = dummy1.Key
	protoDummy.Content = dummy1.Content

	request = protos.DummyObjectRequest{Dummy: &protoDummy}
	dummy, err = client.UpdateDummy(context.TODO(), &request)
	assert.Nil(t, err)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, dummy1.Key)

	// Delete dummy
	idRequest := protos.DummyIdRequest{DummyId: dummy1.Id}
	_, err = client.DeleteDummyById(context.TODO(), &idRequest)
	assert.Nil(t, err)

	// Try to get delete dummy
	idRequest = protos.DummyIdRequest{DummyId: dummy1.Id}
	dummy, err = client.GetDummyById(context.TODO(), &idRequest)
	assert.Nil(t, err)

	callsCnt := controller.GetNumberOfCalls()
	assert.Equal(t, callsCnt, (int64)(6))
}
