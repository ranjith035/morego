package recorder

import (
	"strings"
	"testing"
)

func TestRecorderOptimizations(t *testing.T) {
	rec := NewRecorder()

	// Add consecutive fills on same locator
	rec.AddAction(ActionIR{
		Type:             ActionFill,
		SelectorStrategy: "TEST_ID",
		Selector:         "username",
		Value:            "a",
	})
	rec.AddAction(ActionIR{
		Type:             ActionFill,
		SelectorStrategy: "TEST_ID",
		Selector:         "username",
		Value:            "ab",
	})
	rec.AddAction(ActionIR{
		Type:             ActionFill,
		SelectorStrategy: "TEST_ID",
		Selector:         "username",
		Value:            "abc",
	})

	// Add click on button
	rec.AddAction(ActionIR{
		Type:             ActionClick,
		SelectorStrategy: "ACCESSIBILITY_ID",
		Selector:         "submit",
	})

	if len(rec.actions) != 2 {
		t.Fatalf("Expected consecutive Fills to merge into 1 action, total action count: %d", len(rec.actions))
	}

	if rec.actions[0].Value != "abc" {
		t.Errorf("Expected merged Fill action value to be 'abc', got %q", rec.actions[0].Value)
	}
}

func TestLanguageGenerators(t *testing.T) {
	rec := NewRecorder()

	rec.AddAction(ActionIR{
		Type:             ActionClick,
		SelectorStrategy: "ACCESSIBILITY_ID",
		Selector:         "login_btn",
	})
	rec.AddAction(ActionIR{
		Type:             ActionFill,
		SelectorStrategy: "TEST_ID",
		Selector:         "password_field",
		Value:            "my_secret_pass",
	})
	rec.AddAction(ActionIR{
		Type:       ActionSwipe,
		StartX:     10,
		StartY:     20,
		EndX:       10,
		EndY:       80,
		DurationMS: 250,
	})

	languages := []string{"typescript", "python", "go", "java", "csharp", "kotlin"}

	for _, lang := range languages {
		code, err := rec.Generate(lang)
		if err != nil {
			t.Fatalf("Failed to generate code for language %s: %v", lang, err)
		}

		if len(code) == 0 {
			t.Errorf("Generated code for %s is empty", lang)
		}

		// Verify target selector strategy is preserved
		if !strings.Contains(code, "login_btn") || !strings.Contains(code, "password_field") {
			t.Errorf("Generated code for %s did not contain expected elements", lang)
		}

		// Verify Swipe parameters are contained
		if !strings.Contains(code, "20") && !strings.Contains(code, "250") {
			t.Errorf("Generated code for %s did not contain Swipe parameters", lang)
		}
	}
}
