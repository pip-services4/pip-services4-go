package test_persistence

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

func TestDummyMapFilePersistence(t *testing.T) {
	filename := "../../data/dummies_map.json"

	//cleaning file before testing
	f, err := os.Create(filename)
	if err != nil {
		t.Error("Can't clean file: ", filename)
	}
	_ = f.Close()

	persistence := NewDummyMapFilePersistence(filename)
	persistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	defer persistence.Close(context.Background())

	fixture := NewDummyMapPersistenceFixture(persistence)
	_ = persistence.Open(context.Background())

	t.Run("DummyMapFilePersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyMapFilePersistence:Batch", fixture.TestBatchOperations)

}
