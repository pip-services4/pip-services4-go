package controllers_test

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type DummyAzureFunctionServiceFactory struct {
	cbuild.Factory
	Descriptor                *cref.Descriptor
	ControllerDescriptor      *cref.Descriptor
	AzureServiceDescriptor    *cref.Descriptor
	CmdAzureServiceDescriptor *cref.Descriptor
}

func NewDummyAzureFunctionServiceFactory() *DummyAzureFunctionServiceFactory {

	c := DummyAzureFunctionServiceFactory{
		Factory:                   *cbuild.NewFactory(),
		Descriptor:                cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		AzureServiceDescriptor:    cref.NewDescriptor("pip-services-dummies", "controller", "azurefunc", "*", "1.0"),
		CmdAzureServiceDescriptor: cref.NewDescriptor("pip-services-dummies", "controller", "commandable-azurefunc", "*", "1.0"),
	}

	c.RegisterType(c.AzureServiceDescriptor, NewDummyAzureFunctionController)
	c.RegisterType(c.CmdAzureServiceDescriptor, NewDummyCommandableAzureFunctionController)
	return &c
}
