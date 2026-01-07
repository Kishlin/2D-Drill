package entities

// OreType represents different ore types found in the world
type OreType int

const (
	OreCopper OreType = iota
	OreIron
	OreGold
	OreMythril
	OrePlatinum
	OreDiamond
)

// OreMetadata contains Gaussian distribution parameters for ore generation
type OreMetadata struct {
	PeakDepth float32 // Tile Y coordinate where ore is most common
	Sigma     float32 // Standard deviation (spread of distribution)
	MaxWeight float32 // Weight at peak depth (relative spawn chance)
}

// OreDistributions maps each ore type to its generation parameters
var OreDistributions = map[OreType]OreMetadata{
	OreCopper:   {PeakDepth: -75, Sigma: 120, MaxWeight: 8.0},
	OreIron:     {PeakDepth: 70, Sigma: 90, MaxWeight: 5.0},
	OreGold:     {PeakDepth: 230, Sigma: 80, MaxWeight: 3.0},
	OreMythril:  {PeakDepth: 360, Sigma: 70, MaxWeight: 2.2},
	OrePlatinum: {PeakDepth: 500, Sigma: 80, MaxWeight: 1.8},
	OreDiamond:  {PeakDepth: 600, Sigma: 180, MaxWeight: 0.15},
}

// OreValues maps each ore type to its sell value in currency
var OreValues = map[OreType]int{
	OreCopper:   25,
	OreIron:     75,
	OreGold:     300,
	OreMythril:  1500,
	OrePlatinum: 10000,
	OreDiamond:  30000,
}

// GetAllOreTypes returns all ore types for iteration
func GetAllOreTypes() []OreType {
	return []OreType{
		OreCopper,
		OreIron,
		OreGold,
		OreMythril,
		OrePlatinum,
		OreDiamond,
	}
}

// CalculateInventoryValue calculates total sell value of an ore inventory
func CalculateInventoryValue(inventory [6]int) int {
	total := 0
	for oreType, count := range inventory {
		if count > 0 {
			total += OreValues[OreType(oreType)] * count
		}
	}
	return total
}
