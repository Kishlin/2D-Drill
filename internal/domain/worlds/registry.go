package worlds

// worldRegistry stores all registered world configurations
// Worlds can be registered via RegisterWorld() or via init() functions in example packages
var worldRegistry = make(map[string]*WorldGameConfig)

// RegisterWorld adds a world configuration to the registry
// This is called by world definition packages (e.g., examples/hard_mode.go) via init()
func RegisterWorld(name string, config *WorldGameConfig) {
	worldRegistry[name] = config
}

// InitializeDefaultWorlds registers all built-in world configurations
// Called during package initialization to ensure default world is always available
func init() {
	RegisterWorld("default", NewDefaultWorld())
}
