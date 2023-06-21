package commands

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
)

// IEventListener an interface for listener objects that receive notifications on fired events.
//
//	see IEvent
//	see Event
//	Example:
//		type MyListener struct {
//			msg string
//		}
//
//		func (l *MyListener) OnEvent(ctx context.Context, event IEvent, args Parameters) {
//			fmt.Println("Fired event " + event.Name())
//		}
//
//		var event = NewEvent("myevent")
//		_listener := MyListener{}
//		event.AddListener(_listener)
//		event.Notify(context.Background(), "123", Parameters.FromTuples("param1", "ABC"))
//
//		// Console output: Fired event myevent
type IEventListener interface {
	// OnEvent a method called when events this listener is subscrubed to are fired.
	//	Parameters:
	//		- ctx context.Context execution context to trace execution through call chain.
	//		- e: IEvent a fired evemt
	//		- value: *run.Parameters event arguments.
	OnEvent(ctx context.Context, e IEvent, value *exec.Parameters)
}
