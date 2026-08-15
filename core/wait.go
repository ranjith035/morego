package core

import (
	"context"
	"time"
)

// WaitState represents target lifecycle expectations.
type WaitState string

const (
	WaitStateAttached  WaitState = "ATTACHED"
	WaitStateVisible   WaitState = "VISIBLE"
	WaitStateStable    WaitState = "STABLE"
	WaitStateEnabled   WaitState = "ENABLED"
	WaitStateClickable WaitState = "CLICKABLE"
)

// WaitOptions details pacing rules for polling operations.
type WaitOptions struct {
	Timeout  time.Duration
	Interval time.Duration
}

// WaitEngine acts as the auto-wait sync gateway, waiting for target UI elements to resolve.
type WaitEngine interface {
	Component

	// WaitForState blocks until a locator matches and achieves the target state.
	WaitForState(ctx context.Context, sessionID string, locator *Locator, state WaitState, opts WaitOptions) (*Element, error)

	// WaitForCondition blocks until the specified assertion function evaluates to true.
	WaitForCondition(ctx context.Context, sessionID string, condition func() (bool, error), opts WaitOptions) error
}
