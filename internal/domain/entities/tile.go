package entities

type TileType int

const (
	TileTypeEmpty TileType = iota // Air/empty space
	TileTypeDirt                   // Solid dirt (diggable)
)

type Tile struct {
	Type TileType
}

func NewTile(tileType TileType) *Tile {
	return &Tile{Type: tileType}
}

func (t *Tile) IsSolid() bool {
	return t.Type != TileTypeEmpty
}

func (t *Tile) IsDiggable() bool {
	return t.Type == TileTypeDirt
}
