package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	services "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

// DefaultHttpFactory are creates RPC components by their descriptors
//
//	see Factory
//	see HttpEndpoint
//	see HeartbeatRestController
//	see StatusRestController
type DefaultHttpFactory struct {
	cbuild.Factory
}

// NewDefaultHttpFactory creates a new instance of the factory.
func NewDefaultHttpFactory() *DefaultHttpFactory {
	c := DefaultHttpFactory{}
	c.Factory = *cbuild.NewFactory()

	httpEndpointDescriptor := cref.NewDescriptor("pip-services", "endpoint", "http", "*", "1.0")
	statusControllerDescriptor := cref.NewDescriptor("pip-services", "status-controller", "http", "*", "1.0")
	heartbeatControllerDescriptor := cref.NewDescriptor("pip-services", "heartbeat-controller", "http", "*", "1.0")

	c.RegisterType(httpEndpointDescriptor, services.NewHttpEndpoint)
	c.RegisterType(heartbeatControllerDescriptor, services.NewHeartbeatRestController)
	c.RegisterType(statusControllerDescriptor, services.NewStatusRestController)
	return &c
}
