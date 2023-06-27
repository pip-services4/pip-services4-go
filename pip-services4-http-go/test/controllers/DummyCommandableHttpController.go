package test_controllers

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	services "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"
)

type DummyCommandableHttpController struct {
	*services.CommandableHttpController
}

func NewDummyCommandableHttpController() *DummyCommandableHttpController {
	c := &DummyCommandableHttpController{}
	c.CommandableHttpController = services.InheritCommandableHttpController(c, "dummies")
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return c
}

func (c *DummyCommandableHttpController) Register() {
	if !c.SwaggerAuto && c.SwaggerEnabled {
		c.RegisterOpenApiSpec("swagger yaml content")
	}
	c.CommandableHttpController.Register()
}
