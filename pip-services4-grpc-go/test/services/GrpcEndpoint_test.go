package test_services

import (
	"context"
	"testing"
	"time"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	grpcservices "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
	"github.com/stretchr/testify/assert"
)

func TestGrpcEndpoint(t *testing.T) {
	ctx := context.Background()

	grpcConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", 3005,
	)

	endpoint := grpcservices.NewGrpcEndpoint()
	endpoint.Configure(ctx, grpcConfig)

	endpoint.Open(ctx)

	// wait server start
	<-time.After(100 * time.Millisecond)

	assert.True(t, endpoint.IsOpen())
	endpoint.Close(ctx)
}
