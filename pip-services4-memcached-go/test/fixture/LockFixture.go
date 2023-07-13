package test_fixture

import (
	"context"
	"testing"

	clock "github.com/pip-services4/pip-services4-go/pip-services4-logic-go/lock"
	"github.com/stretchr/testify/assert"
)

const (
	LOCK1 string = "lock_1"
	LOCK2 string = "lock_2"
	LOCK3 string = "lock_3"
)

type LockFixture struct {
	lock clock.ILock
}

func NewLockFixture(lock clock.ILock) *LockFixture {
	c := LockFixture{}
	c.lock = lock
	return &c
}

func (c *LockFixture) TestTryAcquireLock(t *testing.T) {

	ctx := context.Background()

	// Try to acquire lock for the first time
	result, err := c.lock.TryAcquireLock(ctx, LOCK1, 3000)
	assert.Nil(t, err)
	assert.True(t, result)

	// Try to acquire lock for the second time
	result, err = c.lock.TryAcquireLock(ctx, LOCK1, 3000)
	assert.Nil(t, err)
	assert.False(t, result)

	// Release the lock
	err = c.lock.ReleaseLock(ctx, LOCK1)
	assert.Nil(t, err)

	// Try to acquire lock for the third time
	result, err = c.lock.TryAcquireLock(ctx, LOCK1, 3000)
	assert.Nil(t, err)
	assert.True(t, result)

	c.lock.ReleaseLock(ctx, LOCK1)
}

func (c *LockFixture) TestAcquireLock(t *testing.T) {

	ctx := context.Background()

	// Acquire lock for the first time
	err := c.lock.AcquireLock(ctx, LOCK2, 3000, 1000)
	assert.Nil(t, err)

	// Acquire lock for the second time
	err = c.lock.AcquireLock(ctx, LOCK2, 3000, 1000)
	assert.NotNil(t, err)

	// Release the lock
	err = c.lock.ReleaseLock(ctx, LOCK2)
	assert.Nil(t, err)

	// Acquire lock for the third time
	err = c.lock.AcquireLock(ctx, LOCK2, 3000, 1000)
	assert.Nil(t, err)

	c.lock.ReleaseLock(ctx, LOCK2)
}

func (c *LockFixture) TestReleaseLock(t *testing.T) {
	ctx := context.Background()

	// Acquire lock for the first time
	result, err := c.lock.TryAcquireLock(ctx, LOCK3, 3000)
	assert.Nil(t, err)
	assert.True(t, result)

	// Release the lock for the first time
	err = c.lock.ReleaseLock(ctx, LOCK3)
	assert.Nil(t, err)
	// Release the lock for the second time
	err = c.lock.ReleaseLock(ctx, LOCK3)
	assert.Nil(t, err)
}
