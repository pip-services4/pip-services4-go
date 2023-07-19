package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/test/fixtures"
)

type DummyJsonSqlitePersistence struct {
	*persist.IdentifiableJsonSqlitePersistence[fixtures.Dummy, string]
}

func NewDummyJsonSqlitePersistence() *DummyJsonSqlitePersistence {
	c := &DummyJsonSqlitePersistence{}
	c.IdentifiableJsonSqlitePersistence = persist.InheritIdentifiableJsonSqlitePersistence[fixtures.Dummy, string](c, "dummies_json")
	return c
}

func (c *DummyJsonSqlitePersistence) DefineSchema() {
	c.ClearSchema()
	c.IdentifiableJsonSqlitePersistence.DefineSchema()
	c.EnsureTable("", "")
	c.EnsureIndex(c.TableName+"_json_key", map[string]string{"JSON_EXTRACT(data, '$.key')": "1"}, map[string]string{"unique": "true"})
}

func (c *DummyJsonSqlitePersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "JSON_EXTRACT(data, '$.key')='" + key + "'"
	}

	return c.IdentifiableJsonSqlitePersistence.GetPageByFilter(ctx,
		filterObj, paging,
		"", "",
	)
}

func (c *DummyJsonSqlitePersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	filterObj := ""
	if key, ok := filter.GetAsNullableString("Key"); ok && key != "" {
		filterObj += "JSON_EXTRACT(data, '$.key')='" + key + "'"
	}

	return c.IdentifiableJsonSqlitePersistence.GetCountByFilter(ctx, filterObj)
}

func (c *DummyJsonSqlitePersistence) GetOneRandom(ctx context.Context) (item fixtures.Dummy, err error) {
	return c.IdentifiableJsonSqlitePersistence.GetOneRandom(ctx, "")
}
