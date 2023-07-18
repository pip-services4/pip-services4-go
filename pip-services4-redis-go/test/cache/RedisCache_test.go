package test_cache

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	rediscache "github.com/pip-services4/pip-services4-go/pip-services4-redis-go/cache"
	redisfixture "github.com/pip-services4/pip-services4-go/pip-services4-redis-go/test/fixture"
)

func TestRedisCache(t *testing.T) {
	ctx := context.Background()

	var cache *rediscache.RedisCache[any]
	var fixture *redisfixture.CacheFixture

	host := os.Getenv("REDIS_SERVICE_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_SERVICE_PORT")
	if port == "" {
		port = "6379"
	}

	cache = rediscache.NewRedisCache[any]()
	config := cconf.NewConfigParamsFromTuples(
		"connection.host", host,
		"connection.port", port,
	)
	cache.Configure(ctx, config)
	fixture = redisfixture.NewCacheFixture(cache)
	cache.Open(ctx)
	defer cache.Close(ctx)

	t.Run("TestRedisCache:Store and Retrieve", fixture.TestStoreAndRetrieve)
	t.Run("TestRedisCache:Retrieve Expired", fixture.TestRetrieveExpired)
	t.Run("TestRedisCache:Remove", fixture.TestRemove)
}
