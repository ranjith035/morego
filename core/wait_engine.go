package core

import (
	"context"
	"fmt"
	"time"
)

// WaitEngineImpl implements the WaitEngine interface.
type WaitEngineImpl struct {
	container Container
}

// NewWaitEngine constructs a WaitEngine instance.
func NewWaitEngine(c Container) *WaitEngineImpl {
	return &WaitEngineImpl{container: c}
}

func (we *WaitEngineImpl) Name() string {
	return "WaitEngine"
}

func (we *WaitEngineImpl) Init(ctx context.Context) error {
	return nil
}

func (we *WaitEngineImpl) Shutdown(ctx context.Context) error {
	return nil
}

func (we *WaitEngineImpl) WaitForState(ctx context.Context, sessionID string, locator *Locator, state WaitState, opts WaitOptions) (*Element, error) {
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	interval := opts.Interval
	if interval <= 0 {
		interval = 50 * time.Millisecond // Tight default for snappy tests
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var lastBounds *BoundingBox
	var stableTicks int

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for locator %q to achieve state %q: %w", locator.Selector, state, ctx.Err())
		default:
			leVal, ok := we.container.Resolve("LocatorEngine")
			if !ok {
				return nil, fmt.Errorf("LocatorEngine not wired in container")
			}
			le := leVal.(LocatorEngine)

			elem, err := le.FindElement(ctx, sessionID, locator)
			if err != nil {
				// Element is not attached, reset counters
				lastBounds = nil
				stableTicks = 0
				if state == WaitStateAttached {
					time.Sleep(interval)
					continue
				}
				time.Sleep(interval)
				continue
			}

			// Attached check
			if state == WaitStateAttached {
				return elem, nil
			}

			// Visible check
			visible := elem.Bounds.Width > 0 && elem.Bounds.Height > 0 && elem.Attributes["visible"] != "false"
			if !visible {
				lastBounds = nil
				stableTicks = 0
				time.Sleep(interval)
				continue
			}
			if state == WaitStateVisible {
				return elem, nil
			}

			// Stability checks (compare consecutive bounding boxes)
			currentBounds := elem.Bounds
			stable := false
			if lastBounds != nil && *lastBounds == currentBounds {
				stableTicks++
				if stableTicks >= 2 { // Stable for at least two consecutive poll intervals
					stable = true
				}
			} else {
				stableTicks = 0
				lastBounds = &currentBounds
			}

			if state == WaitStateStable {
				if stable {
					return elem, nil
				}
				time.Sleep(interval)
				continue
			}

			// Enabled check
			enabled := elem.Attributes["enabled"] != "false" && elem.Attributes["disabled"] != "true"
			if !enabled {
				time.Sleep(interval)
				continue
			}
			if state == WaitStateEnabled {
				if stable {
					return elem, nil
				}
				time.Sleep(interval)
				continue
			}

			// Clickable check
			clickable := enabled && elem.Attributes["clickable"] != "false"
			if state == WaitStateClickable {
				if clickable && stable {
					return elem, nil
				}
				time.Sleep(interval)
				continue
			}
			time.Sleep(interval)
		}
	}
}

func (we *WaitEngineImpl) WaitForCondition(ctx context.Context, sessionID string, condition func() (bool, error), opts WaitOptions) error {
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	interval := opts.Interval
	if interval <= 0 {
		interval = 50 * time.Millisecond
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for condition: %w", ctx.Err())
		case <-ticker.C:
			ok, err := condition()
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
	}
}
