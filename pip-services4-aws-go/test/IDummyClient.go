package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDummyClient interface {
	GetDummies(ctx context.Context, filter *cquery.FilterParams, paging *cquery.PagingParams) (result *cquery.DataPage[Dummy], err error)
	GetDummyById(ctx context.Context, dummyId string) (result *Dummy, err error)
	CreateDummy(ctx context.Context, dummy Dummy) (result *Dummy, err error)
	UpdateDummy(ctx context.Context, dummy Dummy) (result *Dummy, err error)
	DeleteDummy(ctx context.Context, dummyId string) (result *Dummy, err error)
}
