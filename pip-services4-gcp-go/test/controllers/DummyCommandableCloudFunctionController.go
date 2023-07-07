package services_test

import (
	"context"

	crefer "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	gcpctrl "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/controllers"
)

type DummyCommandableCloudFunctionController struct {
	*gcpctrl.CommandableCloudFunctionController
}

func NewDummyCommandableCloudFunctionController() *DummyCommandableCloudFunctionController {
	c := DummyCommandableCloudFunctionController{}
	c.CommandableCloudFunctionController = gcpctrl.NewCommandableCloudFunctionController("dummies")
	c.DependencyResolver.Put(context.Background(), "service", crefer.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return &c
}
