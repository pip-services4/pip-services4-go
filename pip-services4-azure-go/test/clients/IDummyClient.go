package clients_test

import (
	"context"

	tdata "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/test/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type IDummyClient interface {
	GetDummies(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (result cquery.DataPage[tdata.Dummy], err error)
	GetDummyById(ctx context.Context, dummyId string) (result tdata.Dummy, err error)
	CreateDummy(ctx context.Context, dummy tdata.Dummy) (result tdata.Dummy, err error)
	UpdateDummy(ctx context.Context, dummy tdata.Dummy) (result tdata.Dummy, err error)
	DeleteDummy(ctx context.Context, dummyId string) (result tdata.Dummy, err error)
}
