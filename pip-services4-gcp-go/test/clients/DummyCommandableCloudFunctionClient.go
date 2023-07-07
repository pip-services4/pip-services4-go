package clients_test

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	gcpclient "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/clients"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/data"
	rpcclient "github.com/pip-services4/pip-services4-go/pip-services4-http-go/clients"
)

type DummyCommandableCloudFunctionClient struct {
	*gcpclient.CommandableCloudFunctionClient
}

func NewDummyCommandableCloudFunctionClient() *DummyCommandableCloudFunctionClient {
	return &DummyCommandableCloudFunctionClient{
		CommandableCloudFunctionClient: gcpclient.NewCommandableCloudFunctionClient("dummies"),
	}
}

func (c *DummyCommandableCloudFunctionClient) GetDummies(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (result cquery.DataPage[tdata.Dummy], err error) {
	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, &filter)
	c.AddPagingParams(params, &paging)

	response, err := c.CallCommand(ctx, "dummies.get_dummies", cdata.NewAnyValueMapFromValue(params.Value()))
	if err != nil {
		return *cquery.NewEmptyDataPage[tdata.Dummy](), err
	}

	return rpcclient.HandleHttpResponse[cquery.DataPage[tdata.Dummy]](response, cctx.GetTraceId(ctx))
}

func (c *DummyCommandableCloudFunctionClient) GetDummyById(ctx context.Context, dummyId string) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "dummies.get_dummy_by_id", params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return rpcclient.HandleHttpResponse[tdata.Dummy](response, cctx.GetTraceId(ctx))
}

func (c *DummyCommandableCloudFunctionClient) CreateDummy(ctx context.Context, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "dummies.create_dummy", params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return rpcclient.HandleHttpResponse[tdata.Dummy](response, cctx.GetTraceId(ctx))
}

func (c *DummyCommandableCloudFunctionClient) UpdateDummy(ctx context.Context, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "dummies.update_dummy", params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return rpcclient.HandleHttpResponse[tdata.Dummy](response, cctx.GetTraceId(ctx))
}

func (c *DummyCommandableCloudFunctionClient) DeleteDummy(ctx context.Context, dummyId string) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "dummies.delete_dummy", params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return rpcclient.HandleHttpResponse[tdata.Dummy](response, cctx.GetTraceId(ctx))
}
