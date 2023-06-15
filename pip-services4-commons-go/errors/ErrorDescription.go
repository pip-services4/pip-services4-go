package errors

// ErrorDescription seializeable error description. It is use to pass information about errors between microservices
// implemented in different languages. On the receiving side ErrorDescription is used to recreate exception
// object close to its original type without missing additional details.
//
//	category - Standard error category
//	cause - Original error wrapped by this exception
//	code - A unique error code
//	correlation_id - A unique transaction id to trace execution throug call chain
//	details - A map with additional details that can be used to restore error description in other languages
//	message - A human-readable error description (usually written in English)
//	stack_trace - Stack trace of the exception
//	status - HTTP status code associated with this error type
//	type - Data type of the original error
type ErrorDescription struct {
	Type          string         `json:"type"`
	Category      string         `json:"category"`
	Status        int            `json:"status"`
	Code          string         `json:"code"`
	Message       string         `json:"message"`
	Details       map[string]any `json:"details"`
	CorrelationId string         `json:"correlation_id"`
	Cause         string         `json:"cause"`
	StackTrace    string         `json:"stack_trace"`
}
