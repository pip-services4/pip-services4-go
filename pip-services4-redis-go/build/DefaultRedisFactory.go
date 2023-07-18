package build

import (
	cbuild "github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	rediscache "github.com/pip-services4/pip-services4-go/pip-services4-redis-go/cache"
	redislock "github.com/pip-services4/pip-services4-go/pip-services4-redis-go/lock"
)

/*
DefaultRedisFactory are creates Redis components by their descriptors.

See RedisCache
See RedisLock
*/
type DefaultRedisFactory struct {
	*cbuild.Factory
	Descriptor           *cref.Descriptor
	RedisCacheDescriptor *cref.Descriptor
	RedisLockDescriptor  *cref.Descriptor
}

// NewDefaultRedisFactory method are create a new instance of the factory.
func NewDefaultRedisFactory() *DefaultRedisFactory {

	c := DefaultRedisFactory{}
	c.Factory = cbuild.NewFactory()
	c.Descriptor = cref.NewDescriptor("pip-services", "factory", "redis", "default", "1.0")
	c.RedisCacheDescriptor = cref.NewDescriptor("pip-services", "cache", "redis", "*", "1.0")
	c.RedisLockDescriptor = cref.NewDescriptor("pip-services", "lock", "redis", "*", "1.0")
	c.RegisterType(c.RedisCacheDescriptor, rediscache.NewRedisCache[any])
	c.RegisterType(c.RedisLockDescriptor, redislock.NewRedisLock)
	return &c
}
