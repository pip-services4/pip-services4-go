package log

import (
	"encoding/json"
)

// Standard log levels.
// Logs at debug and trace levels are usually captured only locally for
// troubleshooting and never sent to consolidated log services.
//	Log levels:
//		- None  = 0 Nothing to log
//		- Fatal = 1 Log only fatal errors that cause processes to crash
//		- Error = 2 Log all errors.
//		- Warn  = 3 Log errors and warnings
//		- Info  = 4 Log errors and important information messages
//		- Debug = 5 Log everything except traces
//		- Trace = 6 Log everything.
const (
	LevelNone  LevelType = 0
	LevelFatal LevelType = 1
	LevelError LevelType = 2
	LevelWarn  LevelType = 3
	LevelInfo  LevelType = 4
	LevelDebug LevelType = 5
	LevelTrace LevelType = 6
)

// LevelType is a type to represent log level names
type LevelType uint8

func (l *LevelType) UnmarshalJSON(data []byte) (err error) {
	var result string
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	*l = logLevelFromString(result)
	return
}

func (l LevelType) MarshalJSON() ([]byte, error) {
	return json.Marshal(logLevelToString(l))
}
