package test_controllers

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

type DummyLambdaControllerFactory struct {
	*cbuild.Factory
	Descriptor                    *cref.Descriptor
	ControllerDescriptor          *cref.Descriptor
	LambdaControllerDescriptor    *cref.Descriptor
	CmdLambdaControllerDescriptor *cref.Descriptor
}

func NewDummyLambdaControllerFactory() *DummyLambdaControllerFactory {

	c := DummyLambdaControllerFactory{
		Factory:                       cbuild.NewFactory(),
		Descriptor:                    cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		LambdaControllerDescriptor:    cref.NewDescriptor("pip-services-dummies", "controller", "awslambda", "*", "1.0"),
		CmdLambdaControllerDescriptor: cref.NewDescriptor("pip-services-dummies", "controller", "commandable-awslambda", "*", "1.0"),
	}

	c.RegisterType(c.LambdaControllerDescriptor, NewDummyLambdaController)
	c.RegisterType(c.CmdLambdaControllerDescriptor, NewDummyCommandableLambdaController)
	return &c
}
