package core

import "context"

// Component defines the lifecycle contract for all modular subsystems.
type Component interface {
	// Name returns a unique identifier for the component.
	Name() string

	// Init initializes the component's internal state and resources.
	Init(ctx context.Context) error

	// Shutdown gracefully cleans up allocated resources.
	Shutdown(ctx context.Context) error
}
