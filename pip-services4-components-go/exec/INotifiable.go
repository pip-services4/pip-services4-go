package exec

import "context"

// INotifiable interface for components that can be asynchronously notified.
// The notification may include optional argument that describe the occurred event.
//	see Notifier
//	see IExecutable
//	Example:
//		type MyComponent {}
//		...
//		func (mc *MyComponent)Notify(ctx context.Context, args Parameters) {
//			fmt.Println("Occured event " + args.GetAsString("event"));
//		}
//		myComponent := MyComponent{};
//		myComponent.Notify("123", NewParametersFromTuples("event", "Test Event"));
type INotifiable interface {
	// Notify notifies the component about occured event.
	//	Parameters:
	//		- ctx context.Context a context to trace execution through call chain.
	//		- args *Parameters notification arguments.
	Notify(ctx context.Context, args *Parameters)
}
