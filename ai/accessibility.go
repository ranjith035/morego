package ai

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// AccessibilityViolation records accessibility rule gaps.
type AccessibilityViolation struct {
	ElementID  string
	Class      string
	Rule       string
	Suggestion string
	Severity   string // "high", "medium", "low"
}

// AuditAccessibility parses view hierarchies to find accessibility compliance gaps.
func AuditAccessibility(hierarchyXML string) ([]AccessibilityViolation, error) {
	var root XMLNode
	err := xml.Unmarshal([]byte(hierarchyXML), &root)
	if err != nil {
		return nil, err
	}

	var candidates []XMLNode
	collectNodes(root, &candidates)

	var violations []AccessibilityViolation

	for _, node := range candidates {
		isButton := strings.Contains(strings.ToLower(node.Class), "button") ||
			strings.Contains(strings.ToLower(node.XMLName.Local), "button")

		// Rule 1: Clickable element must have a label or content-description
		if isButton {
			cleanText := strings.TrimSpace(node.Text)
			cleanDesc := strings.TrimSpace(node.ContentDesc)
			if cleanText == "" && cleanDesc == "" {
				violations = append(violations, AccessibilityViolation{
					ElementID:  node.ResourceID,
					Class:      node.Class,
					Rule:       "CLICKABLE_MISSING_LABEL",
					Suggestion: "Provide a 'content-desc' or 'text' attribute to declare label semantics for screen readers.",
					Severity:   "high",
				})
			}
		}

		// Rule 2: Input fields (EditText) should have descriptive text labels/placeholders
		isInput := strings.Contains(strings.ToLower(node.Class), "edittext") ||
			strings.Contains(strings.ToLower(node.XMLName.Local), "input")

		if isInput {
			if strings.TrimSpace(node.Text) == "" && strings.TrimSpace(node.ContentDesc) == "" {
				violations = append(violations, AccessibilityViolation{
					ElementID:  node.ResourceID,
					Class:      node.Class,
					Rule:       "INPUT_MISSING_语义_HINT",
					Suggestion: "Ensure text input fields define accessibility descriptions or placeholder text values.",
					Severity:   "medium",
				})
			}
		}
	}

	return violations, nil
}

func (av AccessibilityViolation) String() string {
	return fmt.Sprintf("[%s] %s on Element: %s (%s)", av.Severity, av.Rule, av.ElementID, av.Class)
}
