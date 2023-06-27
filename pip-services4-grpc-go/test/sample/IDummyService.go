package test_sample

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDummyService interface {
	GetPageByFilter(ctx context.Context, filter *cquery.FilterParams, paging *cquery.PagingParams) (result *DummyDataPage, err error)
	GetOneById(ctx context.Context, id string) (result *Dummy, err error)
	Create(ctx context.Context, entity Dummy) (result *Dummy, err error)
	Update(ctx context.Context, entity Dummy) (result *Dummy, err error)
	DeleteById(ctx context.Context, id string) (result *Dummy, err error)
}
