package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/test/fixtures"
)

type DummySqlServerPersistence struct {
	*persist.IdentifiableSqlServerPersistence[fixtures.Dummy, string]
}

func NewDummySqlServerPersistence() *DummySqlServerPersistence {
	c := &DummySqlServerPersistence{}
	c.IdentifiableSqlServerPersistence = persist.InheritIdentifiableSqlServerPersistence[fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummySqlServerPersistence) DefineSchema() {
	c.ClearSchema()
	c.EnsureSchema("CREATE TABLE [" + c.TableName + "] ([id] VARCHAR(32) PRIMARY KEY, [key] VARCHAR(50), [content] VARCHAR(MAX))")
	c.EnsureIndex(c.IdentifiableSqlServerPersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummySqlServerPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "[key]='" + key + "'"
	}
	sorting := ""

	return c.IdentifiableSqlServerPersistence.GetPageByFilter(ctx,
		filterObj, paging,
		sorting, "",
	)
}

func (c *DummySqlServerPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "[key]='" + key + "'"
	}
	return c.IdentifiableSqlServerPersistence.GetCountByFilter(ctx, filterObj)
}

func (c *DummySqlServerPersistence) GetOneRandom(ctx context.Context) (item fixtures.Dummy, err error) {
	return c.IdentifiableSqlServerPersistence.GetOneRandom(ctx, "")
}
