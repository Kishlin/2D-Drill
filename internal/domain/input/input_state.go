package input

// InputState represents platform-agnostic input state
type InputState struct {
	Left  bool
	Right bool
	Up    bool
	Down  bool
}

func NewInputState() InputState {
	return InputState{
		Left:  false,
		Right: false,
		Up:    false,
		Down:  false,
	}
}

func (is InputState) HasHorizontalInput() bool {
	return is.Left || is.Right
}

func (is InputState) HasVerticalInput() bool {
	return is.Up || is.Down
}
