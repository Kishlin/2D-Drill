package config

type FloorType string

const (
	FloorConcrete FloorType = "concrete" // Solid, walkable
	FloorLava     FloorType = "lava"     // Deals damage on contact
)

type BossRoomConfig struct {
	BossType    string    // e.g., "test_boss", "earth_guardian"
	FloorType   FloorType // Floor behavior
	FloorDamage float32   // Damage per frame when standing on hazardous floor (e.g., lava)
	RoomHeight  float32   // Boss room height in pixels (~720)
	FloorHeight float32   // Floor tiles below room (in tile count, will be multiplied by tile size)
}
