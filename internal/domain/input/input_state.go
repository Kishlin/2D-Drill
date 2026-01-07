package input

// InputState represents platform-agnostic input state
type InputState struct {
	Left  bool
	Right bool
	Up    bool
	Drill bool // Down for drilling
	Sell  bool // E key for selling at market
}

func NewInputState() InputState {
	return InputState{
		Left:  false,
		Right: false,
		Up:    false,
		Drill: false,
		Sell:  false,
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
