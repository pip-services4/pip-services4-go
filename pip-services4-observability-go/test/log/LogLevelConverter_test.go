package test_log

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	"github.com/stretchr/testify/assert"
)

func TestLogLevelConverter(t *testing.T) {
	level := log.LevelConverter.ToLogLevel("info")
	assert.Equal(t, log.LevelInfo, level)

	level = log.LevelConverter.ToLogLevel("4")
	assert.Equal(t, log.LevelInfo, level)

	str := log.LevelConverter.ToString(level)
	assert.Equal(t, "INFO", str)
}
