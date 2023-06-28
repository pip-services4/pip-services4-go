package test_container

import (
	"context"
	"encoding/json"

	awscont "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/containers"
	awstest "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type DummyLambdaFunction struct {
	*awscont.LambdaFunction
	service awstest.IDummyService
}

func NewDummyLambdaFunction() *DummyLambdaFunction {
	c := &DummyLambdaFunction{}
	c.LambdaFunction = awscont.InheriteLambdaFunction(c, "dummy", "Dummy lambda function")

	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	c.AddFactory(awstest.NewDummyFactory())
	return c
}

func (c *DummyLambdaFunction) SetReferences(ctx context.Context, references cref.IReferences) {
	c.LambdaFunction.SetReferences(ctx, references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr == nil && depRes != nil {
		c.service = depRes.(awstest.IDummyService)
	}
}

func (c *DummyLambdaFunction) getCtxWithTraceId(ctx context.Context, params map[string]any) context.Context {
	traceId, _ := params["trace_id"].(string)
	return cctx.NewContextWithTraceId(ctx, traceId)
}

func (c *DummyLambdaFunction) getPageByFilter(ctx context.Context, params map[string]any) (any, error) {
	return c.service.GetPageByFilter(
		c.getCtxWithTraceId(ctx, params),
		cquery.NewFilterParamsFromValue(params["filter"]),
		cquery.NewPagingParamsFromValue(params["paging"]),
	)
}

func (c *DummyLambdaFunction) getOneById(ctx context.Context, params map[string]any) (any, error) {
	return c.service.GetOneById(
		c.getCtxWithTraceId(ctx, params),
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaFunction) create(ctx context.Context, params map[string]any) (any, error) {
	ctx = c.getCtxWithTraceId(ctx, params)
	val, _ := json.Marshal(params["dummy"])
	var entity = awstest.Dummy{}
	json.Unmarshal(val, &entity)

	c.Logger().Debug(ctx, "Create method called Dummy %v", entity)

	res, err := c.service.Create(
		c.getCtxWithTraceId(ctx, params),
		entity,
	)

	c.Logger().Debug(ctx, "Create method called Result: %v Err: %v", res, err)

	return res, err
}

func (c *DummyLambdaFunction) update(ctx context.Context, params map[string]any) (any, error) {
	val, _ := json.Marshal(params["dummy"])
	var entity = awstest.Dummy{}
	json.Unmarshal(val, &entity)
	return c.service.Update(
		c.getCtxWithTraceId(ctx, params),
		entity,
	)
}

func (c *DummyLambdaFunction) deleteById(ctx context.Context, params map[string]any) (any, error) {
	ctx = c.getCtxWithTraceId(ctx, params)

	c.Logger().Debug(ctx, "DeleteById method called Id %v", params["dummy_id"].(string))

	res, err := c.service.DeleteById(
		ctx,
		params["dummy_id"].(string),
	)
	c.Logger().Debug(ctx, "DeleteById method called Result: %v Err: %v", res, err)

	return res, err
}

func (c *DummyLambdaFunction) Register() {

	c.RegisterAction(
		"get_dummies",
		cvalid.NewObjectSchema().
			WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).
			WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()).Schema,
		c.getPageByFilter)

	c.RegisterAction(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().
			WithOptionalProperty("dummy_id", cconv.String).Schema,
		c.getOneById)

	c.RegisterAction(
		"create_dummy",
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", awstest.NewDummySchema()).Schema,
		c.create)

	c.RegisterAction(
		"update_dummy",
		cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", awstest.NewDummySchema()).Schema,
		c.update)

	c.RegisterAction(
		"delete_dummy",
		cvalid.NewObjectSchema().
			WithOptionalProperty("dummy_id", cconv.String).Schema,
		c.deleteById)
}
