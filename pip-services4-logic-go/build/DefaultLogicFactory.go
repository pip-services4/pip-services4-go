package build

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-logic-go/cache"
	"github.com/pip-services4/pip-services4-go/pip-services4-logic-go/lock"
	"github.com/pip-services4/pip-services4-go/pip-services4-logic-go/state"
)

// Creates ICache components by their descriptors.

var NullCacheDescriptor = refer.NewDescriptor("pip-services", "cache", "null", "*", "1.0")
var MemoryCacheDescriptor = refer.NewDescriptor("pip-services", "cache", "memory", "*", "1.0")

var NullLockDescriptor = refer.NewDescriptor("pip-services", "lock", "null", "*", "1.0")
var MemoryLockDescriptor = refer.NewDescriptor("pip-services", "lock", "memory", "*", "1.0")

var NullStateStoreDescriptor = refer.NewDescriptor("pip-services", "state-store", "null", "*", "1.0")
var MemoryStateStoreDescriptor = refer.NewDescriptor("pip-services", "state-store", "memory", "*", "1.0")

// NewDefaultLogicFactory create a new instance of the factory.
//
//	Returns: *build.Factory
func NewDefaultLogicFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(NullCacheDescriptor, cache.NewNullCache[any])
	factory.RegisterType(MemoryCacheDescriptor, cache.NewMemoryCache[any])

	factory.RegisterType(NullLockDescriptor, lock.NewNullLock)
	factory.RegisterType(MemoryLockDescriptor, lock.NewMemoryLock)

	factory.RegisterType(NullStateStoreDescriptor, state.NewEmptyNullStateStore[any])
	factory.RegisterType(MemoryStateStoreDescriptor, state.NewEmptyMemoryStateStore[any])

	return factory
}
