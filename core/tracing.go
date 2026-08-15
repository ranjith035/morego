package core

import "context"

// Span tracks scope durations of single operations.
type Span interface {
	// End stops timing resolution and registers trace results.
	End()
}

// Tracer controls spawning spans and mapping execution headers over network bounds.
type Tracer interface {
	Component

	// StartSpan registers a nested trace span.
	StartSpan(ctx context.Context, name string) (context.Context, Span)

	// InjectContext writes trace context variables into headers map carriers.
	InjectContext(ctx context.Context, carrier map[string]string) error

	// ExtractContext parses active trace IDs out of header maps.
	ExtractContext(ctx context.Context, carrier map[string]string) (context.Context, error)
}
