package test_persistence

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cpersist "github.com/pip-services4/pip-services4-go/pip-services4-persistence-go/persistence"
)

type DummyMapMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence[DummyMap, string]
}

func NewDummyMapMemoryPersistence() *DummyMapMemoryPersistence {
	return &DummyMapMemoryPersistence{
		*cpersist.NewIdentifiableMemoryPersistence[DummyMap, string](),
	}
}

func filterFunc(filter cquery.FilterParams) func(item DummyMap) bool {

	var key string

	if _key, ok := filter.GetAsNullableString("Key"); ok {
		key = _key
	}

	return func(value DummyMap) bool {
		if _val, ok := value["Key"]; ok {
			if _key, ok := _val.(string); !ok && key != "" && _key != key {
				return false
			}
			return true
		}

		return false
	}
}

func sortFunc(a, b DummyMap) bool {
	_val, ok := a["Key"]
	if !ok {
		return false
	}
	_keyA, ok := _val.(string)
	if !ok {
		return false
	}

	_val, ok = b["Key"]
	if !ok {
		return false
	}
	_keyB, ok := _val.(string)
	if !ok {
		return false
	}

	return len(_keyA) < len(_keyB)
}

func (c *DummyMapMemoryPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (result cquery.DataPage[DummyMap], err error) {

	return c.IdentifiableMemoryPersistence.
		GetPageByFilter(ctx, filterFunc(filter), paging, sortFunc, nil)
}

func (c *DummyMapMemoryPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	return c.IdentifiableMemoryPersistence.
		GetCountByFilter(ctx, filterFunc(filter))
}
