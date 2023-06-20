package test_lock

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-logic-go/lock"
)

func newMemoryLockFixture() *LockFixture {
	locker := lock.NewMemoryLock()
	fixture := NewLockFixture(locker)
	return fixture
}

func TestMemoryLockTryAcquireLock(t *testing.T) {
	fixture := newMemoryLockFixture()
	fixture.TestTryAcquireLock(t)
}

func TestMemoryLockAcquireLock(t *testing.T) {
	fixture := newMemoryLockFixture()
	fixture.TestAcquireLock(t)
}

func TestMemoryLockReleaseLock(t *testing.T) {
	fixture := newMemoryLockFixture()
	fixture.TestReleaseLock(t)
}
