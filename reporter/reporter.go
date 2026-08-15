package reporter

import (
	"time"
)

// TestStatus defines possible test outcomes.
type TestStatus string

const (
	StatusPassed  TestStatus = "passed"
	StatusFailed  TestStatus = "failed"
	StatusSkipped TestStatus = "skipped"
)

// LogEntry stores contextual debug outputs printed during execution.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

// MetricPoint samples hardware properties over durations.
type MetricPoint struct {
	Timestamp  time.Time `json:"timestamp"`
	CPUPercent float64   `json:"cpu_percent"`
	RAMMB      float64   `json:"ram_mb"`
}

// TestStep models individual script commands.
type TestStep struct {
	Name       string        `json:"name"`
	Duration   time.Duration `json:"duration"`
	Status     TestStatus    `json:"status"`
	Screenshot string        `json:"screenshot,omitempty"` // Base64 encoded PNG representation
}

// TestResult stores results for a single test case.
type TestResult struct {
	Name         string        `json:"name"`
	Status       TestStatus    `json:"status"`
	Duration     time.Duration `json:"duration"`
	ErrorMessage string        `json:"error_message,omitempty"`
	StackTrace   string        `json:"stack_trace,omitempty"`
	Logs         []LogEntry    `json:"logs"`
	Steps        []TestStep    `json:"steps"`
	Metrics      []MetricPoint `json:"metrics"`
}

// SuiteResult organizes collections of runs.
type SuiteResult struct {
	Name       string        `json:"name"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Duration   time.Duration `json:"duration"`
	Tests      []TestResult  `json:"tests"`
	TotalCount int           `json:"total_count"`
	PassCount  int           `json:"pass_count"`
	FailCount  int           `json:"fail_count"`
}
