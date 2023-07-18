package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-mysql-go/test/fixtures"
)

type DummyRefMySqlPersistence struct {
	persist.IdentifiableMySqlPersistence[*fixtures.Dummy, string]
}

func NewDummyRefMySqlPersistence() *DummyRefMySqlPersistence {
	c := &DummyRefMySqlPersistence{}
	c.IdentifiableMySqlPersistence = *persist.InheritIdentifiableMySqlPersistence[*fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummyRefMySqlPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[*fixtures.Dummy], err error) {

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

func (c *DummyRefMySqlPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "`key`='" + key + "'"
	}
	return c.IdentifiableMySqlPersistence.GetCountByFilter(ctx, filterObj)
}
