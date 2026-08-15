package core

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

// Mock LocatorEngine that changes element states based on poll counts
type transitioningLocatorEngine struct {
	mu        sync.Mutex
	pollCount int
	failFirst int // Number of times to return error (not attached)
	unstable  int // Number of times to return shifting bounds (not stable)
	disabled  int // Number of times to return disabled (not enabled)
}

func (t *transitioningLocatorEngine) Name() string                       { return "LocatorEngine" }
func (t *transitioningLocatorEngine) Init(ctx context.Context) error     { return nil }
func (t *transitioningLocatorEngine) Shutdown(ctx context.Context) error { return nil }

func (t *transitioningLocatorEngine) FindElement(ctx context.Context, sessionID string, locator *Locator) (*Element, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.pollCount++

	if t.pollCount <= t.failFirst {
		return nil, errors.New("element not attached")
	}

	bounds := BoundingBox{X: 100, Y: 100, Width: 200, Height: 50}
	if t.pollCount <= t.failFirst+t.unstable {
		// bounds shift to mock instability
		bounds.X += t.pollCount
	}

	isEnabled := "true"
	if t.pollCount <= t.failFirst+t.unstable+t.disabled {
		isEnabled = "false"
	}

	return &Element{
		ID:      "test_element",
		TagName: "node",
		Bounds:  bounds,
		Attributes: map[string]string{
			"visible":   "true",
			"enabled":   isEnabled,
			"clickable": "true",
		},
	}, nil
}

func (t *transitioningLocatorEngine) FindElements(ctx context.Context, sessionID string, locator *Locator) ([]*Element, error) {
	elem, err := t.FindElement(ctx, sessionID, locator)
	if err != nil {
		return nil, err
	}
	return []*Element{elem}, nil
}

// Test successful wait for attachment
func TestWaitForStateAttached(t *testing.T) {
	container := &mockContainer{components: make(map[string]interface{})}
	mockLE := &transitioningLocatorEngine{failFirst: 3}
	container.Register("LocatorEngine", mockLE)

	we := NewWaitEngine(container)
	ctx := context.Background()

	elem, err := we.WaitForState(ctx, "session_1", &Locator{Selector: "btn"}, WaitStateAttached, WaitOptions{
		Timeout:  1 * time.Second,
		Interval: 10 * time.Millisecond,
	})

	if err != nil {
		t.Fatalf("WaitForStateAttached failed: %v", err)
	}

	if elem == nil || elem.ID != "test_element" {
		t.Errorf("Expected element 'test_element', got %+v", elem)
	}

	mockLE.mu.Lock()
	polls := mockLE.pollCount
	mockLE.mu.Unlock()

	if polls < 4 {
		t.Errorf("Expected at least 4 polls before success, got %d", polls)
	}
}

// Test successful wait for stability
func TestWaitForStateStable(t *testing.T) {
	container := &mockContainer{components: make(map[string]interface{})}
	mockLE := &transitioningLocatorEngine{failFirst: 1, unstable: 3}
	container.Register("LocatorEngine", mockLE)

	we := NewWaitEngine(container)
	ctx := context.Background()

	elem, err := we.WaitForState(ctx, "session_1", &Locator{Selector: "btn"}, WaitStateStable, WaitOptions{
		Timeout:  1 * time.Second,
		Interval: 10 * time.Millisecond,
	})

	if err != nil {
		t.Fatalf("WaitForStateStable failed: %v", err)
	}

	if elem == nil {
		t.Fatal("Expected element, got nil")
	}

	mockLE.mu.Lock()
	polls := mockLE.pollCount
	mockLE.mu.Unlock()

	// Should poll:
	// 1: fail (not attached)
	// 2: bounds X=102 (stableTicks=0, lastBounds=102)
	// 3: bounds X=103 (stableTicks=0, lastBounds=103)
	// 4: bounds X=104 (stableTicks=0, lastBounds=104)
	// 5: bounds X=100 (stableTicks=0, lastBounds=100) (unstable ends)
	// 6: bounds X=100 (stableTicks=1, lastBounds=100)
	// 7: bounds X=100 (stableTicks=2, lastBounds=100) -> returns!
	if polls < 7 {
		t.Errorf("Expected at least 7 polls for stability, got %d", polls)
	}
}

// Test timeout scenario
func TestWaitForStateTimeout(t *testing.T) {
	container := &mockContainer{components: make(map[string]interface{})}
	// never becomes attached
	container.Register("LocatorEngine", &transitioningLocatorEngine{failFirst: 100})

	we := NewWaitEngine(container)
	ctx := context.Background()

	_, err := we.WaitForState(ctx, "session_1", &Locator{Selector: "btn"}, WaitStateAttached, WaitOptions{
		Timeout:  50 * time.Millisecond,
		Interval: 10 * time.Millisecond,
	})

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "timeout waiting for locator") {
		t.Errorf("Expected timeout error message, got: %v", err)
	}
}

// Test context cancellation
func TestWaitForStateCancellation(t *testing.T) {
	container := &mockContainer{components: make(map[string]interface{})}
	container.Register("LocatorEngine", &transitioningLocatorEngine{failFirst: 100})

	we := NewWaitEngine(container)
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context asynchronously after a brief delay
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	_, err := we.WaitForState(ctx, "session_1", &Locator{Selector: "btn"}, WaitStateAttached, WaitOptions{
		Timeout:  1 * time.Second,
		Interval: 5 * time.Millisecond,
	})

	if err == nil {
		t.Fatal("Expected cancellation error, got nil")
	}

	if !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context canceled error, got: %v", err)
	}
}
