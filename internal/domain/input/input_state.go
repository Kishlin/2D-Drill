package input

// InputState represents platform-agnostic input state
type InputState struct {
	Left        bool
	Right       bool
	Up          bool
	Drill       bool // Down for drilling
	Sell        bool // E key for selling at market
	UseTeleport bool // T key for teleport item
	UseRepair   bool // R key for repair item
	UseRefuel   bool // F key for refuel item
	UseBomb     bool // B key for bomb item
	UseBigBomb  bool // G key for big bomb item
}

func NewInputState() InputState {
	return InputState{
		Left:        false,
		Right:       false,
		Up:          false,
		Drill:       false,
		Sell:        false,
		UseTeleport: false,
		UseRepair:   false,
		UseRefuel:   false,
		UseBomb:     false,
		UseBigBomb:  false,
	}
}

func (is InputState) HasHorizontalInput() bool {
	return is.Left || is.Right
}

func (is InputState) HasVerticalInput() bool {
	return is.Up
}

// HasMovementInput returns true if player is actively moving or drilling
// (Left, Right, Up, or Drill - but NOT interactions)
func (is InputState) HasMovementInput() bool {
	return is.Left || is.Right || is.Up || is.Drill
}
