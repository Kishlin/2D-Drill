package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type InventoryUI struct {
	oreConfigs []config.OreConfig
	state      *InventoryState
	active     bool
}

func NewInventoryUI(oreConfigs []config.OreConfig) *InventoryUI {
	return &InventoryUI{
		oreConfigs: oreConfigs,
		state:      NewInventoryState(),
		active:     false,
	}
}

// Process handles input for the inventory UI and returns whether it should close
func (u *InventoryUI) Process(inputState input.InputState) bool {
	// Close on Q/Escape
	if inputState.CloseShop {
		u.active = false
		return true
	}

	// Close on movement input
	if inputState.HasMovementInput() {
		u.active = false
		return true
	}

	// Skip the first frame to avoid processing the I that opened the inventory
	if u.state.IsFirstFrame() {
		u.state.ClearFirstFrame()
		return false
	}

	// Toggle close on I press (after first frame)
	if inputState.Inventory {
		u.active = false
		return true
	}

	return false
}

func (u *InventoryUI) Open() {
	u.state.Reset()
	u.active = true
}

func (u *InventoryUI) IsActive() bool {
	return u.active
}

func (u *InventoryUI) GetRenderState() *InventoryState {
	return u.state
}

func (u *InventoryUI) GetOreConfigs() []config.OreConfig {
	return u.oreConfigs
}

func (u *InventoryUI) ResetState() {
	u.state.Reset()
}
