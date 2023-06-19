package errors

// Errors related to calling operations, which require the component to be in a specific state.
// For instance: business calls when the component is not ready.

// NewInvalidStateError creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- traceId string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewInvalidStateError(traceId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category: InvalidState,
		TraceId:  traceId,
		Code:     code,
		Message:  message,
		Status:   500,
	}
}
