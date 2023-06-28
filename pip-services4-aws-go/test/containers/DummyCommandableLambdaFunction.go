package test_container

import (
	"context"

	awscont "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/containers"
	awstest "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type DummyCommandableLambdaFunction struct {
	*awscont.CommandableLambdaFunction
}

func NewDummyCommandableLambdaFunction() *DummyCommandableLambdaFunction {
	c := &DummyCommandableLambdaFunction{}
	c.CommandableLambdaFunction = awscont.NewCommandableLambdaFunction("dummy", "Dummy lambda function")

	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	c.AddFactory(awstest.NewDummyFactory())
	return c
}
