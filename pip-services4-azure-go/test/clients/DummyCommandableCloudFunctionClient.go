package clients_test

import (
	"context"

	azureclient "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/clients"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-azure-go/test/data"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	httpclient "github.com/pip-services4/pip-services4-go/pip-services4-http-go/clients"
)

type DummyCommandableAzureFunctionClient struct {
	*azureclient.CommandableAzureFunctionClient
}

func NewDummyCommandableAzureFunctionClient() *DummyCommandableAzureFunctionClient {
	return &DummyCommandableAzureFunctionClient{
		CommandableAzureFunctionClient: azureclient.NewCommandableAzureFunctionClient("dummies"),
	}
}

func (c *DummyCommandableAzureFunctionClient) GetDummies(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (result cquery.DataPage[tdata.Dummy], err error) {
	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, &filter)
	c.AddPagingParams(params, &paging)

	response, err := c.CallCommand(ctx, "dummies.get_dummies", cdata.NewAnyValueMapFromValue(params.Value()))
	if err != nil {
		return *cquery.NewEmptyDataPage[tdata.Dummy](), err
	}

	return httpclient.HandleHttpResponse[cquery.DataPage[tdata.Dummy]](response, utils.ContextHelper.GetTraceId(ctx))
}

func (c *DummyCommandableAzureFunctionClient) GetDummyById(ctx context.Context, dummyId string) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "dummies.get_dummy_by_id", params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return httpclient.HandleHttpResponse[tdata.Dummy](response, utils.ContextHelper.GetTraceId(ctx))
}

func (c *DummyCommandableAzureFunctionClient) CreateDummy(ctx context.Context, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "dummies.create_dummy", params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return httpclient.HandleHttpResponse[tdata.Dummy](response, utils.ContextHelper.GetTraceId(ctx))
}

func (c *DummyCommandableAzureFunctionClient) UpdateDummy(ctx context.Context, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "dummies.update_dummy", params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return httpclient.HandleHttpResponse[tdata.Dummy](response, utils.ContextHelper.GetTraceId(ctx))
}

func (c *DummyCommandableAzureFunctionClient) DeleteDummy(ctx context.Context, dummyId string) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "dummies.delete_dummy", params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return httpclient.HandleHttpResponse[tdata.Dummy](response, utils.ContextHelper.GetTraceId(ctx))
}
