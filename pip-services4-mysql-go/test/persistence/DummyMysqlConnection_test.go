package test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	conn "github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/connect"
	tf "github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestDummyMySqlConnection(t *testing.T) {

	var persistence *DummyMySqlPersistence
	var fixture tf.DummyPersistenceFixture
	var connection *conn.MySqlConnection

	mysqlUri := os.Getenv("MYSQL_URI")
	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		mysqlHost = "localhost"
	}

	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	mysqlDatabase := os.Getenv("MYSQL_DB")
	if mysqlDatabase == "" {
		mysqlDatabase = "test"
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "" {
		mysqlUser = "user"
	}
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" {
		mysqlPassword = "password"
	}

	if mysqlUri == "" && mysqlHost == "" {
		t.Skip("Connection params not set")
	}

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", mysqlUri,
		"connection.host", mysqlHost,
		"connection.port", mysqlPort,
		"connection.database", mysqlDatabase,
		"credential.username", mysqlUser,
		"credential.password", mysqlPassword,
	)

	connection = conn.NewMySqlConnection()
	connection.Configure(context.Background(), dbConfig)

	persistence = NewDummyMySqlPersistence()
	descr := cref.NewDescriptor("pip-services", "connection", "mysql", "default", "1.0")
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

	t.Run("DummyMySqlConnection:CRUD", fixture.TestCrudOperations)

	opnErr = persistence.Clear(context.Background())
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummyMySqlConnection:Batch", fixture.TestBatchOperations)

}
