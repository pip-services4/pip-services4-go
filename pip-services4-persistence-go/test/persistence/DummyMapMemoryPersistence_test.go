package test_persistence

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

func TestDummyMapMemoryPersistence(t *testing.T) {
	persister := NewDummyMapMemoryPersistence()
	persister.Configure(context.Background(), cconf.NewEmptyConfigParams())

	fixture := NewDummyMapPersistenceFixture(persister)

	t.Run("DummyMapMemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyMapMemoryPersistence:Batch", fixture.TestBatchOperations)

}
