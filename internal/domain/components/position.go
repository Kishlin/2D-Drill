package components

import "github.com/Kishlin/drill-game/internal/domain/types"

type Position struct {
	AABB types.AABB
}

func NewPosition(x, y, width, height float32) Position {
	return Position{AABB: types.NewAABB(x, y, width, height)}
}

func (p Position) Intersects(other Position) bool {
	return p.AABB.Intersects(other.AABB)
}
