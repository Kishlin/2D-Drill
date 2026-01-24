package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

type TileType int

const (
	TileTypeEmpty TileType = iota // Air/empty space
	TileTypeDirt                  // Solid dirt (drillable)
	TileTypeOre                   // Solid ore (drillable, contains ore)
	TileTypeRock                  // Solid rock (impenetrable, not drillable)
	TileTypeLava                  // Lava (drillable, deals damage on completion)
	TileTypeFloor                 // Floor (solid, not drillable, not nukeable)
)

type Tile struct {
	Type     TileType
	OreID    string // Only meaningful if Type == TileTypeOre (e.g., "copper", "gold")
	HazardID string // Only meaningful if Type == TileTypeRock or TileTypeLava (e.g., "rock", "lava")
}

func NewTile(tileType TileType) *Tile {
	return &Tile{Type: tileType}
}

func NewOreTileByID(oreID string) *Tile {
	return &Tile{Type: TileTypeOre, OreID: oreID}
}

func NewHazardTileByID(hazardID string, hazardCfg *config.HazardConfig) *Tile {
	if hazardCfg != nil && hazardCfg.Drillable == false {
		return &Tile{Type: TileTypeRock, HazardID: hazardID}
	}
	return &Tile{Type: TileTypeLava, HazardID: hazardID}
}

func (t *Tile) IsSolid() bool {
	return t.Type != TileTypeEmpty
}

func (t *Tile) IsDrillable() bool {
	// Rock and Floor are NOT drillable
	// Lava IS drillable (deals damage on completion)
	return t.Type == TileTypeDirt || t.Type == TileTypeOre || t.Type == TileTypeLava
}

func (t *Tile) GetAABB(gridX, gridY int, tileSize float32) types.AABB {
	return types.AABB{
		X:      float32(gridX) * tileSize,
		Y:      float32(gridY) * tileSize,
		Width:  tileSize,
		Height: tileSize,
	}
}
