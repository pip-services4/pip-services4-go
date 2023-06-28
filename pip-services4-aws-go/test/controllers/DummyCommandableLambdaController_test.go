package test_controllers

import (
	"context"
	"encoding/json"
	"testing"

	awstest "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

func TestDummyCommandableLambdaController(t *testing.T) {
	ctx := context.Background()

	restConfig := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"service.descriptor", "pip-services-dummies:service:default:default:1.0",
		"controller.descriptor", "pip-services-dummies:controller:commandable-awslambda:default:1.0",
	)

	var _dummy1 awstest.Dummy
	var _dummy2 awstest.Dummy
	var lambda *DummyLambdaFunction
	srv := awstest.NewDummyService()

	lambda = NewDummyLambdaFunction()
	lambda.Configure(ctx, restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(ctx,
		cref.NewDescriptor("pip-services-dummies", "service", "default", "default", "1.0"), srv,
	)
	lambda.SetReferences(ctx, references)
	opnErr := lambda.Open(ctx)
	assert.Nil(t, opnErr)
	defer lambda.Close(ctx)

	_dummy1 = awstest.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	_dummy2 = awstest.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	var dummy1 awstest.Dummy

	params := make(map[string]any)

	// Create one dummy
	params["dummy"] = _dummy1
	params["cmd"] = "dummy.create_dummy"

	resBody, bodyErr := lambda.Act(params)
	assert.Nil(t, bodyErr)

	var dummy awstest.Dummy
	jsonErr := json.Unmarshal([]byte(resBody), &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy1.Content)
	assert.Equal(t, dummy.Key, _dummy1.Key)

	dummy1 = dummy

	// Create another dummy
	params["dummy"] = _dummy2
	params["cmd"] = "dummy.create_dummy"

	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)

	jsonErr = json.Unmarshal([]byte(resBody), &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy2.Content)
	assert.Equal(t, dummy.Key, _dummy2.Key)
	//dummy2 = dummy

	// Get all dummies
	delete(params, "dummy")
	params["cmd"] = "dummy.get_dummies"
	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)

	var dummies cquery.DataPage[awstest.Dummy]
	jsonErr = json.Unmarshal([]byte(resBody), &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy

	dummy1.Content = "Updated Content 1"

	params["dummy"] = dummy1
	params["cmd"] = "dummy.update_dummy"

	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)
	jsonErr = json.Unmarshal([]byte(resBody), &dummy)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)

	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, _dummy1.Key)
	dummy1 = dummy

	// Delete dummy
	delete(params, "dummy")
	params["dummy_id"] = dummy1.Id
	params["cmd"] = "dummy.delete_dummy"
	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)

	// Try to get delete dummy
	dummies.Data = dummies.Data[:0]
	dummies.Total = 0

	params["dummy_id"] = dummy1.Id
	params["cmd"] = "dummy.get_dummy_by_id"

	resBody, bodyErr = lambda.Act(params)
	assert.Nil(t, bodyErr)
	jsonErr = json.Unmarshal([]byte(resBody), &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 0)
}
