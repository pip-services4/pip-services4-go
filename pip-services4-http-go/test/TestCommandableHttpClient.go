package test

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-http-go/clients"
)

type TestCommandableHttpClient struct {
	*clients.CommandableHttpClient
}

func NewTestCommandableHttpClient(baseRoute string) *TestCommandableHttpClient {
	c := &TestCommandableHttpClient{}
	c.CommandableHttpClient = clients.NewCommandableHttpClient(baseRoute)
	return c
}
