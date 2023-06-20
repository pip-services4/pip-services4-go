package read

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// IFilteredPageReader is interface for data processing components
// that can retrieve a page of data items by a filter.
//
//	Typed params:
//		- T any type
type IFilteredPageReader[T any] interface {

	// GetPageByFilter gets a page of data items using filter
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- filter  data.FilterParams filter parameters
	//		- paging data.PagingParams paging parameters
	//		- sort data.SortParams sort parameters
	//	Returns: data.DataPage[T], error list of items or error.
	GetPageByFilter(ctx context.Context,
		filter cdata.FilterParams, paging cdata.PagingParams, sort cdata.SortParams) (page cdata.DataPage[T], err error)
}
