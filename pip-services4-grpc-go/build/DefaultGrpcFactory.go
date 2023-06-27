package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	grpcservices "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
)

// DefaultGrpcFactory creates GRPC components by their descriptors.
// See Factory
// See GrpcEndpoint
// See HeartbeatGrpcService
// See StatusGrpcService
type DefaultGrpcFactory struct {
	*cbuild.Factory
	Descriptor             *cref.Descriptor
	GrpcEndpointDescriptor *cref.Descriptor
}

// NewDefaultGrpcFactory method are creates a new instance of the factory.
func NewDefaultGrpcFactory() *DefaultGrpcFactory {

	c := DefaultGrpcFactory{
		Factory: cbuild.NewFactory(),
	}
	c.Descriptor = cref.NewDescriptor("pip-services", "factory", "grpc", "default", "1.0")
	c.GrpcEndpointDescriptor = cref.NewDescriptor("pip-services", "endpoint", "grpc", "*", "1.0")

	c.RegisterType(c.GrpcEndpointDescriptor, grpcservices.NewGrpcEndpoint)
	return &c
}
