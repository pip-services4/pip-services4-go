package test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	tf "github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/test/fixtures"
)

func TestDummySqlitePersistence(t *testing.T) {

	var persistence *DummySqlitePersistence
	var fixture tf.DummyPersistenceFixture

	sqliteDatabase := os.Getenv("SQLITE_DB")
	if sqliteDatabase == "" {
		sqliteDatabase = "../../data/test.db"
	}

	if sqliteDatabase == "" {
		panic("Connection params losse")
	}

	dbConfig := cconf.NewConfigParamsFromTuples(
		"connection.database", sqliteDatabase,
	)

	persistence = NewDummySqlitePersistence()
	fixture = *tf.NewDummyPersistenceFixture(persistence)
	persistence.Configure(context.Background(), dbConfig)

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

	t.Run("DummySqlitePersistence:CRUD", fixture.TestCrudOperations)

	opnErr = persistence.Clear(context.Background())
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummySqlitePersistence:Batch", fixture.TestBatchOperations)

	opnErr = persistence.Clear(context.Background())
	if opnErr != nil {
		t.Error("Error cleaned persistence", opnErr)
		return
	}

	t.Run("DummySqlitePersistence:Random", fixture.TestRandomOperation)
}
