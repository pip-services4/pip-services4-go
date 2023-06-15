package errors

// ApplicationError Defines a base class to defive various application exceptions.
//
// Most languages have own definition of base exception (error) types.
// However, this class is implemented symmetrically in all languages supported by
// PipServices toolkit. It allows to create portable implementations and support
// proper error propagation in microservices calls.
//
// Error propagation means that when microservice implemented in one
// language calls microservice(s) implemented in a different language(s),
// errors are returned throught the entire call chain and restored in their original (or close) type.
//
// Since number of potential exception types is endless,
// PipServices toolkit supports only 12 standard categories of exceptions defined in ErrorCategory.
// This ApplicationException class acts as a basis for all other 12 standard exception types.
//
// Most exceptions have just free-form message that describes occured error.
// That may not be sufficient to create meaninful error descriptions.
// The ApplicationException class proposes an extended error definition that has more standard fields:
//
//	message: is a human-readable error description
//	category: one of 12 standard error categories of errors
//	status: numeric HTTP status code for REST invocations
//	code: a unique error code, usually defined as "MY_ERROR_CODE"
//	correlation_id: a unique transaction id to trace execution through a call chain
//	details: map with error parameters that can help to recreate meaningful error description in other languages
//	stack_trace: a stack trace
//	cause: original error that is wrapped by this exception
//	ApplicationException class is not serializable.
//		To pass errors through the wire it is converted into
//		ErrorDescription object and restored on receiving end into identical exception type.
//
//	see ErrorCategory
//	see ErrorDescription
type ApplicationError struct {
	Message       string         `json:"message"`
	Category      string         `json:"category"`
	Status        int            `json:"status"`
	Code          string         `json:"code"`
	Details       map[string]any `json:"details"`
	CorrelationId string         `json:"correlation_id"`
	StackTrace    string         `json:"stack_trace"`
	Cause         string         `json:"cause"`
}

// Error return string error message
func (e *ApplicationError) Error() string {
	return e.Message
}

// WithCode Add code to ApplicationError
//	Parameters: code string a error code
//	Returns: *ApplicationError
func (e *ApplicationError) WithCode(code string) *ApplicationError {
	e.Code = code
	return e
}

// WithStatus add status to ApplicationError
//	Parameters: status int a status code
//	Returns: *ApplicationError
func (e *ApplicationError) WithStatus(status int) *ApplicationError {
	e.Status = status
	return e
}

// WithDetails add error details to ApplicationError
//	Parameters:
//		- key string a detail key word
//		- value any an value of detail object
//	Returns: *ApplicationError
func (e *ApplicationError) WithDetails(key string, value any) *ApplicationError {
	if e.Details == nil {
		e.Details = map[string]any{}
	}
	e.Details[key] = value
	return e
}

// WithCause add cause to ApplicationError
//	Parameters: cause error a cause error object
//	Returns: *ApplicationError
func (e *ApplicationError) WithCause(cause error) *ApplicationError {
	e.Cause = cause.Error()
	return e
}

// WithCauseString add cause to ApplicationError
//	Parameters: cause string a cause string describe an error
//	Returns: *ApplicationError
func (e *ApplicationError) WithCauseString(cause string) *ApplicationError {
	e.Cause = cause
	return e
}

// WithCorrelationId add Correlation Id to ApplicationError
//	Parameters: correlationId string a correlation string
//	Returns: *ApplicationError
func (e *ApplicationError) WithCorrelationId(correlationId string) *ApplicationError {
	e.CorrelationId = correlationId
	return e
}

// Wrap error by ApplicationError struct
//	Parameters: err error an error what neet to wrap
//	Returns: *ApplicationError
func (e *ApplicationError) Wrap(err error) *ApplicationError {
	if er, ok := err.(*ApplicationError); ok {
		return er
	}

	e.WithCause(err)
	return e
}

// WrapError wrap error by ApplicationError struct and sets message
//	Parameters:
//		- err error an error what neet to wrap
//		- message string error message
//	Returns: *ApplicationError
func WrapError(err error, message string) *ApplicationError {
	if e, ok := err.(*ApplicationError); ok {
		return e
	}

	return NewError(message).WithCause(err)
}

// NewError creates a new instance of application error and assigns its message.
//	Parameters: message string an error message
//	Returns: *ApplicationError generated new ApplicationError
func NewError(message string) *ApplicationError {
	if message == "" {
		message = "Unknown error"
	}
	return &ApplicationError{Code: "UNKNOWN", Message: message, Status: 500}
}
