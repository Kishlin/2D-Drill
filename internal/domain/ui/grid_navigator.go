package ui

// GridNavigator handles grid-based selection with wrapping navigation.
type GridNavigator struct {
	Selected int
	cols     int
	rows     int
}

func NewGridNavigator(cols, rows int) GridNavigator {
	return GridNavigator{Selected: 0, cols: cols, rows: rows}
}

func (g *GridNavigator) NavigateUp() {
	row := g.Selected / g.cols
	col := g.Selected % g.cols

	if row == 0 {
		row = g.rows - 1
	} else {
		row--
	}

	g.Selected = row*g.cols + col
}

func (g *GridNavigator) NavigateDown() {
	row := g.Selected / g.cols
	col := g.Selected % g.cols

	if row == g.rows-1 {
		row = 0
	} else {
		row++
	}

	g.Selected = row*g.cols + col
}

func (g *GridNavigator) NavigateLeft() {
	row := g.Selected / g.cols
	col := g.Selected % g.cols

	if col == 0 {
		col = g.cols - 1
	} else {
		col--
	}

	g.Selected = row*g.cols + col
}

func (g *GridNavigator) NavigateRight() {
	row := g.Selected / g.cols
	col := g.Selected % g.cols

	if col == g.cols-1 {
		col = 0
	} else {
		col++
	}

	g.Selected = row*g.cols + col
}

func (g *GridNavigator) GetSelectedRow() int {
	return g.Selected / g.cols
}

func (g *GridNavigator) GetSelectedCol() int {
	return g.Selected % g.cols
}
