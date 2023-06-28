package test_clients

import (
	"context"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	grpcclients "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/clients"
	testproto "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/protos"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/sample"
)

type DummyGrpcClient struct {
	*grpcclients.GrpcClient
}

func NewDummyGrpcClient() *DummyGrpcClient {
	dgc := DummyGrpcClient{}
	dgc.GrpcClient = grpcclients.NewGrpcClient("dummies.Dummies")
	return &dgc
}

func (c *DummyGrpcClient) GetDummies(ctx context.Context, filter *cquery.FilterParams, paging *cquery.PagingParams) (result *tsample.DummyDataPage, err error) {

	req := &testproto.DummiesPageRequest{
		TraceId: cctx.GetTraceId(ctx),
	}
	if filter != nil {
		req.Filter = filter.Value()
	}
	if paging != nil {
		req.Paging = &testproto.PagingParams{
			Skip:  paging.GetSkip(0),
			Take:  paging.GetTake(100),
			Total: paging.Total,
		}
	}
	reply := new(testproto.DummiesPage)
	err = c.CallWithContext(ctx, "get_dummies", req, reply)
	c.Instrument(ctx, "dummy.get_page_by_filter")
	if err != nil {
		return nil, err
	}
	result = toDummiesPage(reply)
	return result, nil
}

func (c *DummyGrpcClient) GetDummyById(ctx context.Context, dummyId string) (result *tsample.Dummy, err error) {

	req := &testproto.DummyIdRequest{
		TraceId: cctx.GetTraceId(ctx),
		DummyId: dummyId,
	}

	reply := new(testproto.Dummy)
	err = c.CallWithContext(ctx, "get_dummy_by_id", req, reply)
	c.Instrument(ctx, "dummy.get_one_by_id")
	if err != nil {
		return nil, err
	}
	result = toDummy(reply)
	if result != nil && result.Id == "" && result.Key == "" {
		result = nil
	}
	return result, nil
}

func (c *DummyGrpcClient) CreateDummy(ctx context.Context, dummy tsample.Dummy) (result *tsample.Dummy, err error) {

	req := &testproto.DummyObjectRequest{
		TraceId: cctx.GetTraceId(ctx),
		Dummy:   fromDummy(&dummy),
	}

	reply := new(testproto.Dummy)
	err = c.CallWithContext(ctx, "create_dummy", req, reply)
	c.Instrument(ctx, "dummy.create")
	if err != nil {
		return nil, err
	}
	result = toDummy(reply)
	if result != nil && result.Id == "" && result.Key == "" {
		result = nil
	}
	return result, nil
}

func (c *DummyGrpcClient) UpdateDummy(ctx context.Context, dummy tsample.Dummy) (result *tsample.Dummy, err error) {
	req := &testproto.DummyObjectRequest{
		TraceId: cctx.GetTraceId(ctx),
		Dummy:   fromDummy(&dummy),
	}
	reply := new(testproto.Dummy)
	err = c.CallWithContext(ctx, "update_dummy", req, reply)
	c.Instrument(ctx, "dummy.update")
	if err != nil {
		return nil, err
	}
	result = toDummy(reply)
	if result != nil && result.Id == "" && result.Key == "" {
		result = nil
	}
	return result, nil
}

func (c *DummyGrpcClient) DeleteDummy(ctx context.Context, dummyId string) (result *tsample.Dummy, err error) {

	req := &testproto.DummyIdRequest{
		TraceId: cctx.GetTraceId(ctx),
		DummyId: dummyId,
	}

	reply := new(testproto.Dummy)
	c.CallWithContext(ctx, "delete_dummy_by_id", req, reply)
	c.Instrument(ctx, "dummy.delete_by_id")
	if err != nil {
		return nil, err
	}
	result = toDummy(reply)
	if result != nil && result.Id == "" && result.Key == "" {
		result = nil
	}
	return result, nil
}
