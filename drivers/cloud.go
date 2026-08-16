package drivers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CloudGridDriver connects to any W3C Appium-compatible cloud provider (BrowserStack, SauceLabs, LambdaTest).
type CloudGridDriver struct {
	provider  string
	gridURL   string
	sessionID string
	client    *http.Client
	mu        sync.Mutex
}

// NewCloudGridDriver constructs a CloudGridDriver for a specific cloud provider.
func NewCloudGridDriver(provider string) *CloudGridDriver {
	return &CloudGridDriver{
		provider: strings.ToLower(provider),
		client: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

// Connect starts a remote session on the target cloud provider using W3C standards.
func (d *CloudGridDriver) Connect(ctx context.Context, capabilities map[string]string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 1. Resolve remote cloud hub URL
	switch d.provider {
	case "browserstack":
		d.gridURL = capabilities["browserstack_url"]
		if d.gridURL == "" {
			user := capabilities["browserstack.user"]
			key := capabilities["browserstack.key"]
			if user == "" || key == "" {
				user = capabilities["username"]
				key = capabilities["access_key"]
			}
			d.gridURL = fmt.Sprintf("https://%s:%s@hub-cloud.browserstack.com/wd/hub", user, key)
		}
	case "saucelabs":
		d.gridURL = capabilities["saucelabs_url"]
		if d.gridURL == "" {
			user := capabilities["username"]
			key := capabilities["access_key"]
			region := capabilities["region"]
			if region == "eu" {
				d.gridURL = fmt.Sprintf("https://%s:%s@ondemand.eu-central-1.saucelabs.com/wd/hub", user, key)
			} else {
				d.gridURL = fmt.Sprintf("https://%s:%s@ondemand.us-west-1.saucelabs.com/wd/hub", user, key)
			}
		}
	case "lambdatest":
		d.gridURL = capabilities["lambdatest_url"]
		if d.gridURL == "" {
			user := capabilities["username"]
			key := capabilities["access_key"]
			d.gridURL = fmt.Sprintf("https://%s:%s@hub.lambdatest.com/wd/hub", user, key)
		}
	default:
		d.gridURL = capabilities["grid_url"]
	}

	d.gridURL = strings.TrimSuffix(d.gridURL, "/")

	// 2. Prepare standardized W3C capabilities payload
	capsMap := make(map[string]interface{})
	for k, v := range capabilities {
		if k != "username" && k != "access_key" && k != "browserstack.key" && k != "browserstack.user" {
			capsMap[k] = v
		}
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"capabilities": map[string]interface{}{
			"alwaysMatch": capsMap,
		},
	})

	resp, err := d.client.Post(d.gridURL+"/session", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to establish session with %s grid: %w", d.provider, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s grid connection rejected (status %d): %s", d.provider, resp.StatusCode, string(body))
	}

	var res struct {
		SessionID string `json:"sessionId"`
		Value     struct {
			SessionID string `json:"sessionId"`
		} `json:"value"`
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(bodyBytes, &res)

	d.sessionID = res.SessionID
	if d.sessionID == "" {
		d.sessionID = res.Value.SessionID
	}

	if d.sessionID == "" {
		return fmt.Errorf("failed to extract session ID from %s response: %s", d.provider, string(bodyBytes))
	}

	return nil
}

func (d *CloudGridDriver) Disconnect(ctx context.Context) error {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

	if sessionID == "" {
		return nil
	}

	req, _ := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/session/%s", gridURL, sessionID), nil)
	resp, err := d.client.Do(req)
	if err == nil {
		resp.Body.Close()
	}

	d.mu.Lock()
	d.sessionID = ""
	d.mu.Unlock()

	return nil
}

func (d *CloudGridDriver) GetSource(ctx context.Context, format string) (string, error) {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

	resp, err := d.client.Get(fmt.Sprintf("%s/session/%s/source", gridURL, sessionID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Value string `json:"value"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	return res.Value, nil
}

func (d *CloudGridDriver) ClickAt(ctx context.Context, x, y int) error {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"x": x,
		"y": y,
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/actions", gridURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *CloudGridDriver) Click(ctx context.Context, elementID string) error {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/element/%s/click", gridURL, sessionID, elementID), "application/json", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *CloudGridDriver) Fill(ctx context.Context, elementID string, value string) error {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

	// W3C standard requires keyboard typing value character arrays
	reqBody, _ := json.Marshal(map[string]interface{}{
		"text": value,
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/element/%s/value", gridURL, sessionID, elementID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *CloudGridDriver) Swipe(ctx context.Context, startX, startY, endX, endY int, duration time.Duration) error {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

	// Dispatch W3C touch gestures action chains
	reqBody, _ := json.Marshal(map[string]interface{}{
		"actions": []interface{}{
			map[string]interface{}{
				"type": "pointer",
				"id":   "finger1",
				"actions": []interface{}{
					map[string]interface{}{"type": "pointerMove", "duration": 0, "x": startX, "y": startY},
					map[string]interface{}{"type": "pointerDown", "button": 0},
					map[string]interface{}{"type": "pointerMove", "duration": duration.Milliseconds(), "x": endX, "y": endY},
					map[string]interface{}{"type": "pointerUp", "button": 0},
				},
			},
		},
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/actions", gridURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *CloudGridDriver) Screenshot(ctx context.Context, elementID string) ([]byte, error) {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

	resp, err := d.client.Get(fmt.Sprintf("%s/session/%s/screenshot", gridURL, sessionID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Value string `json:"value"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)

	return io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(res.Value)))
}

func (d *CloudGridDriver) FindElement(ctx context.Context, strategy string, selector string) (string, error) {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

	wdaStrategy := "xpath"
	switch strings.ToUpper(strategy) {
	case "ACCESSIBILITY_ID":
		wdaStrategy = "accessibility id"
	case "CLASS_NAME", "CLASS":
		wdaStrategy = "class name"
	case "XPATH":
		wdaStrategy = "xpath"
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"using": wdaStrategy,
		"value": selector,
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/element", gridURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Value map[string]string `json:"value"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)

	for _, v := range res.Value {
		return v, nil
	}

	return "", fmt.Errorf("element not found on %s grid", d.provider)
}

func (d *CloudGridDriver) LaunchApp(ctx context.Context, appID string, args []string, env map[string]string) error {
	return nil
}

func (d *CloudGridDriver) TerminateApp(ctx context.Context, appID string) error {
	return nil
}

func (d *CloudGridDriver) InstallApp(ctx context.Context, appPath string) error {
	return errors.New("install app is not supported on cloud driver")
}

func (d *CloudGridDriver) UninstallApp(ctx context.Context, appID string) error {
	return errors.New("uninstall app is not supported on cloud driver")
}

func (d *CloudGridDriver) ExecuteScript(ctx context.Context, script string, arguments []string) (string, error) {
	return "", errors.New("execute script is not supported on cloud driver")
}

func (d *CloudGridDriver) InjectKeyevent(ctx context.Context, keycode int) error {
	return nil
}

func (d *CloudGridDriver) GetContexts(ctx context.Context) ([]string, error) {
	return []string{"NATIVE_APP"}, nil
}

func (d *CloudGridDriver) SetContext(ctx context.Context, name string) error {
	return nil
}

func (d *CloudGridDriver) GetTelemetry(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"cpu_usage":     0.0,
		"ram_usage_mb":  0.0,
		"battery_level": 100,
	}, nil
}

func (d *CloudGridDriver) SetMockLocation(ctx context.Context, latitude, longitude float64) error {
	return nil
}

func (d *CloudGridDriver) MockBiometrics(ctx context.Context, action string, enrollID int) error {
	return nil
}
