package example_logic

import (
	"context"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cexec "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
	data "github.com/pip-services4/pip-services4-go/pip-services4-swagger-go/example/data"
)

type DummyCommandSet struct {
	ccomand.CommandSet
	controller IDummyService
}

func NewDummyCommandSet(controller IDummyService) *DummyCommandSet {
	c := DummyCommandSet{}
	c.CommandSet = *ccomand.NewCommandSet()

	c.controller = controller

	c.AddCommand(c.makeGetPageByFilterCommand())
	c.AddCommand(c.makeGetOneByIdCommand())
	c.AddCommand(c.makeCreateCommand())
	c.AddCommand(c.makeUpdateCommand())
	c.AddCommand(c.makeDeleteByIdCommand())
	return &c
}

func (c *DummyCommandSet) makeGetPageByFilterCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummies",
		cvalid.NewObjectSchema().WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()),
		func(ctx context.Context, args *cexec.Parameters) (result any, err error) {
			var filter *cquery.FilterParams
			var paging *cquery.PagingParams

			if data, contains := args.Get("filter"); contains {
				filter = cquery.NewFilterParamsFromValue(data)
			}

			if data, contains := args.Get("paging"); contains {
				paging = cquery.NewPagingParamsFromValue(data)
			}

			return c.controller.GetPageByFilter(ctx, filter, paging)
		},
	)
}

func (c *DummyCommandSet) makeGetOneByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(ctx context.Context, args *cexec.Parameters) (result any, err error) {
			id := args.GetAsString("dummy_id")
			return c.controller.GetOneById(ctx, id)
		},
	)
}

func (c *DummyCommandSet) makeCreateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"create_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", data.NewDummySchema()),
		func(ctx context.Context, args *cexec.Parameters) (result any, err error) {
			var entity data.Dummy

			if _val, ok := args.Get("dummy"); ok {
				val, _ := cconv.JsonConverter.ToJson(_val)
				obj, err := cconv.JsonConverter.FromJson(val)

				if err != nil {
					return nil, err
				}

				entity = obj.(data.Dummy)
			}

			return c.controller.Create(ctx, entity)
		},
	)
}

func (c *DummyCommandSet) makeUpdateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"update_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", data.NewDummySchema()),
		func(ctx context.Context, args *cexec.Parameters) (result any, err error) {
			var entity data.Dummy

			if _val, ok := args.Get("dummy"); ok {
				val, _ := cconv.JsonConverter.ToJson(_val)
				obj, err := cconv.JsonConverter.FromJson(val)

				if err != nil {
					return nil, err
				}

				entity = obj.(data.Dummy)
			}
			return c.controller.Update(ctx, entity)
		},
	)
}

func (c *DummyCommandSet) makeDeleteByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"delete_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(ctx context.Context, args *cexec.Parameters) (result any, err error) {
			id := args.GetAsString("dummy_id")
			return c.controller.DeleteById(ctx, id)
		},
	)
}
