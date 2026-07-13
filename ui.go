package main

import (
	"NotaborEngine/internal/notaui"
)

// UI provides an immediate-mode UI system with grid-based layout,
// panels, buttons, text, input fields, sliders, checkboxes, and
// progress bars.
type UI struct {
	handle *notaui.UI
}

// NewUI creates a new UI instance bound to the given engine and window.
func NewUI(engine *Engine, window *Window) (*UI, error) {
	h, err := notaui.New(engine.handle, window.handle)
	if err != nil {
		return nil, err
	}
	return &UI{handle: h}, nil
}

// Draw renders all registered UI widgets for the current frame.
func (ui *UI) Draw() {
	ui.handle.Draw()
}

// HasKeyboardFocus returns true if a text input field currently has focus.
func (ui *UI) HasKeyboardFocus() bool {
	return ui.handle.HasKeyboardFocus()
}

// Grid creates or retrieves a grid layout container by ID.
// Grids auto-size their children into a column×row cell grid.
// Grids can be nested and respond to window resize by default.
func (ui *UI) Grid(id string) *Grid {
	return &Grid{handle: ui.handle.Grid(id)}
}

// Rect is an axis-aligned rectangle with X, Y position and W, H size.
type Rect notaui.Rect

// R creates a Rect from position and size.
func R(x, y, w, h float32) Rect {
	return Rect(notaui.R(x, y, w, h))
}

// Inset returns a Rect shrunk by v on all sides.
func (r Rect) Inset(v float32) Rect { return Rect(notaui.Rect(r).Inset(v)) }

// Contains returns true if the point (x, y) is inside the rect.
func (r Rect) Contains(x, y float32) bool { return notaui.Rect(r).Contains(x, y) }

// Anchor specifies a position within a parent rectangle for relative sizing.
type Anchor = notaui.Anchor

const (
	AnchorTopLeft      = notaui.AnchorTopLeft
	AnchorTopCenter    = notaui.AnchorTopCenter
	AnchorTopRight     = notaui.AnchorTopRight
	AnchorCenterLeft   = notaui.AnchorCenterLeft
	AnchorCenter       = notaui.AnchorCenter
	AnchorCenterRight  = notaui.AnchorCenterRight
	AnchorBottomLeft   = notaui.AnchorBottomLeft
	AnchorBottomCenter = notaui.AnchorBottomCenter
	AnchorBottomRight  = notaui.AnchorBottomRight
	AnchorStretch      = notaui.AnchorStretch
)

// ─── Grid ────────────────────────────────────────────────────────────────

// Grid is a layout container that positions child widgets in a
// column×row cell grid.  Grids can be nested and auto-size to
// their parent (or the window for the top-level grid).
type Grid struct {
	handle *notaui.Grid
}

// Columns sets the number of columns in this grid.
func (g *Grid) Columns(n int) *Grid { g.handle.Columns(n); return g }

// Gap sets the spacing between grid cells.
func (g *Grid) Gap(gap float32) *Grid { g.handle.Gap(gap); return g }

// Padding sets the inset from grid bounds to cells.
func (g *Grid) Padding(padding float32) *Grid { g.handle.Padding(padding); return g }

// RemoveById removes a widget from the grid by ID.
func (g *Grid) RemoveById(id string) { g.handle.RemoveById(id) }

// Clear removes all widgets from the grid.
func (g *Grid) Clear() { g.handle.Clear() }

// SetBounds sizes this grid as a fraction of its parent, positioned
// at the given anchor point.  Values 0<relW,relH≤1.  The grid
// re-sizes with its parent on every frame.
func (g *Grid) SetBounds(anchor Anchor, relW, relH float32) *Grid {
	g.handle.SetBounds(anchor, relW, relH)
	return g
}

// SetRelBounds sizes this grid by explicit parent-relative fractions
// (each 0≤value≤1).  (0,0) is the top-left of the parent; (1,1) is
// the bottom-right.
func (g *Grid) SetRelBounds(left, top, right, bottom float32) *Grid {
	g.handle.SetRelBounds(left, top, right, bottom)
	return g
}

// Button creates (or retrieves) a button and adds it to the grid.
func (g *Grid) Button(id, label string) *notaui.Button {
	return g.handle.Button(id, label)
}

// Text creates (or retrieves) a static text label and adds it to the grid.
func (g *Grid) Text(id, text string) *notaui.Text {
	return g.handle.Text(id, text)
}

// TextFunc creates (or retrieves) a dynamic text label and adds it to the grid.
func (g *Grid) TextFunc(id string, fn func() string) *notaui.Text {
	return g.handle.TextFunc(id, fn)
}

// Panel creates (or retrieves) a panel and adds it to the grid.
func (g *Grid) Panel(id string) *notaui.Panel {
	return g.handle.Panel(id)
}

// Input creates (or retrieves) a text input field and adds it to the grid.
func (g *Grid) Input(id string, value *string) *notaui.Input {
	return g.handle.Input(id, value)
}

// Slider creates (or retrieves) a slider and adds it to the grid.
func (g *Grid) Slider(id string, value *float32, min, max float32) *notaui.Slider {
	return g.handle.Slider(id, value, min, max)
}

// Checkbox creates (or retrieves) a checkbox and adds it to the grid.
func (g *Grid) Checkbox(id, label string, value *bool) *notaui.Checkbox {
	return g.handle.Checkbox(id, label, value)
}

// Progress creates (or retrieves) a progress bar and adds it to the grid.
func (g *Grid) Progress(id string, value *float32, min, max float32) *notaui.ProgressBar {
	return g.handle.Progress(id, value, min, max)
}

// Grid creates (or retrieves) a nested grid and adds it as a child.
func (g *Grid) Grid(id string) *Grid {
	return &Grid{handle: g.handle.Grid(id)}
}
