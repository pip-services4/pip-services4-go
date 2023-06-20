package log

import (
	"context"
	"sync"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
)

// ICachedLogSaver abstract logger that caches captured log messages
// in memory and periodically dumps them. Child classes implement
// saving cached messages to their specified destinations.
//
//	Configuration parameters
//		- level: maximum log level to capture
//		- source: source (context) name
//		- options:
//			- interval: interval in milliseconds to save log messages (default: 10 seconds)
//			- max_cache_size: maximum number of messages stored in this cache (default: 100)
//	References:
//		- *:context-info:*:*:1.0 (optional) ContextInfo to detect the context id and specify counters source
type ICachedLogSaver interface {
	Save(ctx context.Context, messages []LogMessage) error
}

type ICachedLoggerOverrides interface {
	ILoggerOverrides
	Save(ctx context.Context, messages []LogMessage) error
}

type CachedLogger struct {
	Logger
	Cache        []LogMessage
	Updated      bool
	LastDumpTime time.Time
	MaxCacheSize int
	Interval     int
	mtx          *sync.Mutex
	Overrides    ICachedLoggerOverrides
}

const (
	DefaultMaxCacheSize                = 100
	DefaultInterval                    = 10000
	ConfigParameterOptionsInterval     = "options.interval"
	ConfigParameterOptionsMaxCacheSize = "options.max_cache_size"
)

// InheritCachedLogger creates a new instance of the logger from ICachedLogSaver
//
//	Parameters:
//		- overrides ICachedLoggerOverrides
//	Returns: CachedLogger
func InheritCachedLogger(overrides ICachedLoggerOverrides) *CachedLogger {
	c := &CachedLogger{
		Cache:        make([]LogMessage, 0, DefaultMaxCacheSize),
		Updated:      false,
		LastDumpTime: time.Now(),
		MaxCacheSize: DefaultMaxCacheSize,
		Interval:     DefaultInterval,
		mtx:          &sync.Mutex{},
		Overrides:    overrides,
	}
	c.Logger = *InheritLogger(overrides)
	return c
}

// Configure configures component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration parameters to be set.
func (c *CachedLogger) Configure(ctx context.Context, cfg *config.ConfigParams) {
	c.Logger.Configure(ctx, cfg)

	c.Interval = cfg.GetAsIntegerWithDefault(ConfigParameterOptionsInterval, c.Interval)
	c.MaxCacheSize = cfg.GetAsIntegerWithDefault(ConfigParameterOptionsMaxCacheSize, c.MaxCacheSize)
	c.Cache = make([]LogMessage, 0, c.MaxCacheSize)
}

// Writes a log message to the logger destination.
// Parameters:
//   - ctx context.Context
//   - level LogLevel a log level.
//   - ctx context.Context execution context to trace execution through call chain.
//   - err error an error object associated with this message.
//   - message string a human-readable message to log.
func (c *CachedLogger) Write(ctx context.Context, level LevelType, err error, message string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	logMessage := LogMessage{
		Time:    time.Now().UTC(),
		Level:   level,
		Source:  c.source,
		Message: message,
		TraceId: utils.ContextHelper.GetTraceId(ctx),
	}

	if err != nil {
		errorDescription := errors.NewErrorDescription(err)
		logMessage.Error = *errorDescription
	}

	c.Cache = append(c.Cache, logMessage)

	c.update(ctx)
}

// Clear (removes) all cached log messages.
//
//	Parameters:
//		- ctx context.Context
func (c *CachedLogger) Clear(ctx context.Context) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.Cache = make([]LogMessage, 0, c.MaxCacheSize)
	c.Updated = false
}

// Dump (writes) the currently cached log messages.
//
//	Parameters:
//		- ctx context.Context
func (c *CachedLogger) Dump(ctx context.Context) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.dump(ctx)
}

// Update makes message cache as updated and dumps it when timeout expires.
//
//	Parameters:
//		- ctx context.Context
func (c *CachedLogger) Update(ctx context.Context) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.update(ctx)
}

func (c *CachedLogger) update(ctx context.Context) {
	c.Updated = true

	elapsed := int(time.Since(c.LastDumpTime).Seconds() * 1000)

	if elapsed > c.Interval {
		// Todo: Decide what to do with the error
		_ = c.dump(ctx)
	}
}

func (c *CachedLogger) dump(ctx context.Context) error {
	if c.Updated {
		if !c.Updated {
			return nil
		}

		messages := c.Cache
		c.Cache = make([]LogMessage, 0, c.MaxCacheSize)

		err := c.Overrides.Save(ctx, messages)
		if err != nil {

			// Put failed messages back to cache
			c.Cache = append(messages, c.Cache...)

			// Truncate cache to max size
			if len(c.Cache) > c.MaxCacheSize {
				c.Cache = c.Cache[len(c.Cache)-c.MaxCacheSize:]
			}

		}

		c.Updated = false
		c.LastDumpTime = time.Now()
		return err
	}
	return nil
}
