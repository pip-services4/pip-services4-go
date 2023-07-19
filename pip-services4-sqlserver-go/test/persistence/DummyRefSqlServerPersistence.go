package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-sqlserver-go/test/fixtures"
)

type DummyRefSqlServerPersistence struct {
	*persist.IdentifiableSqlServerPersistence[*fixtures.Dummy, string]
}

func NewDummyRefSqlServerPersistence() *DummyRefSqlServerPersistence {
	c := &DummyRefSqlServerPersistence{}
	c.IdentifiableSqlServerPersistence = persist.InheritIdentifiableSqlServerPersistence[*fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummyRefSqlServerPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[*fixtures.Dummy], err error) {

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

func (c *DummyRefSqlServerPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "[key]='" + key + "'"
	}
	return c.IdentifiableSqlServerPersistence.GetCountByFilter(ctx, filterObj)
}
