package test_log

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"

	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

type cachedLoggerSaver struct {
	counter uint32
}

func (c *cachedLoggerSaver) Save(ctx context.Context, messages []log.LogMessage) error {
	c.counter += uint32(len(messages))
	return nil
}

func (c *cachedLoggerSaver) Write(ctx context.Context, level log.LevelType, err error, message string) {

}

var saver = &cachedLoggerSaver{}

func newCustomCachedLogger() *log.CachedLogger {
	logger := log.InheritCachedLogger(saver)
	logger.Configure(
		context.Background(),
		config.NewConfigParamsFromTuples(
			log.ConfigParameterOptionsInterval, 100,
			log.ConfigParameterOptionsMaxCacheSize, 1,
		),
	)
	return logger
}

func newCachedLoggerFixture() *LoggerFixture {
	logger := newCustomCachedLogger()
	fixture := NewLoggerFixture(logger)
	return fixture
}

func TestCachedLogLevel(t *testing.T) {
	fixture := newCachedLoggerFixture()
	fixture.TestLogLevel(t)
}

func TestCachedSimpleLogging(t *testing.T) {
	fixture := newCachedLoggerFixture()
	fixture.TestSimpleLogging(t)
}

func TestCachedErrorLogging(t *testing.T) {
	fixture := newCachedLoggerFixture()
	fixture.TestErrorLogging(t)
}
