package errors

import (
	"strconv"

	cerrors "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
)

// Exception that can be thrown by Mustache Template.
func NewMustacheError(traceId, code, message string, line, column int) *cerrors.ApplicationError {
	if line != 0 || column != 0 {
		message = message + " at line " + strconv.Itoa(line) + " and column " + strconv.Itoa(column)
	}
	return &cerrors.ApplicationError{
		Category: cerrors.BadRequest,
		TraceId:  traceId,
		Code:     code,
		Message:  message,
		Status:   400,
	}
}
