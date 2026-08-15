package drivers

import (
	"fmt"
	"sync"
)

// DriverInfo describes registered driver capacities.
type DriverInfo struct {
	Type     string // e.g. "uiautomator2", "xcuitest"
	Platform string // e.g. "android", "ios"
}

// DriverFactory defines instantiation callback functions.
type DriverFactory func() Driver

// Registry manages thread-safe storage and resolution of drivers.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]DriverFactory
}

// NewRegistry constructs a driver registry instance.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]DriverFactory),
	}
}

// Register maps a driver identifier to its constructor factory.
func (r *Registry) Register(driverType string, factory DriverFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[driverType] = factory
}

// CreateDriver instantiates a driver based on its type identifier.
func (r *Registry) CreateDriver(driverType string) (Driver, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	factory, exists := r.factories[driverType]
	if !exists {
		return nil, fmt.Errorf("driver type %q is not registered", driverType)
	}
	return factory(), nil
}

// GetRegisteredDrivers returns the list of all registered driver types.
func (r *Registry) GetRegisteredDrivers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	drivers := make([]string, 0, len(r.factories))
	for k := range r.factories {
		drivers = append(drivers, k)
	}
	return drivers
}

// DefaultRegistry is the default driver registry instance.
var DefaultRegistry = NewRegistry()
