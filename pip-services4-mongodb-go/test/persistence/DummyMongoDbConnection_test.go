package test_persistence

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-mongodb-go/connect"
)

func TestDummyMongoDbConnection(t *testing.T) {

	var persistence *DummyMongoDbPersistence
	var fixture DummyPersistenceFixture
	var connection *conn.MongoDbConnection

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

	connection = conn.NewMongoDbConnection()
	connection.Configure(context.Background(), dbConfig)

	persistence = NewDummyMongoDbPersistence()
	descr := cref.NewDescriptor("pip-services", "connection", "mongodb", "default", "1.0")
	ref := cref.NewReferencesFromTuples(context.Background(), descr, connection)
	persistence.SetReferences(context.Background(), ref)

	fixture = *NewDummyPersistenceFixture(persistence)

	opnErr := connection.Open(context.Background())
	if opnErr != nil {
		t.Error("Error opened connection", opnErr)
		return
	}
	defer connection.Close(context.Background())

	opnErr = persistence.Open(context.Background())
	if opnErr != nil {
		t.Error("Error opened persistence", opnErr)
		return
	}
	defer persistence.Close(context.Background())

	opnErr = persistence.Clear(context.Background())
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummyMongoDbConnection:CRUD", fixture.TestCrudOperations)
	t.Run("DummyMongoDbConnection:Batch", fixture.TestBatchOperations)

}
