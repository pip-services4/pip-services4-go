package write

import (
	"context"
)

// ISaver interface for data processing components that save data items.
//
//	Typed params:
//		- T any type of getting element
type ISaver[T any] interface {

	// Save given data items.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- items []T a list of items to save.
	//	Returns: error or nil for success.
	Save(ctx context.Context, items []T) error
}
