package refer

// Error when required component dependency cannot be found.

import (
	"context"
	"fmt"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
)

// NewReferenceError Creates an error instance and assigns its values.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- locator any the locator to find reference to dependent component.
//	Returns *errors.ApplicationError
func NewReferenceError(ctx context.Context, locator any) *errors.ApplicationError {
	message := fmt.Sprintf("Failed to obtain reference to %v", locator)
	e := errors.NewInternalError(utils.ContextHelper.GetTraceId(ctx), "REF_ERROR", message)
	e.WithDetails("locator", locator)
	return e
}
