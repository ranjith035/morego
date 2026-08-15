package sdk

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ranjith035/morego/proto/v1"
)

// Session represents an active run session delegate.
type Session struct {
	device    *Device
	sessionID string
	appID     string
}

// ID returns the active session identifier.
func (s *Session) ID() string {
	return s.sessionID
}

// Close terminates the session.
func (s *Session) Close(ctx context.Context) error {
	_, err := s.device.sessionClient.CloseSession(ctx, &pb.CloseSessionRequest{
		SessionId: s.sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}
	return nil
}

// Swipe dispatches swipes gestures directly.
func (s *Session) Swipe(ctx context.Context, startX, startY, endX, endY int, duration time.Duration) error {
	_, err := s.device.driverClient.Swipe(ctx, &pb.SwipeRequest{
		DriverId:   s.sessionID,
		Start:      &pb.Point{X: int32(startX), Y: int32(startY)},
		End:        &pb.Point{X: int32(endX), Y: int32(endY)},
		DurationMs: int32(duration.Milliseconds()),
	})
	if err != nil {
		return fmt.Errorf("swipe action failed: %w", err)
	}
	return nil
}

// Locator constructs a fluid element selector handle.
func (s *Session) Locator(strategy string, selector string) *Locator {
	return &Locator{
		session:  s,
		strategy: strategy,
		selector: selector,
		index:    0,
	}
}

// GetSource retrieves the current UI hierarchy tree layout.
func (s *Session) GetSource(ctx context.Context, format string) (string, error) {
	resp, err := s.device.driverClient.GetSource(ctx, &pb.GetSourceRequest{
		DriverId: s.sessionID,
		Format:   format,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch layout source: %w", err)
	}
	return resp.SourceData, nil
}
