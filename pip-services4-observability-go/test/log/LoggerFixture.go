package test_log

import (
	"context"
	"errors"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	"github.com/stretchr/testify/assert"
)

type LoggerFixture struct {
	logger log.ILogger
}

func NewLoggerFixture(logger log.ILogger) *LoggerFixture {
	return &LoggerFixture{
		logger: logger,
	}
}

func (c *LoggerFixture) TestLogLevel(t *testing.T) {
	assert.True(t, c.logger.Level() >= log.LevelNone)
	assert.True(t, c.logger.Level() <= log.LevelTrace)
}

func (c *LoggerFixture) TestSimpleLogging(t *testing.T) {
	c.logger.SetLevel(log.LevelTrace)

	c.logger.Fatal(context.Background(), nil, "Fatal error message")
	c.logger.Error(context.Background(), nil, "Error message")
	c.logger.Warn(context.Background(), "Warning message")
	c.logger.Info(context.Background(), "Information message")
	c.logger.Debug(context.Background(), "Debug message")
	c.logger.Trace(context.Background(), "Trace message")
}

func (c *LoggerFixture) TestErrorLogging(t *testing.T) {
	err := errors.New("Test error")

	c.logger.Fatal(context.Background(), err, "Fatal error")
	c.logger.Error(context.Background(), err, "Recoverable error")
}
