package lock

import "context"

// NullLock Dummy lock implementation that doesn't do anything.
//It can be used in testing or in situations when lock is required but shall be disabled.
type NullLock struct{}

func NewNullLock() *NullLock {
	return &NullLock{}
}

// TryAcquireLock makes a single attempt to acquire a lock by its key.
// It returns immediately a positive or negative result.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique lock key to acquire.
//		- ttl int64 a lock timeout (time to live) in milliseconds.
//	Returns bool, error true if locked. Error object
func (c *NullLock) TryAcquireLock(ctx context.Context,
	key string, ttl int) (bool, error) {
	return true, nil
}

// AcquireLock makes multiple attempts to acquire a lock by its key within give time interval.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string a unique lock key to acquire.
//		- ttl int64 a lock timeout (time to live) in milliseconds.
//		- timeout int64 a lock acquisition timeout.
//	Returns: error
func (c *NullLock) AcquireLock(ctx context.Context,
	key string, ttl int, timeout int) error {
	return nil
}

// ReleaseLock releases the lock with the given key.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- key string the key of the lock that is to be released.
//	Return: error
func (c *NullLock) ReleaseLock(ctx context.Context,
	key string) error {
	return nil
}
