package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/test/fixtures"
)

type DummyMySqlPersistence struct {
	*persist.IdentifiableMySqlPersistence[fixtures.Dummy, string]
}

func NewDummyMySqlPersistence() *DummyMySqlPersistence {
	c := &DummyMySqlPersistence{}
	c.IdentifiableMySqlPersistence = persist.InheritIdentifiableMySqlPersistence[fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummyMySqlPersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableMySqlPersistence.DefineSchema()
	// Row name must be in double quotes for properly case!!!
	c.EnsureSchema("CREATE TABLE `" + c.TableName + "` (id VARCHAR(32) PRIMARY KEY, `key` VARCHAR(50), `content` TEXT)")
	c.EnsureIndex(c.IdentifiableMySqlPersistence.TableName+"_key", map[string]string{"key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyMySqlPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "`key`='" + key + "'"
	}
	sorting := ""

	return c.IdentifiableMySqlPersistence.GetPageByFilter(ctx,
		filterObj, paging,
		sorting, "",
	)
}

func (c *DummyMySqlPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "`key`='" + key + "'"
	}
	return c.IdentifiableMySqlPersistence.GetCountByFilter(ctx, filterObj)
}

func (c *DummyMySqlPersistence) GetOneRandom(ctx context.Context) (item fixtures.Dummy, err error) {
	return c.IdentifiableMySqlPersistence.GetOneRandom(ctx, "")
}
