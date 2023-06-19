package exec

import "context"

// Notifier helper class that notifies components.
var Notifier = &_TNotifier{}

type _TNotifier struct{}

// NotifyOne notifies specific component.
// To be notified components must implement INotifiable interface.
// If they don't the call to this method has no effect.
//	see INotifiable
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component any the component that is to be notified.
//		- args *Parameters notification arguments.
func (c *_TNotifier) NotifyOne(ctx context.Context, component any, args *Parameters) {
	if v, ok := component.(INotifiable); ok {
		v.Notify(ctx, args)
	}
}

// Notify notifies multiple components.
// To be notified components must implement INotifiable interface.
// If they don't the call to this method has no effect.
//	see NotifyOne
//	see INotifiable
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- components []any a list of components that are to be notified.
//		- args *Parameters notification arguments.
func (c *_TNotifier) Notify(ctx context.Context, components []any, args *Parameters) {
	for _, component := range components {
		c.NotifyOne(ctx, component, args)
	}
}
