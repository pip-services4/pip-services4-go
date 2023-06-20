package read

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// IFilteredReader interface for data processing components that can
// retrieve a list of data items by filter.
//
//	Typed params:
//		- T any type
type IFilteredReader[T any] interface {

	// GetListByFilter gets a list of data items using filter
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- filter data.FilterParams filter parameters
	//		- sort  data.SortParams sort parameters
	//	Returns: []T, error receives list of items or error.
	GetListByFilter(ctx context.Context,
		filter cdata.FilterParams, sort cdata.SortParams) (items []T, err error)
}
