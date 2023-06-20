package validate

import "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"

// ISchema validation schema interface
type ISchema interface {
	Validate(value any) []*ValidationResult
	ValidateAndReturnError(traceId string, value any, strict bool) *errors.ApplicationError
	ValidateAndThrowError(traceId string, value any, strict bool)
}
