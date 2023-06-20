package test_persistence

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

type IDummyPersistence interface {
	GetPageByFilter(ctx context.Context, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[Dummy], err error)
	GetListByIds(ctx context.Context, ids []string) (items []Dummy, err error)
	GetOneById(ctx context.Context, id string) (item Dummy, err error)
	Create(ctx context.Context, item Dummy) (result Dummy, err error)
	Update(ctx context.Context, item Dummy) (result Dummy, err error)
	UpdatePartially(ctx context.Context, id string, data cdata.AnyValueMap) (item Dummy, err error)
	DeleteById(ctx context.Context, id string) (item Dummy, err error)
	DeleteByIds(ctx context.Context, ids []string) (err error)
	GetCountByFilter(ctx context.Context, filter cdata.FilterParams) (count int64, err error)
}
