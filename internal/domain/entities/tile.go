package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
)

type TileType int

const (
	TileTypeEmpty TileType = iota // Air/empty space
	TileTypeDirt                   // Solid dirt (drillable)
	TileTypeOre                    // Solid ore (drillable, contains ore)
	TileTypeRock                   // Solid rock (impenetrable, not drillable)
	TileTypeLava                   // Lava (drillable, deals damage on completion)
)

type Tile struct {
	Type       TileType
	OreType    OreType    // Only meaningful if Type == TileTypeOre
	HazardType HazardType // Only meaningful if Type == TileTypeRock or TileTypeLava
}

func NewTile(tileType TileType) *Tile {
	return &Tile{Type: tileType}
}

func NewOreTile(oreType OreType) *Tile {
	return &Tile{Type: TileTypeOre, OreType: oreType}
}

func NewHazardTile(hazardType HazardType) *Tile {
	if hazardType == HazardRock {
		return &Tile{Type: TileTypeRock, HazardType: hazardType}
	}
	return &Tile{Type: TileTypeLava, HazardType: hazardType}
}

func (t *Tile) IsSolid() bool {
	return t.Type != TileTypeEmpty
}

func (t *Tile) IsDrillable() bool {
	// Rock is NOT drillable, but Lava IS drillable (deals damage on completion)
	return t.Type == TileTypeDirt || t.Type == TileTypeOre || t.Type == TileTypeLava
}

// GetAABB returns the tile's bounding box at given grid coordinates
func (t *Tile) GetAABB(gridX, gridY int, tileSize float32) types.AABB {
	return types.AABB{
		X:      float32(gridX) * tileSize,
		Y:      float32(gridY) * tileSize,
		Width:  tileSize,
		Height: tileSize,
	}
}
