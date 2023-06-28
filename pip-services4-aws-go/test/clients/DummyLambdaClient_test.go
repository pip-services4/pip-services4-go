package test

import (
	"context"
	"os"
	"testing"

	awstest "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

func TestDummyLambdaClient(t *testing.T) {
	ctx := context.Background()

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

	var client *DummyLambdaClient
	var fixture *awstest.DummyClientFixture

	client = NewDummyLambdaClient()
	client.Configure(ctx, lambdaConfig)

	fixture = awstest.NewDummyClientFixture(client)

	client.Open(ctx)

	defer client.Close(ctx)

	t.Run("DummyLambdaClient.CrudOperations", fixture.TestCrudOperations)
}
