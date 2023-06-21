package test_sample

import (
	"context"
	"encoding/json"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type DummyCommandSet struct {
	ccomand.CommandSet
	controller IDummyController
}

func NewDummyCommandSet(controller IDummyController) *DummyCommandSet {
	c := DummyCommandSet{}
	c.CommandSet = *ccomand.NewCommandSet()

	c.controller = controller

	c.AddCommand(c.makeGetPageByFilterCommand())
	c.AddCommand(c.makeGetOneByIdCommand())
	c.AddCommand(c.makeCreateCommand())
	c.AddCommand(c.makeUpdateCommand())
	c.AddCommand(c.makeDeleteByIdCommand())
	c.AddCommand(c.makeCheckTraceIdCommand())
	c.AddCommand(c.makeCheckErrorPropagationCommand())
	return &c
}

func (c *DummyCommandSet) makeGetPageByFilterCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummies",
		cvalid.NewObjectSchema().WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			var filter *cquery.FilterParams
			var paging *cquery.PagingParams

			if _val, ok := args.Get("filter"); ok {
				filter = cquery.NewFilterParamsFromValue(_val)
			}
			if _val, ok := args.Get("paging"); ok {
				paging = cquery.NewPagingParamsFromValue(_val)
			}

			return c.controller.GetPageByFilter(ctx, filter, paging)
		},
	)
}

func (c *DummyCommandSet) makeGetOneByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			id := args.GetAsString("dummy_id")
			return c.controller.GetOneById(ctx, id)
		},
	)
}

func (c *DummyCommandSet) makeCreateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"create_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", NewDummySchema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			var entity Dummy

			if _val, ok := args.Get("dummy"); ok {
				val, _ := json.Marshal(_val)
				json.Unmarshal(val, &entity)
			}

			return c.controller.Create(ctx, entity)
		},
	)
}

func (c *DummyCommandSet) makeUpdateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"update_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", NewDummySchema()),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			var entity Dummy

			if _val, ok := args.Get("dummy"); ok {
				val, _ := json.Marshal(_val)
				json.Unmarshal(val, &entity)
			}

			return c.controller.Update(ctx, entity)
		},
	)
}

func (c *DummyCommandSet) makeDeleteByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"delete_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			id := args.GetAsString("dummy_id")
			return c.controller.DeleteById(ctx, id)
		},
	)
}

func (c *DummyCommandSet) makeCheckTraceIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"check_trace_id",
		cvalid.NewObjectSchema(),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			return c.controller.CheckTraceId(ctx)
		},
	)
}

func (c *DummyCommandSet) makeCheckErrorPropagationCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"check_error_propagation",
		cvalid.NewObjectSchema(),
		func(ctx context.Context, args *exec.Parameters) (result any, err error) {
			return nil, c.controller.CheckErrorPropagation(ctx)
		},
	)
}
