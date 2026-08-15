package core

import "context"

// Config aggregates structural setup configurations.
type Config struct {
	Port              int
	DriverSearchPaths []string
	DefaultTimeoutMS  int
	LogLevel          string
	MetricsEnabled    bool
	TraceEnabled      bool
}

// ConfigurationProvider reads, holds, and syncs runtime engine properties.
type ConfigurationProvider interface {
	Component

	// GetConfig returns the active configuration.
	GetConfig() *Config

	// LoadConfig parses property mappings from file locations.
	LoadConfig(ctx context.Context, path string) error

	// SaveConfig serializes current variables back to a file.
	SaveConfig(ctx context.Context, path string) error
}
