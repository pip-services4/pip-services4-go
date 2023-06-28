package test_controllers

import (
	awscont "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/containers"
	awstest "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test"
)

type DummyLambdaFunction struct {
	*awscont.LambdaFunction
}

func NewDummyLambdaFunction() *DummyLambdaFunction {
	c := &DummyLambdaFunction{}
	c.LambdaFunction = awscont.InheriteLambdaFunction(c, "dummy", "Dummy lambda function")
	c.AddFactory(awstest.NewDummyFactory())
	c.AddFactory(NewDummyLambdaControllerFactory())
	return c
}

func (c *DummyLambdaFunction) Register() {}
