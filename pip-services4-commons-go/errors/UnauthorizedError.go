package errors

// Access errors caused by missing user identity (authentication error)
// or incorrect security permissions (authorization error).

// NewUnauthorizedError creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- traceId string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewUnauthorizedError(traceId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category: Unauthorized,
		TraceId:  traceId,
		Code:     code,
		Message:  message,
		Status:   401,
	}
}
