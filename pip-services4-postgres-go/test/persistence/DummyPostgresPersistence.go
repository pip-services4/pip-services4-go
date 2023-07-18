package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/test/fixtures"
)

type DummyPostgresPersistence struct {
	persist.IdentifiablePostgresPersistence[fixtures.Dummy, string]
}

func NewDummyPostgresPersistence() *DummyPostgresPersistence {
	c := &DummyPostgresPersistence{}
	c.IdentifiablePostgresPersistence = *persist.InheritIdentifiablePostgresPersistence[fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummyPostgresPersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiablePostgresPersistence.DefineSchema()
	// Row name must be in double quotes for properly case!!!
	c.EnsureSchema("CREATE TABLE " + c.QuotedTableName() + " (\"id\" TEXT PRIMARY KEY, \"key\" TEXT, \"content\" TEXT)")
	c.EnsureIndex(c.IdentifiablePostgresPersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyPostgresPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	sorting := ""

	return c.IdentifiablePostgresPersistence.GetPageByFilter(ctx,
		filterObj, paging,
		sorting, "",
	)
}

func (c *DummyPostgresPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	return c.IdentifiablePostgresPersistence.GetCountByFilter(ctx, filterObj)
}

func (c *DummyPostgresPersistence) GetOneRandom(ctx context.Context) (item fixtures.Dummy, err error) {
	return c.IdentifiablePostgresPersistence.GetOneRandom(ctx, "")
}
