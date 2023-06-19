package run

import "context"

// IClosable interface for components that require explicit closure.
// For components that require opening as well as closing use IOpenable interface instead.
//	see IOpenable
//	see Closer
//	Example:
//		type MyConnector {
//			_client interface{}
//		}
//		... // The _client can be lazy created
//		func (mc *MyConnector) Close(ctx context.Context) error {
//			if (mc._client != nil) {
//				mc._client.Close(ctx)
//				mc._client = nil
//				return nil
//			}
//		}
type IClosable interface {
	Close(ctx context.Context) error
}
