package core

import (
	"context"
	"strings"
	"testing"
)

func TestAssertionEngine(t *testing.T) {
	xmlSource := `
	<hierarchy>
		<node class="android.widget.FrameLayout" bounds="[0,0][1080,1920]">
			<node class="android.widget.Button" text="Submit" resource-id="com.example:id/submit_btn" content-desc="Submit description" bounds="[100,200][300,300]" id="btn_1" visible="true" enabled="true" checked="false" clickable="true"/>
			<node class="android.widget.Button" text="Cancel" resource-id="com.example:id/cancel_btn" bounds="[100,400][300,500]" id="btn_2" visible="false" enabled="false" disabled="true" checked="true"/>
			<node class="android.widget.EditText" text="John Doe" resource-id="com.example:id/username" bounds="[100,600][300,700]" id="edit_1" visible="true" enabled="true" value="John Doe"/>
		</node>
	</hierarchy>`

	container := &mockContainer{components: make(map[string]interface{})}
	container.Register("ActionEngine", &mockActionEngine{xmlSource: xmlSource})

	// Register real engines
	le := NewLocatorEngine(container)
	container.Register("LocatorEngine", le)

	we := NewWaitEngine(container)
	container.Register("WaitEngine", we)

	ae := NewAssertionEngine(container)

	ctx := context.Background()
	sessionID := "session_1"

	// 1. AssertVisible
	btnLocator := &Locator{Strategy: StrategyResourceID, Selector: "com.example:id/submit_btn"}
	err := ae.AssertVisible(ctx, sessionID, btnLocator, 500)
	if err != nil {
		t.Errorf("AssertVisible failed: %v", err)
	}

	// 2. AssertNotVisible
	cancelLocator := &Locator{Strategy: StrategyResourceID, Selector: "com.example:id/cancel_btn"}
	err = ae.AssertNotVisible(ctx, sessionID, cancelLocator, 500)
	if err != nil {
		t.Errorf("AssertNotVisible failed: %v", err)
	}

	// 3. AssertExist
	err = ae.AssertExist(ctx, sessionID, btnLocator, 500)
	if err != nil {
		t.Errorf("AssertExist failed: %v", err)
	}

	// 4. AssertChecked
	err = ae.AssertChecked(ctx, sessionID, cancelLocator, 500)
	if err != nil {
		t.Errorf("AssertChecked failed: %v", err)
	}

	// 5. AssertNotChecked
	err = ae.AssertNotChecked(ctx, sessionID, btnLocator, 500)
	if err != nil {
		t.Errorf("AssertNotChecked failed: %v", err)
	}

	// 6. AssertText
	err = ae.AssertText(ctx, sessionID, btnLocator, "Sub", 500)
	if err != nil {
		t.Errorf("AssertText failed: %v", err)
	}

	// 7. AssertValue
	userLocator := &Locator{Strategy: StrategyResourceID, Selector: "com.example:id/username"}
	err = ae.AssertValue(ctx, sessionID, userLocator, "John Doe", 500)
	if err != nil {
		t.Errorf("AssertValue failed: %v", err)
	}

	// 8. AssertEnabled
	err = ae.AssertEnabled(ctx, sessionID, btnLocator, 500)
	if err != nil {
		t.Errorf("AssertEnabled failed: %v", err)
	}

	// 9. AssertDisabled
	err = ae.AssertDisabled(ctx, sessionID, cancelLocator, 500)
	if err != nil {
		t.Errorf("AssertDisabled failed: %v", err)
	}

	// 10. AssertCount
	allButtonsLocator := &Locator{Strategy: StrategyRole, Selector: "android.widget.Button"}
	err = ae.AssertCount(ctx, sessionID, allButtonsLocator, 2, 500)
	if err != nil {
		t.Errorf("AssertCount failed: %v", err)
	}
}

func TestAssertionFailures(t *testing.T) {
	xmlSource := `
	<hierarchy>
		<node class="android.widget.FrameLayout" bounds="[0,0][1080,1920]">
			<node class="android.widget.Button" text="Submit" bounds="[100,200][300,300]" id="btn_1" visible="true" enabled="true"/>
		</node>
	</hierarchy>`

	container := &mockContainer{components: make(map[string]interface{})}
	container.Register("ActionEngine", &mockActionEngine{xmlSource: xmlSource})

	le := NewLocatorEngine(container)
	container.Register("LocatorEngine", le)

	we := NewWaitEngine(container)
	container.Register("WaitEngine", we)

	ae := NewAssertionEngine(container)

	ctx := context.Background()
	sessionID := "session_1"

	// Assert text that doesn't match
	btnLocator := &Locator{Strategy: StrategyText, Selector: "Submit"}
	err := ae.AssertText(ctx, sessionID, btnLocator, "InvalidText", 20)
	if err == nil {
		t.Fatal("Expected AssertText to fail, but it passed")
	}

	if !strings.Contains(err.Error(), "assertion failed") {
		t.Errorf("Expected failure message starting with 'assertion failed', got: %v", err)
	}
}
