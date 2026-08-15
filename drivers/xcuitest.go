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

// XCUITestDriver communicates directly with Apple's WebDriverAgent (WDA) server.
type XCUITestDriver struct {
	wdaURL    string
	sessionID string
	client    *http.Client
	mu        sync.Mutex
}

// NewXCUITestDriver constructs an XCUITestDriver.
func NewXCUITestDriver() *XCUITestDriver {
	return &XCUITestDriver{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// Connect dials the WDA server and initializes a driver session.
func (d *XCUITestDriver) Connect(ctx context.Context, capabilities map[string]string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	wdaURL, ok := capabilities["wda_url"]
	if !ok {
		wdaURL = "http://localhost:8100" // Default WebDriverAgent port
	}
	d.wdaURL = strings.TrimSuffix(wdaURL, "/")

	// Request WDA session creation
	reqBody, _ := json.Marshal(map[string]interface{}{
		"capabilities": map[string]interface{}{
			"alwaysMatch": map[string]interface{}{
				"bundleId": capabilities["bundle_id"],
			},
		},
	})

	resp, err := d.client.Post(d.wdaURL+"/session", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to reach WebDriverAgent at %s (is it running?): %w", d.wdaURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WDA returned non-ok status: %d", resp.StatusCode)
	}

	var res struct {
		SessionID string `json:"sessionId"`
		Value     struct {
			SessionID string `json:"sessionId"`
		} `json:"value"`
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(bodyBytes, &res)

	// Session ID can reside either in top level or inside W3C value wrapper
	d.sessionID = res.SessionID
	if d.sessionID == "" {
		d.sessionID = res.Value.SessionID
	}

	if d.sessionID == "" {
		return fmt.Errorf("failed to extract session ID from WDA response: %s", string(bodyBytes))
	}

	return nil
}

func (d *XCUITestDriver) Disconnect(ctx context.Context) error {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	if sessionID == "" {
		return nil
	}

	req, _ := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/session/%s", wdaURL, sessionID), nil)
	resp, err := d.client.Do(req)
	if err == nil {
		resp.Body.Close()
	}

	d.mu.Lock()
	d.sessionID = ""
	d.mu.Unlock()

	return nil
}

func (d *XCUITestDriver) GetSource(ctx context.Context, format string) (string, error) {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	resp, err := d.client.Get(fmt.Sprintf("%s/session/%s/source", wdaURL, sessionID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Value string `json:"value"`
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(bodyBytes, &res)

	return res.Value, nil
}

func (d *XCUITestDriver) ClickAt(ctx context.Context, x, y int) error {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"x": x,
		"y": y,
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/wda/tap/nil", wdaURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *XCUITestDriver) Click(ctx context.Context, elementID string) error {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/element/%s/click", wdaURL, sessionID, elementID), "application/json", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *XCUITestDriver) Fill(ctx context.Context, elementID string, value string) error {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	// WDA requires keyboard typing value array
	reqBody, _ := json.Marshal(map[string]interface{}{
		"value": strings.Split(value, ""),
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/element/%s/value", wdaURL, sessionID, elementID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *XCUITestDriver) Swipe(ctx context.Context, startX, startY, endX, endY int, duration time.Duration) error {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"fromX":    startX,
		"fromY":    startY,
		"toX":      endX,
		"toY":      endY,
		"duration": duration.Seconds(),
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/wda/dragfromtoforduration", wdaURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *XCUITestDriver) Screenshot(ctx context.Context, elementID string) ([]byte, error) {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	// Screen captures endpoint
	url := fmt.Sprintf("%s/session/%s/screenshot", wdaURL, sessionID)
	if elementID != "" {
		url = fmt.Sprintf("%s/session/%s/element/%s/screenshot", wdaURL, sessionID, elementID)
	}

	resp, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Value string `json:"value"` // base64 payload
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)

	return base64.StdEncoding.DecodeString(res.Value)
}

func (d *XCUITestDriver) FindElement(ctx context.Context, strategy string, selector string) (string, error) {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	// W3C Standard locating mappings
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

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/element", wdaURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("locator failed with status code %d", resp.StatusCode)
	}

	var res struct {
		Value map[string]string `json:"value"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)

	// W3C uses unique uuid element keys: "element-6066-11e4-a52e-4f735466cecf"
	for _, v := range res.Value {
		return v, nil
	}

	return "", errors.New("locator did not return element references")
}

func (d *XCUITestDriver) LaunchApp(ctx context.Context, appID string, args []string, env map[string]string) error {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"bundleId": appID,
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/wda/apps/launch", wdaURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *XCUITestDriver) TerminateApp(ctx context.Context, appID string) error {
	d.mu.Lock()
	sessionID := d.sessionID
	wdaURL := d.wdaURL
	d.mu.Unlock()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"bundleId": appID,
	})

	resp, err := d.client.Post(fmt.Sprintf("%s/session/%s/wda/apps/terminate", wdaURL, sessionID), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (d *XCUITestDriver) InstallApp(ctx context.Context, appPath string) error {
	return errors.New("install app is not supported natively via WebDriverAgent. Use Xcode toolchains instead.")
}

func (d *XCUITestDriver) UninstallApp(ctx context.Context, appID string) error {
	return errors.New("uninstall app is not supported natively via WebDriverAgent. Use Xcode toolchains instead.")
}

func (d *XCUITestDriver) ExecuteScript(ctx context.Context, script string, arguments []string) (string, error) {
	return "", errors.New("execute script is not supported natively via WebDriverAgent")
}
