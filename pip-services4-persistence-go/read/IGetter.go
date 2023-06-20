package read

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// IGetter Interface for data processing components that can get data items.
//
//	Typed params:
//		- T cdata.IIdentifiable[T] any type
//		- K any type of id (key)
type IGetter[T cdata.IIdentifiable[T], K any] interface {

	// GetOneById a data items by its unique id.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- id an id of item to be retrieved.
	//	Returns: T, error item or error
	GetOneById(ctx context.Context, id K) (item T, err error)
}
