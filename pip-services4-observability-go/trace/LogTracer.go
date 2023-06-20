package trace

import (
	"context"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	log "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// LogTracer tracer that dumps recorded traces to logger.
//
//	Configuration parameters:
//		- options:
//			- log_level: log level to record traces (default: debug)
//
//	References:
//		- *:logger:*:*:1.0       [[ILogger]] components to dump the captured counters
//		- *:context-info:*:*:1.0 (optional) [[ContextInfo]] to detect the context id and specify counters source
//
//	See Tracer
//	See CachedCounters
//	See CompositeLogger
//
//	Example:
//		tracer = NewLogTracer();
//		tracer.SetReferences(
//			context.Background(),
//			NewReferencesFromTuples(
//				NewDescriptor("pip-services", "logger", "console", "default", "1.0"), NewConsoleLogger()
//			)
//		);
//		timing := trcer.BeginTrace(context.Background(), "123", "mycomponent", "mymethod");
//		...
//		timing.EndTrace(context.Background());
//		if err != nil {
//			timing.EndFailure(context.Background(), err);
//		}
type LogTracer struct {
	logger   *log.CompositeLogger
	logLevel log.LevelType
}

// NewLogTracer creates a new instance of the tracer.
func NewLogTracer() *LogTracer {
	return &LogTracer{
		logger:   log.NewCompositeLogger(),
		logLevel: log.LevelDebug,
	}
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config configuration parameters to be set.
func (c *LogTracer) Configure(ctx context.Context, config *cconf.ConfigParams) {
	logLvl, ok := config.GetAsObject("options.log_level")
	if ok && logLvl == nil {
		logLvl = c.logLevel
	}
	c.logLevel = log.LevelConverter.ToLogLevel(logLvl)
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- references 	references to locate the component dependencies.
func (c *LogTracer) SetReferences(ctx context.Context, references cref.IReferences) {
	c.logger.SetReferences(ctx, references)
}

func (c *LogTracer) logTrace(ctx context.Context, component string, operation string, err error, duration int64) {
	builder := ""

	if err != nil {
		builder += "Failed to execute "
	} else {
		builder += "Executed "
	}

	builder += component
	builder += "."
	builder += operation

	if duration > 0 {
		builder += " in " + cconv.StringConverter.ToString(duration) + " msec"
	}

	if err != nil {
		c.logger.Error(ctx, err, builder)
	} else {
		c.logger.Log(ctx, c.logLevel, nil, builder)
	}
}

// Trace records an operation trace with its name and duration
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//		- duration      execution duration in milliseconds.
func (c *LogTracer) Trace(ctx context.Context, component string, operation string, duration int64) {
	c.logTrace(ctx, component, operation, nil, duration)
}

// Failure records an operation failure with its name, duration and error
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component         a name of called component
//		- operation         a name of the executed operation.
//		- error             an error object associated with this trace.
//		- duration          execution duration in milliseconds.
func (c *LogTracer) Failure(ctx context.Context, component string, operation string, err error, duration int64) {
	c.logTrace(ctx, component, operation, err, duration)
}

// BeginTrace begins recording an operation trace
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//	Returns: a trace timing object.
func (c *LogTracer) BeginTrace(ctx context.Context, component string, operation string) *TraceTiming {
	return NewTraceTiming(ctx, component, operation, c)
}
