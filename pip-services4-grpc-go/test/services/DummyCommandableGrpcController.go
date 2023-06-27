package test_services

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	grpcservices "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
)

type DummyCommandableGrpcController struct {
	*grpcservices.CommandableGrpcController
}

func NewDummyCommandableGrpcController() *DummyCommandableGrpcController {
	c := &DummyCommandableGrpcController{}
	c.CommandableGrpcController = grpcservices.InheritCommandableGrpcController(c, "dummy")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return c
}
