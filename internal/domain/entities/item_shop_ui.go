package entities

const (
	ItemShopGridCols  = 3
	ItemShopGridRows  = 2
	ItemShopEmptyCell = 5 // Position 5 is empty
)

type ItemShopUIState struct {
	Open          bool // Whether the item shop modal is open
	SelectedIndex int  // Currently selected item in the grid (0-4, skip 5)
}

func NewItemShopUIState() *ItemShopUIState {
	return &ItemShopUIState{
		Open:          false,
		SelectedIndex: 0,
	}
}

func (s *ItemShopUIState) OpenShop() {
	s.Open = true
	s.SelectedIndex = 0
}

func (s *ItemShopUIState) CloseShop() {
	s.Open = false
}

// NavigateUp moves selection up one row in the grid
// Skips the empty cell at position 5
func (s *ItemShopUIState) NavigateUp() {
	row := s.SelectedIndex / ItemShopGridCols
	col := s.SelectedIndex % ItemShopGridCols

	// Move up one row, wrap to bottom if at top
	if row == 0 {
		row = ItemShopGridRows - 1
	} else {
		row--
	}

	newIndex := row*ItemShopGridCols + col

	// If we landed on the empty cell (5), move to the cell above it (2)
	if newIndex == ItemShopEmptyCell {
		newIndex = ItemShopEmptyCell - ItemShopGridCols
	}

	s.SelectedIndex = newIndex
}

// NavigateDown moves selection down one row in the grid
// Skips the empty cell at position 5
func (s *ItemShopUIState) NavigateDown() {
	row := s.SelectedIndex / ItemShopGridCols
	col := s.SelectedIndex % ItemShopGridCols

	// Move down one row, wrap to top if at bottom
	if row == ItemShopGridRows-1 {
		row = 0
	} else {
		row++
	}

	newIndex := row*ItemShopGridCols + col

	// If we landed on the empty cell (5), move to the cell above it (2)
	if newIndex == ItemShopEmptyCell {
		newIndex = ItemShopEmptyCell - ItemShopGridCols
	}

	s.SelectedIndex = newIndex
}

// NavigateLeft moves selection left one column in the grid
// Skips the empty cell at position 5
func (s *ItemShopUIState) NavigateLeft() {
	row := s.SelectedIndex / ItemShopGridCols
	col := s.SelectedIndex % ItemShopGridCols

	// Move left one column, wrap to right edge if at left
	if col == 0 {
		col = ItemShopGridCols - 1
	} else {
		col--
	}

	newIndex := row*ItemShopGridCols + col

	// If we landed on the empty cell (5), move to the cell to its left (4)
	if newIndex == ItemShopEmptyCell {
		newIndex = ItemShopEmptyCell - 1
	}

	s.SelectedIndex = newIndex
}

// NavigateRight moves selection right one column in the grid
// Skips the empty cell at position 5
func (s *ItemShopUIState) NavigateRight() {
	row := s.SelectedIndex / ItemShopGridCols
	col := s.SelectedIndex % ItemShopGridCols

	// Move right one column, wrap to left edge if at right
	if col == ItemShopGridCols-1 {
		col = 0
	} else {
		col++
	}

	newIndex := row*ItemShopGridCols + col

	// If we landed on the empty cell (5), move to the cell to its left (4)
	if newIndex == ItemShopEmptyCell {
		newIndex = ItemShopEmptyCell - 1
	}

	s.SelectedIndex = newIndex
}

func (s *ItemShopUIState) GetSelectedRow() int {
	return s.SelectedIndex / ItemShopGridCols
}

func (s *ItemShopUIState) GetSelectedCol() int {
	return s.SelectedIndex % ItemShopGridCols
}
