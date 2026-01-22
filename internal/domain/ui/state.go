package ui

import "github.com/Kishlin/drill-game/internal/domain/entities"

const (
	// Grid layout: 2 rows x 3 columns for Base + Mk1-Mk5
	UpgradeShopGridCols = 3
	UpgradeShopGridRows = 2
)

type UpgradeShopState struct {
	ActiveTab    entities.UpgradeType
	SelectedTier int
}

func NewUpgradeShopState() *UpgradeShopState {
	return &UpgradeShopState{
		ActiveTab:    entities.UpgradeEngine,
		SelectedTier: 0,
	}
}

func (s *UpgradeShopState) Reset() {
	s.ActiveTab = entities.UpgradeEngine
	s.SelectedTier = 0
}

func (s *UpgradeShopState) NextTab() {
	s.ActiveTab = entities.UpgradeType((int(s.ActiveTab) + 1) % int(entities.UpgradeTypeCount))
	s.SelectedTier = 0
}

func (s *UpgradeShopState) PrevTab() {
	s.ActiveTab = entities.UpgradeType((int(s.ActiveTab) - 1 + int(entities.UpgradeTypeCount)) % int(entities.UpgradeTypeCount))
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
