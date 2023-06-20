package log

import "context"

// ILogger for logger components that capture execution log messages.
type ILogger interface {

	// Level gets the maximum log level. Messages with higher log level are filtered out.
	Level() LevelType

	// SetLevel set the maximum log level.
	SetLevel(value LevelType)

	// Log logs a message at specified log level.
	Log(ctx context.Context, level LevelType, err error, message string, args ...any)

	// Fatal logs fatal (unrecoverable) message that caused the process to crash.
	Fatal(ctx context.Context, err error, message string, args ...any)

	// Error logs recoverable application error.
	Error(ctx context.Context, err error, message string, args ...any)

	// Warn logs a warning that may or may not have a negative impact.
	Warn(ctx context.Context, message string, args ...any)

	// Info logs an important information message
	Info(ctx context.Context, message string, args ...any)

	// Debug logs a high-level debug information for troubleshooting.
	Debug(ctx context.Context, message string, args ...any)

	// Trace logs a low-level debug information for troubleshooting.
	Trace(ctx context.Context, message string, args ...any)
}
