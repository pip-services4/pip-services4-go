package services_test

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	logic "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/logic"
)

type DummyCloudFunctionControllerFactory struct {
	cbuild.Factory
	Descriptor                   *cref.Descriptor
	ServiceDescriptor            *cref.Descriptor
	CloudControllerDescriptor    *cref.Descriptor
	CmdCloudControllerDescriptor *cref.Descriptor
}

func NewDummyCloudFunctionControllerFactory() *DummyCloudFunctionControllerFactory {

	c := DummyCloudFunctionControllerFactory{
		Factory:                      *cbuild.NewFactory(),
		Descriptor:                   cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		CloudControllerDescriptor:    cref.NewDescriptor("pip-services-dummies", "controller", "cloudfunc", "*", "1.0"),
		CmdCloudControllerDescriptor: cref.NewDescriptor("pip-services-dummies", "controller", "commandable-cloudfunc", "*", "1.0"),
		ServiceDescriptor:            cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "1.0"),
	}

	c.RegisterType(c.CloudControllerDescriptor, NewDummyCloudFunctionService)
	c.RegisterType(c.CmdCloudControllerDescriptor, NewDummyCommandableCloudFunctionController)
	c.RegisterType(c.ServiceDescriptor, logic.NewDummyService)
	return &c
}
