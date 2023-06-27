package test

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-http-go/clients"
)

type TestRestClient struct {
	*clients.RestClient
}

func NewTestRestClient(baseRoute string) *TestRestClient {
	c := &TestRestClient{}
	c.RestClient = clients.NewRestClient()
	c.BaseRoute = baseRoute
	return c
}
