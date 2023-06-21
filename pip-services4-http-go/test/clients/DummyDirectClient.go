package test_clients

import (
	"context"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/clients"
	test_sample "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/test/sample"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/test/sample"
)

type DummyDirectClient struct {
	clients.DirectClient
	specificController test_sample.IDummyController
}

func NewDummyDirectClient() *DummyDirectClient {
	ddc := DummyDirectClient{}
	ddc.DirectClient = *clients.NewDirectClient()
	ddc.DependencyResolver.Put(context.Background(), "controller", cref.NewDescriptor("pip-services-dummies", "controller", "*", "*", "*"))
	return &ddc
}

func (c *DummyDirectClient) SetReferences(ctx context.Context, references cref.IReferences) {
	c.DirectClient.SetReferences(ctx, references)

	specificController, ok := c.Controller.(test_sample.IDummyController)
	if !ok {
		panic("DummyDirectClient: Cant't resolv dependency 'controller' to IDummyController")
	}
	c.specificController = specificController

}

func (c *DummyDirectClient) GetDummies(ctx context.Context, filter cdata.FilterParams, paging cdata.PagingParams) (cdata.DataPage[tsample.Dummy], error) {

	timing := c.Instrument(ctx, "dummy.get_page_by_filter")
	result, err := c.specificController.GetPageByFilter(ctx, &filter, &paging)
	timing.EndTiming(ctx, err)
	return *result, err

}

func (c *DummyDirectClient) GetDummyById(ctx context.Context, dummyId string) (tsample.Dummy, error) {

	timing := c.Instrument(ctx, "dummy.get_one_by_id")
	result, err := c.specificController.GetOneById(ctx, dummyId)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) CreateDummy(ctx context.Context, dummy tsample.Dummy) (tsample.Dummy, error) {

	timing := c.Instrument(ctx, "dummy.create")
	result, err := c.specificController.Create(ctx, dummy)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) UpdateDummy(ctx context.Context, dummy tsample.Dummy) (tsample.Dummy, error) {

	timing := c.Instrument(ctx, "dummy.update")
	result, err := c.specificController.Update(ctx, dummy)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) DeleteDummy(ctx context.Context, dummyId string) (tsample.Dummy, error) {

	timing := c.Instrument(ctx, "dummy.delete_by_id")
	result, err := c.specificController.DeleteById(ctx, dummyId)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) CheckTraceId(ctx context.Context) (map[string]string, error) {

	timing := c.Instrument(ctx, "dummy.delete_by_id")
	result, err := c.specificController.CheckTraceId(ctx)
	timing.EndTiming(ctx, err)
	return result, err
}

func (c *DummyDirectClient) CheckErrorPropagation(ctx context.Context) error {
	timing := c.Instrument(ctx, "dummy.check_error_propagation")
	err := c.specificController.CheckErrorPropagation(ctx)
	timing.EndTiming(ctx, err)
	return err
}
