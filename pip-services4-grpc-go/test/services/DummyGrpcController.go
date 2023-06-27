package test_services

import (
	"context"
	"encoding/json"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	grpcservices "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/controllers"
	tsample "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/sample"
	"google.golang.org/grpc"

	grpcproto "github.com/pip-services4/pip-services4-go/pip-services4-grpc-go/test/protos"
)

type DummyGrpcController struct {
	*grpcservices.GrpcController
	service       tsample.IDummyService
	numberOfCalls int64

	grpcproto.DummiesServer
}

func NewDummyGrpcController() *DummyGrpcController {
	c := &DummyGrpcController{}
	c.GrpcController = grpcservices.InheritGrpcService(c, "dummies.Dummies")
	c.numberOfCalls = 0
	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return c
}

func (c *DummyGrpcController) SetReferences(ctx context.Context, references cref.IReferences) {
	c.GrpcController.SetReferences(ctx, references)
	resolv, err := c.DependencyResolver.GetOneRequired("service")
	if err == nil && resolv != nil {
		c.service = resolv.(tsample.IDummyService)
		return
	}
	panic("Can't resolve 'service' reference")
}

func (c *DummyGrpcController) GetNumberOfCalls() int64 {
	return c.numberOfCalls
}

func (c *DummyGrpcController) incrementNumberOfCalls(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	m, err := handler(ctx, req)
	if err != nil {
		c.Logger.Error(ctx, err, "RPC failed with error %v", err.Error())
	}
	c.numberOfCalls++
	return m, err
}

func (c *DummyGrpcController) Open(ctx context.Context) error {

	// Add interceptors
	c.RegisterUnaryInterceptor(c.incrementNumberOfCalls)
	return c.GrpcController.Open(ctx)
}

func (c *DummyGrpcController) GetDummies(ctx context.Context, req *grpcproto.DummiesPageRequest) (*grpcproto.DummiesPage, error) {

	validateErr := c.ValidateRequest(req,
		cvalid.NewObjectSchema().
			WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()).
			WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).Schema)

	if validateErr != nil {
		return nil, validateErr
	}

	filter := cquery.NewFilterParamsFromValue(req.GetFilter())
	paging := cquery.NewEmptyPagingParams()
	if req.Paging != nil {
		paging = cquery.NewPagingParams(req.Paging.GetSkip(), req.Paging.GetTake(), req.Paging.GetTotal())
	}
	data, err := c.service.GetPageByFilter(
		utils.ContextHelper.NewContextWithTraceId(context.Background(), req.TraceId),
		filter,
		paging,
	)
	if err != nil || data == nil {
		return nil, err
	}

	result := grpcproto.DummiesPage{}
	result.Total = *data.Total
	for _, v := range data.Data {
		buf := grpcproto.Dummy{}
		bytes, _ := json.Marshal(v)
		json.Unmarshal(bytes, &buf)
		result.Data = append(result.Data, &buf)
	}
	return &result, err
}

func (c *DummyGrpcController) GetDummyById(ctx context.Context, req *grpcproto.DummyIdRequest) (*grpcproto.Dummy, error) {

	// validation
	validateErr := c.ValidateRequest(req,
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema)

	if validateErr != nil {
		return nil, validateErr
	}
	// ==================================

	data, err := c.service.GetOneById(
		utils.ContextHelper.NewContextWithTraceId(context.Background(), req.TraceId),
		req.DummyId,
	)
	if err != nil {
		return nil, err
	}
	result := grpcproto.Dummy{}
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &result)
	return &result, nil
}

func (c *DummyGrpcController) CreateDummy(ctx context.Context, req *grpcproto.DummyObjectRequest) (*grpcproto.Dummy, error) {

	// validation
	validateErr := c.ValidateRequest(req,
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", tsample.NewDummySchema()).Schema)

	if validateErr != nil {
		return nil, validateErr
	}

	dummy := tsample.Dummy{}
	bytes, _ := json.Marshal(req.Dummy)
	json.Unmarshal(bytes, &dummy)

	data, err := c.service.Create(
		utils.ContextHelper.NewContextWithTraceId(context.Background(), req.TraceId),
		dummy,
	)

	if err != nil || data == nil {
		return nil, err
	}
	result := grpcproto.Dummy{}
	bytes, _ = json.Marshal(data)
	json.Unmarshal(bytes, &result)
	return &result, nil
}

func (c *DummyGrpcController) UpdateDummy(ctx context.Context, req *grpcproto.DummyObjectRequest) (*grpcproto.Dummy, error) {

	validateErr := c.ValidateRequest(req,
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", tsample.NewDummySchema()).Schema)

	if validateErr != nil {
		return nil, validateErr
	}

	dummy := tsample.Dummy{}
	bytes, _ := json.Marshal(req.Dummy)
	json.Unmarshal(bytes, &dummy)

	data, err := c.service.Update(
		utils.ContextHelper.NewContextWithTraceId(context.Background(), req.TraceId),
		dummy,
	)

	if err != nil || data == nil {
		return nil, err
	}
	result := grpcproto.Dummy{}
	bytes, _ = json.Marshal(data)
	json.Unmarshal(bytes, &result)
	return &result, nil
}

func (c *DummyGrpcController) DeleteDummyById(ctx context.Context, req *grpcproto.DummyIdRequest) (*grpcproto.Dummy, error) {

	validateErr := c.ValidateRequest(req,
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy_id", cconv.String).Schema)

	if validateErr != nil {
		return nil, validateErr
	}

	data, err := c.service.DeleteById(
		utils.ContextHelper.NewContextWithTraceId(context.Background(), req.TraceId),
		req.DummyId,
	)
	if err != nil || data == nil {
		return nil, err
	}
	result := grpcproto.Dummy{}
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &result)
	return &result, nil
}

func (c *DummyGrpcController) Register() {
	grpcproto.RegisterDummiesServer(c.Endpoint.GetServer(), c)
}
