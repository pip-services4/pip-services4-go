package examples

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

var ControllerDescriptor = refer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "1.0")

func NewDummyFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(ControllerDescriptor, NewDummyController)

	return factory
}
