package build

// Error raised when factory is not able to create requested component.

import (
	"fmt"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
)

// NewCreateError creates an error instance and assigns its values.
//
//	Parameters:
//		- traceId string
//		- message string human-readable error of the component that cannot be created.
//	Returns: *errors.ApplicationError
func NewCreateError(traceId, message string) *errors.ApplicationError {
	return errors.NewInternalError(traceId, "CANNOT_CREATE", message)
}

// NewCreateErrorByLocator creates an error instance and assigns its values.
//
//	Parameters:
//		- trace_id string
//		- locator any human-readable locator of the component that cannot be created.
//	Returns: *errors.ApplicationError
func NewCreateErrorByLocator(traceId string, locator any) *errors.ApplicationError {
	message := fmt.Sprintf("Requested component %v cannot be created", locator)
	return errors.NewInternalError(traceId, "CANNOT_CREATE", message).
		WithDetails("locator", locator)
}
