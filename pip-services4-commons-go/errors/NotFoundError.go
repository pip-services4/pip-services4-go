package errors

// Errors caused by attempts to access missing objects.

// NewNotFoundError creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- traceId string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewNotFoundError(traceId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category: NotFound,
		TraceId:  traceId,
		Code:     code,
		Message:  message,
		Status:   404,
	}
}
