package core

import (
	"context"
	"time"
)

// Point details 2D screen coordinates.
type Point struct {
	X int
	Y int
}

// ActionEngine coordinates physical touch and state inputs sent down to target platforms.
type ActionEngine interface {
	Component

	// Click targets an element resolved by a locator, performing auto-wait checks beforehand.
	Click(ctx context.Context, sessionID string, locator *Locator) error

	// ClickAt dispatches a click tap directly at coordinate positions.
	ClickAt(ctx context.Context, sessionID string, pt Point) error

	// Fill clears input areas and types the given value.
	Fill(ctx context.Context, sessionID string, locator *Locator, text string) error

	// Swipe performs linear translation movements.
	Swipe(ctx context.Context, sessionID string, start, end Point, duration time.Duration) error

	// Screenshot captures layout graphics. If locator is nil, returns full screen.
	Screenshot(ctx context.Context, sessionID string, locator *Locator) ([]byte, error)

	// GetSource downloads current layout source code tree.
	GetSource(ctx context.Context, sessionID string, format string) (string, error)

	// ExecuteScript evaluates script macros inside system execution sandboxes.
	ExecuteScript(ctx context.Context, sessionID string, script string, args []string) (string, error)
}
