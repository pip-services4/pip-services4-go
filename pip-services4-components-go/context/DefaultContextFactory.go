package context

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// Creates information components by their descriptors.

var ContextInfoDescriptor = refer.NewDescriptor("pip-services", "context-info", "default", "*", "1.0")
var ContainerInfoDescriptor = refer.NewDescriptor("pip-services", "container-info", "default", "*", "1.0")
var ContainerInfoDescriptor2 = refer.NewDescriptor("pip-services-container", "container-info", "default", "*", "1.0")

// NewDefaultContextFactory create a new instance of the factory.
//
//	Returns: *build.Factory
func NewDefaultContextFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(ContextInfoDescriptor, NewContextInfo)
	factory.RegisterType(ContainerInfoDescriptor, NewContextInfo)
	factory.RegisterType(ContainerInfoDescriptor2, NewContextInfo)

	return factory
}
