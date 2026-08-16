package core

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// UINode represents a parsed layout element node from XML view hierarchies.
type UINode struct {
	TagName            string
	Text               string
	ResourceID         string
	Class              string
	ContentDesc        string
	AccessibilityLabel string
	Placeholder        string
	Bounds             BoundingBox
	Attributes         map[string]string
	Children           []*UINode
	Parent             *UINode
}

// ScoredLocator stores ranking score weights.
type ScoredLocator struct {
	Locator *Locator
	Score   float64
}

// LocatorEngineImpl implements the LocatorEngine interface.
type LocatorEngineImpl struct {
	container Container
}

func NewLocatorEngine(c Container) *LocatorEngineImpl {
	return &LocatorEngineImpl{container: c}
}

func (le *LocatorEngineImpl) Name() string {
	return "LocatorEngine"
}

func (le *LocatorEngineImpl) Init(ctx context.Context) error {
	return nil
}

func (le *LocatorEngineImpl) Shutdown(ctx context.Context) error {
	return nil
}

func (le *LocatorEngineImpl) FindElement(ctx context.Context, sessionID string, locator *Locator) (*Element, error) {
	elements, err := le.FindElements(ctx, sessionID, locator)
	if err != nil {
		return nil, err
	}
	if len(elements) == 0 {
		return nil, fmt.Errorf("no elements matched locator selector %q", locator.Selector)
	}
	return elements[0], nil
}

func (le *LocatorEngineImpl) FindElements(ctx context.Context, sessionID string, locator *Locator) ([]*Element, error) {
	// 1. Resolve Driver from Session or Action engine (Mocked or stubbed here)
	// For interface implementation, we fetch active XML source dump using ActionEngine
	actionEngineVal, exists := le.container.Resolve("ActionEngine")
	if !exists {
		return nil, fmt.Errorf("ActionEngine not wired in container")
	}
	actionEngine := actionEngineVal.(ActionEngine)

	xmlSource, err := actionEngine.GetSource(ctx, sessionID, "xml")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch UI hierarchy source from driver: %w", err)
	}

	// 2. Parse UI XML tree
	root, err := ParseXML(strings.NewReader(xmlSource))
	if err != nil {
		return nil, fmt.Errorf("failed to parse UI hierarchy XML: %w", err)
	}

	// 3. Find matching nodes
	nodes := FindNodes(root, locator)

	// 4. Map matching nodes to elements
	elements := make([]*Element, 0, len(nodes))
	for _, n := range nodes {
		id := n.Attributes["id"]
		if id == "" {
			id = generateXPath(n)
		}
		elements = append(elements, &Element{
			ID:         id,
			TagName:    n.TagName,
			Text:       n.Text,
			Bounds:     n.Bounds,
			Attributes: n.Attributes,
		})
	}

	return elements, nil
}

// Helpers

func parseBounds(boundsStr string) BoundingBox {
	// Format: [x1,y1][x2,y2]
	boundsStr = strings.ReplaceAll(boundsStr, "][", ",")
	boundsStr = strings.ReplaceAll(boundsStr, "[", "")
	boundsStr = strings.ReplaceAll(boundsStr, "]", "")
	parts := strings.Split(boundsStr, ",")
	if len(parts) == 4 {
		x1, _ := strconv.Atoi(parts[0])
		y1, _ := strconv.Atoi(parts[1])
		x2, _ := strconv.Atoi(parts[2])
		y2, _ := strconv.Atoi(parts[3])
		return BoundingBox{X: x1, Y: y1, Width: x2 - x1, Height: y2 - y1}
	}
	return BoundingBox{}
}

// ParseXML parses XML hierarchy layouts into tree structures.
func ParseXML(r io.Reader) (*UINode, error) {
	decoder := xml.NewDecoder(r)
	var root *UINode
	var current *UINode

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch se := t.(type) {
		case xml.StartElement:
			node := &UINode{
				TagName:    se.Name.Local,
				Attributes: make(map[string]string),
			}
			for _, attr := range se.Attr {
				node.Attributes[attr.Name.Local] = attr.Value
				switch attr.Name.Local {
				case "text", "value":
					node.Text = attr.Value
				case "resource-id", "name":
					node.ResourceID = attr.Value
				case "class", "type":
					node.Class = attr.Value
				case "content-desc", "accessibility-id":
					node.ContentDesc = attr.Value
				case "label":
					node.AccessibilityLabel = attr.Value
				case "placeholder":
					node.Placeholder = attr.Value
				case "bounds":
					node.Bounds = parseBounds(attr.Value)
				}
			}

			// iOS coordinate split fallback
			if xVal, ok := node.Attributes["x"]; ok {
				x, _ := strconv.Atoi(xVal)
				y, _ := strconv.Atoi(node.Attributes["y"])
				w, _ := strconv.Atoi(node.Attributes["width"])
				h, _ := strconv.Atoi(node.Attributes["height"])
				node.Bounds = BoundingBox{X: x, Y: y, Width: w, Height: h}
			}

			if current != nil {
				node.Parent = current
				current.Children = append(current.Children, node)
			} else {
				root = node
			}
			current = node

		case xml.EndElement:
			if current != nil {
				current = current.Parent
			}
		}
	}
	return root, nil
}

