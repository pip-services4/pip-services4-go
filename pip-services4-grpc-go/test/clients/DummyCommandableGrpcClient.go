package test_clients

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	grpcclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/sample"
)

type DummyCommandableGrpcClient struct {
	*grpcclients.CommandableGrpcClient
}

func NewDummyCommandableGrpcClient() *DummyCommandableGrpcClient {
	dcgc := DummyCommandableGrpcClient{}
	dcgc.CommandableGrpcClient = grpcclients.NewCommandableGrpcClient("dummy")
	return &dcgc
}

func (c *DummyCommandableGrpcClient) GetDummies(ctx context.Context, filter *cquery.FilterParams, paging *cquery.PagingParams) (result *tsample.DummyDataPage, err error) {

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, filter)
	c.AddPagingParams(params, paging)

	response, calErr := c.CallCommand(ctx, "get_dummies", cdata.NewAnyValueMapFromValue(params.Value()))
	if calErr != nil {
		return nil, calErr
	}

	return grpcclients.HandleHttpResponse[*tsample.DummyDataPage](response)
}

func (c *DummyCommandableGrpcClient) GetDummyById(ctx context.Context, dummyId string) (result *tsample.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, calErr := c.CallCommand(ctx, "get_dummy_by_id", params)
	if calErr != nil {
		return nil, calErr
	}

	return grpcclients.HandleHttpResponse[*tsample.Dummy](response)
}

func (c *DummyCommandableGrpcClient) CreateDummy(ctx context.Context, dummy tsample.Dummy) (result *tsample.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, calErr := c.CallCommand(ctx, "create_dummy", params)
	if calErr != nil {
		return nil, calErr
	}

	return grpcclients.HandleHttpResponse[*tsample.Dummy](response)
}

func (c *DummyCommandableGrpcClient) UpdateDummy(ctx context.Context, dummy tsample.Dummy) (result *tsample.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, calErr := c.CallCommand(ctx, "update_dummy", params)
	if calErr != nil {
		return nil, calErr
	}

	return grpcclients.HandleHttpResponse[*tsample.Dummy](response)
}

func (c *DummyCommandableGrpcClient) DeleteDummy(ctx context.Context, dummyId string) (result *tsample.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, calErr := c.CallCommand(ctx, "delete_dummy_by_id", params)
	if calErr != nil {
		return nil, calErr
	}

	return grpcclients.HandleHttpResponse[*tsample.Dummy](response)
}
