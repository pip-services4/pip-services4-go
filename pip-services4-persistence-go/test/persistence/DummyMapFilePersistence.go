package test_persistence

import (
	"context"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type DummyMapFilePersistence struct {
	DummyMapMemoryPersistence
	persister *cpersist.JsonFilePersister[DummyMap]
}

func NewDummyMapFilePersistence(path string) *DummyMapFilePersistence {
	c := &DummyMapFilePersistence{
		DummyMapMemoryPersistence: *NewDummyMapMemoryPersistence(),
	}

	persister := cpersist.NewJsonFilePersister[DummyMap](path)
	c.persister = persister
	c.Loader = persister
	c.Saver = persister

	return c
}

func (c *DummyMapFilePersistence) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.DummyMapMemoryPersistence.Configure(ctx, config)
	c.persister.Configure(ctx, config)
}
