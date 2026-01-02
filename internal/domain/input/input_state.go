package input

// InputState represents platform-agnostic input state
type InputState struct {
	Left  bool
	Right bool
	Up    bool
	Dig   bool // Down for digging
	Sell  bool // E key for selling at shop
}

func NewInputState() InputState {
	return InputState{
		Left:  false,
		Right: false,
		Up:    false,
		Dig:   false,
		Sell:  false,
	}
}

func (is InputState) HasHorizontalInput() bool {
	return is.Left || is.Right
}

func (is InputState) HasVerticalInput() bool {
	return is.Up
}

// HasMovementInput returns true if player is actively moving or digging
// (Left, Right, Up, or Dig - but NOT interactions)
func (is InputState) HasMovementInput() bool {
	return is.Left || is.Right || is.Up || is.Dig
}
