package build

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-messaging-go/queues"
)

// Creates message queue componens.
type IMessageQueueFactory interface {
	// Creates a message queue component and assigns its name.
	// Parameters:
	//   - name: a name of the created message queue.
	CreateQueue(name string) queues.IMessageQueue
}
