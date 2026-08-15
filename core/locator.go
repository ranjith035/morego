package core

import "context"

// LocatorStrategy specifies the element selection method.
type LocatorStrategy string

const (
	StrategyAccessibilityID LocatorStrategy = "ACCESSIBILITY_ID"
	StrategyTestID          LocatorStrategy = "TEST_ID"
	StrategyRole            LocatorStrategy = "ROLE"
	StrategyText            LocatorStrategy = "TEXT"
	StrategyPlaceholder     LocatorStrategy = "PLACEHOLDER"
	StrategyLabel           LocatorStrategy = "LABEL"
	StrategyResourceID      LocatorStrategy = "RESOURCE_ID"
	StrategyXPath           LocatorStrategy = "XPATH"
)

// RelativeDirection defines positional queries relative to target elements.
type RelativeDirection string

const (
	DirectionAbove   RelativeDirection = "ABOVE"
	DirectionBelow   RelativeDirection = "BELOW"
	DirectionLeftOf  RelativeDirection = "LEFT_OF"
	DirectionRightOf RelativeDirection = "RIGHT_OF"
)

// BoundingBox holds structural layouts.
type BoundingBox struct {
	X      int
	Y      int
	Width  int
	Height int
}

// RelativeConstraint structures relative queries.
type RelativeConstraint struct {
	Direction RelativeDirection
	Target    *Locator
	Distance  int
}

// Locator aggregates search instructions.
type Locator struct {
	Strategy    LocatorStrategy
	Selector    string
	Parent      *Locator
	Index       int
	Constraints []RelativeConstraint
}

// Element models a matched node in the target application's active viewport hierarchy.
type Element struct {
	ID         string
	TagName    string
	Text       string
	Bounds     BoundingBox
	Attributes map[string]string
}

// LocatorEngine resolves locators to Element matches.
type LocatorEngine interface {
	Component

	// FindElement searches the active DOM viewport for the first element matching a locator.
	FindElement(ctx context.Context, sessionID string, locator *Locator) (*Element, error)

	// FindElements finds all matching elements.
	FindElements(ctx context.Context, sessionID string, locator *Locator) ([]*Element, error)
}
