package build

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/trace"
)

// Creates ICounters components by their descriptors.

var NullCountersDescriptor = refer.NewDescriptor("pip-services", "counters", "null", "*", "1.0")
var LogCountersDescriptor = refer.NewDescriptor("pip-services", "counters", "log", "*", "1.0")
var CompositeCountersDescriptor = refer.NewDescriptor("pip-services", "counters", "composite", "*", "1.0")

var NullLoggerDescriptor = refer.NewDescriptor("pip-services", "logger", "null", "*", "1.0")
var ConsoleLoggerDescriptor = refer.NewDescriptor("pip-services", "logger", "console", "*", "1.0")
var CompositeLoggerDescriptor = refer.NewDescriptor("pip-services", "logger", "composite", "*", "1.0")

var NullTracerDescriptor = refer.NewDescriptor("pip-services", "tracer", "null", "*", "1.0")
var LogTracerDescriptor = refer.NewDescriptor("pip-services", "tracer", "log", "*", "1.0")
var CompositeTracerDescriptor = refer.NewDescriptor("pip-services", "tracer", "composite", "*", "1.0")

// NewDefaultObservabilityFactory Creates observability components by their descriptors.
//
//	Returns: *build.Factory
func NewDefaultObservabilityFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(NullCountersDescriptor, count.NewNullCounters)
	factory.RegisterType(LogCountersDescriptor, count.NewLogCounters)
	factory.RegisterType(CompositeCountersDescriptor, count.NewCompositeCounters)

	factory.RegisterType(NullLoggerDescriptor, log.NewNullLogger)
	factory.RegisterType(ConsoleLoggerDescriptor, log.NewConsoleLogger)
	factory.RegisterType(CompositeLoggerDescriptor, log.NewCompositeLogger)

	factory.RegisterType(NullTracerDescriptor, trace.NewNullTracer)
	factory.RegisterType(LogTracerDescriptor, trace.NewLogTracer)
	factory.RegisterType(CompositeTracerDescriptor, trace.NewCompositeTracer)

	return factory
}
