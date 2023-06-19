package run

import "context"

// IOpenable interface for components that require explicit opening and closing.
// For components that perform opening on demand consider using ICloseable interface instead.
//	see IOpenable
//	see Opener
//	Example:
//		type MyPersistence {
//			_client any
//		}
//
//		func (mp* MyPersistence)IsOpen() bool {
//			return mp._client != nil;
//		}
//
//		func (mp* MyPersistence) Open(ctx context.Context) error {
//			if (mp.isOpen()) {
//				return nil;
//			}
//		}
//
//		func (mp* MyPersistence) Close(ctx context.Context) {
//			if (mp._client != nil) {
//				mp._client.Close(ctx);
//				mp._client = nil;
//			}
//		}
type IOpenable interface {
	IClosable

	// IsOpen Checks if the component is opened.
	//	Returns: bool true if the component has been opened and false otherwise.
	IsOpen() bool

	// Open opens the component.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//	Return: error
	Open(ctx context.Context) error
}
