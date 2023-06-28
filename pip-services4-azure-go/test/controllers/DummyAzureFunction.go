package controllers_test

import (
	azuresrv "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/containers"
	tbuild "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/test/build"
)

type DummyAzureFunction struct {
	*azuresrv.AzureFunction
}

func NewDummyAzureFunction() *DummyAzureFunction {
	c := DummyAzureFunction{AzureFunction: azuresrv.NewAzureFunctionWithParams("dummy", "Dummy azure function")}
	c.AddFactory(tbuild.NewDummyFactory())
	c.AddFactory(NewDummyAzureFunctionServiceFactory())

	return &c
}
