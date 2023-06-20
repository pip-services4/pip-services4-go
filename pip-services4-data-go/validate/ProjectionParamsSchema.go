package validate

// Schema to validate ProjectionParams

import "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"

// NewProjectionParamsSchema creates a new instance of validation schema.
//	Returns: *ArraySchema
func NewProjectionParamsSchema() *ArraySchema {
	return NewArraySchema(convert.String)
}
