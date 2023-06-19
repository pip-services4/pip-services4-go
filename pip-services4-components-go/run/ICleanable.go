package run

import "context"

// ICleanable Interface for components that should clean their state.
// Cleaning state most often is used during testing.
// But there may be situations when it can be done in production.
//	see Cleaner
//	Example:
//		type MyObjectWithState {
//			_state any
//		}
//		...
//		func (mo *MyObjectWithState) clear(ctx context.Context) {
//			mo._state = any
//		}
type ICleanable interface {
	// Clear clears component state.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//  transaction id to trace execution through call chain.
	Clear(ctx context.Context) error
}
