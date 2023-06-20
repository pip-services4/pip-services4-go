package log

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// LevelConverter Helper class to convert log level values.
var LevelConverter = &_TLogLevelConverter{}

type _TLogLevelConverter struct{}

// ToLogLevel converts numbers and strings to standard log level values.
//
//	Parameters: value any a value to be converted
//	Returns LevelType converted log level
func (c *_TLogLevelConverter) ToLogLevel(value any) LevelType {
	return logLevelFromString(value)
}

// ToString converts log level to a string.
// see LevelType
//
//	Parameters: level int a log level to convert
//	Returns: string log level name string.
func (c *_TLogLevelConverter) ToString(level LevelType) string {
	return logLevelToString(level)
}

// LogLevelFromString converts log level to a LogLevel.
// see LevelType
//
//	Parameters: value any a log level string to convert
//	Returns: int log level value.
func logLevelFromString(value any) LevelType {
	if value == nil {
		return LevelInfo
	}

	str := convert.StringConverter.ToString(value)
	str = strings.ToUpper(str)
	if "0" == str || "NOTHING" == str || "NONE" == str {
		return LevelNone
	} else if "1" == str || "FATAL" == str {
		return LevelFatal
	} else if "2" == str || "ERROR" == str {
		return LevelError
	} else if "3" == str || "WARN" == str || "WARNING" == str {
		return LevelWarn
	} else if "4" == str || "INFO" == str {
		return LevelInfo
	} else if "5" == str || "DEBUG" == str {
		return LevelDebug
	} else if "6" == str || "TRACE" == str {
		return LevelTrace
	} else {
		return LevelInfo
	}
}

// LogLevelToString converts log level to a string.
//
//	see LogLevel
//	Parameters:
//		- level LevelType a log level to convert
//	Returns string log level name string.
func logLevelToString(level LevelType) string {
	switch level {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	default:
		return "UNDEF"
	}
}
