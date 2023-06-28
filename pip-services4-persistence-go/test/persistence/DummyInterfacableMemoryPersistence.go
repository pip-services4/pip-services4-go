package test_persistence

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type DummyInterfacableMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence[DummyInterfacable, string]
}

func NewDummyInterfacableMemoryPersistence() *DummyInterfacableMemoryPersistence {
	return &DummyInterfacableMemoryPersistence{
		*cpersist.NewIdentifiableMemoryPersistence[DummyInterfacable, string](),
	}
}

func (c *DummyInterfacableMemoryPersistence) GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[DummyInterfacable], error) {

	var key string

	if _key, ok := filter.GetAsNullableString("Key"); ok {
		key = _key
	}

	return c.IdentifiableMemoryPersistence.
		GetPageByFilter(ctx,
			func(item DummyInterfacable) bool {
				if key != "" && item.Key != key {
					return false
				}
				return true
			},
			paging,
			func(a, b DummyInterfacable) bool {
				return len(a.Key) < len(b.Key)
			},
			nil,
		)
}

func (c *DummyInterfacableMemoryPersistence) GetCountByFilter(ctx context.Context, filter cquery.FilterParams) (count int64, err error) {

	var key string

	if _key, ok := filter.GetAsNullableString("Key"); ok {
		key = _key
	}

	return c.IdentifiableMemoryPersistence.
		GetCountByFilter(ctx,
			func(item DummyInterfacable) bool {
				if key != "" && item.Key != key {
					return false
				}
				return true
			})
}
