package write

import (
	"context"
)

// IWriter interface for data processing components
// that can create, update and delete data items.
//
//	Typed params:
//		- T any type
//		- K any type of id (key)
type IWriter[T any, K any] interface {

	// Create creates a data item.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- item T an item to be created.
	//	Returns: T, error created item or error.
	Create(ctx context.Context, item T) (value T, err error)

	// Update a data item.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- item T an item to be updated.
	//	Returns: T, error updated item or error.
	Update(ctx context.Context, item T) (value T, err error)

	// DeleteById a data item by it's unique id.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- id K an id of the item to be deleted
	//	Returns: T, error deleted item or error.
	DeleteById(ctx context.Context, id K) (value T, err error)
}
