package grid

import (
	"sync"
)

// FrameBroadcaster manages subscription registries distributing visual frames (MJPEG screenshots).
type FrameBroadcaster struct {
	mu          sync.Mutex
	subscribers map[chan []byte]bool
}

// NewFrameBroadcaster constructs a FrameBroadcaster.
func NewFrameBroadcaster() *FrameBroadcaster {
	return &FrameBroadcaster{
		subscribers: make(map[chan []byte]bool),
	}
}

// Subscribe returns a buffered channel for receiving frame bytes.
func (fb *FrameBroadcaster) Subscribe() chan []byte {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	ch := make(chan []byte, 10)
	fb.subscribers[ch] = true
	return ch
}

// Unsubscribe removes a client channel subscription.
func (fb *FrameBroadcaster) Unsubscribe(ch chan []byte) {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	if _, exists := fb.subscribers[ch]; exists {
		delete(fb.subscribers, ch)
		close(ch)
	}
}

// Broadcast dispatches frame payloads to all active subscriber channels non-blockingly.
func (fb *FrameBroadcaster) Broadcast(frame []byte) {
	fb.mu.Lock()
	defer fb.mu.Unlock()
	for ch := range fb.subscribers {
		select {
		case ch <- frame:
		default: // Drop frame if client buffer is full (slow reader protection)
		}
	}
}
