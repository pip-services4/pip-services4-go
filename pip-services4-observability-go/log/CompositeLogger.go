package log

import (
	"context"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
)

// CompositeLogger aggregates all loggers from component references under a single component.
// It allows logging messages and conveniently send them to multiple destinations.
//
//	References:
//		- *:logger:*:*:1.0 (optional) ILogger components to pass log messages
//
//	see ILogger
//	Example:
//		type MyComponent {
//			_logger CompositeLogger
//		}
//		func (mc* MyComponent) Configure(ctx context.Context, config ConfigParams) {
//			mc._logger.Configure(ctx, config)
//			...
//		}
//
//		func (mc* MyComponent) SetReferences(ctx context.Context, references IReferences) {
//			mc._logger.SetReferences(ctx, references)
//			...
//		}
//
//		func (mc* MyComponent) myMethod(ctx context.Context) {
//			mc._logger.Debug(ctx context.Context, "Called method mycomponent.mymethod")
//			...
//		}
//		var mc MyComponent = MyComponent{}
//		mc._logger = NewCompositeLogger()
type CompositeLogger struct {
	*Logger
	loggers []ILogger
}

// NewCompositeLogger creates a new instance of the logger.
//
//	Returns: *CompositeLogger
func NewCompositeLogger() *CompositeLogger {
	c := &CompositeLogger{
		loggers: []ILogger{},
	}
	c.Logger = InheritLogger(c)
	c.SetLevel(LevelTrace)
	return c
}

// NewCompositeLoggerFromReferences creates a new instance of the logger.
//
//	Parameters:
//		- ctx context.Context
//		- refer.IReferences references to locate the component dependencies.
//	Returns: CompositeLogger
func NewCompositeLoggerFromReferences(ctx context.Context, references refer.IReferences) *CompositeLogger {
	c := NewCompositeLogger()
	c.SetReferences(ctx, references)
	return c
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- ctx context.Context
//		- refer.IReferences references to locate the component dependencies.
func (c *CompositeLogger) SetReferences(ctx context.Context, references refer.IReferences) {
	c.Logger.SetReferences(ctx, references)

	if c.loggers == nil {
		c.loggers = []ILogger{}
	}

	loggers := references.GetOptional(
		refer.NewDescriptor("*", "logger", "*", "*", "*"),
	)
	for _, l := range loggers {
		if l == c {
			continue
		}

		if logger, ok := l.(ILogger); ok {
			c.loggers = append(c.loggers, logger)
		}
	}
}

// Writes a log message to the logger destination(s).
// Parameters:
//   - ctx context.Context
//   - level LogLevel a log level.
//   - ctx context.Context execution context to trace execution through call chain.
//   - err error an error object associated with this message.
//   - message string a human-readable message to log.
func (c *CompositeLogger) Write(ctx context.Context, level LevelType, err error, message string) {
	if c.loggers == nil && len(c.loggers) == 0 {
		return
	}

	for _, logger := range c.loggers {
		logger.Log(ctx, level, err, message)
	}
}
