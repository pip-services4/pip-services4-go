package errors

// Errors that occur during connections to remote services.
// They can be related to misconfiguration, network issues, or the remote service itself.

// NewConnectionError creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- correlation_id string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewConnectionError(correlationId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category:      NoResponse,
		CorrelationId: correlationId,
		Code:          code,
		Message:       message,
		Status:        500,
	}
}
