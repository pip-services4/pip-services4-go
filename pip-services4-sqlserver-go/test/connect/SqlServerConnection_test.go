package test_connect

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestSqlServerConnection(t *testing.T) {
	var connection *conn.SqlServerConnection

	sqlserverUri := os.Getenv("SQLSERVER_URI")
	sqlserverHost := os.Getenv("SQLSERVER_HOST")
	if sqlserverHost == "" {
		sqlserverHost = "localhost"
	}
	sqlserverPort := os.Getenv("SQLSERVER_PORT")
	if sqlserverPort == "" {
		sqlserverPort = "1433"
	}
	sqlserverDatabase := os.Getenv("SQLSERVER_DB")
	if sqlserverDatabase == "" {
		sqlserverDatabase = "master"
	}
	sqlserverUser := os.Getenv("SQLSERVER_USER")
	if sqlserverUser == "" {
		sqlserverUser = "sa"
	}
	sqlserverPassword := os.Getenv("SQLSERVER_PASSWORD")
	if sqlserverPassword == "" {
		sqlserverPassword = "sqlserver_123"
	}

	if sqlserverUri == "" && sqlserverHost == "" {
		t.Skip("No SqlServer credentials")
	}

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", sqlserverUri,
		"connection.host", sqlserverHost,
		"connection.port", sqlserverPort,
		"connection.database", sqlserverDatabase,
		"credential.username", sqlserverUser,
		"credential.password", sqlserverPassword,
		"options.max_pool_size", 10,
		"options.connect_timeout", 100,
		"options.idle_timeout", 100,
	)

	connection = conn.NewSqlServerConnection()
	connection.Configure(context.Background(), dbConfig)
	err := connection.Open(context.Background())
	assert.Nil(t, err)

	assert.NotNil(t, connection.GetConnection())
	assert.NotEmpty(t, connection.GetDatabaseName())
	assert.NotNil(t, connection.GetDatabaseName())

	err = connection.Close(context.Background())
	assert.Nil(t, err)
}
