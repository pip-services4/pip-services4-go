package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/test/fixtures"
)

type DummyJsonMySqlPersistence struct {
	*persist.IdentifiableJsonMySqlPersistence[fixtures.Dummy, string]
}

func NewDummyJsonMySqlPersistence() *DummyJsonMySqlPersistence {
	c := &DummyJsonMySqlPersistence{}
	c.IdentifiableJsonMySqlPersistence = persist.InheritIdentifiableJsonMySqlPersistence[fixtures.Dummy, string](c, "dummies_json")
	return c
}

func (c *DummyJsonMySqlPersistence) DefineSchema() {
	c.ClearSchema()
	c.EnsureTable("", "")
	c.EnsureSchema("ALTER TABLE `" + c.TableName + "` ADD `data_key` VARCHAR(50) AS (JSON_UNQUOTE(`data`->\"$.key\"))")
	c.EnsureIndex(c.TableName+"_json_key", map[string]string{"data_key": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyJsonMySqlPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "data->'$.key'='" + key + "'"
	}

	return c.IdentifiableJsonMySqlPersistence.GetPageByFilter(ctx,
		filterObj, paging,
		"", "",
	)
}

func (c *DummyJsonMySqlPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	filterObj := ""
	if key, ok := filter.GetAsNullableString("Key"); ok && key != "" {
		filterObj += "data->'$.key'='" + key + "'"
	}

	return c.IdentifiableJsonMySqlPersistence.GetCountByFilter(ctx, filterObj)
}

func (c *DummyJsonMySqlPersistence) GetOneRandom(ctx context.Context) (item fixtures.Dummy, err error) {
	return c.IdentifiableJsonMySqlPersistence.GetOneRandom(ctx, "")
}
