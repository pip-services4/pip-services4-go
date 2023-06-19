package test

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

var ShutdownDescriptor = refer.NewDescriptor("pip-services", "shutdown", "*", "*", "1.0")

func NewDefaultTestFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(ShutdownDescriptor, NewShutdown)

	return factory
}
