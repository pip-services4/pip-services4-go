package test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	tf "github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/test/fixtures"
)

func TestDummyRefMySqlPersistence(t *testing.T) {

	var persistence *DummyRefMySqlPersistence
	var fixture tf.DummyRefPersistenceFixture

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
	persistence = NewDummyRefMySqlPersistence()
	persistence.Configure(context.Background(), dbConfig)

	fixture = *tf.NewDummyRefPersistenceFixture(persistence)

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

	t.Run("DummyRefMySqlPersistence:CRUD", fixture.TestCrudOperations)

	opnErr = persistence.Clear(context.Background())
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummyRefMySqlPersistence:Batch", fixture.TestBatchOperations)

}
