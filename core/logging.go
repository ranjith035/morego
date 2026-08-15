package core

// Fields structures metadata logs context.
type Fields map[string]interface{}

// Logger provides structured diagnostics logging contracts.
type Logger interface {
	Component

	// Debug writes verbose logging traces.
	Debug(msg string, args ...interface{})

	// Info writes general system metrics notices.
	Info(msg string, args ...interface{})

	// Warn writes recoverable issue alerts.
	Warn(msg string, args ...interface{})

	// Error writes execution failures.
	Error(msg string, args ...interface{})

	// WithFields spawns a nested logger contextualized with fields metadata.
	WithFields(fields Fields) Logger
}
