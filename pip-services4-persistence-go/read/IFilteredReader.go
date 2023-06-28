package read

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
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
		filter cquery.FilterParams, sort cquery.SortParams) (items []T, err error)
}
