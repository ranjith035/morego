package ai

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
)

// LocatorSuggestion describes a recommended locator and why it is stable.
type LocatorSuggestion struct {
	Strategy       string
	Selector       string
	Reason         string
	StabilityScore float32
}

// SuggestLocators ranks likely-stable locators from the current hierarchy.
func SuggestLocators(hierarchyXML string) ([]LocatorSuggestion, error) {
	var root XMLNode
	if err := xml.Unmarshal([]byte(hierarchyXML), &root); err != nil {
		return nil, err
	}

	var nodes []XMLNode
	collectNodes(root, &nodes)

	suggestions := make([]LocatorSuggestion, 0)
	seen := make(map[string]struct{})

	for _, node := range nodes {
		for _, suggestion := range candidateSuggestions(node) {
			key := suggestion.Strategy + "\x00" + suggestion.Selector
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			suggestions = append(suggestions, suggestion)
		}
	}

	sort.SliceStable(suggestions, func(i, j int) bool {
		if suggestions[i].StabilityScore == suggestions[j].StabilityScore {
			return suggestions[i].Selector < suggestions[j].Selector
		}
		return suggestions[i].StabilityScore > suggestions[j].StabilityScore
	})

	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	return suggestions, nil
}

func candidateSuggestions(node XMLNode) []LocatorSuggestion {
	suggestions := make([]LocatorSuggestion, 0, 5)

	if value := strings.TrimSpace(node.ContentDesc); value != "" {
		suggestions = append(suggestions, LocatorSuggestion{
			Strategy:       "ACCESSIBILITY_ID",
			Selector:       value,
			Reason:         "Accessibility labels are typically stable and user-facing.",
			StabilityScore: 0.98,
		})
	}

	if value := stableTestID(node.ResourceID); value != "" {
		suggestions = append(suggestions, LocatorSuggestion{
			Strategy:       "TEST_ID",
			Selector:       value,
			Reason:         "Test identifiers are explicit automation hooks and usually survive UI copy changes.",
			StabilityScore: 0.95,
		})
	}

	if value := strings.TrimSpace(node.ResourceID); value != "" {
		suggestions = append(suggestions, LocatorSuggestion{
			Strategy:       "RESOURCE_ID",
			Selector:       value,
			Reason:         "Resource identifiers are usually more stable than visual text.",
			StabilityScore: 0.9,
		})
	}

	if value := strings.TrimSpace(node.Text); value != "" {
		suggestions = append(suggestions, LocatorSuggestion{
			Strategy:       "TEXT",
			Selector:       value,
			Reason:         "Visible text is readable but more likely to change than explicit identifiers.",
			StabilityScore: 0.72,
		})
	}

	if class := strings.TrimSpace(node.Class); class != "" {
		suggestions = append(suggestions, LocatorSuggestion{
			Strategy:       "ROLE",
			Selector:       class,
			Reason:         "Class-based locators can help when semantic attributes are missing.",
			StabilityScore: 0.45,
		})
	}

	return suggestions
}

func stableTestID(resourceID string) string {
	resourceID = strings.TrimSpace(resourceID)
	if resourceID == "" {
		return ""
	}

	parts := strings.Split(resourceID, "/")
	last := parts[len(parts)-1]
	if strings.HasPrefix(last, "id/") {
		return strings.TrimPrefix(last, "id/")
	}
	return last
}

func (ls LocatorSuggestion) String() string {
	return fmt.Sprintf("%s=%q (%.2f)", ls.Strategy, ls.Selector, ls.StabilityScore)
}
