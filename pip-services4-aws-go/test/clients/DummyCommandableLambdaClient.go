package test

import (
	"context"

	awsclient "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/clients"
	awstest "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
)

type DummyCommandableLambdaClient struct {
	*awsclient.CommandableLambdaClient
}

func NewDummyCommandableLambdaClient() *DummyCommandableLambdaClient {
	c := &DummyCommandableLambdaClient{
		CommandableLambdaClient: awsclient.NewCommandableLambdaClient("dummy"),
	}
	return c
}
func (c *DummyCommandableLambdaClient) GetDummies(ctx context.Context, filter *cquery.FilterParams,
	paging *cquery.PagingParams) (result *cquery.DataPage[awstest.Dummy], err error) {
	timing := c.Instrument(ctx, "dummy.get_dummies")

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("filter", filter.Value())
	params.SetAsObject("paging", paging)

	calValue, calErr := c.CallCommand(ctx, "get_dummies", params.Value())
	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)

	return awsclient.HandleLambdaResponse[*cquery.DataPage[awstest.Dummy]](calValue)
}

func (c *DummyCommandableLambdaClient) GetDummyById(ctx context.Context, dummyId string) (result *awstest.Dummy, err error) {
	timing := c.Instrument(ctx, "dummy.get_one_by_id")

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy_id", dummyId)

	calValue, calErr := c.CallCommand(ctx, "get_dummy_by_id", params.Value())

	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)

	return awsclient.HandleLambdaResponse[*awstest.Dummy](calValue)
}

func (c *DummyCommandableLambdaClient) CreateDummy(ctx context.Context, dummy awstest.Dummy) (result *awstest.Dummy, err error) {
	timing := c.Instrument(ctx, "dummy.create_dummy")
	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy", dummy)

	calValue, calErr := c.CallCommand(ctx, "create_dummy", params.Value())
	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)
	return awsclient.HandleLambdaResponse[*awstest.Dummy](calValue)
}

func (c *DummyCommandableLambdaClient) UpdateDummy(ctx context.Context, dummy awstest.Dummy) (result *awstest.Dummy, err error) {
	timing := c.Instrument(ctx, "dummy.update_dummy")
	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy", dummy)

	calValue, calErr := c.CallCommand(ctx, "update_dummy", params.Value())
	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)
	return awsclient.HandleLambdaResponse[*awstest.Dummy](calValue)
}

func (c *DummyCommandableLambdaClient) DeleteDummy(ctx context.Context, dummyId string) (result *awstest.Dummy, err error) {
	timing := c.Instrument(ctx, "dummy.delete_dummy")

	params := cdata.NewEmptyAnyValueMap()
	params.SetAsObject("dummy_id", dummyId)
	calValue, calErr := c.CallCommand(ctx, "delete_dummy", params.Value())
	if calErr != nil {
		return nil, calErr
	}

	defer timing.EndTiming(ctx, err)

	return awsclient.HandleLambdaResponse[*awstest.Dummy](calValue)
}
