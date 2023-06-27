package test

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
)

type TestGrpcClient struct {
	*clients.GrpcClient
}

func NewTestGrpcClient(name string) *TestGrpcClient {
	c := &TestGrpcClient{}
	c.GrpcClient = clients.NewGrpcClient(name)
	return c
}
