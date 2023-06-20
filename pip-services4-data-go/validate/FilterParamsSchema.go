package validate

// Schema to validate FilterParams.

import "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"

// NewFilterParamsSchema creates a new instance of validation schema.
//	Returns: *MapSchema
func NewFilterParamsSchema() *MapSchema {
	return NewMapSchema(convert.String, nil)
}
