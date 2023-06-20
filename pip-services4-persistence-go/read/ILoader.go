package read

import "context"

// ILoader interface for data processing components that load data items.
//	Typed params:
//		- T any type of getting element
type ILoader[T any] interface {

	// Load data items.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//	Returns: []T, error a list of data items or error.
	Load(ctx context.Context) (items []T, err error)
}
