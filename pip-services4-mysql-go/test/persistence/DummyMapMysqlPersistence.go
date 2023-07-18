package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/persistence"
)

type DummyMapMySqlPersistence struct {
	persist.IdentifiableMySqlPersistence[map[string]any, string]
}

func NewDummyMapMySqlPersistence() *DummyMapMySqlPersistence {
	c := &DummyMapMySqlPersistence{}
	c.IdentifiableMySqlPersistence = *persist.InheritIdentifiableMySqlPersistence[map[string]any, string](c, "dummies")
	return c
}

func (c *DummyMapMySqlPersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableMySqlPersistence.DefineSchema()
	c.EnsureSchema("CREATE TABLE `" + c.TableName + "` (id VARCHAR(32) PRIMARY KEY, `key` VARCHAR(50), `content` TEXT)")
	c.EnsureIndex(c.IdentifiableMySqlPersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyMapMySqlPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[map[string]any], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "`key`='" + key + "'"
	}
	sorting := ""

	return c.IdentifiableMySqlPersistence.GetPageByFilter(ctx,
		filterObj, paging, sorting, "",
	)
}

func (c *DummyMapMySqlPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "`key`='" + key + "'"
	}
	return c.IdentifiableMySqlPersistence.GetCountByFilter(ctx, filterObj)
}
