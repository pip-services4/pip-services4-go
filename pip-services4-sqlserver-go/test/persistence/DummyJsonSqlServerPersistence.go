package test

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/test/fixtures"
)

type DummyJsonSqlServerPersistence struct {
	*persist.IdentifiableJsonSqlServerPersistence[fixtures.Dummy, string]
}

func NewDummyJsonSqlServerPersistence() *DummyJsonSqlServerPersistence {
	c := &DummyJsonSqlServerPersistence{}
	c.IdentifiableJsonSqlServerPersistence = persist.InheritIdentifiableJsonSqlServerPersistence[fixtures.Dummy, string](c, "dummies_json")
	return c
}

func (c *DummyJsonSqlServerPersistence) DefineSchema() {
	c.ClearSchema()
	c.EnsureTable("", "")
	c.EnsureSchema("ALTER TABLE [" + c.TableName + "] ADD [data_key] AS JSON_VALUE([data],'$.key')")
	c.EnsureIndex(c.TableName+"_key", map[string]string{"data_key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyJsonSqlServerPersistence) GetPageByFilter(ctx context.Context,
	filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "JSON_VALUE([data],'$.key')='" + key + "'"
	}

	return c.IdentifiableJsonSqlServerPersistence.GetPageByFilter(ctx,
		filterObj, paging,
		"", "",
	)
}

func (c *DummyJsonSqlServerPersistence) GetCountByFilter(ctx context.Context,
	filter cdata.FilterParams) (count int64, err error) {

	filterObj := ""
	if key, ok := filter.GetAsNullableString("Key"); ok && key != "" {
		filterObj += "JSON_VALUE([data],'$.key')='" + key + "'"
	}

	return c.IdentifiableJsonSqlServerPersistence.GetCountByFilter(ctx, filterObj)
}

func (c *DummyJsonSqlServerPersistence) GetOneRandom(ctx context.Context) (item fixtures.Dummy, err error) {
	return c.IdentifiableJsonSqlServerPersistence.GetOneRandom(ctx, "")
}
