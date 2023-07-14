package test_connect

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-mongodb-go/connect"
	"github.com/stretchr/testify/assert"
)

func TestMongoDBConnection(t *testing.T) {
	var connection conn.MongoDbConnection

	mongoUri := os.Getenv("MONGO_URI")
	mongoHost := os.Getenv("MONGO_HOST")
	if mongoHost == "" {
		mongoHost = "localhost"
	}
	mongoPort := os.Getenv("MONGO_PORT")
	if mongoPort == "" {
		mongoPort = "27017"
	}
	mongoDatabase := os.Getenv("MONGO_DB")
	if mongoDatabase == "" {
		mongoDatabase = "test"
	}
	if mongoUri == "" && mongoHost == "" {
		return
	}

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", mongoUri,
		"connection.host", mongoHost,
		"connection.port", mongoPort,
		"connection.database", mongoDatabase,
	)

	connection = *conn.NewMongoDbConnection()
	connection.Configure(context.Background(), dbConfig)

	err := connection.Open(context.Background())
	assert.Nil(t, err)

	defer connection.Close(context.Background())

	assert.NotNil(t, connection.GetConnection())
	assert.NotNil(t, connection.GetDatabase())
	assert.NotEqual(t, "", connection.GetDatabaseName())

}
