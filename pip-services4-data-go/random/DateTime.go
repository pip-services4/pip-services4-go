package random

import (
	"time"
)

// DateTime random generator for Date time values.
//	Example:
//		value1 := random.DateTime.NextDate(time.Parse(shortForm, "2007-Jan-01"), time.Parse(shortForm, "2010-Jan-01"));    // Possible result: 2008-01-03
//		value2 := random.DateTime.NextDateTime(time.Parse(shortForm, "2006-Jan-01"), time.Parse(shortForm, "2017-Jan-01"));// Possible result: 2007-03-11 11:20:32
//		value3 := random.DateTime.UpdateDateTime(time.Parse(shortForm, "2010-Jan-01"), );// Possible result: 2010-02-05 11:33:23
var DateTime = &TRandomDateTime{}

type TRandomDateTime struct{}

// NextDate generates a random Date in the range ['minYear', 'maxYear'].
// This method generate dates without time (or time set to 00:00:00)
//	Parameters:
//		- min: time.Time - minimum range value
//		- max: time.Time - maximum range value
//	Returns: time.Time  - a random Date value.
func (c *TRandomDateTime) NextDate(min time.Time, max time.Time) time.Time {
	diff := max.Unix() - min.Unix()
	if diff <= 0 {
		return min
	}

	diff = Long.Next(0, diff)
	date := min.Add(time.Duration(diff) * time.Second)

	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, min.Location())
}

// NextDateTime generates a random Date and time in the range ['minYear', 'maxYear'].
// This method generate dates without time (or time set to 00:00:00)
//	Parameters:
//		- min: time.Time  minimum range value
// 		- max: time.Time - maximum range value
//	Returns: time.Time a random Date and Time value.
func (c *TRandomDateTime) NextDateTime(min time.Time, max time.Time) time.Time {
	diff := max.Unix() - min.Unix()
	if diff <= 0 {
		return min
	}

	diff = Long.Next(0, diff)
	return min.Add(time.Duration(diff) * time.Second)
}

// UpdateDateTime updates (drifts) a Date value within specified range defined
//	Parameters:
//		- value: time.Time - value to drift.
//		- interval: int64 - a range in milliseconds. Default: 10 days
//	Returns: time.Time Date value.
func (c *TRandomDateTime) UpdateDateTime(value time.Time, interval int64) time.Time {
	if interval == 0 {
		interval = 10 * 24 * 3600
	}
	if interval < 0 {
		return value
	}

	// Days to milliseconds
	unixTime := value.Unix() + Long.Next(-interval, interval)
	return time.Unix(unixTime, 0)
}
