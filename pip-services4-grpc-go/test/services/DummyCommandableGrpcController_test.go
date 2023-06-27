package test_services

import (
	"context"
	"encoding/json"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cmdproto "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/protos"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/sample"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func TestDummyCommandableGrpcController(t *testing.T) {
	ctx := context.Background()

	grpcConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3001",
	)

	var Dummy1 tsample.Dummy
	var Dummy2 tsample.Dummy
	var service *DummyCommandableGrpcController
	var client cmdproto.CommandableClient

	srv := tsample.NewDummyService()
	service = NewDummyCommandableGrpcController()
	service.Configure(ctx, grpcConfig)

	references := cref.NewReferencesFromTuples(ctx,
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
		cref.NewDescriptor("pip-services-dummies", "controller", "grpc", "default", "1.0"), service,
	)
	service.SetReferences(ctx, references)
	service.Open(ctx)
	defer service.Close(ctx)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("localhost:3001", opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client = cmdproto.NewCommandableClient(conn)

	Dummy1 = tsample.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	Dummy2 = tsample.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	// Test CRUD Operations
	var dummy, dummy1 tsample.Dummy

	request := cmdproto.InvokeRequest{}

	requestParams := make(map[string]any)
	requestParams["dummy"] = Dummy1
	jsonBuf, _ := json.Marshal(requestParams)

	request.Method = "dummy.create_dummy"
	request.TraceId = "dummy.create_dummy123"
	request.ArgsEmpty = false
	request.ArgsJson = string(jsonBuf)
	response, err := client.Invoke(context.TODO(), &request)

	assert.Nil(t, err)
	assert.False(t, response.ResultEmpty)
	assert.NotEqual(t, response.ResultJson, "")
	json.Unmarshal([]byte(response.ResultJson), &dummy)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, Dummy1.Content)
	assert.Equal(t, dummy.Key, Dummy1.Key)
	dummy1 = dummy

	// Create another dummy
	requestParams["dummy"] = Dummy2
	jsonBuf, _ = json.Marshal(requestParams)

	request.Method = "dummy.create_dummy"
	request.ArgsEmpty = false
	request.ArgsJson = string(jsonBuf)
	response, err = client.Invoke(context.TODO(), &request)

	assert.Nil(t, err)
	assert.False(t, response.ResultEmpty)
	assert.NotEqual(t, response.ResultJson, "")
	json.Unmarshal([]byte(response.ResultJson), &dummy)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, Dummy2.Content)
	assert.Equal(t, dummy.Key, Dummy2.Key)

	// Get all dummies
	request.Method = "dummy.get_dummies"
	request.ArgsEmpty = false
	request.ArgsJson = "{}"
	response, err = client.Invoke(context.TODO(), &request)

	assert.Nil(t, err)
	assert.False(t, response.ResultEmpty)
	assert.NotEqual(t, response.ResultJson, "")
	var dummies tsample.DummyDataPage
	json.Unmarshal([]byte(response.ResultJson), &dummies)

	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	requestParams["dummy"] = dummy1
	jsonBuf, _ = json.Marshal(requestParams)

	request.Method = "dummy.update_dummy"
	request.ArgsEmpty = false
	request.ArgsJson = string(jsonBuf)
	response, err = client.Invoke(context.TODO(), &request)

	assert.Nil(t, err)
	assert.False(t, response.ResultEmpty)
	assert.NotEqual(t, response.ResultJson, "")
	json.Unmarshal([]byte(response.ResultJson), &dummy)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, Dummy1.Key)

	dummy1 = dummy

	// Delete dummy
	delParam := make(map[string]string, 0)
	delParam["dummy_id"] = dummy1.Id
	jsonBuf, _ = json.Marshal(delParam)

	request.Method = "dummy.delete_dummy_by_id"
	request.ArgsEmpty = false
	request.ArgsJson = string(jsonBuf)
	response, err = client.Invoke(context.TODO(), &request)

	assert.Nil(t, err)
	assert.Nil(t, response.Error)

	// Try to get delete dummy
	request.Method = "dummy.get_dummy_by_id"
	request.ArgsEmpty = false
	request.ArgsJson = string(jsonBuf)
	response, err = client.Invoke(context.TODO(), &request)

	assert.Nil(t, err)
	assert.Nil(t, response.Error)
	assert.True(t, response.ResultEmpty)
}
