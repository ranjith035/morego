package core

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AssertionEngineImpl implements the AssertionEngine interface.
type AssertionEngineImpl struct {
	container Container
}

// NewAssertionEngine constructs an AssertionEngine instance.
func NewAssertionEngine(c Container) *AssertionEngineImpl {
	return &AssertionEngineImpl{container: c}
}

func (ae *AssertionEngineImpl) Name() string {
	return "AssertionEngine"
}

func (ae *AssertionEngineImpl) Init(ctx context.Context) error {
	return nil
}

func (ae *AssertionEngineImpl) Shutdown(ctx context.Context) error {
	return nil
}

func (ae *AssertionEngineImpl) getWaitEngine() (WaitEngine, error) {
	val, ok := ae.container.Resolve("WaitEngine")
	if !ok {
		return nil, fmt.Errorf("WaitEngine not wired in container")
	}
	return val.(WaitEngine), nil
}

func (ae *AssertionEngineImpl) getLocatorEngine() (LocatorEngine, error) {
	val, ok := ae.container.Resolve("LocatorEngine")
	if !ok {
		return nil, fmt.Errorf("LocatorEngine not wired in container")
	}
	return val.(LocatorEngine), nil
}

func (ae *AssertionEngineImpl) AssertVisible(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	_, err = we.WaitForState(ctx, sessionID, locator, WaitStateVisible, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})
	if err != nil {
		return fmt.Errorf("assertion failed: element is not visible: %w", err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertNotVisible(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	le, err := ae.getLocatorEngine()
	if err != nil {
		return err
	}

	err = we.WaitForCondition(ctx, sessionID, func() (bool, error) {
		elem, err := le.FindElement(ctx, sessionID, locator)
		if err != nil {
			// Element is not attached, which counts as not visible
			return true, nil
		}
		visible := elem.Bounds.Width > 0 && elem.Bounds.Height > 0 && elem.Attributes["visible"] != "false"
		return !visible, nil
	}, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})

	if err != nil {
		return fmt.Errorf("assertion failed: element remained visible: %w", err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertExist(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	_, err = we.WaitForState(ctx, sessionID, locator, WaitStateAttached, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})
	if err != nil {
		return fmt.Errorf("assertion failed: element does not exist: %w", err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertChecked(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	le, err := ae.getLocatorEngine()
	if err != nil {
		return err
	}

	err = we.WaitForCondition(ctx, sessionID, func() (bool, error) {
		elem, err := le.FindElement(ctx, sessionID, locator)
		if err != nil {
			return false, nil // Wait until attached
		}
		checked := elem.Attributes["checked"] == "true" || elem.Attributes["selected"] == "true"
		return checked, nil
	}, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})

	if err != nil {
		return fmt.Errorf("assertion failed: element is not checked: %w", err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertNotChecked(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	le, err := ae.getLocatorEngine()
	if err != nil {
		return err
	}

	err = we.WaitForCondition(ctx, sessionID, func() (bool, error) {
		elem, err := le.FindElement(ctx, sessionID, locator)
		if err != nil {
			return false, nil // Wait until attached
		}
		checked := elem.Attributes["checked"] == "true" || elem.Attributes["selected"] == "true"
		return !checked, nil
	}, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})

	if err != nil {
		return fmt.Errorf("assertion failed: element remained checked: %w", err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertText(ctx context.Context, sessionID string, locator *Locator, expected string, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	le, err := ae.getLocatorEngine()
	if err != nil {
		return err
	}

	err = we.WaitForCondition(ctx, sessionID, func() (bool, error) {
		elem, err := le.FindElement(ctx, sessionID, locator)
		if err != nil {
			return false, nil // Wait until attached
		}
		return strings.Contains(elem.Text, expected), nil
	}, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})

	if err != nil {
		return fmt.Errorf("assertion failed: element text does not contain %q: %w", expected, err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertValue(ctx context.Context, sessionID string, locator *Locator, expected string, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	le, err := ae.getLocatorEngine()
	if err != nil {
		return err
	}

	err = we.WaitForCondition(ctx, sessionID, func() (bool, error) {
		elem, err := le.FindElement(ctx, sessionID, locator)
		if err != nil {
			return false, nil // Wait until attached
		}
		val, ok := elem.Attributes["value"]
		if !ok {
			val = elem.Text // Fallback to raw text if value is not set
		}
		return val == expected, nil
	}, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})

	if err != nil {
		return fmt.Errorf("assertion failed: element value is not %q: %w", expected, err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertEnabled(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	_, err = we.WaitForState(ctx, sessionID, locator, WaitStateEnabled, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})
	if err != nil {
		return fmt.Errorf("assertion failed: element is not enabled: %w", err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertDisabled(ctx context.Context, sessionID string, locator *Locator, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	le, err := ae.getLocatorEngine()
	if err != nil {
		return err
	}

	err = we.WaitForCondition(ctx, sessionID, func() (bool, error) {
		elem, err := le.FindElement(ctx, sessionID, locator)
		if err != nil {
			return false, nil // Wait until attached
		}
		disabled := elem.Attributes["enabled"] == "false" || elem.Attributes["disabled"] == "true"
		return disabled, nil
	}, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})

	if err != nil {
		return fmt.Errorf("assertion failed: element is not disabled: %w", err)
	}
	return nil
}

func (ae *AssertionEngineImpl) AssertCount(ctx context.Context, sessionID string, locator *Locator, expected int, timeoutMS int) error {
	we, err := ae.getWaitEngine()
	if err != nil {
		return err
	}
	le, err := ae.getLocatorEngine()
	if err != nil {
		return err
	}

	err = we.WaitForCondition(ctx, sessionID, func() (bool, error) {
		elems, err := le.FindElements(ctx, sessionID, locator)
		if err != nil {
			return false, nil
		}
		return len(elems) == expected, nil
	}, WaitOptions{
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	})

	if err != nil {
		return fmt.Errorf("assertion failed: elements count is not %d: %w", expected, err)
	}
	return nil
}
