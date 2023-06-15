package errors

// Errors in read/write local disk operations.

// NewFileError creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- correlation_id string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewFileError(correlationId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category:      FileError,
		CorrelationId: correlationId,
		Code:          code,
		Message:       message,
		Status:        500,
	}
}
