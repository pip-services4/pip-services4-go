package test

import (
	"context"
	"errors"
	"testing"
	"time"

	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	clog "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	"github.com/stretchr/testify/assert"
)

type LoggerFixture struct {
	logger *clog.CachedLogger
}

func NewLoggerFixture(logger *clog.CachedLogger) *LoggerFixture {
	return &LoggerFixture{
		logger: logger,
	}
}

func (c *LoggerFixture) TestLogLevel(t *testing.T) {
	assert.True(t, c.logger.Level() >= clog.LevelNone)
	assert.True(t, c.logger.Level() <= clog.LevelTrace)
}

func (c *LoggerFixture) TestSimpleLogging(t *testing.T) {
	ctx := context.Background()
	c.logger.SetLevel(clog.LevelTrace)

	c.logger.Fatal(ctx, nil, "Fatal error message")
	c.logger.Error(ctx, nil, "Error message")
	c.logger.Warn(ctx, "Warning message")
	c.logger.Info(ctx, "Information message")
	c.logger.Debug(ctx, "Debug message")
	c.logger.Trace(ctx, "Trace message")

	c.logger.Dump(ctx)
	time.Sleep(1000 * time.Millisecond)
}

func (c *LoggerFixture) TestErrorLogging(t *testing.T) {
	ctx := cctx.NewContextWithTraceId(context.Background(), "123")
	var ex error = errors.New("Testing error throw")
	c.logger.Fatal(ctx, ex, "Fatal error")
	c.logger.Error(ctx, ex, "Recoverable error")

	assert.NotNil(t, ex)

	c.logger.Dump(ctx)
	time.Sleep(1000 * time.Millisecond)
}
