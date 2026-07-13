package main

import (
	"NotaborEngine/internal/notaui"
)

// UI provides an immediate-mode UI system with panels, buttons, text,
// input fields, sliders, checkboxes, and grid layouts.
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

// Panel creates or retrieves a panel widget by ID.
func (ui *UI) Panel(id string) *notaui.Panel {
	return ui.handle.Panel(id)
}

// Text creates or retrieves a static text label by ID.
func (ui *UI) Text(id string, text string) *notaui.Text {
	return ui.handle.Text(id, text)
}

// TextFunc creates or retrieves a dynamic text label (updated each frame via fn) by ID.
func (ui *UI) TextFunc(id string, fn func() string) *notaui.Text {
	return ui.handle.TextFunc(id, fn)
}

// Button creates or retrieves a clickable button by ID.
func (ui *UI) Button(id string, label string) *notaui.Button {
	return ui.handle.Button(id, label)
}

// Input creates or retrieves a text input field by ID.
func (ui *UI) Input(id string, value *string) *notaui.Input {
	return ui.handle.Input(id, value)
}

// Slider creates or retrieves a draggable slider by ID.
func (ui *UI) Slider(id string, value *float32, min float32, max float32) *notaui.Slider {
	return ui.handle.Slider(id, value, min, max)
}

// Checkbox creates or retrieves a toggleable checkbox by ID.
func (ui *UI) Checkbox(id string, label string, value *bool) *notaui.Checkbox {
	return ui.handle.Checkbox(id, label, value)
}

// Grid creates or retrieves a grid layout helper by ID.
func (ui *UI) Grid(id string, bounds Rect, columns, rows int) *notaui.Grid {
	return ui.handle.Grid(id, notaui.Rect(bounds), columns, rows)
}

// Rect is an axis-aligned rectangle with X, Y position and W, H size.
type Rect notaui.Rect

// R creates a Rect from position and size.
func R(x, y, w, h float32) Rect {
	return Rect(notaui.R(x, y, w, h))
}
