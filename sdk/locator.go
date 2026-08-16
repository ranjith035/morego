package sdk

import (
	"context"
	"fmt"
	"strings"

	pb "github.com/ranjith035/morego/proto/v1"
)

type relativeConstraint struct {
	direction string
	target    *Locator
}

// Locator models visual elements in the device layout tree.
type Locator struct {
	session     *Session
	strategy    string
	selector    string
	index       int
	parent      *Locator
	constraints []relativeConstraint
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

func mapDirection(d string) pb.RelativeDirection {
	switch strings.ToUpper(d) {
	case "ABOVE":
		return pb.RelativeDirection_RELATIVE_DIRECTION_ABOVE
	case "BELOW":
		return pb.RelativeDirection_RELATIVE_DIRECTION_BELOW
	case "LEFT_OF":
		return pb.RelativeDirection_RELATIVE_DIRECTION_LEFT_OF
	case "RIGHT_OF":
		return pb.RelativeDirection_RELATIVE_DIRECTION_RIGHT_OF
	default:
		return pb.RelativeDirection_RELATIVE_DIRECTION_UNSPECIFIED
	}
}

func (l *Locator) toProto() *pb.Locator {
	if l == nil {
		return nil
	}
	p := &pb.Locator{
		Strategy: mapStrategy(l.strategy),
		Selector: l.selector,
		Index:    int32(l.index),
	}
	if l.parent != nil {
		p.Parent = l.parent.toProto()
	}
	for _, c := range l.constraints {
		p.Constraints = append(p.Constraints, &pb.RelativeConstraint{
			Direction: mapDirection(c.direction),
			Target:    c.target.toProto(),
		})
	}
	return p
}

// Nth filters multiple selections by zero-indexed values.
func (l *Locator) Nth(index int) *Locator {
	return &Locator{
		session:     l.session,
		strategy:    l.strategy,
		selector:    l.selector,
		index:       index,
		parent:      l.parent,
		constraints: l.constraints,
	}
}

// Locator creates a chained sub-locator.
func (l *Locator) Locator(strategy string, selector string) *Locator {
	return &Locator{
		session:  l.session,
		strategy: strategy,
		selector: selector,
		parent:   l,
	}
}

// GetByText locates elements containing matching text inside this locator.
func (l *Locator) GetByText(text string) *Locator {
	return l.Locator("TEXT", text)
}

// GetByRole locates elements matching the class role.
func (l *Locator) GetByRole(role string) *Locator {
	return l.Locator("ROLE", role)
}

// GetByLabel locates elements matching the accessibility label.
func (l *Locator) GetByLabel(label string) *Locator {
	return l.Locator("LABEL", label)
}

// GetByPlaceholder locates elements matching the input placeholder attribute.
func (l *Locator) GetByPlaceholder(placeholder string) *Locator {
	return l.Locator("PLACEHOLDER", placeholder)
}

// GetByAccessibilityID locates elements by content description or accessibility ID.
func (l *Locator) GetByAccessibilityID(id string) *Locator {
	return l.Locator("ACCESSIBILITY_ID", id)
}

// GetByTestID locates elements by test identifier.
func (l *Locator) GetByTestID(id string) *Locator {
	return l.Locator("TEST_ID", id)
}

// Above adds a positional constraint that target must be above the reference locator.
func (l *Locator) Above(other *Locator) *Locator {
	l.constraints = append(l.constraints, relativeConstraint{
		direction: "ABOVE",
		target:    other,
	})
	return l
}

// Below adds a positional constraint that target must be below the reference locator.
func (l *Locator) Below(other *Locator) *Locator {
	l.constraints = append(l.constraints, relativeConstraint{
		direction: "BELOW",
		target:    other,
	})
	return l
}

// LeftOf adds a positional constraint that target must be left of the reference locator.
func (l *Locator) LeftOf(other *Locator) *Locator {
	l.constraints = append(l.constraints, relativeConstraint{
		direction: "LEFT_OF",
		target:    other,
	})
	return l
}

// RightOf adds a positional constraint that target must be right of the reference locator.
func (l *Locator) RightOf(other *Locator) *Locator {
	l.constraints = append(l.constraints, relativeConstraint{
		direction: "RIGHT_OF",
		target:    other,
	})
	return l
}

// Click locates and triggers touch clicks on the target locator.
func (l *Locator) Click(ctx context.Context) error {
	// 1. Resolve element via FindElement RPC
	findResp, err := l.session.device.driverClient.FindElement(ctx, &pb.FindElementRequest{
		SessionId: l.session.sessionID,
		Locator:   l.toProto(),
		TimeoutMs: 5000,
	})
	if err != nil {
		return fmt.Errorf("failed to locate element for Click: %w", err)
	}

	if findResp.Element == nil {
		return fmt.Errorf("element not found for locator: %+v", l)
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
		Locator:   l.toProto(),
		TimeoutMs: 5000,
	})
	if err != nil {
		return fmt.Errorf("failed to locate element for Fill: %w", err)
	}

	if findResp.Element == nil {
		return fmt.Errorf("element not found for locator: %+v", l)
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
