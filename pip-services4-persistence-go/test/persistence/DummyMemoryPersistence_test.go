package test_persistence

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

func TestDummyMemoryPersistence(t *testing.T) {
	persistence := NewDummyMemoryPersistence()
	persistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	fixture := NewDummyPersistenceFixture(persistence)

	t.Run("DummyMemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyMemoryPersistence:Batch", fixture.TestBatchOperations)

}
