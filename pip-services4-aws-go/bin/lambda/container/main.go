package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	awscont "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test/containers"
)

func main() {
	ctx := context.Background()
	var container *awscont.DummyLambdaFunction

	container = awscont.NewDummyLambdaFunction()

	defer container.Close(ctx)
	err := container.Run(ctx)
	if err != nil {
		panic(err)
	}
	lambda.Start(container.GetHandler())
}
