package containers

import cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"

type ILambdaFunctionOverrides interface {
	cref.IReferenceable
	// Perform required registration steps.
	Register()
}
