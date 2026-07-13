package notaui

import "testing"

func TestRectContains(t *testing.T) {
	r := R(10, 20, 100, 50)

	if !r.Contains(10, 20) {
		t.Fatal("expected top-left edge to be inside")
	}
	if !r.Contains(110, 70) {
		t.Fatal("expected bottom-right edge to be inside")
	}
	if r.Contains(111, 70) {
		t.Fatal("expected point past right edge to be outside")
	}
	if r.Contains(10, 71) {
		t.Fatal("expected point past bottom edge to be outside")
	}
}

func TestSliderSetFromMouseClamps(t *testing.T) {
	value := float32(0)
	slider := &Slider{
		value: &value,
		min:   10,
		max:   20,
		rect:  R(100, 0, 200, 20),
	}

	slider.setFromMouse(50)
	if value != 10 {
		t.Fatalf("expected min value, got %f", value)
	}

	slider.setFromMouse(400)
	if value != 20 {
		t.Fatalf("expected max value, got %f", value)
	}

	slider.setFromMouse(200)
	if value != 15 {
		t.Fatalf("expected midpoint value, got %f", value)
	}
}

func TestSliderOnChangeCallback(t *testing.T) {
	value := float32(0)
	got := float32(-1)
	slider := (&Slider{
		value: &value,
		min:   0,
		max:   100,
		rect:  R(0, 0, 100, 20),
	}).OnChange(func(v float32) {
		got = v
	})

	callback := slider.setFromMouse(25)
	if callback == nil {
		t.Fatal("expected callback")
	}
	callback()

	if got != 25 {
		t.Fatalf("expected callback value 25, got %f", got)
	}
}

func TestGridAutoFlow(t *testing.T) {
	g := &Grid{
		name:    "test",
		columns: 3,
		gap:     6,
		padding: 5,
	}

	if g.columns != 3 {
		t.Fatalf("expected 3 columns, got %d", g.columns)
	}
	if g.gap != 6 {
		t.Fatalf("expected gap 6, got %f", g.gap)
	}
	if g.id() != "test" {
		t.Fatalf("expected id 'test', got '%s'", g.id())
	}

	// Set bounds and verify they're stored
	g.setBounds(R(0, 0, 400, 300))
	b := g.bounds()
	if b.W != 400 || b.H != 300 {
		t.Fatalf("expected bounds (0,0,400,300), got %+v", b)
	}

	// Verify default row count for empty grid
	if g.rows() != 0 {
		t.Fatalf("expected 0 rows for empty grid, got %d", g.rows())
	}
}

func TestGridAddWidget(t *testing.T) {
	g := &Grid{
		name:    "test",
		columns: 2,
		gap:     4,
		padding: 2,
	}

	w := &Button{name: "btn", rect: R(0, 0, 0, 0), meta: defaultGridMeta()}
	g.Add(w)

	if len(g.cells) != 1 {
		t.Fatalf("expected 1 cell, got %d", len(g.cells))
	}

	cell := g.cells[0]
	if cell.col != 0 || cell.row != 0 {
		t.Fatalf("expected cell at (0,0), got (%d,%d)", cell.col, cell.row)
	}

	// Add second widget — should auto-flow to (1, 0)
	w2 := &Button{name: "btn2", rect: R(0, 0, 0, 0), meta: defaultGridMeta()}
	g.Add(w2)

	if len(g.cells) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(g.cells))
	}
	cell2 := g.cells[1]
	if cell2.col != 1 || cell2.row != 0 {
		t.Fatalf("expected cell at (1,0), got (%d,%d)", cell2.col, cell2.row)
	}
}

func TestGridAtPosition(t *testing.T) {
	g := &Grid{
		name:    "test",
		columns: 4,
		gap:     4,
		padding: 2,
	}

	w := &Button{name: "btn", rect: R(0, 0, 0, 0), meta: defaultGridMeta()}
	w.setGridPos(2, 3)
	g.Add(w)

	if len(g.cells) != 1 {
		t.Fatalf("expected 1 cell, got %d", len(g.cells))
	}

	cell := g.cells[0]
	if cell.col != 2 || cell.row != 3 {
		t.Fatalf("expected cell at (2,3), got (%d,%d)", cell.col, cell.row)
	}
}
