package core

import (
	"context"
	"time"
)

// SessionStatus represents the operational status of an active automation session.
type SessionStatus string

const (
	SessionStarting SessionStatus = "STARTING"
	SessionActive   SessionStatus = "ACTIVE"
	SessionClosed   SessionStatus = "CLOSED"
	SessionError    SessionStatus = "ERROR"
)

// Session holds connection metadata representing a running test execution.
type Session struct {
	ID           string
	Status       SessionStatus
	Capabilities map[string]string
	CreatedAt    time.Time
	DeviceID     string
	AppID        string
}

// SessionManager orchestrates the creation, tracking, and clean release of automation sessions.
type SessionManager interface {
	Component

	// CreateSession starts a new session with capabilities on a target device.
	CreateSession(ctx context.Context, deviceID string, appID string, caps map[string]string) (*Session, error)

	// CloseSession terminates an active session and releases underlying drivers.
	CloseSession(ctx context.Context, sessionID string) error

	// GetSession retrieves details of an active session.
	GetSession(ctx context.Context, sessionID string) (*Session, error)

	// ListSessions lists all sessions managed by the engine.
	ListSessions(ctx context.Context) ([]*Session, error)
}
