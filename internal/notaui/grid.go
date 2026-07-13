package notaui

// Anchor specifies a position within a parent rectangle for relative sizing.
type Anchor int

const (
	AnchorTopLeft Anchor = iota
	AnchorTopCenter
	AnchorTopRight
	AnchorCenterLeft
	AnchorCenter
	AnchorCenterRight
	AnchorBottomLeft
	AnchorBottomCenter
	AnchorBottomRight
	AnchorStretch
)

// Grid is a layout container that positions child Widgets in a
// column×row cell grid.  Grid implements the Widget interface so
// grids can be nested.
type Grid struct {
	name    string
	ui      *UI
	rect    Rect
	columns int
	gap     float32
	padding float32
	cells   []gridCell
	meta    gridMeta

	// Relative sizing — when relBounds is true the grid's pixel
	// bounds are computed from its parent bounds each frame.
	relBounds bool
	relLeft   float32
	relTop    float32
	relRight  float32
	relBottom float32
}

type gridCell struct {
	w       Widget
	col     int
	row     int
	colSpan int
	rowSpan int
}

func (ui *UI) Grid(id string) *Grid {
	if w, ok := ui.byID[id].(*Grid); ok {
		w.columns = 1
		w.cells = nil
		return w
	}
	g := &Grid{
		name:    id,
		ui:      ui,
		columns: 1,
		gap:     6,
		meta:    defaultGridMeta(),
	}
	ui.add(g)
	return g
}

func (g *Grid) id() string   { return g.name }
func (g *Grid) bounds() Rect { return g.rect }

// setBounds implements the Widget interface.  When the grid has
// relative sizing enabled the pixel bounds are computed from the
// parent rect + stored fractions; otherwise parent rect is used
// directly.
func (g *Grid) setBounds(r Rect) {
	if g.relBounds {
		pw, ph := r.W, r.H
		g.rect = R(
			r.X+g.relLeft*pw,
			r.Y+g.relTop*ph,
			(g.relRight-g.relLeft)*pw,
			(g.relBottom-g.relTop)*ph,
		)
	} else {
		g.rect = r
	}
}

// SetBounds sizes this grid as a fraction of its parent, positioned
// at the given anchor point.  Values 0<relW,relH≤1.  Calling this
// enables relative sizing — the grid re-sizes with its parent on
// every frame.
func (g *Grid) SetBounds(anchor Anchor, relW, relH float32) *Grid {
	g.relBounds = true
	switch anchor {
	case AnchorTopLeft:
		g.relLeft, g.relTop, g.relRight, g.relBottom = 0, 0, relW, relH
	case AnchorTopCenter:
		g.relLeft, g.relTop, g.relRight, g.relBottom = (1-relW)/2, 0, (1+relW)/2, relH
	case AnchorTopRight:
		g.relLeft, g.relTop, g.relRight, g.relBottom = 1-relW, 0, 1, relH
	case AnchorCenterLeft:
		g.relLeft, g.relTop, g.relRight, g.relBottom = 0, (1-relH)/2, relW, (1+relH)/2
	case AnchorCenter:
		g.relLeft, g.relTop, g.relRight, g.relBottom = (1-relW)/2, (1-relH)/2, (1+relW)/2, (1+relH)/2
	case AnchorCenterRight:
		g.relLeft, g.relTop, g.relRight, g.relBottom = 1-relW, (1-relH)/2, 1, (1+relH)/2
	case AnchorBottomLeft:
		g.relLeft, g.relTop, g.relRight, g.relBottom = 0, 1-relH, relW, 1
	case AnchorBottomCenter:
		g.relLeft, g.relTop, g.relRight, g.relBottom = (1-relW)/2, 1-relH, (1+relW)/2, 1
	case AnchorBottomRight:
		g.relLeft, g.relTop, g.relRight, g.relBottom = 1-relW, 1-relH, 1, 1
	default:
		g.relBounds = false
	}
	return g
}

// SetRelBounds sizes this grid by explicit parent-relative fractions
// (each 0≤value≤1).  (0,0) is the top-left of the parent; (1,1) is
// the bottom-right.
func (g *Grid) SetRelBounds(left, top, right, bottom float32) *Grid {
	g.relBounds = true
	g.relLeft, g.relTop, g.relRight, g.relBottom = left, top, right, bottom
	return g
}

// Widget interface methods – Grid is also positionable in a parent grid.
func (g *Grid) gridCol() int        { return g.meta.col }
func (g *Grid) gridRow() int        { return g.meta.row }
func (g *Grid) colSpan() int        { return g.meta.colSpan }
func (g *Grid) rowSpan() int        { return g.meta.rowSpan }
func (g *Grid) setGridPos(c, r int) { g.meta.setGridPos(c, r) }
func (g *Grid) setColSpan(n int)    { g.meta.setColSpan(n) }
func (g *Grid) setRowSpan(n int)    { g.meta.setRowSpan(n) }

// Columns sets the number of columns for this grid. Triggers re-layout.
func (g *Grid) Columns(n int) *Grid {
	if n < 1 {
		n = 1
	}
	g.columns = n
	return g
}

// Gap sets the spacing between grid cells.
func (g *Grid) Gap(gap float32) *Grid {
	if gap < 0 {
		gap = 0
	}
	g.gap = gap
	return g
}

// Padding sets the inset from the grid bounds to the cells.
func (g *Grid) Padding(padding float32) *Grid {
	if padding < 0 {
		padding = 0
	}
	g.padding = padding
	return g
}

