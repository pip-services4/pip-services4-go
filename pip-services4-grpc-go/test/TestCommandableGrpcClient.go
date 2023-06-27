package test

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
)

type TestCommandableGrpcClient struct {
	*clients.CommandableGrpcClient
}

func NewTestCommandableGrpcClient(name string) *TestCommandableGrpcClient {
	c := &TestCommandableGrpcClient{}
	c.CommandableGrpcClient = clients.NewCommandableGrpcClient(name)
	return c
}
