package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	awsserv "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test/controllers"
)

func main() {
	ctx := context.Background()
	var container *awsserv.DummyLambdaFunction

	container = awsserv.NewDummyLambdaFunction()

	defer container.Close(ctx)
	opnErr := container.Run(ctx)
	if opnErr == nil {
		lambda.Start(container.GetHandler())
	}

}