// Add inserts a Widget into the grid using auto-flow (row-major).
// Returns the Widget for chaining.  Use Widget.At/ColSpan/RowSpan
// to override auto-flow position before passing to Add, or set
// them after adding.
func (g *Grid) Add(w Widget) Widget {
	g.cells = append(g.cells, g.flowCell(w))
	if g.ui != nil {
		g.ui.managed[w.id()] = true
	}
	return w
}

// flowCell determines the auto-flow position for w and returns a gridCell.
func (g *Grid) flowCell(w Widget) gridCell {
	col, row := w.gridCol(), w.gridRow()
	if col >= 0 && row >= 0 {
		// Explicit position – do not affect auto-flow cursor.
		return gridCell{w: w, col: col, row: row, colSpan: w.colSpan(), rowSpan: w.rowSpan()}
	}

	// Auto-flow: find next free cell.
	cols := g.columns
	occupied := g.occupiedMap()

	nextCol := 0
	nextRow := 0
	for occupied[nextRow*gridStride+nextCol] {
		nextCol++
		if nextCol >= cols {
			nextCol = 0
			nextRow++
		}
		// Safety – should never loop forever with valid input.
		if nextRow > 10000 {
			break
		}
	}

	cs := w.colSpan()
	if cs < 1 {
		cs = 1
	}
	rs := w.rowSpan()
	if rs < 1 {
		rs = 1
	}
	return gridCell{w: w, col: nextCol, row: nextRow, colSpan: cs, rowSpan: rs}
}

const gridStride = 10000

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// occupiedMap builds a set of occupied cell positions, accounting for spans.
func (g *Grid) occupiedMap() map[int]bool {
	m := make(map[int]bool)
	for _, c := range g.cells {
		cs := c.w.colSpan()
		rs := c.w.rowSpan()
		for r := 0; r < maxInt(rs, 1); r++ {
			for co := 0; co < maxInt(cs, 1); co++ {
				m[(c.row+r)*gridStride+(c.col+co)] = true
			}
		}
	}
	return m
}

// RemoveById removes a Widget by its ID from the grid.
func (g *Grid) RemoveById(id string) {
	filtered := make([]gridCell, 0, len(g.cells))
	for _, c := range g.cells {
		if c.w.id() == id {
			if g.ui != nil {
				delete(g.ui.managed, id)
			}
			continue
		}
		filtered = append(filtered, c)
	}
	g.cells = filtered
}

// Clear removes all Widgets from the grid.
func (g *Grid) Clear() {
	if g.ui != nil {
		for _, c := range g.cells {
			delete(g.ui.managed, c.w.id())
		}
	}
	g.cells = nil
}

// rows returns the number of rows needed to fit all cells.
func (g *Grid) rows() int {
	maxRow := 0
	for _, c := range g.cells {
		rs := maxInt(c.w.rowSpan(), 1)
		r := c.row + rs - 1
		if r > maxRow {
			maxRow = r
		}
	}
	return maxRow + 1
}

// draw implements the Widget interface: calculates cell bounds from
// the grid's own bounds and delegates drawing to each child.
func (g *Grid) draw(ui *UI) {
	if g.rect.IsEmpty() {
		return
	}

	rows := g.rows()
	if rows == 0 {
		return
	}
	cols := g.columns

	inner := g.rect.Inset(g.padding)
	totalGapW := float32(cols-1) * g.gap
	totalGapH := float32(rows-1) * g.gap
	cellW := (inner.W - totalGapW) / float32(cols)
	cellH := (inner.H - totalGapH) / float32(rows)

	for _, c := range g.cells {
		cs := maxInt(c.w.colSpan(), 1)
		rs := maxInt(c.w.rowSpan(), 1)
		cx := inner.X + float32(c.col)*(cellW+g.gap)
		cy := inner.Y + float32(c.row)*(cellH+g.gap)
		cw := cellW*float32(cs) + g.gap*float32(cs-1)
		ch := cellH*float32(rs) + g.gap*float32(rs-1)
		c.w.setBounds(R(cx, cy, cw, ch))
		c.w.draw(ui)
	}
}

// Convenience shorthands – create Widget and Add in one call.

func (g *Grid) Button(id, label string) *Button {
	w := g.ui.Button(id, label)
	g.Add(w)
	return w
}

func (g *Grid) Text(id, text string) *Text {
	w := g.ui.Text(id, text)
	g.Add(w)
	return w
}

func (g *Grid) TextFunc(id string, fn func() string) *Text {
	w := g.ui.TextFunc(id, fn)
	g.Add(w)
	return w
}

func (g *Grid) Panel(id string) *Panel {
	w := g.ui.Panel(id)
	g.Add(w)
	return w
}

func (g *Grid) Input(id string, value *string) *Input {
	w := g.ui.Input(id, value)
	g.Add(w)
	return w
}

func (g *Grid) Slider(id string, value *float32, min, max float32) *Slider {
	w := g.ui.Slider(id, value, min, max)
	g.Add(w)
	return w
}

func (g *Grid) Checkbox(id, label string, value *bool) *Checkbox {
	w := g.ui.Checkbox(id, label, value)
	g.Add(w)
	return w
}

func (g *Grid) Progress(id string, value *float32, min, max float32) *ProgressBar {
	w := g.ui.Progress(id, value, min, max)
	g.Add(w)
	return w
}

// Grid creates a nested grid and adds it as a child. Returns the child grid.
func (g *Grid) Grid(id string) *Grid {
	child := g.ui.Grid(id)
	g.Add(child)
	return child
}
