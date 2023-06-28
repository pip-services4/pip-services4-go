package test_persistence

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDummyMapPersistence interface {
	GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[DummyMap], err error)
	GetListByIds(ctx context.Context, ids []string) (items []DummyMap, err error)
	GetOneById(ctx context.Context, id string) (item DummyMap, err error)
	Create(ctx context.Context, item DummyMap) (result DummyMap, err error)
	Update(ctx context.Context, item DummyMap) (result DummyMap, err error)
	UpdatePartially(ctx context.Context, id string, data cdata.AnyValueMap) (item DummyMap, err error)
	DeleteById(ctx context.Context, id string) (item DummyMap, err error)
	DeleteByIds(ctx context.Context, ids []string) (err error)
	GetCountByFilter(ctx context.Context, filter cquery.FilterParams) (count int64, err error)
}
