package clients

import "time"

type DataDogMetricPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}
