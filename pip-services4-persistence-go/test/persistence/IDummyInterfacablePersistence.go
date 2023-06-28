package test_persistence

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDummyInterfacablePersistence interface {
	GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[DummyInterfacable], err error)
	GetListByIds(ctx context.Context, ids []string) (items []DummyInterfacable, err error)
	GetOneById(ctx context.Context, id string) (item DummyInterfacable, err error)
	Create(ctx context.Context, item DummyInterfacable) (result DummyInterfacable, err error)
	Update(ctx context.Context, item DummyInterfacable) (result DummyInterfacable, err error)
	UpdatePartially(ctx context.Context, id string, data cdata.AnyValueMap) (item DummyInterfacable, err error)
	DeleteById(ctx context.Context, id string) (item DummyInterfacable, err error)
	DeleteByIds(ctx context.Context, ids []string) (err error)
	GetCountByFilter(ctx context.Context, filter cquery.FilterParams) (count int64, err error)
}
