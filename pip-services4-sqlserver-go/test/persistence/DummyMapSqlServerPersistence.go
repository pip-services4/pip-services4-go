package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/persistence"
)

type DummyMapSqlServerPersistence struct {
	*persist.IdentifiableSqlServerPersistence[map[string]any, string]
}

func NewDummyMapSqlServerPersistence() *DummyMapSqlServerPersistence {
	c := &DummyMapSqlServerPersistence{}
	c.IdentifiableSqlServerPersistence = persist.InheritIdentifiableSqlServerPersistence[map[string]any, string](c, "dummies")
	return c
}

func (c *DummyMapSqlServerPersistence) DefineSchema() {
	c.ClearSchema()
	c.EnsureSchema("CREATE TABLE [" + c.TableName + "] ([id] VARCHAR(32) PRIMARY KEY, [key] VARCHAR(50), [content] VARCHAR(MAX))")
	c.EnsureIndex(c.IdentifiableSqlServerPersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyMapSqlServerPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[map[string]any], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "[key]='" + key + "'"
	}
	sorting := ""

	return c.IdentifiableSqlServerPersistence.GetPageByFilter(ctx,
		filterObj, paging, sorting, "",
	)
}

func (c *DummyMapSqlServerPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "[key]='" + key + "'"
	}
	return c.IdentifiableSqlServerPersistence.GetCountByFilter(ctx, filterObj)
}
