package ai

import (
	"testing"
)

func TestSelfHealing(t *testing.T) {
	// A new hierarchy where submit button has resource-id "com.example:id/submit_btn_new" instead of "submit_btn_old"
	hierarchyXML := `
	<hierarchy>
		<node class="android.widget.FrameLayout" bounds="[0,0][1080,1920]">
			<node class="android.widget.Button" text="Login" resource-id="com.example:id/submit_btn_new" content-desc="Submit form" bounds="[100,200][300,300]" />
		</node>
	</hierarchy>`

	// History points to the old resource ID
	history := PastNodeHistory{
		Class:       "android.widget.Button",
		Text:        "Login",
		ResourceID:  "com.example:id/submit_btn_old",
		ContentDesc: "Submit form",
	}

	result, err := HealLocator("com.example:id/submit_btn_old", history, hierarchyXML)
	if err != nil {
		t.Fatalf("Healing failed: %v", err)
	}

	if result == nil {
		t.Fatal("Self-healing returned nil candidate, expected a healed match")
	}

	if result.HealedSelector != "com.example:id/submit_btn_new" {
		t.Errorf("Expected healed selector to be 'com.example:id/submit_btn_new', got %q", result.HealedSelector)
	}

	if result.Confidence < 0.5 {
		t.Errorf("Expected high healing match confidence, got %f", result.Confidence)
	}
}

func TestAccessibilityAudits(t *testing.T) {
	// Node 1: Button with empty label and desc
	// Node 2: Input field with empty label/placeholder
	hierarchyXML := `
	<hierarchy>
		<node class="android.widget.FrameLayout" bounds="[0,0][1080,1920]">
			<node class="android.widget.Button" text="" resource-id="com.example:id/btn_icon" content-desc="" bounds="[100,200][150,250]" />
			<node class="android.widget.EditText" text="" resource-id="com.example:id/empty_input" content-desc="" bounds="[100,300][300,350]" />
		</node>
	</hierarchy>`

	violations, err := AuditAccessibility(hierarchyXML)
	if err != nil {
		t.Fatalf("Accessibility audit failed: %v", err)
	}

	if len(violations) != 2 {
		t.Fatalf("Expected 2 accessibility violations, got %d", len(violations))
	}

	if violations[0].Rule != "CLICKABLE_MISSING_LABEL" {
		t.Errorf("Expected first violation rule to be CLICKABLE_MISSING_LABEL, got %s", violations[0].Rule)
	}

	if violations[1].Rule != "INPUT_MISSING_语义_HINT" {
		t.Errorf("Expected second violation rule to be INPUT_MISSING_语义_HINT, got %s", violations[1].Rule)
	}
}

func TestVisualLocator(t *testing.T) {
	hierarchyXML := `
	<hierarchy>
		<node class="android.widget.FrameLayout">
			<node class="android.widget.ImageView" resource-id="com.example:id/shopping_cart" content-desc="cart icon" bounds="[800,100][900,200]" />
		</node>
	</hierarchy>`

	bounds, err := ResolveVisualLocator("cart icon", hierarchyXML)
	if err != nil {
		t.Fatalf("Visual resolving failed: %v", err)
	}

	if bounds == nil {
		t.Fatal("Expected resolved bounds visual hit, got nil")
	}

	if bounds.X <= 0 || bounds.Y <= 0 {
		t.Errorf("Expected positive resolved coordinates, got X=%d, Y=%d", bounds.X, bounds.Y)
	}
}
