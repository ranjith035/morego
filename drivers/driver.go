package drivers

import (
	"context"
	"time"
)

// Driver defines the execution interface that all native/remote automation agents must implement.
type Driver interface {
	// Connect establishes the automation session based on capabilities.
	Connect(ctx context.Context, capabilities map[string]string) error

	// Disconnect shuts down the session and releases target processes.
	Disconnect(ctx context.Context) error

	// Click executes a click gesture at the center of the resolved element.
	Click(ctx context.Context, elementID string) error

	// ClickAt dispatches a click gesture directly at coordinate offsets.
	ClickAt(ctx context.Context, x, y int) error

	// Fill inputs keyboard characters into the target input element.
	Fill(ctx context.Context, elementID string, value string) error

	// Swipe executes linear swipe gestures between start and end bounds.
	Swipe(ctx context.Context, startX, startY, endX, endY int, duration time.Duration) error

	// Screenshot captures the screen or specific element layout coordinates.
	Screenshot(ctx context.Context, elementID string) ([]byte, error)

	// FindElement searches the active UI hierarchy and returns a unique element reference.
	FindElement(ctx context.Context, strategy string, selector string) (string, error)

	// LaunchApp runs the target app package on the device.
	LaunchApp(ctx context.Context, appID string, args []string, env map[string]string) error

	// TerminateApp forces execution stop of the package.
	TerminateApp(ctx context.Context, appID string) error

	// InstallApp deploys app binaries to the device target.
	InstallApp(ctx context.Context, appPath string) error

	// UninstallApp removes target app binaries.
	UninstallApp(ctx context.Context, appID string) error

	// GetSource dumps layout nodes representation.
	GetSource(ctx context.Context, format string) (string, error)

	// ExecuteScript evaluates macros or scripts inside the runtime driver context.
	ExecuteScript(ctx context.Context, script string, arguments []string) (string, error)
}
