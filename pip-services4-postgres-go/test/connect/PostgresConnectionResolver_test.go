package test_connect

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConnectionResolver(t *testing.T) {

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.host", "localhost",
		"connection.port", 5432,
		"connection.database", "test",
		"connection.sslmode", "verify-ca",
		"credential.username", "postgres",
		"credential.password", "postgres",
	)

	resolver := conn.NewPostgresConnectionResolver()
	resolver.Configure(context.Background(), dbConfig)

	config, err := resolver.Resolve(context.Background())
	assert.Nil(t, err)

	assert.NotNil(t, config)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/test?sslmode=verify-ca", config)
}
