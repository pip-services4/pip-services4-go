package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-sqlite-go/test/fixtures"
)

type DummyRefSqlitePersistence struct {
	*persist.IdentifiableSqlitePersistence[*fixtures.Dummy, string]
}

func NewDummyRefSqlitePersistence() *DummyRefSqlitePersistence {
	c := &DummyRefSqlitePersistence{}
	c.IdentifiableSqlitePersistence = persist.InheritIdentifiableSqlitePersistence[*fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummyRefSqlitePersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[*fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	sorting := ""

	return c.IdentifiableSqlitePersistence.GetPageByFilter(ctx,
		filterObj, paging,
		sorting, "",
	)
}

func (c *DummyRefSqlitePersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	return c.IdentifiableSqlitePersistence.GetCountByFilter(ctx, filterObj)
}
