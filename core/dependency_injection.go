package core

// Container acts as the dependency injection registration book.
type Container interface {
	// Register binds a module implementation instance to a name registry key.
	Register(name string, value interface{})

	// Resolve searches for registered modules matching a key.
	Resolve(name string) (interface{}, bool)
}
