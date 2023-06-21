package commands

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
)

// Event concrete implementation of IEvent interface. It allows to send asynchronous
// notifications to multiple subscribed listeners.
//
//	Example:
//		event: = NewEvent("my_event");
//		event.AddListener(myListener);
//		event.Notify(cpntext.Backgroudn(),
//			Parameters.fromTuples(
//				"param1", "ABC",
//				"param2", 123,
//			)
//		);
type Event struct {
	name      string
	listeners []IEventListener
}

// NewEvent creates a new event and assigns its name.
// Throws an Error if the name is null.
//
//	Parameters: name: string the name of the event that is to be created.
//	Returns: Event
func NewEvent(name string) *Event {
	if name == "" {
		panic("Name cannot be empty")
	}

	return &Event{
		name:      name,
		listeners: []IEventListener{},
	}
}

// Name gets the name of the event.
//
//	Returns: string the name of this event.
func (c *Event) Name() string {
	return c.name
}

// Listeners gets all listeners registered in this event.
//
//	Returns: []IEventListener a list of listeners.
func (c *Event) Listeners() []IEventListener {
	return c.listeners
}

// AddListener adds a listener to receive notifications when this event is fired.
//
//	Parameters: listener: IEventListener the listener reference to add.
func (c *Event) AddListener(listener IEventListener) {
	c.listeners = append(c.listeners, listener)
}

// RemoveListener removes a listener, so that it no longer receives notifications for this event.
//
//	Parameters: listener: IEventListener the listener reference to remove.
func (c *Event) RemoveListener(listener IEventListener) {
	for i, l := range c.listeners {
		if listener == l {
			c.listeners = append(c.listeners[:i], c.listeners[i+1:]...)
			break
		}
	}
}

// Notify fires this event and notifies all registred listeners.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- args: Parameters the parameters to raise this event with.
func (c *Event) Notify(ctx context.Context, args *exec.Parameters) {
	for _, listener := range c.listeners {
		listener.OnEvent(ctx, c, args)
	}
}
