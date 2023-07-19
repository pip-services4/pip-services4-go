package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/controllers"
)

// DefaultSwaggerFactory are creates RPC components by their descriptors.

// See Factory
// See HttpEndpoint
// See HeartbeatRestController
// See StatusRestController
type DefaultSwaggerFactory struct {
	cbuild.Factory
	Descriptor               *cref.Descriptor
	SwaggerServiceDescriptor *cref.Descriptor
}

// NewDefaultSwaggerFactorymethod create a new instance of the factory.
func NewDefaultSwaggerFactory() *DefaultSwaggerFactory {
	c := DefaultSwaggerFactory{}
	c.Factory = *cbuild.NewFactory()
	c.Descriptor = cref.NewDescriptor("pip-services", "factory", "swagger", "default", "1.0")
	c.SwaggerServiceDescriptor = cref.NewDescriptor("pip-services", "swagger-controller", "*", "*", "1.0")

	c.RegisterType(c.SwaggerServiceDescriptor, controllers.NewSwaggerController)
	return &c
}
