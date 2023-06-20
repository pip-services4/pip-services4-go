package validate

//Schema to validate TokenizedPagingParams.

import "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"

// NewTokenizedPagingParamsSchema creates a new instance of validation schema.
//	Returns: *TokenizedPagingParamsSchema
func NewTokenizedPagingParamsSchema() *ObjectSchema {
	return NewObjectSchema().
		WithOptionalProperty("token", convert.String).
		WithOptionalProperty("take", convert.Long).
		WithOptionalProperty("total", convert.Boolean)
}
