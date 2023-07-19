package test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/connect"
	tf "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestDummySqlServerConnection(t *testing.T) {

	var persistence *DummySqlServerPersistence
	var fixture tf.DummyPersistenceFixture
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
	)

	connection = conn.NewSqlServerConnection()
	connection.Configure(context.Background(), dbConfig)

	persistence = NewDummySqlServerPersistence()
	descr := cref.NewDescriptor("pip-services", "connection", "sqlserver", "default", "1.0")
	ref := cref.NewReferencesFromTuples(context.Background(), descr, connection)
	persistence.SetReferences(context.Background(), ref)

	fixture = *tf.NewDummyPersistenceFixture(persistence)

	opnErr := connection.Open(context.Background())
	if opnErr != nil {
		t.Error("Error opened connection", opnErr)
		return
	}
	defer func() {
		err := connection.Close(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	opnErr = persistence.Open(context.Background())
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

	t.Run("Connection", func(t *testing.T) {
		assert.NotNil(t, connection.GetConnection())
		assert.NotNil(t, connection.GetDatabaseName())
		assert.NotEqual(t, "", connection.GetDatabaseName())
	})

	t.Run("DummySqlServerConnection:CRUD", fixture.TestCrudOperations)

	opnErr = persistence.Clear(context.Background())
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummySqlServerConnection:Batch", fixture.TestBatchOperations)

}
