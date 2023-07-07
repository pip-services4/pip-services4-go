package services_test

import (
	gcpsrv "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/containers"
	tbuild "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/build"
)

type DummyCloudFunction struct {
	*gcpsrv.CloudFunction
}

func NewDummyCloudFunction() *DummyCloudFunction {
	c := DummyCloudFunction{CloudFunction: gcpsrv.NewCloudFunctionWithParams("dummy", "Dummy cloud function")}
	c.AddFactory(tbuild.NewDummyFactory())
	c.AddFactory(NewDummyCloudFunctionControllerFactory())

	return &c
}
