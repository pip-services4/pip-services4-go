package test_persistence

import (
	"context"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

// extends DummyMemoryPersistence
type DummyRefFilePersistence struct {
	DummyRefMemoryPersistence
	persister *cpersist.JsonFilePersister[*DummyRef]
}

func NewDummyRefFilePersistence(path string) *DummyRefFilePersistence {
	c := &DummyRefFilePersistence{
		DummyRefMemoryPersistence: *NewDummyRefMemoryPersistence(),
	}
	persister := cpersist.NewJsonFilePersister[*DummyRef](path)
	c.persister = persister
	c.Loader = persister
	c.Saver = persister
	return c
}

func (c *DummyRefFilePersistence) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.DummyRefMemoryPersistence.Configure(ctx, config)
	c.persister.Configure(ctx, config)
}
