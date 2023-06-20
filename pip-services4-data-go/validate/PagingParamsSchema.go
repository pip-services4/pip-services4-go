package validate

// Schema to validate PagingParams.

import "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"

// NewPagingParamsSchema creates a new instance of validation schema.
//	Returns: *PagingParamsSchema
func NewPagingParamsSchema() *ObjectSchema {
	return NewObjectSchema().
		WithOptionalProperty("skip", convert.Long).
		WithOptionalProperty("take", convert.Long).
		WithOptionalProperty("total", convert.Boolean)
}
