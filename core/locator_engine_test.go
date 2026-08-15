package core

import (
	"context"
	"strings"
	"testing"
	"time"
)

// Mock ActionEngine implementation for testing
type mockActionEngine struct {
	xmlSource string
}

func (m *mockActionEngine) Name() string                       { return "ActionEngine" }
func (m *mockActionEngine) Init(ctx context.Context) error     { return nil }
func (m *mockActionEngine) Shutdown(ctx context.Context) error { return nil }

func (m *mockActionEngine) Click(ctx context.Context, sessionID string, locator *Locator) error {
	return nil
}

func (m *mockActionEngine) ClickAt(ctx context.Context, sessionID string, pt Point) error {
	return nil
}

func (m *mockActionEngine) Fill(ctx context.Context, sessionID string, locator *Locator, text string) error {
	return nil
}

func (m *mockActionEngine) Swipe(ctx context.Context, sessionID string, start, end Point, duration time.Duration) error {
	return nil
}

func (m *mockActionEngine) Screenshot(ctx context.Context, sessionID string, locator *Locator) ([]byte, error) {
	return []byte("mock_screenshot"), nil
}

func (m *mockActionEngine) GetSource(ctx context.Context, sessionID string, format string) (string, error) {
	return m.xmlSource, nil
}

func (m *mockActionEngine) ExecuteScript(ctx context.Context, sessionID string, script string, args []string) (string, error) {
	return "mock_script_result", nil
}

// Mock DI Container implementation for testing
type mockContainer struct {
	components map[string]interface{}
}

func (m *mockContainer) Register(name string, value interface{}) {
	m.components[name] = value
}

func (m *mockContainer) Resolve(name string) (interface{}, bool) {
	val, ok := m.components[name]
	return val, ok
}

// Test XML parser
func TestParseXML(t *testing.T) {
	xmlSource := `
	<hierarchy rotation="0">
		<node class="android.widget.FrameLayout" bounds="[0,0][1080,1920]" id="root_id">
			<node class="android.widget.Button" text="Submit" resource-id="com.example:id/submit_btn" content-desc="Submit description" bounds="[100,200][300,300]" id="btn_1"/>
			<node class="android.widget.EditText" text="" placeholder="Enter Username" bounds="[100,400][300,500]" test-id="username_field" id="edit_1"/>
			<node class="android.widget.TextView" text="Username Label" bounds="[100,100][300,150]" id="label_1"/>
		</node>
	</hierarchy>`

	root, err := ParseXML(strings.NewReader(xmlSource))
	if err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}

	if root == nil {
		t.Fatal("Parsed root node is nil")
	}

	if root.TagName != "hierarchy" {
		t.Errorf("Expected root tag name 'hierarchy', got %q", root.TagName)
	}

	if len(root.Children) != 1 {
		t.Fatalf("Expected 1 child for root, got %d", len(root.Children))
	}

	frameLayout := root.Children[0]
	if frameLayout.TagName != "node" || frameLayout.Class != "android.widget.FrameLayout" {
		t.Errorf("Expected FrameLayout node, got TagName: %s, Class: %s", frameLayout.TagName, frameLayout.Class)
	}

	// Verify bounds parsing
	expectedBounds := BoundingBox{X: 0, Y: 0, Width: 1080, Height: 1920}
	if frameLayout.Bounds != expectedBounds {
		t.Errorf("Expected bounds %+v, got %+v", expectedBounds, frameLayout.Bounds)
	}

	if len(frameLayout.Children) != 3 {
		t.Fatalf("Expected 3 children for FrameLayout, got %d", len(frameLayout.Children))
	}

	button := frameLayout.Children[0]
	if button.Text != "Submit" || button.ResourceID != "com.example:id/submit_btn" || button.ContentDesc != "Submit description" {
		t.Errorf("Button parsing attributes mismatched: %+v", button)
	}

	buttonExpectedBounds := BoundingBox{X: 100, Y: 200, Width: 200, Height: 100}
	if button.Bounds != buttonExpectedBounds {
		t.Errorf("Expected button bounds %+v, got %+v", buttonExpectedBounds, button.Bounds)
	}
}

