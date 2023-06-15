package errors

// Errors caused by attempts to access missing objects.

// NewNotFoundError creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- correlation_id string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewNotFoundError(correlationId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category:      NotFound,
		CorrelationId: correlationId,
		Code:          code,
		Message:       message,
		Status:        404,
	}
}
