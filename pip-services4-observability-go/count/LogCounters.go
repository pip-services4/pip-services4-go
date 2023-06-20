package count

import (
	"context"
	"math"
	"sort"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/pip-services4/pip-services4-go/pip-services4-observability-go/log"
)

// LogCounters performance counters that periodically dumps counters measurements to logger.
//
//	Configuration parameters:
//		- options:
//			- interval: interval in milliseconds to save current counters measurements (default: 5 mins)
//			- reset_timeout: timeout in milliseconds to reset the counters. 0 disables the reset (default: 0)
//	References:
//		- *:logger:*:*:1.0 ILogger components to dump the captured counters
//		- *:context-info:*:*:1.0 (optional) ContextInfo to detect the context id and specify counters source
//
//	see Counter
//	see CachedCounters
//	see CompositeLogger
//
//	Example:
//		counters := NewLogCounters();
//		counters.SetReferences(context.Background(), NewReferencesFromTuples(
//			NewDescriptor("pip-services", "logger", "console", "default", "1.0"), NewConsoleLogger()
//		));
//		counters.IncrementOne(context.Background(), "mycomponent.mymethod.calls")
//		timing := counters.BeginTiming(context.Background(),"mycomponent.mymethod.exec_time")
//		defer timing.EndTiming(context.Background())
//
//		// do something
//		counters.Dump(context.Background())
type LogCounters struct {
	*CachedCounters
	logger *log.CompositeLogger
}

// NewLogCounters creates a new instance of the counters.
//
//	Returns: *LogCounters
func NewLogCounters() *LogCounters {
	c := &LogCounters{
		logger: log.NewCompositeLogger(),
	}
	c.CachedCounters = InheritCacheCounters(c)
	return c
}

// SetReferences sets references to dependent components.
//
//	Parameters:
//		- references refer.IReferences references to locate the component dependencies.
func (c *LogCounters) SetReferences(ctx context.Context, references refer.IReferences) {
	c.logger.SetReferences(ctx, references)
}

func (c *LogCounters) counterToString(counter Counter) string {
	result := "Counter " + counter.Name + " { "
	result = result + "\"type\": " + counter.Type.ToString()

	switch counter.Type {
	case Increment:
		result = result + ", \"count\": " + convert.StringConverter.ToString(counter.Count)
	case LastValue:
		result = result + ", \"last\": " + convert.StringConverter.ToString(counter.Last)
	case Timestamp:
		result = result + ", \"time\": " + convert.StringConverter.ToString(counter.Time)
	default:

		result = result + ", \"last\": " + convert.StringConverter.ToString(counter.Last)

		if counter.Count > 0 {
			result = result + ", \"count\": " + convert.StringConverter.ToString(counter.Count)
		}

		if counter.Min != math.MaxFloat32 {
			result = result + ", \"min\": " + convert.StringConverter.ToString(counter.Min)
		}

		if counter.Max != -math.MaxFloat32 {
			result = result + ", \"max\": " + convert.StringConverter.ToString(counter.Max)
		}

		result = result + ", \"avg\": " + convert.StringConverter.ToString(counter.Average)
	}

	result = result + " }"
	return result
}

// Save the current counters measurements.
//
//	Parameters:
//		- ctx context.Context
//		- counters []*Counter current counters measurements to be saves.
func (c *LogCounters) Save(ctx context.Context, counters []Counter) error {
	if len(counters) == 0 {
		return nil
	}

	sort.Slice(counters, func(i, j int) bool {
		return counters[i].Name < counters[j].Name
	})

	for _, counter := range counters {
		c.logger.Info(ctx, "", c.counterToString(counter))
	}

	return nil
}
