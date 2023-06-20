package lock

import (
	"context"
	"sync"
	"time"
)

// MemoryLock Lock that is used to synchronize execution within one process using shared memory.
//
//	Configuration parameters:
//		options:
//		retry_timeout: timeout in milliseconds to retry lock acquisition. (Default: 100)
//	see ILock
//	see Lock
//	Example:
//		lock := NewMemoryLock()
//		err = lock.AcquireLock(context.Background(), "123", "key1")
//		if err == nil {
//			 _ = lock.ReleaseLock(context.Background(), "123", "key1")
//			// Processing...
//		}
type MemoryLock struct {
	*Lock
	mux   sync.Mutex
	locks map[string]time.Time
}

// NewMemoryLock create new memory lock
//
//	Returns: *MemoryLock
func NewMemoryLock() *MemoryLock {
	c := &MemoryLock{
		locks: map[string]time.Time{},
	}
	c.Lock = InheritLock(c)

	return c
}

// TryAcquireLock makes a single attempt to acquire a lock by its key.
// It returns immediately a positive or negative result.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique lock key to acquire.
//		- ttl int64 a lock timeout (time to live) in milliseconds.
//	Returns bool, error true if locked. Error object
func (c *MemoryLock) TryAcquireLock(ctx context.Context,
	key string, ttl int64) (bool, error) {

	c.mux.Lock()
	defer c.mux.Unlock()

	expireTime, ok := c.locks[key]
	if ok {
		if expireTime.After(time.Now()) {
			return false, nil
		}
	}

	expireTime = time.Now().Add(time.Duration(ttl) * time.Millisecond)
	c.locks[key] = expireTime

	return true, nil
}

// ReleaseLock releases the lock with the given key.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string the key of the lock that is to be released.
//	Return: error
func (c *MemoryLock) ReleaseLock(ctx context.Context,
	key string) error {

	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.locks, key)

	return nil
}
