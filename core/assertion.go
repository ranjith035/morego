package core

import "context"

// AssertionEngine performs auto-retrying assertions on target element nodes.
type AssertionEngine interface {
	Component

	// AssertVisible confirms that the matching node resolves and stays visible.
	AssertVisible(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error

	// AssertNotVisible asserts that the element is hidden or detached.
	AssertNotVisible(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error

	// AssertChecked asserts that checkbox/radio nodes are checked.
	AssertChecked(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error

	// AssertNotChecked asserts that target checkbox/radio nodes are unchecked.
	AssertNotChecked(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error

	// AssertText asserts that elements match the expected text representation.
	AssertText(ctx context.Context, sessionID string, locator *Locator, expected string, timeoutMS int) error

	// AssertEnabled asserts that target elements are interactive.
	AssertEnabled(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error

	// AssertDisabled asserts that target elements are disabled.
	AssertDisabled(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error

	// AssertCount asserts that the locator matches a specific collection length.
	AssertCount(ctx context.Context, sessionID string, locator *Locator, expected int, timeoutMS int) error

	// AssertValue asserts that input element value matches expected.
	AssertValue(ctx context.Context, sessionID string, locator *Locator, expected string, timeoutMS int) error

	// AssertExist asserts that an element is attached in the hierarchy.
	AssertExist(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error
}
