package entities

// OreType represents different ore types found in the world
type OreType int

const (
	OreCopper OreType = iota
	OreIron
	OreSilver
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
	OreCopper:   {PeakDepth: -50, Sigma: 150, MaxWeight: 10.0},
	OreIron:     {PeakDepth: 0, Sigma: 200, MaxWeight: 8.0},
	OreSilver:   {PeakDepth: 150, Sigma: 180, MaxWeight: 6.0},
	OreGold:     {PeakDepth: 300, Sigma: 150, MaxWeight: 4.0},
	OreMythril:  {PeakDepth: 500, Sigma: 200, MaxWeight: 3.0},
	OrePlatinum: {PeakDepth: 700, Sigma: 180, MaxWeight: 2.0},
	OreDiamond:  {PeakDepth: 900, Sigma: 150, MaxWeight: 1.0},
}

// OreValues maps each ore type to its sell value in currency
var OreValues = map[OreType]int{
	OreCopper:   10,
	OreIron:     25,
	OreSilver:   75,
	OreGold:     250,
	OreMythril:  1000,
	OrePlatinum: 5000,
	OreDiamond:  30000,
}

// GetAllOreTypes returns all ore types for iteration
func GetAllOreTypes() []OreType {
	return []OreType{
		OreCopper,
		OreIron,
		OreSilver,
		OreGold,
		OreMythril,
		OrePlatinum,
		OreDiamond,
	}
}

// CalculateInventoryValue calculates total sell value of an ore inventory
func CalculateInventoryValue(inventory [7]int) int {
	total := 0
	for oreType, count := range inventory {
		if count > 0 {
			total += OreValues[OreType(oreType)] * count
		}
	}
	return total
}
