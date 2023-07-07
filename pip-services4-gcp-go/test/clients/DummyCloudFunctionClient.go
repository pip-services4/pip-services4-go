package clients_test

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	gcpclient "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/clients"
	tdata "github.com/pip-services4/pip-services4-go/pip-services4-gcp-go/test/data"
	httpclients "github.com/pip-services4/pip-services4-go/pip-services4-http-go/clients"
)

type DummyCloudFunctionClient struct {
	*gcpclient.CloudFunctionClient
}

func NewDummyCloudFunctionClient() *DummyCloudFunctionClient {
	return &DummyCloudFunctionClient{
		CloudFunctionClient: gcpclient.NewCloudFunctionClient(),
	}
}

func (c *DummyCloudFunctionClient) GetDummies(ctx context.Context, filter cquery.FilterParams, paging cquery.PagingParams) (result cquery.DataPage[tdata.Dummy], err error) {
	timing := c.Instrument(ctx, "dummies.get_dummies")

	response, err := c.Call(ctx, "dummies.get_dummies", nil)
	if err != nil {
		return cquery.DataPage[tdata.Dummy]{}, err
	}

	defer timing.EndTiming(ctx, err)
	return httpclients.HandleHttpResponse[cquery.DataPage[tdata.Dummy]](response, cctx.GetTraceId(ctx))
}

func (c *DummyCloudFunctionClient) GetDummyById(ctx context.Context, dummyId string) (result tdata.Dummy, err error) {
	timing := c.Instrument(ctx, "dummies.get_dummy_by_id")

	response, err := c.Call(ctx, "dummies.get_dummy_by_id", cdata.NewAnyValueMapFromTuples("dummy_id", dummyId))
	if err != nil {
		return tdata.Dummy{}, err
	}

	defer timing.EndTiming(ctx, err)
	return httpclients.HandleHttpResponse[tdata.Dummy](response, cctx.GetTraceId(ctx))
}

func (c *DummyCloudFunctionClient) CreateDummy(ctx context.Context, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	timing := c.Instrument(ctx, "dummies.create_dummy")

	response, err := c.Call(ctx, "dummies.create_dummy", cdata.NewAnyValueMapFromTuples("dummy", dummy))
	if err != nil {
		return tdata.Dummy{}, err
	}

	defer timing.EndTiming(ctx, err)
	return httpclients.HandleHttpResponse[tdata.Dummy](response, cctx.GetTraceId(ctx))
}

func (c *DummyCloudFunctionClient) UpdateDummy(ctx context.Context, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	timing := c.Instrument(ctx, "dummies.update_dummy")

	response, err := c.Call(ctx, "dummies.update_dummy", cdata.NewAnyValueMapFromTuples("dummy", dummy))
	if err != nil {
		return tdata.Dummy{}, err
	}

	defer timing.EndTiming(ctx, err)
	return httpclients.HandleHttpResponse[tdata.Dummy](response, cctx.GetTraceId(ctx))
}

func (c *DummyCloudFunctionClient) DeleteDummy(ctx context.Context, dummyId string) (result tdata.Dummy, err error) {
	timing := c.Instrument(ctx, "dummies.delete_dummy")

	response, err := c.Call(ctx, "dummies.delete_dummy", cdata.NewAnyValueMapFromTuples("dummy_id", dummyId))
	if err != nil {
		return tdata.Dummy{}, err
	}

	defer timing.EndTiming(ctx, err)
	return httpclients.HandleHttpResponse[tdata.Dummy](response, cctx.GetTraceId(ctx))
}
