package count

import (
	"strings"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	ccount "github.com/pip-services4/pip-services4-go/pip-services4-observability-go/count"
)

// PrometheusCounterConverter is helper class that converts performance counter values into
// a response from Prometheus metrics service.
var PrometheusCounterConverter TPrometheusCounterConverter = TPrometheusCounterConverter{}

type TPrometheusCounterConverter struct {
}

// ToString method converts the given counters to a string that is returned by Prometheus metrics service.
//
//	Parameters:
//		- counters  a list of counters to convert.
//		- source    a source (context) name.
//		- instance  a unique instance name (usually a host name).
//
// Returns string
// string view of counter
func (c *TPrometheusCounterConverter) ToString(counters []ccount.Counter, source string, instance string) string {

	if len(counters) == 0 {
		return ""
	}

	var builder string = ""

	for _, counter := range counters {
		counterName := c.parseCounterName(counter)
		labels := c.generateCounterLabel(counter, source, instance)

		switch counter.Type {
		case ccount.Increment:
			builder += "# TYPE " + counterName + " gauge\n"
			builder += counterName + labels + " " + cconv.StringConverter.ToString(counter.Count) + "\n"
		case ccount.Interval:
			builder += "# TYPE " + counterName + "_max gauge\n"
			builder += counterName + "_max" + labels + " " + cconv.StringConverter.ToString(counter.Max) + "\n"
			builder += "# TYPE " + counterName + "_min gauge\n"
			builder += counterName + "_min" + labels + " " + cconv.StringConverter.ToString(counter.Min) + "\n"
			builder += "# TYPE " + counterName + "_average gauge\n"
			builder += counterName + "_average" + labels + " " + cconv.StringConverter.ToString(counter.Average) + "\n"
			builder += "# TYPE " + counterName + "_count gauge\n"
			builder += counterName + "_count" + labels + " " + cconv.StringConverter.ToString(counter.Count) + "\n"
		case ccount.LastValue:
			builder += "# TYPE " + counterName + " gauge\n"
			builder += counterName + labels + " " + cconv.StringConverter.ToString(counter.Last) + "\n"
		case ccount.Statistics:
			builder += "# TYPE " + counterName + "_max gauge\n"
			builder += counterName + "_max" + labels + " " + cconv.StringConverter.ToString(counter.Max) + "\n"
			builder += "# TYPE " + counterName + "_min gauge\n"
			builder += counterName + "_min" + labels + " " + cconv.StringConverter.ToString(counter.Min) + "\n"
			builder += "# TYPE " + counterName + "_average gauge\n"
			builder += counterName + "_average" + labels + " " + cconv.StringConverter.ToString(counter.Average) + "\n"
			builder += "# TYPE " + counterName + "_count gauge\n"
			builder += counterName + "_count" + labels + " " + cconv.StringConverter.ToString(counter.Count) + "\n"
		case ccount.Timestamp: // Prometheus doesn't support non-numeric metrics
			builder += "# TYPE " + counterName + " gauge\n" //" untyped\n"
			builder += counterName + labels + " " + cconv.StringConverter.ToString(counter.Time.Unix()) + "\n"
		}
	}

	return builder
}

func (c *TPrometheusCounterConverter) AtomicCountersToCounters(atomicCounters []*ccount.AtomicCounter) []ccount.Counter {
	counters := make([]ccount.Counter, len(atomicCounters))

	for _, atomicCounter := range atomicCounters {
		counter := ccount.Counter{
			Name:    atomicCounter.Name(),
			Type:    atomicCounter.Type(),
			Last:    atomicCounter.Last(),
			Count:   atomicCounter.Count(),
			Min:     atomicCounter.Min(),
			Max:     atomicCounter.Max(),
			Average: atomicCounter.Average(),
			Time:    atomicCounter.Time(),
		}

		counters = append(counters, counter)
	}

	return counters

}

func (c *TPrometheusCounterConverter) generateCounterLabel(counter ccount.Counter, source string, instance string) string {

	labels := make(map[string]string, 0)

	if source != "" {
		labels["source"] = source
	}

	if instance != "" {
		labels["instance"] = instance
	}

	nameParts := strings.Split(counter.Name, ".")

	// If there are other predictable names from which we can parse labels, we can add them below
	if len(nameParts) >= 3 && nameParts[2] == "exec_time" {
		labels["service"] = nameParts[0]
		labels["command"] = nameParts[1]
	}

	if len(labels) == 0 {
		return ""
	}

	builder := "{"
	for key := range labels {
		if len(builder) > 1 {
			builder += ","
		}
		builder += key + `="` + labels[key] + `"`
	}
	builder += "}"

	return builder
}

func (c *TPrometheusCounterConverter) parseCounterName(counter ccount.Counter) string {
	if counter.Name == "" {
		return ""
	}

	nameParts := strings.Split(counter.Name, ".")

	// If there are other predictable names from which we can parse labels, we can add them below
	if len(nameParts) >= 3 && nameParts[2] == "exec_time" {
		return nameParts[2]
	}

	// TODO: are there other assumptions we can make?
	// Or just return as a single, valid name
	result := strings.ToLower(counter.Name)
	result = strings.Replace(result, ".", "_", -1)
	result = strings.Replace(result, "/", "_", -1)

	return result
}

func (c *TPrometheusCounterConverter) parseCounterLabels(counter ccount.Counter, source string, instance string) interface{} {
	labels := make(map[string]string, 0)

	if source != "" {
		labels["source"] = source
	}

	if instance != "" {
		labels["instance"] = instance
	}

	nameParts := strings.Split(counter.Name, ".")

	// If there are other predictable names from which we can parse labels, we can add them below
	if len(nameParts) >= 3 && nameParts[2] == "exec_time" {
		labels["service"] = nameParts[0]
		labels["command"] = nameParts[1]
	}

	return labels
}
