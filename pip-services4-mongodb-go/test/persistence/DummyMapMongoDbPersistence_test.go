package test_persistence

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

func TestDummyMapMongoDbPersistence(t *testing.T) {

	var persistence *DummyMapMongoDbPersistence
	var fixture DummyMapPersistenceFixture

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

	persistence = NewDummyMapMongoDbPersistence()
	persistence.Configure(context.Background(), dbConfig)

	fixture = *NewDummyMapPersistenceFixture(persistence)

	opnErr := persistence.Open(context.Background())
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

	t.Run("DummyMapMongoDbPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyMapMongoDbPersistence:Batch", fixture.TestBatchOperations)

}
