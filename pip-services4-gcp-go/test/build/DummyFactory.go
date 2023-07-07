package build_test

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	tlogic "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/logic"
)

type DummyFactory struct {
	cbuild.Factory
	Descriptor        *cref.Descriptor
	ServiceDescriptor *cref.Descriptor
}

// NewDefaultRpcFactory creates a new instance of the factory.
func NewDummyFactory() *DummyFactory {
	c := DummyFactory{
		Factory:           *cbuild.NewFactory(),
		Descriptor:        cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		ServiceDescriptor: cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "1.0"),
	}

	c.RegisterType(c.ServiceDescriptor, tlogic.NewDummyService)
	return &c
}
