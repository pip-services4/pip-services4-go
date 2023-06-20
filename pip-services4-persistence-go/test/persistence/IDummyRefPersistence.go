package test_persistence

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

type IDummyRefPersistence interface {
	GetPageByFilter(ctx context.Context, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[*DummyRef], err error)
	GetListByIds(ctx context.Context, ids []string) (items []*DummyRef, err error)
	GetOneById(ctx context.Context, id string) (item *DummyRef, err error)
	Create(ctx context.Context, item *DummyRef) (result *DummyRef, err error)
	Update(ctx context.Context, item *DummyRef) (result *DummyRef, err error)
	UpdatePartially(ctx context.Context, id string, data cdata.AnyValueMap) (item *DummyRef, err error)
	DeleteById(ctx context.Context, id string) (item *DummyRef, err error)
	DeleteByIds(ctx context.Context, ids []string) (err error)
	GetCountByFilter(ctx context.Context, filter cdata.FilterParams) (count int64, err error)
}
