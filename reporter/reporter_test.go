package reporter

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestReportersExport(t *testing.T) {
	now := time.Now()

	// Mock suite result
	suite := &SuiteResult{
		Name:       "Login and Registration Suite",
		StartTime:  now,
		EndTime:    now.Add(12 * time.Second),
		Duration:   12 * time.Second,
		TotalCount: 2,
		PassCount:  1,
		FailCount:  1,
		Tests: []TestResult{
			{
				Name:     "Should login successfully with valid admin credentials",
				Status:   StatusPassed,
				Duration: 5 * time.Second,
				Logs: []LogEntry{
					{Timestamp: now.Add(1 * time.Second), Level: "info", Message: "Connecting to Android device..."},
					{Timestamp: now.Add(2 * time.Second), Level: "info", Message: "App launched successfully"},
					{Timestamp: now.Add(4 * time.Second), Level: "success", Message: "Login button clicked"},
				},
				Steps: []TestStep{
					{Name: "Launch App", Duration: 1 * time.Second, Status: StatusPassed},
					{Name: "Fill Username", Duration: 2 * time.Second, Status: StatusPassed},
					{Name: "Click Login", Duration: 1 * time.Second, Status: StatusPassed, Screenshot: "base64_png_placeholder"},
				},
				Metrics: []MetricPoint{
					{Timestamp: now.Add(1 * time.Second), CPUPercent: 12.5, RAMMB: 128.0},
					{Timestamp: now.Add(3 * time.Second), CPUPercent: 45.0, RAMMB: 256.0},
					{Timestamp: now.Add(5 * time.Second), CPUPercent: 22.0, RAMMB: 180.0},
				},
			},
			{
				Name:         "Should display validation errors for empty password",
				Status:       StatusFailed,
				Duration:     7 * time.Second,
				ErrorMessage: "AssertionError: expected error message 'Password is required' to be visible",
				StackTrace:   "at main.go:82\ncontext.go:204",
				Logs: []LogEntry{
					{Timestamp: now.Add(6 * time.Second), Level: "info", Message: "Entering admin username"},
					{Timestamp: now.Add(10 * time.Second), Level: "err", Message: "Validation banner not found"},
				},
				Steps: []TestStep{
					{Name: "Launch App", Duration: 1 * time.Second, Status: StatusPassed},
					{Name: "Fill Username", Duration: 2 * time.Second, Status: StatusPassed},
					{Name: "Click Submit", Duration: 4 * time.Second, Status: StatusFailed, Screenshot: "error_screenshot_png"},
				},
				Metrics: []MetricPoint{
					{Timestamp: now.Add(6 * time.Second), CPUPercent: 8.0, RAMMB: 120.0},
					{Timestamp: now.Add(9 * time.Second), CPUPercent: 60.0, RAMMB: 380.0},
				},
			},
		},
	}

	// 1. JSON Export
	jsonBytes, err := ExportToJSON(suite)
	if err != nil {
		t.Fatalf("JSON Export failed: %v", err)
	}
	if len(jsonBytes) == 0 {
		t.Error("Generated JSON report is empty")
	}
	if !strings.Contains(string(jsonBytes), "Login and Registration Suite") {
		t.Error("JSON report does not contain suite name")
	}

	// 2. JUnit XML Export
	xmlBytes, err := ExportToJUnit(suite)
	if err != nil {
		t.Fatalf("JUnit Export failed: %v", err)
	}
	if len(xmlBytes) == 0 {
		t.Error("Generated JUnit XML report is empty")
	}
	if !bytes.HasPrefix(xmlBytes, []byte("<?xml")) {
		t.Error("JUnit report does not have XML header prefix")
	}
	if !strings.Contains(string(xmlBytes), "<failure message=") {
		t.Error("JUnit report does not contain failure message tag")
	}

	// 3. HTML Dashboard Export
	htmlBytes, err := GenerateHTML(suite)
	if err != nil {
		t.Fatalf("HTML Export failed: %v", err)
	}
	if len(htmlBytes) == 0 {
		t.Error("Generated HTML report is empty")
	}
	if !strings.Contains(string(htmlBytes), "<polyline fill=\"none\" stroke=\"var(--primary)\"") {
		t.Error("HTML report does not contain CPU polyline chart tags")
	}
	if !strings.Contains(string(htmlBytes), "expected error message &#39;Password is required&#39;") {
		t.Error("HTML report does not contain failure message details")
	}
}
