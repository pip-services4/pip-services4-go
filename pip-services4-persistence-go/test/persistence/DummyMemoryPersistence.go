package test_persistence

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type DummyMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence[Dummy, string]
}

func NewDummyMemoryPersistence() *DummyMemoryPersistence {
	return &DummyMemoryPersistence{
		*cpersist.NewIdentifiableMemoryPersistence[Dummy, string](),
	}
}

func (c *DummyMemoryPersistence) GetPageByFilter(ctx context.Context, filter cdata.FilterParams, paging cdata.PagingParams) (cdata.DataPage[Dummy], error) {

	var key string

	if _key, ok := filter.GetAsNullableString("Key"); ok {
		key = _key
	}

	return c.IdentifiableMemoryPersistence.
		GetPageByFilter(ctx,
			func(item Dummy) bool {
				if key != "" && item.Key != key {
					return false
				}
				return true
			},
			paging,
			func(a, b Dummy) bool {
				return len(a.Key) < len(b.Key)
			},
			nil,
		)
}

func (c *DummyMemoryPersistence) GetCountByFilter(ctx context.Context, filter cdata.FilterParams) (count int64, err error) {

	var key string

	if _key, ok := filter.GetAsNullableString("Key"); ok {
		key = _key
	}

	return c.IdentifiableMemoryPersistence.
		GetCountByFilter(ctx,
			func(item Dummy) bool {
				if key != "" && item.Key != key {
					return false
				}
				return true
			})
}
