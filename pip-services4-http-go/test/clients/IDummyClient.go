package test_clients

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"
)

type IDummyClient interface {
	GetDummies(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (result cquery.DataPage[tsample.Dummy], err error)
	GetDummyById(ctx context.Context, dummyId string) (result tsample.Dummy, err error)
	CreateDummy(ctx context.Context, dummy tsample.Dummy) (result tsample.Dummy, err error)
	UpdateDummy(ctx context.Context, dummy tsample.Dummy) (result tsample.Dummy, err error)
	DeleteDummy(ctx context.Context, dummyId string) (result tsample.Dummy, err error)

	CheckTraceId(ctx context.Context) (result map[string]string, err error)

	CheckErrorPropagation(ctx context.Context) error
}
