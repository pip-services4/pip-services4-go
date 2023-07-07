package test_logic

import (
	cconvert "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cvalid "github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
)

type DummySchema struct {
	cvalid.ObjectSchema
}

func NewDummySchema() *DummySchema {
	c := DummySchema{
		ObjectSchema: *cvalid.NewObjectSchema(),
	}

	c.WithOptionalProperty("id", cconvert.String)
	c.WithRequiredProperty("key", cconvert.String)
	c.WithOptionalProperty("content", cconvert.String)

	return &c
}
