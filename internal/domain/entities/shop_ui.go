package entities

const (
	// Grid layout: 2 rows x 3 columns for Base + Mk1-Mk5
	ShopGridCols = 3
	ShopGridRows = 2
	ShopTierCount = 6 // Base (0) + Mk1-Mk5 (1-5)
)

// ShopUIState tracks the current state of the shop modal UI
type ShopUIState struct {
	Open         bool        // Whether the shop modal is open
	ActiveTab    UpgradeType // Currently selected upgrade category (0-5)
	SelectedTier int         // Currently selected tier in the grid (0-5)
}

// NewShopUIState creates a new shop UI state (closed by default)
func NewShopUIState() *ShopUIState {
	return &ShopUIState{
		Open:         false,
		ActiveTab:    UpgradeEngine,
		SelectedTier: 0,
	}
}

// OpenShop opens the modal and resets selection to first available upgrade
func (s *ShopUIState) OpenShop() {
	s.Open = true
	s.ActiveTab = UpgradeEngine
	s.SelectedTier = 0
}

// CloseShop closes the modal
func (s *ShopUIState) CloseShop() {
	s.Open = false
}

// NextTab cycles to the next upgrade type (wraps around)
func (s *ShopUIState) NextTab() {
	s.ActiveTab = UpgradeType((int(s.ActiveTab) + 1) % int(UpgradeTypeCount))
	// Reset selection to first tier when changing tabs
	s.SelectedTier = 0
}

// PrevTab cycles to the previous upgrade type (wraps around)
func (s *ShopUIState) PrevTab() {
	s.ActiveTab = UpgradeType((int(s.ActiveTab) - 1 + int(UpgradeTypeCount)) % int(UpgradeTypeCount))
	// Reset selection to first tier when changing tabs
	s.SelectedTier = 0
}

// NavigateUp moves selection up one row in the grid
func (s *ShopUIState) NavigateUp() {
	row := s.SelectedTier / ShopGridCols
	col := s.SelectedTier % ShopGridCols

	// Move up one row, wrap to bottom if at top
	if row == 0 {
		row = ShopGridRows - 1
	} else {
		row--
	}

	s.SelectedTier = row*ShopGridCols + col
}

// NavigateDown moves selection down one row in the grid
func (s *ShopUIState) NavigateDown() {
	row := s.SelectedTier / ShopGridCols
	col := s.SelectedTier % ShopGridCols

	// Move down one row, wrap to top if at bottom
	if row == ShopGridRows-1 {
		row = 0
	} else {
		row++
	}

	s.SelectedTier = row*ShopGridCols + col
}

// NavigateLeft moves selection left one column in the grid
func (s *ShopUIState) NavigateLeft() {
	row := s.SelectedTier / ShopGridCols
	col := s.SelectedTier % ShopGridCols

	// Move left one column, wrap to right edge if at left
	if col == 0 {
		col = ShopGridCols - 1
	} else {
		col--
	}

	s.SelectedTier = row*ShopGridCols + col
}

// NavigateRight moves selection right one column in the grid
func (s *ShopUIState) NavigateRight() {
	row := s.SelectedTier / ShopGridCols
	col := s.SelectedTier % ShopGridCols

	// Move right one column, wrap to left edge if at right
	if col == ShopGridCols-1 {
		col = 0
	} else {
		col++
	}

	s.SelectedTier = row*ShopGridCols + col
}

// GetSelectedRow returns the row (0 or 1) of the current selection
func (s *ShopUIState) GetSelectedRow() int {
	return s.SelectedTier / ShopGridCols
}

// GetSelectedCol returns the column (0, 1, or 2) of the current selection
func (s *ShopUIState) GetSelectedCol() int {
	return s.SelectedTier % ShopGridCols
}
