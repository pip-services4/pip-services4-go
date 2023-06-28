package controllers_test

import (
	"context"

	azureserv "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/controllers"
	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type DummyCommandableAzureFunctionController struct {
	*azureserv.CommandableAzureFunctionController
}

func NewDummyCommandableAzureFunctionController() *DummyCommandableAzureFunctionController {
	c := DummyCommandableAzureFunctionController{}
	c.CommandableAzureFunctionController = azureserv.NewCommandableAzureFunctionController("dummies")
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return &c
}
