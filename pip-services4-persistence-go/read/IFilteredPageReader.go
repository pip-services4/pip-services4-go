package read

import (
	"context"

	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
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
		filter cquery.FilterParams, paging cquery.PagingParams, sort cquery.SortParams) (page cquery.DataPage[T], err error)
}
