package entities

import (
	"github.com/Kishlin/drill-game/internal/domain/types"
)

type TileType int

const (
	TileTypeEmpty TileType = iota // Air/empty space
	TileTypeDirt                   // Solid dirt (diggable)
	TileTypeOre                    // Solid ore (diggable, contains ore)
)

type Tile struct {
	Type    TileType
	OreType OreType // Only meaningful if Type == TileTypeOre
}

func NewTile(tileType TileType) *Tile {
	return &Tile{Type: tileType}
}

func NewOreTile(oreType OreType) *Tile {
	return &Tile{Type: TileTypeOre, OreType: oreType}
}

func (t *Tile) IsSolid() bool {
	return t.Type != TileTypeEmpty
}

func (t *Tile) IsDiggable() bool {
	return t.Type == TileTypeDirt || t.Type == TileTypeOre
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
