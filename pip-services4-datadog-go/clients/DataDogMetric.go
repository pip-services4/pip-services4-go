package clients

type DataDogMetric struct {
	Metric   string               `json:"metric"`
	Service  string               `json:"service"`
	Host     string               `json:"host"`
	Tags     map[string]string    `json:"tags"`
	Type     string               `json:"type"`
	Interval int64                `json:"interval"`
	Points   []DataDogMetricPoint `json:"points"`
}
