package worlds

import (
	"fmt"
	"log/slog"
)

// GetWorld retrieves a world configuration by name
// Returns an error if the world is not registered
func GetWorld(name string) (*WorldGameConfig, error) {
	config, exists := worldRegistry[name]
	if !exists {
		return nil, fmt.Errorf("world '%s' not found in registry. Available worlds: %v", name, ListWorlds())
	}
	// Return a copy to prevent accidental modifications
	return config, nil
}

// MustGetWorld retrieves a world configuration by name, panicking if not found
// Useful for application initialization when the world must exist
func MustGetWorld(name string) *WorldGameConfig {
	config, err := GetWorld(name)
	if err != nil {
		slog.Error("Failed to load world", "name", name, "error", err)
		panic(err)
	}
	return config
}

// ListWorlds returns a list of all registered world names
func ListWorlds() []string {
	names := make([]string, 0, len(worldRegistry))
	for name := range worldRegistry {
		names = append(names, name)
	}
	return names
}
