package reporter

import (
	"bytes"
	"fmt"
	"html/template"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Automation Report: {{.Name}}</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800&family=Fira+Code:wght@400;500&display=swap" rel="stylesheet">
  
  <style>
    :root {
      --bg-darker: #0c0e12;
      --bg-dark: #121620;
      --bg-panel: #1b202e;
      --border-color: #2b3247;
      --text-main: #f3f4f6;
      --text-muted: #9ca3af;
      --primary: #7c3aed;
      --accent-green: #10b981;
      --accent-red: #ef4444;
      --accent-cyan: #06b6d4;
      --font-main: 'Outfit', sans-serif;
      --font-code: 'Fira Code', monospace;
    }

    * { box-sizing: border-box; margin: 0; padding: 0; }

    body {
      background-color: var(--bg-darker);
      color: var(--text-main);
      font-family: var(--font-main);
      padding: 40px 24px;
      line-height: 1.5;
    }

    .container {
      max-width: 1200px;
      margin: 0 auto;
      display: flex;
      flex-direction: column;
      gap: 32px;
    }

    header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 20px;
    }

    .title-group h1 {
      font-size: 28px;
      font-weight: 800;
      letter-spacing: -0.5px;
      background: linear-gradient(135deg, #ffffff, #9ca3af);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    .title-group p {
      font-size: 13px;
      color: var(--text-muted);
      margin-top: 4px;
    }

    .badge {
      padding: 6px 12px;
      border-radius: 12px;
      font-size: 12px;
      font-weight: 700;
      text-transform: uppercase;
    }

    .badge.passed { background-color: rgba(16, 185, 129, 0.15); border: 1px solid var(--accent-green); color: var(--accent-green); }
    .badge.failed { background-color: rgba(239, 68, 68, 0.15); border: 1px solid var(--accent-red); color: var(--accent-red); }

    /* Summary Card Grid */
    .summary-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 20px;
    }

    .summary-card {
      background-color: var(--bg-dark);
      border: 1px solid var(--border-color);
      border-radius: 12px;
      padding: 20px;
      display: flex;
      flex-direction: column;
      gap: 6px;
    }

    .summary-card-label {
      font-size: 12px;
      font-weight: 600;
      color: var(--text-muted);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    .summary-card-val {
      font-size: 28px;
      font-weight: 800;
    }

    .summary-card-val.passed { color: var(--accent-green); }
    .summary-card-val.failed { color: var(--accent-red); }

    /* Test List Panel */
    .section-title {
      font-size: 18px;
      font-weight: 700;
      color: var(--text-muted);
      margin-bottom: 16px;
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .test-list {
      display: flex;
      flex-direction: column;
      gap: 16px;
    }

    .test-item {
      background-color: var(--bg-dark);
      border: 1px solid var(--border-color);
      border-radius: 12px;
      overflow: hidden;
      transition: border-color 0.2s;
    }

    .test-item:hover {
      border-color: var(--primary);
    }

    .test-header {
      padding: 16px 24px;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: space-between;
      user-select: none;
    }

    .test-info {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .test-title {
      font-weight: 600;
      font-size: 16px;
    }

    .test-meta {
      font-size: 12px;
      color: var(--text-muted);
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .test-details {
      padding: 24px;
      border-top: 1px solid var(--border-color);
      background-color: var(--bg-panel);
      display: none;
    }

    .test-details.open {
      display: block;
    }

    .error-box {
      background-color: rgba(239, 68, 68, 0.08);
      border: 1px solid rgba(239, 68, 68, 0.3);
      padding: 16px;
      border-radius: 8px;
      margin-bottom: 24px;
    }

    .error-title {
      font-size: 13px;
      font-weight: 700;
      color: var(--accent-red);
      margin-bottom: 6px;
    }

    .error-message {
      font-family: var(--font-code);
      font-size: 12px;
      white-space: pre-wrap;
      color: #fca5a5;
    }

    /* Steps and Timeline */
    .steps-container {
      margin-bottom: 24px;
    }

    .step-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 10px 0;
      border-bottom: 1px solid rgba(255,255,255,0.05);
      font-size: 13px;
    }

    .step-item:last-child {
      border-bottom: none;
    }

    .step-name {
      font-family: var(--font-code);
      font-weight: 500;
    }

    .step-duration {
      font-size: 11px;
      color: var(--text-muted);
    }

    /* Embedded Screenshot link */
    .screenshot-link {
      color: var(--accent-cyan);
      text-decoration: none;
      font-weight: 600;
      cursor: pointer;
    }

    .screenshot-preview {
      margin-top: 12px;
      max-width: 100%;
      border: 2px solid var(--border-color);
      border-radius: 8px;
      display: none;
    }

    /* SVG Performance Chart */
    .chart-container {
      margin-top: 24px;
      border-top: 1px solid var(--border-color);
      padding-top: 20px;
    }

    .chart-title {
      font-size: 13px;
      font-weight: 600;
      color: var(--text-muted);
      margin-bottom: 12px;
    }

    .chart-svg {
      width: 100%;
      height: 120px;
      background-color: var(--bg-darker);
      border: 1px solid var(--border-color);
      border-radius: 8px;
    }

    /* Log line output */
    .logs-console {
      background-color: var(--bg-darker);
      border: 1px solid var(--border-color);
      border-radius: 8px;
      padding: 16px;
      font-family: var(--font-code);
      font-size: 11px;
      max-height: 200px;
      overflow-y: auto;
      color: #98c379;
    }

    .log-row { margin-bottom: 4px; }
    .log-row.err { color: var(--accent-red); }
    .log-row.info { color: var(--accent-cyan); }
  </style>
</head>
<body>

  <div class="container">
    <!-- Header -->
    <header>
      <div class="title-group">
        <h1>{{.Name}}</h1>
        <p>Execution Summary: Finished at {{.EndTime.Format "2006-01-02 15:04:05 MST"}}</p>
      </div>
      <div>
        {{if gt .FailCount 0}}
        <span class="badge failed">FAIL</span>
        {{else}}
        <span class="badge passed">PASS</span>
        {{end}}
      </div>
    </header>

    <!-- Summary Statistics -->
    <div class="summary-grid">
      <div class="summary-card">
        <span class="summary-card-label">Total Tests</span>
        <span class="summary-card-val">{{.TotalCount}}</span>
      </div>
      <div class="summary-card">
        <span class="summary-card-label">Passed</span>
        <span class="summary-card-val passed">{{.PassCount}}</span>
      </div>
      <div class="summary-card">
        <span class="summary-card-label">Failed</span>
        <span class="summary-card-val failed">{{.FailCount}}</span>
      </div>
      <div class="summary-card">
        <span class="summary-card-label">Duration</span>
        <span class="summary-card-val">{{.Duration}}</span>
      </div>
    </div>

    <!-- Test List -->
    <div>
      <div class="section-title">
        <span>📋</span> Test Execution Details
      </div>
      
      <div class="test-list">
        {{range $index, $test := .Tests}}
        <div class="test-item">
          <!-- Accordion Header -->
          <div class="test-header" onclick="toggleDetails(this)">
            <div class="test-info">
              <span class="badge {{$test.Status}}">{{$test.Status}}</span>
              <span class="test-title">{{$test.Name}}</span>
            </div>
            <div class="test-meta">
              <span>⏱️ {{$test.Duration}}</span>
              <span>▼</span>
            </div>
          </div>
          
          <!-- Accordion Details -->
          <div class="test-details">
            {{if $test.ErrorMessage}}
            <div class="error-box">
              <div class="error-title">Failure Exception</div>
              <div class="error-message">{{$test.ErrorMessage}}</div>
              {{if $test.StackTrace}}
              <div class="error-message" style="margin-top:8px; opacity:0.8; font-size:11px;">{{$test.StackTrace}}</div>
              {{end}}
            </div>
            {{end}}

            <!-- Steps timeline -->
            {{if $test.Steps}}
            <div class="steps-container">
              <div class="chart-title">Test Steps Timeline</div>
              {{range $step := $test.Steps}}
              <div class="step-item">
                <span class="step-name">🔹 {{$step.Name}}</span>
                <div style="display: flex; align-items: center; gap: 12px;">
                  <span class="step-duration">⏱️ {{$step.Duration}}</span>
                  {{if $step.Screenshot}}
                  <span class="screenshot-link" onclick="toggleScreenshot(this, event)">📷 Screenshot</span>
                  <img class="screenshot-preview" src="data:image/png;base64,{{$step.Screenshot}}" />
                  {{end}}
                </div>
              </div>
              {{end}}
            </div>
            {{end}}

            <!-- Console logs -->
            {{if $test.Logs}}
            <div class="steps-container">
              <div class="chart-title">Diagnostic Logs</div>
              <div class="logs-console">
                {{range $log := $test.Logs}}
                <div class="log-row {{$log.Level}}">
                  [{{$log.Timestamp.Format "15:04:05"}}] [{{$log.Level}}] {{$log.Message}}
                </div>
                {{end}}
              </div>
            </div>
            {{end}}

            <!-- Metrics Chart -->
            {{if $test.Metrics}}
            <div class="chart-container">
              <div class="chart-title">Hardware Telemetry (CPU % & RAM MB)</div>
              <svg class="chart-svg" viewBox="0 0 1000 120" preserveAspectRatio="none">
                <!-- Grid Lines -->
                <line x1="0" y1="30" x2="1000" y2="30" stroke="#1f2430" stroke-width="1" />
                <line x1="0" y1="60" x2="1000" y2="60" stroke="#1f2430" stroke-width="1" />
                <line x1="0" y1="90" x2="1000" y2="90" stroke="#1f2430" stroke-width="1" />
                
                <!-- CPU Polyline (Violet) -->
                <polyline fill="none" stroke="var(--primary)" stroke-width="3" points="
                  {{range $i, $m := $test.Metrics}}{{multiply $i 100}},{{cpuY $m.CPUPercent}} {{end}}
                " />

                <!-- RAM Polyline (Cyan) -->
                <polyline fill="none" stroke="var(--accent-cyan)" stroke-width="2" stroke-dasharray="4,4" points="
                  {{range $i, $m := $test.Metrics}}{{multiply $i 100}},{{ramY $m.RAMMB}} {{end}}
                " />
              </svg>
              <div style="display:flex; gap: 16px; margin-top: 8px; font-size: 11px; color: var(--text-muted);">
                <div style="display:flex; align-items:center; gap:6px;">
                  <div style="width:12px; height:4px; background-color:var(--primary);"></div> CPU Usage (%)
                </div>
                <div style="display:flex; align-items:center; gap:6px;">
                  <div style="width:12px; height:4px; border-bottom:2px dashed var(--accent-cyan);"></div> Memory Usage (MB)
                </div>
              </div>
            </div>
            {{end}}
          </div>
        </div>
        {{end}}
      </div>
    </div>
  </div>

  <script>
    function toggleDetails(header) {
      const details = header.nextElementSibling;
      const isOpen = details.classList.contains('open');
      
      // Close all other panels
      document.querySelectorAll('.test-details').forEach(el => el.classList.remove('open'));
      
      if (!isOpen) {
        details.classList.add('open');
      }
    }

    function toggleScreenshot(link, event) {
      event.stopPropagation(); // Avoid triggering accordion close
      const img = link.nextElementSibling;
      const isVisible = img.style.display === 'block';
      img.style.display = isVisible ? 'none' : 'block';
    }
  </script>
</body>
</html>`

// GenerateHTML outputs interactive dashboard strings.
func GenerateHTML(result *SuiteResult) ([]byte, error) {
	// Setup template helper funcs for SVG mapping
	funcMap := template.FuncMap{
		"multiply": func(i int, multiplier int) int {
			return i * multiplier
		},
		"cpuY": func(percent float64) int {
			// Percent from 0-100 mapped to SVG height (120-10)
			return int(110.0 - (percent * 0.9))
		},
		"ramY": func(mb float64) int {
			// RAM from 0-1000 mapped to SVG height (120-10)
			val := mb / 10.0
			if val > 100.0 {
				val = 100.0
			}
			return int(110.0 - (val * 0.9))
		},
	}

	tmpl, err := template.New("html_report").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTML report: %w", err)
	}

	return buf.Bytes(), nil
}
