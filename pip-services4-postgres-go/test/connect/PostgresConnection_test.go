package test_connect

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConnection(t *testing.T) {
	var connection *conn.PostgresConnection

	postgresUri := os.Getenv("POSTGRES_URI")
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = "localhost"
	}
	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresPort == "" {
		postgresPort = "5432"
	}
	postgresDatabase := os.Getenv("POSTGRES_DB")
	if postgresDatabase == "" {
		postgresDatabase = "test"
	}
	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		postgresUser = "postgres"
	}
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		postgresPassword = "postgres#"
	}

	if postgresUri == "" && postgresHost == "" {
		panic("Connection params losse")
	}

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", postgresUri,
		"connection.host", postgresHost,
		"connection.port", postgresPort,
		"connection.database", postgresDatabase,
		"credential.username", postgresUser,
		"credential.password", postgresPassword,
		"options.max_pool_size", 10,
		"options.connect_timeout", 100,
		"options.idle_timeout", 100,
	)

	connection = conn.NewPostgresConnection()
	connection.Configure(context.Background(), dbConfig)
	err := connection.Open(context.Background())
	assert.Nil(t, err)

	defer connection.Close(context.Background())

	assert.NotNil(t, connection.GetConnection())
	assert.NotNil(t, connection.GetDatabaseName())
	assert.NotEqual(t, "", connection.GetDatabaseName())
}
