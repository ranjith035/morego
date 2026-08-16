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

// SeeTestDriver connects to Digital.ai SeeTest Cloud Grid endpoints.
type SeeTestDriver struct {
	gridURL   string
	sessionID string
	client    *http.Client
	mu        sync.Mutex
}

// NewSeeTestDriver constructs a SeeTestDriver.
func NewSeeTestDriver() *SeeTestDriver {
	return &SeeTestDriver{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Connect starts a remote session on SeeTest Cloud using accessKey capabilities.
func (d *SeeTestDriver) Connect(ctx context.Context, capabilities map[string]string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	gridURL, ok := capabilities["seetest_url"]
	if !ok {
		gridURL = "https://cloud.seetest.io/wd/hub" // Default SeeTest endpoint
	}
	d.gridURL = strings.TrimSuffix(gridURL, "/")

	accessKey, ok := capabilities["access_key"]
	if !ok {
		accessKey = capabilities["seetest_access_key"]
	}

	caps := map[string]interface{}{
		"capabilities": map[string]interface{}{
			"alwaysMatch": map[string]interface{}{
				"deviceQuery": capabilities["device_query"],
				"bundleId":    capabilities["bundle_id"],
				"accessKey":   accessKey,
			},
		},
	}
	reqBody, _ := json.Marshal(caps)

	resp, err := d.client.Post(d.gridURL+"/session", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to connect to SeeTest cloud: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("SeeTest grid connection rejected (status %d): %s", resp.StatusCode, string(body))
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
		return fmt.Errorf("failed to obtain session ID from SeeTest response: %s", string(bodyBytes))
	}

	return nil
}

func (d *SeeTestDriver) Disconnect(ctx context.Context) error {
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

func (d *SeeTestDriver) GetSource(ctx context.Context, format string) (string, error) {
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

func (d *SeeTestDriver) ClickAt(ctx context.Context, x, y int) error {
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

func (d *SeeTestDriver) Click(ctx context.Context, elementID string) error {
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

func (d *SeeTestDriver) Fill(ctx context.Context, elementID string, value string) error {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

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

func (d *SeeTestDriver) Swipe(ctx context.Context, startX, startY, endX, endY int, duration time.Duration) error {
	d.mu.Lock()
	sessionID := d.sessionID
	gridURL := d.gridURL
	d.mu.Unlock()

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

func (d *SeeTestDriver) Screenshot(ctx context.Context, elementID string) ([]byte, error) {
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

func (d *SeeTestDriver) FindElement(ctx context.Context, strategy string, selector string) (string, error) {
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

	return "", errors.New("element not found on SeeTest grid")
}

func (d *SeeTestDriver) LaunchApp(ctx context.Context, appID string, args []string, env map[string]string) error {
	return nil
}

func (d *SeeTestDriver) TerminateApp(ctx context.Context, appID string) error {
	return nil
}

func (d *SeeTestDriver) InstallApp(ctx context.Context, appPath string) error {
	return errors.New("install app is not supported on SeeTest cloud driver")
}

func (d *SeeTestDriver) UninstallApp(ctx context.Context, appID string) error {
	return errors.New("uninstall app is not supported on SeeTest cloud driver")
}

func (d *SeeTestDriver) ExecuteScript(ctx context.Context, script string, arguments []string) (string, error) {
	return "", errors.New("execute script is not supported on SeeTest cloud driver")
}

func (d *SeeTestDriver) InjectKeyevent(ctx context.Context, keycode int) error {
	return nil
}
