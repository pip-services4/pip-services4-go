package errors

// Errors returned by remote services or by the network during call attempts.

// NewInvocationError creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- correlation_id string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewInvocationError(correlationId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category:      FailedInvocation,
		CorrelationId: correlationId,
		Code:          code,
		Message:       message,
		Status:        500,
	}
}
