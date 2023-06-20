package test_log

import (
	"context"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

func newCompositeLoggerFixture() *LoggerFixture {
	logger := log.NewCompositeLogger()

	refs := refer.NewReferencesFromTuples(
		context.Background(),
		refer.NewDescriptor("pip-services", "logger", "console", "default", "1.0"), log.NewConsoleLogger(),
		//log.CompositeLoggerDescriptor, logger,
	)
	logger.SetReferences(context.Background(), refs)

	fixture := NewLoggerFixture(logger)
	return fixture
}

func TestCompositeLogLevel(t *testing.T) {
	fixture := newCompositeLoggerFixture()
	fixture.TestLogLevel(t)
}

func TestCompositeSimpleLogging(t *testing.T) {
	fixture := newCompositeLoggerFixture()
	fixture.TestSimpleLogging(t)
}

func TestCompositeErrorLogging(t *testing.T) {
	fixture := newCompositeLoggerFixture()
	fixture.TestErrorLogging(t)
}
