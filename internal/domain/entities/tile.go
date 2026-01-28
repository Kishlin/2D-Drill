package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

type TileType int

const (
	TileTypeEmpty  TileType = iota // Air/empty space
	TileTypeDirt                   // Solid dirt (drillable)
	TileTypeOre                    // Solid ore (drillable, contains ore)
	TileTypeHazard                 // Hazard tile (drillability from config)
	TileTypeFloor                  // Floor (solid, not drillable, not nukeable)
)

type Tile struct {
	Type      TileType
	OreID     string // Only meaningful if Type == TileTypeOre (e.g., "copper", "gold")
	HazardID  string // Only meaningful if Type == TileTypeHazard (e.g., "rock", "lava")
	Drillable bool   // Only meaningful if Type == TileTypeHazard
}

func NewTile(tileType TileType) *Tile {
	return &Tile{Type: tileType}
}

func NewOreTileByID(oreID string) *Tile {
	return &Tile{Type: TileTypeOre, OreID: oreID}
}

func NewHazardTileByID(hazardID string, hazardCfg *config.HazardConfig) *Tile {
	drillable := true
	if hazardCfg != nil {
		drillable = hazardCfg.Drillable
	}
	return &Tile{Type: TileTypeHazard, HazardID: hazardID, Drillable: drillable}
}

func (t *Tile) IsSolid() bool {
	return t.Type != TileTypeEmpty
}

func (t *Tile) IsDrillable() bool {
	switch t.Type {
	case TileTypeDirt, TileTypeOre:
		return true
	case TileTypeHazard:
		return t.Drillable
	default:
		return false
	}
}

func (t *Tile) GetAABB(gridX, gridY int, tileSize float32) types.AABB {
	return types.AABB{
		X:      float32(gridX) * tileSize,
		Y:      float32(gridY) * tileSize,
		Width:  tileSize,
		Height: tileSize,
	}
}
