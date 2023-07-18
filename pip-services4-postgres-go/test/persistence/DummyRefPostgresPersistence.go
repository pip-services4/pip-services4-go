package test

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	persist "github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/persistence"
	"github.com/pip-services4/pip-services4-go/pip-services4-postgres-go/test/fixtures"
)

type DummyRefPostgresPersistence struct {
	persist.IdentifiablePostgresPersistence[*fixtures.Dummy, string]
}

func NewDummyRefPostgresPersistence() *DummyRefPostgresPersistence {
	c := &DummyRefPostgresPersistence{}
	c.IdentifiablePostgresPersistence = *persist.InheritIdentifiablePostgresPersistence[*fixtures.Dummy, string](c, "dummies")
	return c
}

func (c *DummyRefPostgresPersistence) GetPageByFilter(ctx context.Context,
	filter cquery.FilterParams, paging cquery.PagingParams) (page cquery.DataPage[*fixtures.Dummy], err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	sorting := ""

	return c.IdentifiablePostgresPersistence.GetPageByFilter(ctx,
		filterObj, paging,
		sorting, "",
	)
}

func (c *DummyRefPostgresPersistence) GetCountByFilter(ctx context.Context,
	filter cquery.FilterParams) (count int64, err error) {

	key, ok := filter.GetAsNullableString("Key")
	filterObj := ""
	if ok && key != "" {
		filterObj += "key='" + key + "'"
	}
	return c.IdentifiablePostgresPersistence.GetCountByFilter(ctx, filterObj)
}
