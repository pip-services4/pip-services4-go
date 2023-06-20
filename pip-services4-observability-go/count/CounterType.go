package count

import "encoding/json"

type CounterType uint8

//	Types of counters that measure different types of metrics
//	Interval: = 0 Counters that measure execution time intervals
//	LastValue: = 1 Counters that keeps the latest measured value
//	Statistics: = 2 Counters that measure min/average/max statistics
//	Timestamp: = 3 Counter that record timestamps
//	Increment: = 4 Counter that increment counters
const (
	Interval   CounterType = 0
	LastValue  CounterType = 1
	Statistics CounterType = 2
	Timestamp  CounterType = 3
	Increment  CounterType = 4
)

// ToString method converting counter type to string
func (c CounterType) ToString() string {
	name := ""

	switch c {
	case Interval:
		name = "interval"
	case LastValue:
		name = "lastvalue"
	case Statistics:
		name = "statistics"
	case Timestamp:
		name = "timestamp"
	case Increment:
		name = "increment"
	}

	return name
}

// NewCounterTypeFromString creates new CounterType from string
func NewCounterTypeFromString(value string) CounterType {
	switch value {
	case "interval":
		return Interval
	case "lastvalue":
		return LastValue
	case "statistics":
		return Statistics
	case "timestamp":
		return Timestamp
	case "increment":
		return Increment
	}
	return Interval
}

func (c *CounterType) UnmarshalJSON(data []byte) (err error) {
	var result string
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	*c = NewCounterTypeFromString(result)
	return
}

func (c CounterType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ToString())
}
