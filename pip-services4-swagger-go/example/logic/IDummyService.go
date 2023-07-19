package example_logic

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	data "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/example/data"
)

type IDummyService interface {
	GetPageByFilter(ctx context.Context, filter *cquery.FilterParams, paging *cquery.PagingParams) (result *cquery.DataPage[data.Dummy], err error)
	GetOneById(ctx context.Context, id string) (result *data.Dummy, err error)
	Create(ctx context.Context, entity data.Dummy) (result *data.Dummy, err error)
	Update(ctx context.Context, entity data.Dummy) (result *data.Dummy, err error)
	DeleteById(ctx context.Context, id string) (result *data.Dummy, err error)
}
