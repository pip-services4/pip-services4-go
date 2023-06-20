package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
)

// ConsoleLogger is a logger that writes log messages to console.
// Errors are written to standard err stream and all other messages to standard out stream.
//
//	Configuration parameters:
//		- level: maximum log level to capture
//		- source: source (context) name
//
//	References:
//		- *:context-info:*:*:1.0 (optional) ContextInfo to detect the context id and specify counters source
//	see Logger
//
//	Example:
//		logger = NewConsoleLogger()
//		logger.SetLevel(LogLevel.Debug)
//		logger.Error(context.Background(), "123", ex, "Error occured: %s", ex.message)
//		logger.Debug(context.Background(), "123", "Everything is OK.")
type ConsoleLogger struct {
	*Logger
}

// NewConsoleLogger creates a new instance of the logger.
//
//	Returns: ConsoleLogger
func NewConsoleLogger() *ConsoleLogger {
	c := &ConsoleLogger{}
	c.Logger = InheritLogger(c)
	return c
}

// Write a log message to the logger destination.
//
//	Parameters:
//		- ctx context.Context
//		- level LevelType a log level.
//		- ctx context.Context execution context to trace execution through call chain.
//		- err error an error object associated with this message.
//		- message string a human-readable message to log.
func (c *ConsoleLogger) Write(ctx context.Context, level LevelType, err error, message string) {
	if c.Level() < level {
		return
	}
	traceId := utils.ContextHelper.GetTraceId(ctx)

	if traceId == "" {
		traceId = "---"
	}
	levelStr := LevelConverter.ToString(level)
	dateStr := convert.StringConverter.ToString(time.Now().UTC())

	build := strings.Builder{}
	build.WriteString("[")
	build.WriteString(traceId)
	build.WriteString(":")
	build.WriteString(levelStr)
	build.WriteString(":")
	build.WriteString(dateStr)
	build.WriteString("] ")

	build.WriteString(message)

	if err != nil {
		if len(message) == 0 {
			build.WriteString("Error: ")
		} else {
			build.WriteString(": ")
		}

		build.WriteString(c.ComposeError(err))
	}

	build.WriteString("\n")
	output := build.String()

	if level == LevelFatal || level == LevelError || level == LevelWarn {
		fmt.Fprintln(os.Stderr, output)
	} else {
		fmt.Println(output)
	}
}