// Matches checks if the node satisfies selector properties.
func (n *UINode) Matches(l *Locator) bool {
	switch l.Strategy {
	case StrategyAccessibilityID:
		return n.ContentDesc == l.Selector || n.AccessibilityLabel == l.Selector || n.Attributes["accessibility-id"] == l.Selector
	case StrategyTestID:
		return n.Attributes["test-id"] == l.Selector || n.ResourceID == l.Selector || strings.HasSuffix(n.ResourceID, "/"+l.Selector) || strings.HasSuffix(n.ResourceID, ":id/"+l.Selector)
	case StrategyRole:
		return strings.EqualFold(n.Class, l.Selector) || strings.HasSuffix(strings.ToLower(n.Class), "."+strings.ToLower(l.Selector))
	case StrategyText:
		return n.Text == l.Selector || strings.Contains(n.Text, l.Selector)
	case StrategyPlaceholder:
		return n.Placeholder == l.Selector || n.Attributes["placeholder"] == l.Selector
	case StrategyLabel:
		return n.AccessibilityLabel == l.Selector || n.Attributes["label"] == l.Selector
	case StrategyResourceID:
		return n.ResourceID == l.Selector
	case StrategyXPath:
		return strings.Contains(strings.ToLower(n.TagName), strings.ToLower(strings.TrimPrefix(l.Selector, "//")))
	default:
		return false
	}
}

// FindNodes locates matched nodes inside trees recursively.
func FindNodes(root *UINode, l *Locator) []*UINode {
	if root == nil {
		return nil
	}
	var matches []*UINode

	if l.Parent != nil {
		parentMatches := FindNodes(root, l.Parent)
		for _, parent := range parentMatches {
			matches = append(matches, findNodesInSubtree(parent, l)...)
		}
	} else {
		matches = findNodesInSubtree(root, l)
	}

	// Filter by relative constraints
	if len(l.Constraints) > 0 {
		matches = filterRelative(root, matches, l.Constraints)
	}

	// Filter by index
	if l.Index > 0 && l.Index < len(matches) {
		matches = []*UINode{matches[l.Index]}
	}

	return matches
}

func findNodesInSubtree(root *UINode, l *Locator) []*UINode {
	if root == nil {
		return nil
	}
	var matches []*UINode

	tempLocator := *l
	tempLocator.Parent = nil // isolate target node match
	if root.Matches(&tempLocator) {
		matches = append(matches, root)
	}

	for _, child := range root.Children {
		matches = append(matches, findNodesInSubtree(child, l)...)
	}
	return matches
}

func filterRelative(root *UINode, candidates []*UINode, constraints []RelativeConstraint) []*UINode {
	var filtered []*UINode
	for _, c := range candidates {
		ok := true
		for _, rc := range constraints {
			targets := FindNodes(root, rc.Target)
			if len(targets) == 0 {
				ok = false
				break
			}
			anyMatch := false
			for _, t := range targets {
				if isSpatialMatch(c.Bounds, t.Bounds, rc.Direction) {
					anyMatch = true
					break
				}
			}
			if !anyMatch {
				ok = false
				break
			}
		}
		if ok {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func isSpatialMatch(c, t BoundingBox, dir RelativeDirection) bool {
	switch dir {
	case DirectionAbove:
		return (c.Y + c.Height) <= t.Y
	case DirectionBelow:
		return c.Y >= (t.Y + t.Height)
	case DirectionLeftOf:
		return (c.X + c.Width) <= t.X
	case DirectionRightOf:
		return c.X >= (t.X + t.Width)
	default:
		return false
	}
}

// RankLocators scores possible locator choices for a node.
func RankLocators(root *UINode, target *UINode) []*ScoredLocator {
	var results []*ScoredLocator

	// 1. Accessibility ID
	if desc := target.ContentDesc; desc != "" {
		results = append(results, &ScoredLocator{
			Locator: &Locator{Strategy: StrategyAccessibilityID, Selector: desc},
			Score:   100.0,
		})
	}
	if label := target.AccessibilityLabel; label != "" {
		results = append(results, &ScoredLocator{
			Locator: &Locator{Strategy: StrategyLabel, Selector: label},
			Score:   95.0,
		})
	}

	// 2. Test ID
	if tid := target.Attributes["test-id"]; tid != "" {
		results = append(results, &ScoredLocator{
			Locator: &Locator{Strategy: StrategyTestID, Selector: tid},
			Score:   90.0,
		})
	}

	// 3. Text
	if target.Text != "" {
		results = append(results, &ScoredLocator{
			Locator: &Locator{Strategy: StrategyText, Selector: target.Text},
			Score:   75.0,
		})
	}

	// 4. Resource ID
	if target.ResourceID != "" {
		results = append(results, &ScoredLocator{
			Locator: &Locator{Strategy: StrategyResourceID, Selector: target.ResourceID},
			Score:   60.0,
		})
	}

	// 5. XPath fallback
	xpath := generateXPath(target)
	results = append(results, &ScoredLocator{
		Locator: &Locator{Strategy: StrategyXPath, Selector: xpath},
		Score:   10.0,
	})

	return results
}

func generateXPath(node *UINode) string {
	if node == nil {
		return ""
	}
	if node.Parent == nil {
		return "/" + node.TagName
	}
	idx := 1
	for _, sibling := range node.Parent.Children {
		if sibling == node {
			break
		}
		if sibling.TagName == node.TagName {
			idx++
		}
	}
	return generateXPath(node.Parent) + "/" + node.TagName + "[" + strconv.Itoa(idx) + "]"
}
