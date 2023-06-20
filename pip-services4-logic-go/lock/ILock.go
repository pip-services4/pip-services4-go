package lock

import "context"

// ILock Interface for locks to synchronize work or parallel processes and to prevent collisions.
// The lock allows managing multiple locks identified by unique keys.
type ILock interface {
	// TryAcquireLock Makes a single attempt to acquire a lock by its key.
	// It returns immediately a positive or negative result.
	TryAcquireLock(ctx context.Context, key string, ttl int64) (bool, error)

	// AcquireLock makes multiple attempts to acquire a lock by its key within
	// give time interval.
	AcquireLock(ctx context.Context, key string, ttl int64, timeout int64) error

	// ReleaseLock releases previously acquired lock by its key.
	ReleaseLock(ctx context.Context, key string) error
}
