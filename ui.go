package main

import (
	"NotaborEngine/internal/notaui"
)

type UI struct {
	handle *notaui.UI
}

func NewUI(engine *Engine, window *Window) (*UI, error) {
	h, err := notaui.New(engine.handle, window.handle)
	if err != nil {
		return nil, err
	}
	return &UI{handle: h}, nil
}

func (ui *UI) Draw() {
	ui.handle.Draw()
}

func (ui *UI) HasKeyboardFocus() bool {
	return ui.handle.HasKeyboardFocus()
}

func (ui *UI) Panel(id string) *notaui.Panel {
	return ui.handle.Panel(id)
}

func (ui *UI) Text(id string, text string) *notaui.Text {
	return ui.handle.Text(id, text)
}

func (ui *UI) TextFunc(id string, fn func() string) *notaui.Text {
	return ui.handle.TextFunc(id, fn)
}

func (ui *UI) Button(id string, label string) *notaui.Button {
	return ui.handle.Button(id, label)
}

func (ui *UI) Input(id string, value *string) *notaui.Input {
	return ui.handle.Input(id, value)
}

func (ui *UI) Slider(id string, value *float32, min float32, max float32) *notaui.Slider {
	return ui.handle.Slider(id, value, min, max)
}

func (ui *UI) Checkbox(id string, label string, value *bool) *notaui.Checkbox {
	return ui.handle.Checkbox(id, label, value)
}

func (ui *UI) Grid(id string, bounds Rect, columns, rows int) *notaui.Grid {
	return ui.handle.Grid(id, notaui.Rect(bounds), columns, rows)
}

type Rect notaui.Rect

func R(x, y, w, h float32) Rect {
	return Rect(notaui.R(x, y, w, h))
}
