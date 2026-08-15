package core

import "context"

// ExecutionEngine acts as the central orchestration daemon uniting all components.
type ExecutionEngine interface {
	Component

	// Container returns the DI container wiring internal engines.
	Container() Container

	// StartServer runs the engine gRPC listener.
	StartServer(ctx context.Context, port int) error
}
