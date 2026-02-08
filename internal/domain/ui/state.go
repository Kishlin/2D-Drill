package ui

import "github.com/Kishlin/drill-game/internal/domain/upgrades"

const (
	// Grid layout: 2 rows x 3 columns for Base + Mk1-Mk5
	UpgradeShopGridCols = 3
	UpgradeShopGridRows = 2
)

type UpgradeShopState struct {
	GridNavigator
	ActiveTab upgrades.UpgradeType
}

func NewUpgradeShopState() *UpgradeShopState {
	return &UpgradeShopState{
		GridNavigator: NewGridNavigator(UpgradeShopGridCols, UpgradeShopGridRows),
		ActiveTab:     upgrades.TypeEngine,
	}
}

func (s *UpgradeShopState) Reset() {
	s.ActiveTab = upgrades.TypeEngine
	s.Selected = 0
}

func (s *UpgradeShopState) NextTab() {
	s.ActiveTab = upgrades.UpgradeType((int(s.ActiveTab) + 1) % int(upgrades.TypeCount))
	s.Selected = 0
}

func (s *UpgradeShopState) PrevTab() {
	s.ActiveTab = upgrades.UpgradeType((int(s.ActiveTab) - 1 + int(upgrades.TypeCount)) % int(upgrades.TypeCount))
	s.Selected = 0
}

const (
	ItemShopGridCols  = 3
	ItemShopGridRows  = 2
	ItemShopEmptyCell = 5
)

type ItemShopState struct {
	GridNavigator
}

func NewItemShopState() *ItemShopState {
	return &ItemShopState{
		GridNavigator: NewGridNavigator(ItemShopGridCols, ItemShopGridRows),
	}
}

func (s *ItemShopState) Reset() {
	s.Selected = 0
}

func (s *ItemShopState) NavigateUp() {
	s.GridNavigator.NavigateUp()
	if s.Selected == ItemShopEmptyCell {
		s.Selected = ItemShopEmptyCell - ItemShopGridCols
	}
}

func (s *ItemShopState) NavigateDown() {
	s.GridNavigator.NavigateDown()
	if s.Selected == ItemShopEmptyCell {
		s.Selected = ItemShopEmptyCell - ItemShopGridCols
	}
}

func (s *ItemShopState) NavigateLeft() {
	s.GridNavigator.NavigateLeft()
	if s.Selected == ItemShopEmptyCell {
		s.Selected = ItemShopEmptyCell - 1
	}
}

func (s *ItemShopState) NavigateRight() {
	s.GridNavigator.NavigateRight()
	if s.Selected == ItemShopEmptyCell {
		s.Selected = ItemShopEmptyCell - 1
	}
}

type MarketState struct {
	FirstFrameTracker
}

func NewMarketState() *MarketState {
	return &MarketState{FirstFrameTracker: NewFirstFrameTracker()}
}

func (s *MarketState) Reset() {
	s.ResetFirstFrame()
}

const ModalServiceOptionCount = 4

type ModalServiceState struct {
	SelectedIndex int
	FirstFrameTracker
}

func NewModalServiceState() *ModalServiceState {
	return &ModalServiceState{
		SelectedIndex:     0,
		FirstFrameTracker: NewFirstFrameTracker(),
	}
}

func (s *ModalServiceState) Reset() {
	s.SelectedIndex = 0
	s.ResetFirstFrame()
}

func (s *ModalServiceState) NavigateUp() {
	s.SelectedIndex--
	if s.SelectedIndex < 0 {
		s.SelectedIndex = ModalServiceOptionCount - 1
	}
}

func (s *ModalServiceState) NavigateDown() {
	s.SelectedIndex++
	if s.SelectedIndex >= ModalServiceOptionCount {
		s.SelectedIndex = 0
	}
}

type InventoryState struct {
	FirstFrameTracker
}

func NewInventoryState() *InventoryState {
	return &InventoryState{FirstFrameTracker: NewFirstFrameTracker()}
}

func (s *InventoryState) Reset() {
	s.ResetFirstFrame()
}
