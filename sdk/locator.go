package sdk

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/ranjith035/morego/proto/v1"
)

// Locator models visual elements in the device layout tree.
type Locator struct {
	session  *Session
	strategy string
	selector string
	index    int
}

func mapStrategy(s string) pb.LocatorStrategy {
	switch strings.ToUpper(s) {
	case "ACCESSIBILITY_ID":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_ACCESSIBILITY_ID
	case "TEST_ID":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_TEST_ID
	case "ROLE":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_ROLE
	case "TEXT":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_TEXT
	case "PLACEHOLDER":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_PLACEHOLDER
	case "LABEL":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_LABEL
	case "RESOURCE_ID":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_RESOURCE_ID
	case "XPATH":
		return pb.LocatorStrategy_LOCATOR_STRATEGY_XPATH
	default:
		return pb.LocatorStrategy_LOCATOR_STRATEGY_UNSPECIFIED
	}
}

// Nth filters multiple selections by zero-indexed values.
func (l *Locator) Nth(index int) *Locator {
	return &Locator{
		session:  l.session,
		strategy: l.strategy,
		selector: l.selector,
		index:    index,
	}
}

// Click locates and triggers touch clicks on the target locator.
func (l *Locator) Click(ctx context.Context) error {
	// 1. Resolve element via FindElement RPC
	findResp, err := l.session.device.driverClient.FindElement(ctx, &pb.FindElementRequest{
		SessionId: l.session.sessionID,
		Locator: &pb.Locator{
			Strategy: mapStrategy(l.strategy),
			Selector: l.selector,
		},
		TimeoutMs: 5000,
	})
	if err != nil {
		return fmt.Errorf("failed to locate element for Click: %w", err)
	}

	if findResp.Element == nil {
		return fmt.Errorf("element not found for selector %q", l.selector)
	}

	// 2. Click the resolved element ID
	_, err = l.session.device.driverClient.Click(ctx, &pb.ClickRequest{
		DriverId: l.session.sessionID,
		Target: &pb.ClickRequest_ElementId{
			ElementId: findResp.Element.ElementId,
		},
	})
	if err != nil {
		return fmt.Errorf("click request failed: %w", err)
	}

	return nil
}

// Fill locates and enters text inside the target input element locator.
func (l *Locator) Fill(ctx context.Context, value string) error {
	// 1. Resolve element via FindElement RPC
	findResp, err := l.session.device.driverClient.FindElement(ctx, &pb.FindElementRequest{
		SessionId: l.session.sessionID,
		Locator: &pb.Locator{
			Strategy: mapStrategy(l.strategy),
			Selector: l.selector,
		},
		TimeoutMs: 5000,
	})
	if err != nil {
		return fmt.Errorf("failed to locate element for Fill: %w", err)
	}

	if findResp.Element == nil {
		return fmt.Errorf("element not found for selector %q", l.selector)
	}

	// 2. Fill the resolved element ID
	_, err = l.session.device.driverClient.Fill(ctx, &pb.FillRequest{
		DriverId:  l.session.sessionID,
		ElementId: findResp.Element.ElementId,
		Value:     value,
	})
	if err != nil {
		return fmt.Errorf("fill request failed: %w", err)
	}

	return nil
}
