package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/test/fixtures"
)

type DummyJsonPostgresPersistence struct {
	*persist.IdentifiableJsonPostgresPersistence[fixtures.Dummy, string]
}

func NewDummyJsonPostgresPersistence() *DummyJsonPostgresPersistence {
	c := &DummyJsonPostgresPersistence{}
	c.IdentifiableJsonPostgresPersistence = persist.InheritIdentifiableJsonPostgresPersistence[fixtures.Dummy, string](c, "dummies_json")
	return c
}

func (c *DummyJsonPostgresPersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableJsonPostgresPersistence.DefineSchema()
	c.EnsureTable("", "")
	c.EnsureIndex(c.TableName+"_key", map[string]string{"(data->'key')": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyJsonPostgresPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "data->>'key'='" + key + "'"
	}

	return c.IdentifiableJsonPostgresPersistence.GetPageByFilter(ctx,
		filterObj, paging,
		"", "",
	)
}

func (c *DummyJsonPostgresPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	filterObj := ""
	if key, ok := filter.GetAsNullableString("Key"); ok && key != "" {
		filterObj += "data->>'key'='" + key + "'"
	}

	return c.IdentifiableJsonPostgresPersistence.GetCountByFilter(ctx, filterObj)
}

func (c *DummyJsonPostgresPersistence) GetOneRandom(ctx context.Context) (item fixtures.Dummy, err error) {
	return c.IdentifiableJsonPostgresPersistence.GetOneRandom(ctx, "")
}
