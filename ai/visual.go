package ai

import (
	"encoding/xml"
	"errors"
	"strings"
)

// VisualBoundingBox stores coordinates.
type VisualBoundingBox struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ResolveVisualLocator simulates mapping abstract image labels to screen coordinates.
func ResolveVisualLocator(label string, hierarchyXML string) (*VisualBoundingBox, error) {
	var root XMLNode
	err := xml.Unmarshal([]byte(hierarchyXML), &root)
	if err != nil {
		return nil, err
	}

	var candidates []XMLNode
	collectNodes(root, &candidates)

	cleanLabel := strings.ToLower(label)

	for _, node := range candidates {
		// Look for matching semantic descriptions
		if strings.Contains(strings.ToLower(node.ContentDesc), cleanLabel) ||
			strings.Contains(strings.ToLower(node.ResourceID), cleanLabel) ||
			strings.Contains(strings.ToLower(node.Text), cleanLabel) {

			// Return dummy bounds or parsed bounds
			return &VisualBoundingBox{
				X:      100, // mock offset
				Y:      200,
				Width:  80,
				Height: 80,
			}, nil
		}
	}

	return nil, errors.New("visual element matching label not found on screen")
}
