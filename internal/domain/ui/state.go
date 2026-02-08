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
	firstFrame bool
}

func NewMarketState() *MarketState {
	return &MarketState{firstFrame: true}
}

func (s *MarketState) Reset() {
	s.firstFrame = true
}

func (s *MarketState) IsFirstFrame() bool {
	return s.firstFrame
}

func (s *MarketState) ClearFirstFrame() {
	s.firstFrame = false
}

const HospitalOptionCount = 4

type HospitalState struct {
	SelectedIndex int
	firstFrame    bool
}

func NewHospitalState() *HospitalState {
	return &HospitalState{
		SelectedIndex: 0,
		firstFrame:    true,
	}
}

func (s *HospitalState) Reset() {
	s.SelectedIndex = 0
	s.firstFrame = true
}

func (s *HospitalState) NavigateUp() {
	s.SelectedIndex--
	if s.SelectedIndex < 0 {
		s.SelectedIndex = HospitalOptionCount - 1
	}
}

func (s *HospitalState) NavigateDown() {
	s.SelectedIndex++
	if s.SelectedIndex >= HospitalOptionCount {
		s.SelectedIndex = 0
	}
}

func (s *HospitalState) IsFirstFrame() bool {
	return s.firstFrame
}

func (s *HospitalState) ClearFirstFrame() {
	s.firstFrame = false
}

const FuelStationOptionCount = 4

type FuelStationState struct {
	SelectedIndex int
	firstFrame    bool
}

func NewFuelStationState() *FuelStationState {
	return &FuelStationState{
		SelectedIndex: 0,
		firstFrame:    true,
	}
}

func (s *FuelStationState) Reset() {
	s.SelectedIndex = 0
	s.firstFrame = true
}

func (s *FuelStationState) NavigateUp() {
	s.SelectedIndex--
	if s.SelectedIndex < 0 {
		s.SelectedIndex = FuelStationOptionCount - 1
	}
}

func (s *FuelStationState) NavigateDown() {
	s.SelectedIndex++
	if s.SelectedIndex >= FuelStationOptionCount {
		s.SelectedIndex = 0
	}
}

func (s *FuelStationState) IsFirstFrame() bool {
	return s.firstFrame
}

func (s *FuelStationState) ClearFirstFrame() {
	s.firstFrame = false
}

type InventoryState struct {
	firstFrame bool
}

func NewInventoryState() *InventoryState {
	return &InventoryState{firstFrame: true}
}

func (s *InventoryState) Reset() {
	s.firstFrame = true
}

func (s *InventoryState) IsFirstFrame() bool {
	return s.firstFrame
}

func (s *InventoryState) ClearFirstFrame() {
	s.firstFrame = false
}
