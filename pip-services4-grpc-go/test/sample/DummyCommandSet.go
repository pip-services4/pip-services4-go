package test_sample

import (
	"context"
	"encoding/json"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	crun "github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	cdata "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	ccomand "github.com/pip-services4/pip-services4-go/pip-services4-rpc-go/commands"
)

type DummyCommandSet struct {
	ccomand.CommandSet
	controller IDummyService
}

func NewDummyCommandSet(controller IDummyService) *DummyCommandSet {
	dcs := DummyCommandSet{}
	dcs.CommandSet = *ccomand.NewCommandSet()

	dcs.controller = controller

	dcs.AddCommand(dcs.makeGetPageByFilterCommand())
	dcs.AddCommand(dcs.makeGetOneByIdCommand())
	dcs.AddCommand(dcs.makeCreateCommand())
	dcs.AddCommand(dcs.makeUpdateCommand())
	dcs.AddCommand(dcs.makeDeleteByIdCommand())
	return &dcs
}

func (c *DummyCommandSet) makeGetPageByFilterCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummies",
		cvalid.NewObjectSchema().WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()),
		func(ctx context.Context, args *crun.Parameters) (result any, err error) {
			var filter *cdata.FilterParams
			var paging *cdata.PagingParams

			if _val, ok := args.Get("filter"); ok {
				filter = cdata.NewFilterParamsFromValue(_val)
			}
			if _val, ok := args.Get("paging"); ok {
				paging = cdata.NewPagingParamsFromValue(_val)
			}

			return c.controller.GetPageByFilter(ctx, filter, paging)
		},
	)
}

func (c *DummyCommandSet) makeGetOneByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(ctx context.Context, args *crun.Parameters) (result any, err error) {
			id := args.GetAsString("dummy_id")
			return c.controller.GetOneById(ctx, id)
		},
	)
}

func (c *DummyCommandSet) makeCreateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"create_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", NewDummySchema()),
		func(ctx context.Context, args *crun.Parameters) (result any, err error) {
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
		func(ctx context.Context, args *crun.Parameters) (result any, err error) {
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
		"delete_dummy_by_id",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(ctx context.Context, args *crun.Parameters) (result any, err error) {
			id := args.GetAsString("dummy_id")
			return c.controller.DeleteById(ctx, id)
		},
	)
}
