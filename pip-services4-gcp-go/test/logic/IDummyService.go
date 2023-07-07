package test_logic

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/data"
)

type IDummyService interface {
	GetPageByFilter(ctx context.Context, filter *cquery.FilterParams, paging *cquery.PagingParams) (result *cquery.DataPage[tdata.Dummy], err error)
	GetOneById(ctx context.Context, id string) (result tdata.Dummy, err error)
	Create(ctx context.Context, entity tdata.Dummy) (result tdata.Dummy, err error)
	Update(ctx context.Context, entity tdata.Dummy) (result tdata.Dummy, err error)
	DeleteById(ctx context.Context, id string) (result tdata.Dummy, err error)
}
