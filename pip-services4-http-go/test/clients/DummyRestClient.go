package test_clients

import (
	"context"
	"net/http"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/pip-services4/pip-services4-go/pip-services4-http-go/clients"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-http-go/test/sample"
)

type DummyRestClient struct {
	clients.RestClient
}

func NewDummyRestClient() *DummyRestClient {
	drc := DummyRestClient{}
	drc.RestClient = *clients.NewRestClient()
	return &drc
}

func (c *DummyRestClient) GetDummies(ctx context.Context, filter cquery.FilterParams,
	paging cquery.PagingParams) (result cquery.DataPage[tsample.Dummy], err error) {

	defer c.Instrument(ctx, "dummy.get_page_by_filter")

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, &filter)
	c.AddPagingParams(params, &paging)

	response, err := c.Call(ctx, http.MethodGet, "/dummies", params, nil)
	if err != nil {
		return *cquery.NewEmptyDataPage[tsample.Dummy](), err
	}

	return clients.HandleHttpResponse[cquery.DataPage[tsample.Dummy]](response, utils.ContextHelper.GetClient(ctx))
}

func (c *DummyRestClient) GetDummyById(ctx context.Context, dummyId string) (result tsample.Dummy, err error) {
	defer c.Instrument(ctx, "dummy.get_one_by_id")

	response, err := c.Call(ctx, http.MethodGet, "/dummies/"+dummyId, nil, nil)
	if err != nil {
		return tsample.Dummy{}, err
	}

	return clients.HandleHttpResponse[tsample.Dummy](response, utils.ContextHelper.GetClient(ctx))
}

func (c *DummyRestClient) CreateDummy(ctx context.Context, dummy tsample.Dummy) (result tsample.Dummy, err error) {

	defer c.Instrument(ctx, "dummy.create")

	response, err := c.Call(ctx, http.MethodPost, "/dummies", nil, dummy)
	if err != nil {
		return tsample.Dummy{}, err
	}

	return clients.HandleHttpResponse[tsample.Dummy](response, utils.ContextHelper.GetClient(ctx))
}

func (c *DummyRestClient) UpdateDummy(ctx context.Context, dummy tsample.Dummy) (result tsample.Dummy, err error) {

	defer c.Instrument(ctx, "dummy.update")

	response, err := c.Call(ctx, http.MethodPut, "/dummies", nil, dummy)
	if err != nil {
		return tsample.Dummy{}, err
	}

	return clients.HandleHttpResponse[tsample.Dummy](response, utils.ContextHelper.GetTraceId(ctx))
}

func (c *DummyRestClient) DeleteDummy(ctx context.Context, dummyId string) (result tsample.Dummy, err error) {

	defer c.Instrument(ctx, "dummy.delete_by_id")

	response, err := c.Call(ctx, http.MethodDelete, "/dummies/"+dummyId, nil, nil)
	if err != nil {
		return tsample.Dummy{}, err
	}

	return clients.HandleHttpResponse[tsample.Dummy](response, utils.ContextHelper.GetTraceId(ctx))
}

func (c *DummyRestClient) CheckTraceId(ctx context.Context) (result map[string]string, err error) {

	defer c.Instrument(ctx, "dummy.check_trace_id")

	response, err := c.Call(ctx, http.MethodGet, "/dummies/check/trace_id", nil, nil)
	if err != nil {
		return nil, err
	}

	return clients.HandleHttpResponse[map[string]string](response, utils.ContextHelper.GetTraceId(ctx))
}

func (c *DummyRestClient) CheckErrorPropagation(ctx context.Context) error {

	c.Instrument(ctx, "dummy.check_error_propagation")

	_, err := c.Call(ctx, http.MethodGet, "/dummies/check/error_propagation", nil, nil)
	return err
}
