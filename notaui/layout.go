package notaui

// Grid is a tiny layout helper for placing widgets in predictable cells.
// Use UI.Grid(...).Button(...), or call Cell when you want to place manually.
type Grid struct {
	ui      *UI
	bounds  Rect
	columns int
	rows    int
	gap     float32
	padding float32
}

func (ui *UI) Grid(id string, bounds Rect, columns, rows int) *Grid {
	return &Grid{
		ui:      ui,
		bounds:  bounds,
		columns: maxInt(columns, 1),
		rows:    maxInt(rows, 1),
		gap:     6,
	}
}

func (g *Grid) Gap(gap float32) *Grid {
	if gap < 0 {
		gap = 0
	}
	g.gap = gap
	return g
}

func (g *Grid) Padding(padding float32) *Grid {
	if padding < 0 {
		padding = 0
	}
	g.padding = padding
	return g
}

func (g *Grid) Cell(column, row int) Rect {
	columns := maxInt(g.columns, 1)
	rows := maxInt(g.rows, 1)

	column = maxInt(column, 0)
	row = maxInt(row, 0)
	if column >= columns {
		column = columns - 1
	}
	if row >= rows {
		row = rows - 1
	}

	inner := g.bounds.Inset(g.padding)
	totalGapX := float32(columns-1) * g.gap
	totalGapY := float32(rows-1) * g.gap
	cellW := (inner.W - totalGapX) / float32(columns)
	cellH := (inner.H - totalGapY) / float32(rows)

	if cellW < 0 {
		cellW = 0
	}
	if cellH < 0 {
		cellH = 0
	}

	return R(
		inner.X+float32(column)*(cellW+g.gap),
		inner.Y+float32(row)*(cellH+g.gap),
		cellW,
		cellH,
	)
}

func (g *Grid) Text(id, text string, column, row int) *Text {
	cell := g.Cell(column, row)
	return g.ui.Text(id, text).At(cell.X, cell.Y+(cell.H-7)*0.5)
}

func (g *Grid) Button(id, label string, column, row int) *Button {
	return g.ui.Button(id, label).RectR(g.Cell(column, row))
}

func (g *Grid) Input(id string, value *string, column, row int) *Input {
	return g.ui.Input(id, value).RectR(g.Cell(column, row))
}

func (g *Grid) Slider(id string, value *float32, min, max float32, column, row int) *Slider {
	return g.ui.Slider(id, value, min, max).RectR(g.Cell(column, row))
}

func (g *Grid) Checkbox(id, label string, value *bool, column, row int) *Checkbox {
	return g.ui.Checkbox(id, label, value).RectR(g.Cell(column, row))
}

func (g *Grid) Progress(id string, value *float32, min, max float32, column, row int) *ProgressBar {
	return g.ui.Progress(id, value, min, max).RectR(g.Cell(column, row))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
