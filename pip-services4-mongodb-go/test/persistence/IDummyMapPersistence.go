package test_persistence

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDummyMapPersistence interface {
	GetPageByFilter(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[map[string]any], err error)
	GetListByIds(ctx context.Context, ids []string) (items []map[string]any, err error)
	GetOneById(ctx context.Context, id string) (item map[string]any, err error)
	Create(ctx context.Context, item map[string]any) (result map[string]any, err error)
	Update(ctx context.Context, item map[string]any) (result map[string]any, err error)
	UpdatePartially(ctx context.Context, id string, data cdata.AnyValueMap) (item map[string]any, err error)
	DeleteById(ctx context.Context, id string) (item map[string]any, err error)
	DeleteByIds(ctx context.Context, ids []string) (err error)
	GetCountByFilter(ctx context.Context, filter cquery.FilterParams) (count int64, err error)
}
