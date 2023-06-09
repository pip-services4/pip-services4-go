package test_persistence

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
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

func (c *DummyMemoryPersistence) GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (cquery.DataPage[Dummy], error) {

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

func (c *DummyMemoryPersistence) GetCountByFilter(ctx context.Context, filter cquery.FilterParams) (count int64, err error) {

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
