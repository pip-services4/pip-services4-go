package containers_test

import (
	"context"

	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	gcpcont "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/containers"
	tbuild "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/build"
)

type DummyCommandableCloudFunction struct {
	*gcpcont.CommandableCloudFunction
}

func NewDummyCommandableCloudFunction() *DummyCommandableCloudFunction {
	c := DummyCommandableCloudFunction{}
	c.CommandableCloudFunction = gcpcont.NewCommandableCloudFunctionWithParams("dummy", "Dummy commandable cloud function")
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))

	c.AddFactory(tbuild.NewDummyFactory())

	return &c
}
