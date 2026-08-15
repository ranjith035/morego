package core

import "context"

// Plugin describes interceptor hooks executed at test lifecycle steps.
type Plugin interface {
	Name() string
	Version() string

	// OnBeforeAction intercepts automation commands prior to execution. Returns false to halt.
	OnBeforeAction(ctx context.Context, actionName string, payload map[string]string) (bool, error)

	// OnAfterAction intercepts updates after actions complete.
	OnAfterAction(ctx context.Context, actionName string, payload map[string]string) error
}

// PluginLoader dynamically loads and binds hooks extensions.
type PluginLoader interface {
	Component

	// LoadPlugin mounts dynamic plugins from file locations.
	LoadPlugin(ctx context.Context, path string) (Plugin, error)

	// UnloadPlugin releases a loaded plugin.
	UnloadPlugin(ctx context.Context, name string) error

	// GetPlugins returns references to all loaded plugins.
	GetPlugins() []Plugin
}
