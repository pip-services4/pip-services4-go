package test_data

import (
	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type DummySchema struct {
	*cvalid.ObjectSchema
}

func NewDummySchema() *DummySchema {
	ds := DummySchema{}
	ds.ObjectSchema = cvalid.NewObjectSchema()
	ds.WithOptionalProperty("id", cconv.String)
	ds.WithRequiredProperty("key", cconv.String)
	ds.WithOptionalProperty("content", cconv.String)
	return &ds
}
