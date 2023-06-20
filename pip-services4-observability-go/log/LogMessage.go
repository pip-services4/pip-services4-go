package log

import (
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
)

// LogMessage Data object to store captured log messages.
// This object is used by CachedLogger.
type LogMessage struct {
	Time    time.Time               `json:"time"`
	Source  string                  `json:"source"`
	Level   LevelType               `json:"level"`
	TraceId string                  `json:"trace_id"`
	Error   errors.ErrorDescription `json:"error"`
	Message string                  `json:"message"`
}

// NewLogMessage create new log message object
//
//	Parameters:
//		- level LevelType a log level
//		- source string a source
//		- traceId string transaction id to trace execution through call chain.
//		- err errors.ErrorDescription an error object associated with this message.
//		- message string a human-readable message to log.
//	Returns: LogMessage
func NewLogMessage(traceId string, level LevelType, source string,
	err errors.ErrorDescription, message string) LogMessage {
	return LogMessage{
		Time:    time.Now().UTC(),
		Source:  source,
		Level:   level,
		TraceId: traceId,
		Error:   err,
		Message: message,
	}
}
