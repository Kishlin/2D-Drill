package world

type World struct {
	GroundLevel float32
	Width       float32
	Height      float32
}

func NewWorld(width, height, groundLevel float32) *World {
	return &World{
		Width:       width,
		Height:      height,
		GroundLevel: groundLevel,
	}
}

func (w *World) GetGroundLevel() float32 {
	return w.GroundLevel
}

func (w *World) IsInBounds(x, y float32) bool {
	return x >= 0 && x <= w.Width && y >= 0 && y <= w.Height
}
