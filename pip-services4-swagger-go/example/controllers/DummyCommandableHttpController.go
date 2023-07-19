package example_services

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cctrl "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type DummyCommandableHttpController struct {
	*cctrl.CommandableHttpController
}

func NewDummyCommandableHttpController() *DummyCommandableHttpController {
	c := DummyCommandableHttpController{}
	c.CommandableHttpController = cctrl.InheritCommandableHttpController(&c, "dummies2")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return &c
}
