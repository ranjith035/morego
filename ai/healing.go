package ai

import (
	"encoding/xml"
	"strings"
)

// XMLNode represents nodes scanned in XML tree hierarchies.
type XMLNode struct {
	XMLName     xml.Name
	Class       string    `xml:"class,attr"`
	Text        string    `xml:"text,attr"`
	ResourceID  string    `xml:"resource-id,attr"`
	ContentDesc string    `xml:"content-desc,attr"`
	Bounds      string    `xml:"bounds,attr"`
	Nodes       []XMLNode `xml:"node"`
}

// PastNodeHistory defines references metadata recorded during previous successful runs.
type PastNodeHistory struct {
	Class       string
	Text        string
	ResourceID  string
	ContentDesc string
}

// HealingResult holds properties of healed selectors candidates.
type HealingResult struct {
	HealedSelector string
	Strategy       string // e.g. "RESOURCE_ID", "TEXT", "XPATH"
	Confidence     float64
	Node           XMLNode
}

// HealLocator evaluates element trees to propose healed queries.
func HealLocator(brokenSelector string, history PastNodeHistory, hierarchyXML string) (*HealingResult, error) {
	var root XMLNode
	err := xml.Unmarshal([]byte(hierarchyXML), &root)
	if err != nil {
		return nil, err
	}

	var candidates []XMLNode
	collectNodes(root, &candidates)

	var bestNode XMLNode
	bestScore := 0.0

	for _, node := range candidates {
		score := 0.0

		// 1. Class match (weight: 0.3)
		if node.Class == history.Class && node.Class != "" {
			score += 0.3
		}

		// 2. Text match/similarity (weight: 0.3)
		if node.Text != "" && history.Text != "" {
			if node.Text == history.Text {
				score += 0.3
			} else if strings.Contains(node.Text, history.Text) || strings.Contains(history.Text, node.Text) {
				score += 0.15
			}
		}

		// 3. Resource ID match/similarity (weight: 0.3)
		if node.ResourceID != "" && history.ResourceID != "" {
			if node.ResourceID == history.ResourceID {
				score += 0.3
			} else {
				// Suffix matching (e.g. "submit_btn_old" vs "submit_btn_new")
				s1 := getIDSuffix(node.ResourceID)
				s2 := getIDSuffix(history.ResourceID)
				if s1 == s2 && s1 != "" {
					score += 0.2
				}
			}
		}

		// 4. Content Desc match (weight: 0.1)
		if node.ContentDesc != "" && history.ContentDesc != "" {
			if node.ContentDesc == history.ContentDesc {
				score += 0.1
			}
		}

		if score > bestScore {
			bestScore = score
			bestNode = node
		}
	}

	if bestScore < 0.25 {
		return nil, nil // Similarity is too low to claim a healed locator match
	}

	// Propose best strategy based on matched properties
	healedSel := bestNode.ResourceID
	strategy := "RESOURCE_ID"
	if healedSel == "" {
		healedSel = bestNode.Text
		strategy = "TEXT"
	}
	if healedSel == "" {
		healedSel = bestNode.ContentDesc
		strategy = "ACCESSIBILITY_ID"
	}

	return &HealingResult{
		HealedSelector: healedSel,
		Strategy:       strategy,
		Confidence:     bestScore,
		Node:           bestNode,
	}, nil
}

func collectNodes(node XMLNode, list *[]XMLNode) {
	*list = append(*list, node)
	for _, child := range node.Nodes {
		collectNodes(child, list)
	}
}

func getIDSuffix(id string) string {
	parts := strings.Split(id, "/")
	return parts[len(parts)-1]
}
