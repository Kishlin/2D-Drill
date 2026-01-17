package input

// InputState represents platform-agnostic input state
type InputState struct {
	// Continuous inputs (held down for movement)
	Left  bool
	Right bool
	Up    bool
	Drill bool // Down for drilling

	// Discrete inputs (single press actions)
	Sell        bool // E key for selling at market / interact / purchase
	UseTeleport bool // T key for teleport item
	UseRepair   bool // R key for repair item
	UseRefuel   bool // F key for refuel item
	UseBomb     bool // B key for bomb item
	UseBigBomb  bool // G key for big bomb item
	PrevTab     bool // Z key for previous tab in shop
	NextTab     bool // X key for next tab in shop
	CloseShop   bool // Q or Escape key to close shop

	// Discrete navigation (for UI, single press)
	NavLeft  bool // Left arrow or A key (discrete)
	NavRight bool // Right arrow or D key (discrete)
	NavUp    bool // Up arrow or W key (discrete)
	NavDown  bool // Down arrow or S key (discrete)
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
		PrevTab:     false,
		NextTab:     false,
		CloseShop:   false,
		NavLeft:     false,
		NavRight:    false,
		NavUp:       false,
		NavDown:     false,
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
