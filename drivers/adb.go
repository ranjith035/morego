package drivers

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ADBXMLNode represents nodes returned by Android's uiautomator dump.
type ADBXMLNode struct {
	XMLName     xml.Name
	Class       string       `xml:"class,attr"`
	Text        string       `xml:"text,attr"`
	ResourceID  string       `xml:"resource-id,attr"`
	ContentDesc string       `xml:"content-desc,attr"`
	Bounds      string       `xml:"bounds,attr"`
	Nodes       []ADBXMLNode `xml:"node"`
}

// ADBDriver implements a concrete, usable Driver using command line adb calls.
type ADBDriver struct {
	deviceID string
	mu       sync.Mutex
	cache    map[string]ADBXMLNode // Map elementID -> node
	idGen    int
}

// NewADBDriver constructs an ADBDriver instance.
func NewADBDriver() *ADBDriver {
	return &ADBDriver{
		cache: make(map[string]ADBXMLNode),
	}
}

// Connect selects the target Android device ID from capabilities or auto-detects.
func (d *ADBDriver) Connect(ctx context.Context, capabilities map[string]string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	udid, ok := capabilities["udid"]
	if !ok {
		udid = capabilities["device_id"]
	}

	if udid == "" {
		// Auto-detect the first connected device
		cmd := exec.CommandContext(ctx, "adb", "devices")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run adb devices: %w", err)
		}

		lines := strings.Split(out.String(), "\r\n")
		if len(lines) <= 1 {
			lines = strings.Split(out.String(), "\n")
		}

		var firstDevice string
		for _, line := range lines[1:] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) >= 2 && parts[1] == "device" {
				firstDevice = parts[0]
				break
			}
		}

		if firstDevice == "" {
			return errors.New("no connected Android devices found via adb")
		}
		udid = firstDevice
	}

	d.deviceID = udid
	return nil
}

func (d *ADBDriver) Disconnect(ctx context.Context) error {
	return nil // Stateless command line operations require no teardown
}

func (d *ADBDriver) GetSource(ctx context.Context, format string) (string, error) {
	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	if deviceID == "" {
		return "", errors.New("adb driver is not connected")
	}

	// 1. Trigger hierarchy layout dump on device
	dumpCmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "shell", "uiautomator", "dump", "/data/local/tmp/window_dump.xml")
	if err := dumpCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to dump XML hierarchy via adb: %w", err)
	}

	// 2. Read layout file from device tmp
	readCmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "shell", "cat", "/data/local/tmp/window_dump.xml")
	var out bytes.Buffer
	readCmd.Stdout = &out
	if err := readCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read dumped XML hierarchy: %w", err)
	}

	xmlStr := out.String()
	// Strip out standard terminal progress warnings like "UI hierchary dumped to:..." if they exist
	if idx := strings.Index(xmlStr, "<?xml"); idx != -1 {
		xmlStr = xmlStr[idx:]
	}

	return xmlStr, nil
}

func (d *ADBDriver) ClickAt(ctx context.Context, x, y int) error {
	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	if deviceID == "" {
		return errors.New("adb driver is not connected")
	}

	cmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tap coordinates [%d, %d]: %w", x, y, err)
	}
	return nil
}

func (d *ADBDriver) Click(ctx context.Context, elementID string) error {
	d.mu.Lock()
	node, exists := d.cache[elementID]
	d.mu.Unlock()

	if !exists {
		return fmt.Errorf("element ID %s not found in locator cache", elementID)
	}

	x, y, err := parseBoundsCenter(node.Bounds)
	if err != nil {
		return fmt.Errorf("failed to parse bounds center for %s: %w", node.Bounds, err)
	}

	return d.ClickAt(ctx, x, y)
}

func (d *ADBDriver) Fill(ctx context.Context, elementID string, value string) error {
	// 1. Focus element
	err := d.Click(ctx, elementID)
	if err != nil {
		return err
	}

	// Sleep briefly to let keyboard state settle
	time.Sleep(200 * time.Millisecond)

	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	// Replace space characters with Android escape token %s
	escapedVal := strings.ReplaceAll(value, " ", "%s")

	cmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "shell", "input", "text", escapedVal)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to type text %q: %w", value, err)
	}
	return nil
}

func (d *ADBDriver) Swipe(ctx context.Context, startX, startY, endX, endY int, duration time.Duration) error {
	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	if deviceID == "" {
		return errors.New("adb driver is not connected")
	}

	durationMS := duration.Milliseconds()
	if durationMS <= 0 {
		durationMS = 200
	}

	cmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "shell", "input", "swipe",
		strconv.Itoa(startX), strconv.Itoa(startY), strconv.Itoa(endX), strconv.Itoa(endY), strconv.FormatInt(durationMS, 10))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("swipe action failed: %w", err)
	}
	return nil
}