// Test Locator strategies
func TestLocatorMatching(t *testing.T) {
	xmlSource := `
	<hierarchy>
		<node class="android.widget.FrameLayout" bounds="[0,0][1080,1920]">
			<node class="android.widget.Button" text="Submit" resource-id="com.example:id/submit_btn" content-desc="Submit description" bounds="[100,200][300,300]" id="btn_1"/>
			<node class="android.widget.EditText" text="" placeholder="Enter Username" bounds="[100,400][300,500]" test-id="username_field" id="edit_1"/>
		</node>
	</hierarchy>`

	container := &mockContainer{components: make(map[string]interface{})}
	container.Register("ActionEngine", &mockActionEngine{xmlSource: xmlSource})

	engine := NewLocatorEngine(container)

	ctx := context.Background()

	// 1. GetByText
	elem, err := engine.FindElement(ctx, "session_123", &Locator{Strategy: StrategyText, Selector: "Submit"})
	if err != nil {
		t.Fatalf("GetByText failed: %v", err)
	}
	if elem.TagName != "node" || elem.Attributes["id"] != "btn_1" {
		t.Errorf("Expected btn_1, got element %+v", elem)
	}

	// 2. GetByAccessibilityId
	elem, err = engine.FindElement(ctx, "session_123", &Locator{Strategy: StrategyAccessibilityID, Selector: "Submit description"})
	if err != nil {
		t.Fatalf("GetByAccessibilityId failed: %v", err)
	}
	if elem.Attributes["id"] != "btn_1" {
		t.Errorf("Expected btn_1, got element %+v", elem)
	}

	// 3. GetByTestId
	elem, err = engine.FindElement(ctx, "session_123", &Locator{Strategy: StrategyTestID, Selector: "username_field"})
	if err != nil {
		t.Fatalf("GetByTestId failed: %v", err)
	}
	if elem.Attributes["id"] != "edit_1" {
		t.Errorf("Expected edit_1, got element %+v", elem)
	}

	// 4. GetByRole
	elems, err := engine.FindElements(ctx, "session_123", &Locator{Strategy: StrategyRole, Selector: "android.widget.Button"})
	if err != nil {
		t.Fatalf("GetByRole failed: %v", err)
	}
	if len(elems) != 1 || elems[0].Attributes["id"] != "btn_1" {
		t.Errorf("Expected only btn_1 matching button role, got %+v", elems)
	}
}

// Test Relative constraint filtering
func TestRelativeLocators(t *testing.T) {
	xmlSource := `
	<hierarchy>
		<node class="android.widget.FrameLayout" bounds="[0,0][1080,1920]">
			<node class="android.widget.Button" text="Submit" bounds="[100,200][300,300]" id="btn_1"/>
			<node class="android.widget.TextView" text="Username Label" bounds="[100,100][300,150]" id="label_1"/>
		</node>
	</hierarchy>`

	container := &mockContainer{components: make(map[string]interface{})}
	container.Register("ActionEngine", &mockActionEngine{xmlSource: xmlSource})

	engine := NewLocatorEngine(container)

	ctx := context.Background()

	// Query: Find node where strategy=StrategyText, selector="Username Label" and is ABOVE target button
	targetLocator := &Locator{Strategy: StrategyText, Selector: "Submit"}
	relativeConstraint := RelativeConstraint{
		Direction: DirectionAbove,
		Target:    targetLocator,
	}

	queryLocator := &Locator{
		Strategy:    StrategyText,
		Selector:    "Username Label",
		Constraints: []RelativeConstraint{relativeConstraint},
	}

	elem, err := engine.FindElement(ctx, "session_123", queryLocator)
	if err != nil {
		t.Fatalf("Find relative locator failed: %v", err)
	}

	if elem.Attributes["id"] != "label_1" {
		t.Errorf("Expected label_1 relative above target button, got %s", elem.Attributes["id"])
	}
}

// Test Locator ranking score calculation
func TestLocatorRanking(t *testing.T) {
	root := &UINode{
		TagName: "node",
		Class:   "android.widget.Button",
	}

	targetNode := &UINode{
		TagName:     "node",
		Class:       "android.widget.Button",
		Text:        "Save changes",
		ResourceID:  "com.example:id/save_btn",
		ContentDesc: "Click to save changes",
		Attributes: map[string]string{
			"test-id": "save_action_button",
		},
		Parent: root,
	}
	root.Children = append(root.Children, targetNode)

	rankings := RankLocators(root, targetNode)

	if len(rankings) == 0 {
		t.Fatal("Expected locator rankings to be generated, got 0")
	}

	// Verify accessibility id is first and has the highest score
	bestRank := rankings[0]
	if bestRank.Locator.Strategy != StrategyAccessibilityID || bestRank.Score != 100.0 {
		t.Errorf("Expected best rank AccessibilityID (Score 100), got strategy: %s, score: %v", bestRank.Locator.Strategy, bestRank.Score)
	}

	// Verify XPath generation
	var xpathRank *ScoredLocator
	for _, r := range rankings {
		if r.Locator.Strategy == StrategyXPath {
			xpathRank = r
			break
		}
	}

	if xpathRank == nil {
		t.Fatal("Expected XPath locator in rankings, got none")
	}

	expectedXPath := "/node/node[1]"
	if xpathRank.Locator.Selector != expectedXPath {
		t.Errorf("Expected XPath selector %q, got %q", expectedXPath, xpathRank.Locator.Selector)
	}
}
