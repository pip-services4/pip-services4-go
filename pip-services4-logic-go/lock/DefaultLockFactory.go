package lock

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// Creates ILock components by their descriptors.

var NullLockDescriptor = refer.NewDescriptor("pip-services", "lock", "null", "*", "1.0")
var MemoryLockDescriptor = refer.NewDescriptor("pip-services", "lock", "memory", "*", "1.0")

// NewDefaultLockFactory create a new instance of the factory.
//
//	Returns: *build.Factory
func NewDefaultLockFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(NullLockDescriptor, NewNullLock)
	factory.RegisterType(MemoryLockDescriptor, NewMemoryLock)

	return factory
}
