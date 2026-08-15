package core

// MetricsCollector tracks operational telemetry variables (counters, gauges, durations).
type MetricsCollector interface {
	Component

	// CounterInc increments a named counter metric.
	CounterInc(name string, labelNames []string, labelValues []string)

	// GaugeSet sets a metric value representing specific states.
	GaugeSet(name string, value float64, labelNames []string, labelValues []string)

	// HistogramRecord logs execution timings or sizes against range buckets.
	HistogramRecord(name string, value float64, labelNames []string, labelValues []string)
}
