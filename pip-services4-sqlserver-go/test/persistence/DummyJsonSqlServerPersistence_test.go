package test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	tf "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/test/fixtures"
)

func TestDummyJsonSqlServerPersistence(t *testing.T) {

	var persistence *DummyJsonSqlServerPersistence
	var fixture tf.DummyPersistenceFixture

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
	)

	persistence = NewDummyJsonSqlServerPersistence()
	fixture = *tf.NewDummyPersistenceFixture(persistence)
	persistence.Configure(context.Background(), dbConfig)

	opnErr := persistence.Open(context.Background())
	if opnErr != nil {
		t.Error("Error opened persistence", opnErr)
		return
	}

	defer func() {
		err := persistence.Close(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	opnErr = persistence.Clear(context.Background())
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummySqlServerConnection:CRUD", fixture.TestCrudOperations)

	opnErr = persistence.Clear(context.Background())
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummySqlServerConnection:Batch", fixture.TestBatchOperations)

}
