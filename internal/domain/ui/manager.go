package ui

import (
	"github.com/Kishlin/drill-game/internal/domain/components"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type Manager struct {
	uis        map[components.InteractableType]UI
	activeUI   UI
	activeType components.InteractableType
	isActive   bool
}

func NewManager() *Manager {
	return &Manager{
		uis:      make(map[components.InteractableType]UI),
		isActive: false,
	}
}

func (m *Manager) Register(t components.InteractableType, ui UI) {
	m.uis[t] = ui
}

func (m *Manager) OpenUI(t components.InteractableType) bool {
	if ui, ok := m.uis[t]; ok {
		m.activeUI = ui
		m.activeType = t
		m.isActive = true
		return true
	}
	return false
}

func (m *Manager) Process(player *entities.Player, inputState input.InputState) Result {
	if !m.isActive || m.activeUI == nil {
		return NoChange()
	}

	result := m.activeUI.Process(player, inputState)
	if result.ShouldClose {
		m.activeUI = nil
		m.isActive = false
	}
	return result
}

func (m *Manager) HasActiveUI() bool {
	return m.isActive
}

func (m *Manager) GetActiveUI() UI {
	return m.activeUI
}

func (m *Manager) GetActiveType() components.InteractableType {
	return m.activeType
}

func (m *Manager) GetRegisteredUI(t components.InteractableType) UI {
	return m.uis[t]
}
