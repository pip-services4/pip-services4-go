package write

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// IPartialUpdater interface for data processing components to update data items partially.
//
//	Typed params:
//		- T any type
//		- K type of id (key)
type IPartialUpdater[T any, K any] interface {

	// UpdatePartially updates only few selected fields in a data item.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- id K an id of data item to be updated.
	//		- data data.AnyValueMap a map with fields to be updated.
	//	Returns: T, error updated item or error.
	UpdatePartially(ctx context.Context, id K, data cdata.AnyValueMap) (item T, err error)
}
