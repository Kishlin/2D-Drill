package ui

import "github.com/Kishlin/drill-game/internal/domain/upgrades"

const (
	// Grid layout: 2 rows x 3 columns for Base + Mk1-Mk5
	UpgradeShopGridCols = 3
	UpgradeShopGridRows = 2
)

type UpgradeShopState struct {
	ActiveTab    upgrades.UpgradeType
	SelectedTier int
}

func NewUpgradeShopState() *UpgradeShopState {
	return &UpgradeShopState{
		ActiveTab:    upgrades.TypeEngine,
		SelectedTier: 0,
	}
}

func (s *UpgradeShopState) Reset() {
	s.ActiveTab = upgrades.TypeEngine
	s.SelectedTier = 0
}

func (s *UpgradeShopState) NextTab() {
	s.ActiveTab = upgrades.UpgradeType((int(s.ActiveTab) + 1) % int(upgrades.TypeCount))
	s.SelectedTier = 0
}

func (s *UpgradeShopState) PrevTab() {
	s.ActiveTab = upgrades.UpgradeType((int(s.ActiveTab) - 1 + int(upgrades.TypeCount)) % int(upgrades.TypeCount))
	s.SelectedTier = 0
}

func (s *UpgradeShopState) NavigateUp() {
	row := s.SelectedTier / UpgradeShopGridCols
	col := s.SelectedTier % UpgradeShopGridCols

	if row == 0 {
		row = UpgradeShopGridRows - 1
	} else {
		row--
	}

	s.SelectedTier = row*UpgradeShopGridCols + col
}

func (s *UpgradeShopState) NavigateDown() {
	row := s.SelectedTier / UpgradeShopGridCols
	col := s.SelectedTier % UpgradeShopGridCols

	if row == UpgradeShopGridRows-1 {
		row = 0
	} else {
		row++
	}

	s.SelectedTier = row*UpgradeShopGridCols + col
}

func (s *UpgradeShopState) NavigateLeft() {
	row := s.SelectedTier / UpgradeShopGridCols
	col := s.SelectedTier % UpgradeShopGridCols

	if col == 0 {
		col = UpgradeShopGridCols - 1
	} else {
		col--
	}

	s.SelectedTier = row*UpgradeShopGridCols + col
}

func (s *UpgradeShopState) NavigateRight() {
	row := s.SelectedTier / UpgradeShopGridCols
	col := s.SelectedTier % UpgradeShopGridCols

	if col == UpgradeShopGridCols-1 {
		col = 0
	} else {
		col++
	}

	s.SelectedTier = row*UpgradeShopGridCols + col
}

func (s *UpgradeShopState) GetSelectedRow() int {
	return s.SelectedTier / UpgradeShopGridCols
}

func (s *UpgradeShopState) GetSelectedCol() int {
	return s.SelectedTier % UpgradeShopGridCols
}

const (
	ItemShopGridCols  = 3
	ItemShopGridRows  = 2
	ItemShopEmptyCell = 5
)

type ItemShopState struct {
	SelectedIndex int
}

func NewItemShopState() *ItemShopState {
	return &ItemShopState{
		SelectedIndex: 0,
	}
}

func (s *ItemShopState) Reset() {
	s.SelectedIndex = 0
}

func (s *ItemShopState) NavigateUp() {
	row := s.SelectedIndex / ItemShopGridCols
	col := s.SelectedIndex % ItemShopGridCols

	if row == 0 {
		row = ItemShopGridRows - 1
	} else {
		row--
	}

	newIndex := row*ItemShopGridCols + col

	if newIndex == ItemShopEmptyCell {
		newIndex = ItemShopEmptyCell - ItemShopGridCols
	}

	s.SelectedIndex = newIndex
}

func (s *ItemShopState) NavigateDown() {
	row := s.SelectedIndex / ItemShopGridCols
	col := s.SelectedIndex % ItemShopGridCols

	if row == ItemShopGridRows-1 {
		row = 0
	} else {
		row++
	}

	newIndex := row*ItemShopGridCols + col

	if newIndex == ItemShopEmptyCell {
		newIndex = ItemShopEmptyCell - ItemShopGridCols
	}

	s.SelectedIndex = newIndex
}

func (s *ItemShopState) NavigateLeft() {
	row := s.SelectedIndex / ItemShopGridCols
	col := s.SelectedIndex % ItemShopGridCols

	if col == 0 {
		col = ItemShopGridCols - 1
	} else {
		col--
	}

	newIndex := row*ItemShopGridCols + col

	if newIndex == ItemShopEmptyCell {
		newIndex = ItemShopEmptyCell - 1
	}

	s.SelectedIndex = newIndex
}

func (s *ItemShopState) NavigateRight() {
	row := s.SelectedIndex / ItemShopGridCols
	col := s.SelectedIndex % ItemShopGridCols

	if col == ItemShopGridCols-1 {
		col = 0
	} else {
		col++
	}

	newIndex := row*ItemShopGridCols + col

	if newIndex == ItemShopEmptyCell {
		newIndex = ItemShopEmptyCell - 1
	}

	s.SelectedIndex = newIndex
}

func (s *ItemShopState) GetSelectedRow() int {
	return s.SelectedIndex / ItemShopGridCols
}

func (s *ItemShopState) GetSelectedCol() int {
	return s.SelectedIndex % ItemShopGridCols
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
