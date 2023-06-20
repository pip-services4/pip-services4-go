package test_persistence

import (
	"context"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type DummyFilePersistence struct {
	DummyMemoryPersistence
	persister *cpersist.JsonFilePersister[Dummy]
}

func NewDummyFilePersistence(path string) *DummyFilePersistence {
	c := &DummyFilePersistence{
		DummyMemoryPersistence: *NewDummyMemoryPersistence(),
	}
	persister := cpersist.NewJsonFilePersister[Dummy](path)
	c.persister = persister
	c.Loader = persister
	c.Saver = persister
	return c
}

func (c *DummyFilePersistence) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.DummyMemoryPersistence.Configure(ctx, config)
	c.persister.Configure(ctx, config)
}
