package errors

// Errors raised by conflicts between object versions
// that were posted by the user and those that are stored on the server.

// NewConflictError creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- correlation_id string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewConflictError(correlationId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category:      Conflict,
		CorrelationId: correlationId,
		Code:          code,
		Message:       message,
		Status:        409,
	}
}
