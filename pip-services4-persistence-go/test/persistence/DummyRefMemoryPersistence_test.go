package test_persistence

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

func TestDummyRefMemoryPersistence(t *testing.T) {
	persister := NewDummyRefMemoryPersistence()
	persister.Configure(context.Background(), cconf.NewEmptyConfigParams())

	fixture := NewDummyRefPersistenceFixture(persister)

	t.Run("DummyRefMemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyRefMemoryPersistence:Batch", fixture.TestBatchOperations)

}
