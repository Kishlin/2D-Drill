package bosses

import "fmt"

// BossConstructor creates a boss instance given room parameters
type BossConstructor func(roomStartY, worldWidth float32) Boss

// registry holds all registered boss constructors
var registry = make(map[string]BossConstructor)

// Register adds a boss constructor to the registry.
// Call this in your boss package's init() function.
func Register(bossType string, constructor BossConstructor) {
	if _, exists := registry[bossType]; exists {
		panic(fmt.Sprintf("boss type already registered: %s", bossType))
	}
	registry[bossType] = constructor
}

// Create instantiates a boss by type name using the registry.
func Create(bossType string, roomStartY, worldWidth float32) (Boss, error) {
	constructor, exists := registry[bossType]
	if exists == false {
		return nil, fmt.Errorf("unknown boss type: %s", bossType)
	}
	return constructor(roomStartY, worldWidth), nil
}
