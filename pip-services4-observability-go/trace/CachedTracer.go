package trace

import (
	"context"
	"sync"
	"time"

	cerr "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	cconf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	cctx "github.com/pip-services4/pip-services4-go/pip-services4-components-go/context"
	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/utils"
)

// ICachedTraceSaver Abstract tracer that caches recorded traces in memory and periodically dumps them.
// Child classes implement saving cached traces to their specified destinations.
//
//	Configuration parameters:
//		- source:         source (context) name
//		- options:
//		- interval:       interval in milliseconds to save log messages (default: 10 seconds)
//		- maxcache_size:  maximum number of messages stored in this cache (default: 100)
//
//	References:
//		- *:context-info:*:*:1.0 (optional) [[ContextInfo]] to detect the context id and specify counters source
//
//	See ITracer
//	See OperationTrace
type ICachedTraceSaver interface {
	Save(ctx context.Context, operations []OperationTrace) error
}

type CachedTracer struct {
	source       string
	Cache        []OperationTrace
	updated      bool
	lastDumpTime time.Time
	maxCacheSize int
	interval     int64
	saver        ICachedTraceSaver
	mtx          *sync.Mutex
}

// InheritCachedTracer creates a new instance of the logger.
func InheritCachedTracer(saver ICachedTraceSaver) *CachedTracer {
	return &CachedTracer{
		Cache:        make([]OperationTrace, 0),
		updated:      false,
		lastDumpTime: time.Now().UTC(),
		maxCacheSize: 100,
		interval:     10000,
		saver:        saver,
	}
}

// Configure component by passing configuration parameters.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- config configuration parameters to be set.
func (c *CachedTracer) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.interval = config.GetAsLongWithDefault("options.interval", c.interval)
	c.maxCacheSize = config.GetAsIntegerWithDefault("options.maxcache_size", c.maxCacheSize)
	c.source = config.GetAsStringWithDefault("source", c.source)
}

// SetReferences references to dependent components.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- references references to locate the component dependencies.
func (c *CachedTracer) SetReferences(ctx context.Context, references cref.IReferences) {
	ref := references.GetOneOptional(
		cref.NewDescriptor(
			"pip-services",
			"context-info",
			"*", "*", "1.0",
		),
	)
	if ref != nil {
		if contextInfo, ok := ref.(*cctx.ContextInfo); ok && contextInfo != nil && c.source == "" {
			c.source = contextInfo.Name
		}
	}
}

// Writes a log message to the logger destination.
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//		- error         an error object associated with this trace.
//		- duration      execution duration in milliseconds.
func (c *CachedTracer) Write(ctx context.Context, component string, operation string, err error, duration int64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var errorDesc *cerr.ErrorDescription

	if err != nil {
		errorDesc = cerr.NewErrorDescription(err)
	}

	trace := OperationTrace{
		Time:      time.Now().UTC(),
		Source:    c.source,
		Component: component,
		Operation: operation,
		TraceId:   utils.ContextHelper.GetTraceId(ctx),
		Duration:  duration,
		Error:     *errorDesc,
	}
	c.Cache = append(c.Cache, trace)
	c.update(ctx)
}

// Trace records an operation trace with its name and duration
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//		- duration      execution duration in milliseconds.
func (c *CachedTracer) Trace(ctx context.Context, component string, operation string, duration int64) {
	c.Write(ctx, component, operation, nil, duration)
}

// Failure records an operation failure with its name, duration and error
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//		- error         an error object associated with this trace.
//		- duration      execution duration in milliseconds.
func (c *CachedTracer) Failure(ctx context.Context, component string, operation string, err error, duration int64) {
	c.Write(ctx, component, operation, err, duration)
}

// BeginTrace begins recording an operation trace
//
//	Parameters:
//		- ctx context.Context execution context to trace execution through call chain.
//		- component     a name of called component
//		- operation     a name of the executed operation.
//	Returns: a trace timing object.
func (c *CachedTracer) BeginTrace(ctx context.Context, component string, operation string) *TraceTiming {
	return NewTraceTiming(ctx, component, operation, c)
}

// Clear (removes) all cached log messages.
func (c *CachedTracer) Clear() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.Cache = make([]OperationTrace, 0)
	c.updated = false
}

// Dump (writes) the currently cached log messages.
//
//	See [[Write]]
func (c *CachedTracer) Dump(ctx context.Context) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.dump(ctx)
}

// Update makes trace cache as updated
// and dumps it when timeout expires.
//
//	See Dump
func (c *CachedTracer) Update(ctx context.Context) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.update(ctx)
}

func (c *CachedTracer) update(ctx context.Context) {
	c.updated = true
	elapsed := int64(time.Since(c.lastDumpTime).Seconds() * 1000)
	if elapsed > c.interval {
		// Todo: Decide what to do with the error
		c.Dump(ctx)
	}
}

func (c *CachedTracer) dump(ctx context.Context) {
	if c.updated {
		if !c.updated {
			return
		}

		traces := c.Cache
		c.Cache = make([]OperationTrace, 0)

		err := c.saver.Save(ctx, traces)

		if err != nil {
			// Adds traces back to the cache
			traces = append(traces, c.Cache...)
			c.Cache = traces

			// Truncate cache to max size
			if len(c.Cache) > c.maxCacheSize {
				c.Cache = c.Cache[len(c.Cache)-c.maxCacheSize:]
			}
		}

		c.updated = false
		c.lastDumpTime = time.Now().UTC()
	}
}
