package test_lock

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	redislock "github.com/pip-services4/pip-services4-go/pip-services4-redis-go/lock"
	redisfixture "github.com/pip-services4/pip-services4-go/pip-services4-redis-go/test/fixture"
)

func TestRedisLock(t *testing.T) {
	var lock *redislock.RedisLock
	var fixture *redisfixture.LockFixture

	ctx := context.Background()

	host := os.Getenv("REDIS_SERVICE_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_SERVICE_PORT")
	if port == "" {
		port = "6379"
	}

	lock = redislock.NewRedisLock()

	config := cconf.NewConfigParamsFromTuples(
		"connection.host", host,
		"connection.port", port,
	)
	lock.Configure(ctx, config)
	fixture = redisfixture.NewLockFixture(lock)

	lock.Open(ctx)
	defer lock.Close(ctx)

	t.Run("Try Acquire Lock", fixture.TestTryAcquireLock)
	t.Run("Acquire Lock", fixture.TestAcquireLock)
	t.Run("Release Lock", fixture.TestReleaseLock)
}
