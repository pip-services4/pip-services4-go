package test_persistence

import (
	"context"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type DummyInterfacableFilePersistence struct {
	DummyInterfacableMemoryPersistence
	persister *cpersist.JsonFilePersister[DummyInterfacable]
}

func NewDummyInterfacableFilePersistence(path string) *DummyInterfacableFilePersistence {
	c := &DummyInterfacableFilePersistence{
		DummyInterfacableMemoryPersistence: *NewDummyInterfacableMemoryPersistence(),
	}
	persister := cpersist.NewJsonFilePersister[DummyInterfacable](path)
	c.persister = persister
	c.Loader = persister
	c.Saver = persister
	return c
}

func (c *DummyInterfacableFilePersistence) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.DummyInterfacableMemoryPersistence.Configure(ctx, config)
	c.persister.Configure(ctx, config)
}
