package count

import (
	"time"
)

// Counter data object to store measurement for a performance counter.
// This object is used by CachedCounters to store counters.
type Counter struct {
	Name    string      `json:"name"`
	Type    CounterType `json:"type"`
	Last    float64     `json:"last"`
	Count   int64       `json:"count"`
	Min     float64     `json:"min"`
	Max     float64     `json:"max"`
	Average float64     `json:"average"`
	Time    time.Time   `json:"time"`
}
