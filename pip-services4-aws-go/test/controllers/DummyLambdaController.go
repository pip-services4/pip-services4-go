package test_controllers

import (
	"context"
	"encoding/json"

	awsserv "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/controllers"
	awstest "github.com/pip-services4/pip-services4-go/pip-services4-aws-go/test"
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type DummyLambdaController struct {
	*awsserv.LambdaController
	controller awstest.IDummyService
}

func NewDummyLambdaController() *DummyLambdaController {
	c := &DummyLambdaController{}
	c.LambdaController = awsserv.InheritLambdaController(c, "dummy")

	c.DependencyResolver.Put(context.Background(), "service", cref.NewDescriptor("pip-services-dummies", "service", "default", "*", "*"))
	return c
}

func (c *DummyLambdaController) SetReferences(ctx context.Context, references cref.IReferences) {
	c.LambdaController.SetReferences(ctx, references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("service")
	if depErr == nil && depRes != nil {
		c.controller = depRes.(awstest.IDummyService)
	}
}

func (c *DummyLambdaController) getCtxWithTraceId(ctx context.Context, params map[string]any) context.Context {
	traceId, _ := params["trace_id"].(string)
	return utils.ContextHelper.NewContextWithTraceId(ctx, traceId)
}

func (c *DummyLambdaController) getPageByFilter(ctx context.Context, params map[string]any) (any, error) {
	return c.controller.GetPageByFilter(c.getCtxWithTraceId(ctx, params),
		cquery.NewFilterParamsFromValue(params["filter"]),
		cquery.NewPagingParamsFromValue(params["paging"]),
	)
}

func (c *DummyLambdaController) getOneById(ctx context.Context, params map[string]any) (any, error) {
	return c.controller.GetOneById(c.getCtxWithTraceId(ctx, params),
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaController) create(ctx context.Context, params map[string]any) (any, error) {
	val, _ := json.Marshal(params["dummy"])
	var entity awstest.Dummy
	json.Unmarshal(val, &entity)
	return c.controller.Create(c.getCtxWithTraceId(ctx, params),
		entity,
	)
}

func (c *DummyLambdaController) update(ctx context.Context, params map[string]any) (any, error) {
	val, _ := json.Marshal(params["dummy"])
	var entity awstest.Dummy
	json.Unmarshal(val, &entity)
	return c.controller.Update(c.getCtxWithTraceId(ctx, params),
		entity,
	)
}

func (c *DummyLambdaController) deleteById(ctx context.Context, params map[string]any) (any, error) {
	return c.controller.DeleteById(c.getCtxWithTraceId(ctx, params),
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaController) Register() {

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
