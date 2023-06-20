package log

import "context"

// NullLogger dummy implementation of logger that doesn't do anything.
// It can be used in testing or in situations when logger is required but shall be disabled.
type NullLogger struct{}

// NewNullLogger creates a new instance of the logger.
//	Returns *NullLogger
func NewNullLogger() *NullLogger {
	c := &NullLogger{}
	return c
}

// Level gets the maximum log level. Messages with higher log level are filtered out.
//	Returns: LevelType the maximum log level.
func (c *NullLogger) Level() LevelType {
	return LevelNone
}

// SetLevel set the maximum log level.
//	Parameters:
//		- value int a new maximum log level.
func (c *NullLogger) SetLevel(value LevelType) {
}

// Log a message at specified log level.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- level LevelType a log level.
//		- err error an error object associated with this message.
//		- message string a human-readable message to log.
//		- args ...any arguments to parameterize the message.
func (c *NullLogger) Log(ctx context.Context, level LevelType, err error, message string, args ...any) {
}

// Fatal logs fatal (unrecoverable) message that caused the process to crash.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- err error an error object associated with this message.
//		- message string a human-readable message to log.
//		- args ...any arguments to parameterize the message.
func (c *NullLogger) Fatal(ctx context.Context, err error, message string, args ...any) {
}

// Logs recoverable application error.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- err error an error object associated with this message.
//		- message string a human-readable message to log.
//		- args ...any arguments to parameterize the message.
func (c *NullLogger) Error(ctx context.Context, err error, message string, args ...any) {
}

// Warn logs a warning that may or may not have a negative impact.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- message string a human-readable message to log.
//		- args ...any arguments to parameterize the message.
func (c *NullLogger) Warn(ctx context.Context, message string, args ...any) {
}

// Info logs an important information message
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- message string a human-readable message to log.
//		- args ...any arguments to parameterize the message.
func (c *NullLogger) Info(ctx context.Context, message string, args ...any) {
}

// Debug logs a high-level debug information for troubleshooting.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- message string a human-readable message to log.
//		- args ...any arguments to parameterize the message.
func (c *NullLogger) Debug(ctx context.Context, message string, args ...any) {
}

// Trace logs a low-level debug information for troubleshooting.
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- message string a human-readable message to log.
//		- args ...any arguments to parameterize the message.
func (c *NullLogger) Trace(ctx context.Context, message string, args ...any) {
}