func (d *ADBDriver) Screenshot(ctx context.Context, elementID string) ([]byte, error) {
	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	if deviceID == "" {
		return nil, errors.New("adb driver is not connected")
	}

	// Run screencap returning raw png to standard output stream
	cmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "exec-out", "screencap", "-p")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("screenshot capture failed: %w", err)
	}

	return out.Bytes(), nil
}

func (d *ADBDriver) FindElement(ctx context.Context, strategy string, selector string) (string, error) {
	xmlSource, err := d.GetSource(ctx, "xml")
	if err != nil {
		return "", err
	}

	var root ADBXMLNode
	err = xml.Unmarshal([]byte(xmlSource), &root)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal current XML: %w", err)
	}

	var nodes []ADBXMLNode
	collectADBNodes(root, &nodes)

	var matchedNode *ADBXMLNode
	for _, node := range nodes {
		match := false
		switch strings.ToUpper(strategy) {
		case "ACCESSIBILITY_ID", "TEST_ID":
			match = node.ContentDesc == selector || node.ResourceID == selector
		case "RESOURCE_ID":
			match = node.ResourceID == selector
		case "TEXT":
			match = node.Text == selector || strings.Contains(node.Text, selector)
		case "CLASS_NAME", "CLASS":
			match = node.Class == selector
		case "BOUNDS":
			match = node.Bounds == selector
		}

		if match {
			matchedNode = &node
			break
		}
	}

	if matchedNode == nil {
		return "", fmt.Errorf("element not found matching strategy %s with selector %s", strategy, selector)
	}

	d.mu.Lock()
	d.idGen++
	elementID := fmt.Sprintf("element_%d", d.idGen)
	d.cache[elementID] = *matchedNode
	d.mu.Unlock()

	return elementID, nil
}

func (d *ADBDriver) LaunchApp(ctx context.Context, appID string, args []string, env map[string]string) error {
	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	if deviceID == "" {
		return errors.New("adb driver is not connected")
	}

	// Trigger monkey starter packet
	cmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "shell", "monkey", "-p", appID, "-c", "android.intent.category.LAUNCHER", "1")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to launch app %s: %w", appID, err)
	}
	return nil
}

func (d *ADBDriver) TerminateApp(ctx context.Context, appID string) error {
	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	if deviceID == "" {
		return errors.New("adb driver is not connected")
	}

	cmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "shell", "am", "force-stop", appID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to terminate app %s: %w", appID, err)
	}
	return nil
}

func (d *ADBDriver) InstallApp(ctx context.Context, appPath string) error {
	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	if deviceID == "" {
		return errors.New("adb driver is not connected")
	}

	cmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "install", appPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install app package %s: %w", appPath, err)
	}
	return nil
}

func (d *ADBDriver) UninstallApp(ctx context.Context, appID string) error {
	d.mu.Lock()
	deviceID := d.deviceID
	d.mu.Unlock()

	if deviceID == "" {
		return errors.New("adb driver is not connected")
	}

	cmd := exec.CommandContext(ctx, "adb", "-s", deviceID, "uninstall", appID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall app package %s: %w", appID, err)
	}
	return nil
}

func (d *ADBDriver) ExecuteScript(ctx context.Context, script string, arguments []string) (string, error) {
	return "", errors.New("execute script is not supported on native adb driver")
}

// Helpers
func collectADBNodes(node ADBXMLNode, list *[]ADBXMLNode) {
	*list = append(*list, node)
	for _, child := range node.Nodes {
		collectADBNodes(child, list)
	}
}

// parseBoundsCenter reads Android's bounds attribute format "[x1,y1][x2,y2]"
func parseBoundsCenter(boundsStr string) (int, int, error) {
	r := regexp.MustCompile(`\[(\d+),(\d+)\]\[(\d+),(\d+)\]`)
	matches := r.FindStringSubmatch(boundsStr)
	if len(matches) < 5 {
		return 0, 0, fmt.Errorf("invalid bounds format %q", boundsStr)
	}

	x1, _ := strconv.Atoi(matches[1])
	y1, _ := strconv.Atoi(matches[2])
	x2, _ := strconv.Atoi(matches[3])
	y2, _ := strconv.Atoi(matches[4])

	centerX := x1 + (x2-x1)/2
	centerY := y1 + (y2-y1)/2

	return centerX, centerY, nil
}
