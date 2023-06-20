package read

import (
	"context"

	cdata "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// IQuerableReader interface for data processing components that can query a list of data items.
//
//	Typed params:
//		- T any type
type IQuerableReader[T any] interface {

	// GetListByQuery gets a list of data items using a query string.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- query string a query string
	//		- sort data.SortParams sort parameters
	// Returns []T, error list of items or error.
	GetListByQuery(ctx context.Context,
		query string, sort cdata.SortParams) (items []T, err error)
}
