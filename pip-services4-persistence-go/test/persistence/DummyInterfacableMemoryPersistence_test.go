package test_persistence

import (
	"context"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
)

func TestDummyInterfacableMemoryPersistence(t *testing.T) {
	persistence := NewDummyInterfacableMemoryPersistence()
	persistence.Configure(context.Background(), cconf.NewEmptyConfigParams())

	fixture := NewDummyInterfacablePersistenceFixture(persistence)

	t.Run("DummyInterfacableMemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyInterfacableMemoryPersistence:Batch", fixture.TestBatchOperations)

}
