package entities

const (
	// Grid layout: 2 rows x 3 columns for Base + Mk1-Mk5
	ShopGridCols = 3
	ShopGridRows = 2
)

type UpgradeShopUIState struct {
	Open         bool        // Whether the shop modal is open
	ActiveTab    UpgradeType // Currently selected upgrade category (0-5)
	SelectedTier int         // Currently selected tier in the grid (0-5)
}

func NewUpgradeShopUIState() *UpgradeShopUIState {
	return &UpgradeShopUIState{
		Open:         false,
		ActiveTab:    UpgradeEngine,
		SelectedTier: 0,
	}
}

func (s *UpgradeShopUIState) OpenShop() {
	s.Open = true
	s.ActiveTab = UpgradeEngine
	s.SelectedTier = 0
}

func (s *UpgradeShopUIState) CloseShop() {
	s.Open = false
}

func (s *UpgradeShopUIState) NextTab() {
	s.ActiveTab = UpgradeType((int(s.ActiveTab) + 1) % int(UpgradeTypeCount))
	// Reset selection to first tier when changing tabs
	s.SelectedTier = 0
}

func (s *UpgradeShopUIState) PrevTab() {
	s.ActiveTab = UpgradeType((int(s.ActiveTab) - 1 + int(UpgradeTypeCount)) % int(UpgradeTypeCount))
	// Reset selection to first tier when changing tabs
	s.SelectedTier = 0
}

func (s *UpgradeShopUIState) NavigateUp() {
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

func (s *UpgradeShopUIState) NavigateDown() {
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

func (s *UpgradeShopUIState) NavigateLeft() {
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

func (s *UpgradeShopUIState) NavigateRight() {
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

func (s *UpgradeShopUIState) GetSelectedRow() int {
	return s.SelectedTier / ShopGridCols
}

func (s *UpgradeShopUIState) GetSelectedCol() int {
	return s.SelectedTier % ShopGridCols
}
