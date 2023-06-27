package test_clients

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/pip-services4/pip-services4-go/pip-services4-http-go/clients"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"
)

type DummyCommandableHttpClient struct {
	clients.CommandableHttpClient
}

func NewDummyCommandableHttpClient() *DummyCommandableHttpClient {
	dchc := DummyCommandableHttpClient{}
	dchc.CommandableHttpClient = *clients.NewCommandableHttpClient("dummies")
	return &dchc
}

func (c *DummyCommandableHttpClient) GetDummies(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (result cquery.DataPage[tsample.Dummy], err error) {
	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, &filter)
	c.AddPagingParams(params, &paging)

	response, err := c.CallCommand(ctx, "get_dummies", cdata.NewAnyValueMapFromValue(params.Value()))
	if err != nil {
		return *cquery.NewEmptyDataPage[tsample.Dummy](), err
	}

	return clients.HandleHttpResponse[cquery.DataPage[tsample.Dummy]](response, utils.ContextHelper.GetTraceId(ctx))
}

func (c *DummyCommandableHttpClient) GetDummyById(ctx context.Context, dummyId string) (tsample.Dummy, error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "get_dummy_by_id", params)
	if err != nil {
		return tsample.Dummy{}, err
	}

	return clients.HandleHttpResponse[tsample.Dummy](response, utils.ContextHelper.GetClient(ctx))
}

func (c *DummyCommandableHttpClient) CreateDummy(ctx context.Context, dummy tsample.Dummy) (result tsample.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "create_dummy", params)
	if err != nil {
		return tsample.Dummy{}, err
	}

	return clients.HandleHttpResponse[tsample.Dummy](response, utils.ContextHelper.GetClient(ctx))
}

func (c *DummyCommandableHttpClient) UpdateDummy(ctx context.Context, dummy tsample.Dummy) (result tsample.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "update_dummy", params)
	if err != nil {
		return tsample.Dummy{}, err
	}

	return clients.HandleHttpResponse[tsample.Dummy](response, utils.ContextHelper.GetClient(ctx))
}

func (c *DummyCommandableHttpClient) DeleteDummy(ctx context.Context, dummyId string) (result tsample.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "delete_dummy", params)
	if err != nil {
		return tsample.Dummy{}, err
	}

	return clients.HandleHttpResponse[tsample.Dummy](response, utils.ContextHelper.GetClient(ctx))
}

func (c *DummyCommandableHttpClient) CheckTraceId(ctx context.Context) (result map[string]string, err error) {

	params := cdata.NewEmptyAnyValueMap()

	response, err := c.CallCommand(ctx, "check_trace_id", params)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[map[string]string](response, utils.ContextHelper.GetClient(ctx))
}

func (c *DummyCommandableHttpClient) CheckErrorPropagation(ctx context.Context) error {
	params := cdata.NewEmptyAnyValueMap()
	_, calErr := c.CallCommand(ctx, "check_error_propagation", params)
	return calErr
}
