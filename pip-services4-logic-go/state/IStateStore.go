package state

import "context"

// IStateStore interface for state storages that are used to store and retrieve transaction states.
type IStateStore[T any] interface {

	// Load state from the store using its key.
	// If value is missing in the store it returns nil.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- key           a unique state key.
	//	Returns: the state value or nil if value wasn't found.
	Load(ctx context.Context, key string) T

	// LoadBulk loads an array of states from the store using their keys.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- keys          unique state keys.
	//	Returns: an array with state values and their corresponding keys.
	LoadBulk(ctx context.Context, keys []string) []StateValue[T]

	// Save state into the store.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- key           a unique state key.
	//		- value         a state value.
	//	Returns: the state that was stored in the store.
	Save(ctx context.Context, key string, value T) T

	// Delete a state from the store by its key.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- key           a unique value key.
	//	Returns: the state that was deleted in the store.
	Delete(ctx context.Context, key string) T
}
