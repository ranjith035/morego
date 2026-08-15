package core

import "context"

// Event wraps publisher payloads passing through internal message routes.
type Event struct {
	Topic     string
	Payload   interface{}
	Timestamp int64
}

// Subscription tracks channel routing rules for subscribers.
type Subscription interface {
	// Topic returning the subscribed path.
	Topic() string

	// Channel reads messages incoming.
	Channel() <-chan Event

	// Close drops subscription hooks.
	Close() error
}

// EventBus is the core decoupled event broker.
type EventBus interface {
	Component

	// Publish broadcasts an event message.
	Publish(ctx context.Context, event Event) error

	// Subscribe requests a channel handle listening for topic occurrences.
	Subscribe(ctx context.Context, topic string) (Subscription, error)
}
