package containers_test

import (
	"context"

	azurecont "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/containers"
	tbuild "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/test/build"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type DummyCommandableAzureFunction struct {
	*azurecont.CommandableAzureFunction
}

func NewDummyCommandableAzureFunction() *DummyCommandableAzureFunction {
	c := DummyCommandableAzureFunction{}
	c.CommandableAzureFunction = azurecont.NewCommandableAzureFunctionWithParams("dummy", "Dummy commandable azure function")
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))

	c.AddFactory(tbuild.NewDummyFactory())

	return &c
}
