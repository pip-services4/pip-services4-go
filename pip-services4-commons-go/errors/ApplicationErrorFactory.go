package errors

// ApplicationErrorFactory is a factory to recreate exceptions from ErrorDescription values passed through the wire.
//	see ErrorDescription
//	see ApplicationError
var ApplicationErrorFactory *_TApplicationErrorFactory = &_TApplicationErrorFactory{}

type _TApplicationErrorFactory struct{}

// Create recreates ApplicationError object from serialized ErrorDescription.
// It tries to restore original exception type using type or error category fields.
//	Parameters: description: ErrorDescription a serialized error description received as a result of remote call
//	Returns: *ApplicationError
func (c *_TApplicationErrorFactory) Create(description *ErrorDescription) *ApplicationError {
	return NewErrorFromDescription(description)
}

// NewErrorFromDescription Recreates ApplicationError object from description.
// It tries to restore original exception type using type or error category fields.
//	Parameters: description: ErrorDescription a serialized error description received as a result of remote call
//	Returns: *ApplicationError
func NewErrorFromDescription(description *ErrorDescription) *ApplicationError {
	if description == nil {
		return nil
	}

	var err *ApplicationError = nil
	category := description.Category
	code := description.Code
	message := description.Message
	correlationId := description.CorrelationId

	// Create well-known exception type based on error category
	switch category {
	case Unknown:
		err = NewUnknownError(correlationId, code, message)
		break
	case Internal:
		err = NewInternalError(correlationId, code, message)
		break
	case Misconfiguration:
		err = NewConfigError(correlationId, code, message)
		break
	case NoResponse:
		err = NewConnectionError(correlationId, code, message)
		break
	case FailedInvocation:
		err = NewInvocationError(correlationId, code, message)
		break
	case FileError:
		err = NewFileError(correlationId, code, message)
		break
	case BadRequest:
		err = NewBadRequestError(correlationId, code, message)
		break
	case Unauthorized:
		err = NewUnauthorizedError(correlationId, code, message)
		break
	case Conflict:
		err = NewConflictError(correlationId, code, message)
		break
	case NotFound:
		err = NewNotFoundError(correlationId, code, message)
		break
	case InvalidState:
		err = NewInvalidStateError(correlationId, code, message)
		break
	case Unsupported:
		err = NewUnsupportedError(correlationId, code, message)
		break
	default:
		err = NewUnknownError(correlationId, code, message)
		err.Category = category
		err.Status = description.Status
	}

	// Fill error with details
	err.Details = description.Details
	err.Cause = description.Cause
	err.StackTrace = description.StackTrace

	return err
}
