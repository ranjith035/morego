package core

import (
	"context"
)

// Platform defines the target mobile operating system.
type Platform string

const (
	PlatformAndroid Platform = "ANDROID"
	PlatformIOS     Platform = "IOS"
)

// Device details physical or virtual execution hardware.
type Device struct {
	ID         string
	Name       string
	Platform   Platform
	OSVersion  string
	IsEmulator bool
	Status     string
}

// DeviceManager controls connecting, installing, and executing apps on hardware targets.
type DeviceManager interface {
	Component

	// ListDevices lists connected mobile targets matching filters.
	ListDevices(ctx context.Context, platform Platform) ([]*Device, error)

	// GetDevice returns properties of a specific target.
	GetDevice(ctx context.Context, deviceID string) (*Device, error)

	// InstallApp deploys a binary app package to a target.
	InstallApp(ctx context.Context, deviceID string, appPath string) error

	// UninstallApp removes an application identifier package from the target.
	UninstallApp(ctx context.Context, deviceID string, appID string) error

	// LaunchApp executes an application process on a target device.
	LaunchApp(ctx context.Context, deviceID string, appID string, args []string, env map[string]string) (int, error)

	// TerminateApp halts the execution process of an application.
	TerminateApp(ctx context.Context, deviceID string, appID string) error
}
