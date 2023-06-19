package errors

import (
	"fmt"
)

// ErrorDescriptionFactory is a factory to create serializeable ErrorDescription
// from ApplicationException or from arbitrary errors.
// The ErrorDescriptions are used to pass errors through the wire
// between microservices implemented in different languages.
// They allow to restore exceptions on the receiving side close to
// the original type and preserve additional information.
//
//	see ErrorDescription
//	see ApplicationError
var ErrorDescriptionFactory = &_TErrorDescriptionFactory{}

type _TErrorDescriptionFactory struct{}

// Create creates a serializable ErrorDescription from error object.
//
//	Parameters: err error an error object
//	Returns: *ErrorDescription a serializeable ErrorDescription object that describes the error.
func (c *_TErrorDescriptionFactory) Create(err any) *ErrorDescription {
	return NewErrorDescription(err)
}

// NewErrorDescription creates a serializable ErrorDescription from error object.
//
//	Parameters: err any an error object
//	Returns: *ErrorDescription a serializeable ErrorDescription object that describes the error.
func NewErrorDescription(err any) *ErrorDescription {
	description := &ErrorDescription{
		Category: Unknown,
		Code:     "UNKNOWN",
		Status:   500,
		Message:  "Unknown error",
	}

	if ex, ok := err.(*ApplicationError); ok {
		description.Category = ex.Category
		description.Status = ex.Status
		description.Code = ex.Code
		description.Message = ex.Message
		description.Details = ex.Details
		description.TraceId = ex.TraceId
		description.Cause = ex.Cause
		description.StackTrace = ex.StackTrace
	} else if err != nil {
		description.Message = fmt.Sprintf("%v", err)
	}

	return description
}
