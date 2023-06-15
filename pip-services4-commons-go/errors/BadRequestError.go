package errors

//Errors due to improper user requests.
//
//For example: missing or incorrect parameters.

// NewBadRequestError Creates an error instance and assigns its values.
//	see ErrorCategory
//	Parameters:
//		- correlation_id string a unique transaction id to trace execution through call chain.
//		- code string a unique error code.
//		- message string a human-readable description of the error.
//	Returns: *ApplicationError
func NewBadRequestError(correlationId, code, message string) *ApplicationError {
	return &ApplicationError{
		Category:      BadRequest,
		CorrelationId: correlationId,
		Code:          code,
		Message:       message,
		Status:        400,
	}
}
