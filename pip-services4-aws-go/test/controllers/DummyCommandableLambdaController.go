package test_controllers

import (
	"context"

	awsserv "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/controllers"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type DummyCommandableLambdaController struct {
	*awsserv.CommandableLambdaController
}

func NewDummyCommandableLambdaController() *DummyCommandableLambdaController {
	c := &DummyCommandableLambdaController{}
	c.CommandableLambdaController = awsserv.InheritCommandableLambdaController(c, "dummy")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return c
}
