package test

import (
	"context"
	"os"
	"testing"

	awstest "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

func TestDummyCommandableLambdaClient(t *testing.T) {

	lambdaArn := os.Getenv("LAMBDA_ARN")
	awsAccessId := os.Getenv("AWS_ACCESS_ID")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")

	if lambdaArn == "" || awsAccessId == "" || awsAccessKey == "" {
		t.Skip("AWS credentials is empty")
	}

	lambdaConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "aws",
		"connection.arn", lambdaArn,
		"credential.access_id", awsAccessId,
		"credential.access_key", awsAccessKey,
		"options.connection_timeout", 30000,
	)

	var client *DummyCommandableLambdaClient
	var fixture *awstest.DummyClientFixture

	client = NewDummyCommandableLambdaClient()
	client.Configure(context.Background(), lambdaConfig)

	fixture = awstest.NewDummyClientFixture(client)

	err := client.Open(context.Background())
	assert.Nil(t, err)

	defer client.Close(context.Background())

	t.Run("DummyCommandableLambdaClient.CrudOperations", fixture.TestCrudOperations)
}
