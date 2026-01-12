package entities

type HazardType int

const (
	HazardRock HazardType = iota
	HazardLava
)

type HazardMetadata struct {
	PeakDepth float32 // Tile Y coordinate where hazard is most common
	Sigma     float32 // Standard deviation (spread of distribution)
	MaxWeight float32 // Weight at peak depth (relative spawn chance)
}

var HazardDistributions = map[HazardType]HazardMetadata{
	HazardRock: {
		PeakDepth: 650.0, // ~80% depth
		Sigma:     200.0, // Wide spread to start at ~40%
		MaxWeight: 15.0,  // High to dominate at depth
	},
	HazardLava: {
		PeakDepth: 750.0, // ~85% depth
		Sigma:     150.0, // Narrower, starts ~60-65%
		MaxWeight: 12.0,  // High to dominate at depth
	},
}

func GetAllHazardTypes() []HazardType {
	return []HazardType{HazardRock, HazardLava}
}

func GetHazardName(hazardType HazardType) string {
	switch hazardType {
	case HazardRock:
		return "Rock"
	case HazardLava:
		return "Lava"
	default:
		return "Unknown"
	}
}
