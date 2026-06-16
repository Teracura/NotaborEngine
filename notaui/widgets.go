package notaui

import (
	"NotaborEngine/internal/notacolor"
	"fmt"
)

type Text struct {
	name   string
	x, y   float32
	scale  float32
	color  notacolor.Color
	text   string
	textFn func() string
}

func (ui *UI) Text(id, text string) *Text {
	if w, ok := ui.byID[id].(*Text); ok {
		w.text = text
		w.textFn = nil
		return w
	}

	w := &Text{name: id, text: text, scale: 1, color: ui.Theme.Text}
	ui.add(w)
	return w
}

func (ui *UI) TextFunc(id string, text func() string) *Text {
	if w, ok := ui.byID[id].(*Text); ok {
		w.textFn = text
		return w
	}

	w := &Text{name: id, textFn: text, scale: 1, color: ui.Theme.Text}
	ui.add(w)
	return w
}

func (w *Text) id() string { return w.name }

func (w *Text) At(x, y float32) *Text {
	w.x, w.y = x, y
	return w
}

func (w *Text) Scale(scale float32) *Text {
	w.scale = scale
	return w
}

func (w *Text) Color(color notacolor.Color) *Text {
	w.color = color
	return w
}

func (w *Text) draw(ui *UI) {
	text := w.text
	if w.textFn != nil {
		text = w.textFn()
	}
	ui.drawText(text, w.x, w.y, w.scale, w.color)
}

type Panel struct {
	name   string
	bounds Rect
	fill   notacolor.Color
	border notacolor.Color
}

func (ui *UI) Panel(id string) *Panel {
	if w, ok := ui.byID[id].(*Panel); ok {
		return w
	}

	w := &Panel{name: id, fill: ui.Theme.Panel, border: ui.Theme.PanelBorder}
	ui.add(w)
	return w
}

func (w *Panel) id() string { return w.name }

func (w *Panel) Rect(x, y, width, height float32) *Panel {
	w.bounds = R(x, y, width, height)
	return w
}

func (w *Panel) RectR(r Rect) *Panel {
	w.bounds = r
	return w
}

func (w *Panel) Colors(fill, border notacolor.Color) *Panel {
	w.fill = fill
	w.border = border
	return w
}

func (w *Panel) draw(ui *UI) {
	ui.drawRect(w.bounds, w.fill)
	ui.drawBorder(w.bounds, 1, w.border)
}

type Button struct {
	name    string
	label   string
	bounds  Rect
	onClick func()
}

func (ui *UI) Button(id, label string) *Button {
	if w, ok := ui.byID[id].(*Button); ok {
		w.label = label
		return w
	}

	w := &Button{name: id, label: label}
	ui.add(w)
	return w
}

func (w *Button) id() string { return w.name }

func (w *Button) Rect(x, y, width, height float32) *Button {
	w.bounds = R(x, y, width, height)
	return w
}

func (w *Button) RectR(r Rect) *Button {
	w.bounds = r
	return w
}

func (w *Button) OnClick(fn func()) *Button {
	w.onClick = fn
	return w
}

func (w *Button) draw(ui *UI) {
	bg := ui.Theme.Surface
	if w.bounds.Contains(ui.mouseX, ui.mouseY) {
		bg = ui.Theme.SurfaceHover
	}
	if ui.activeID == w.name && ui.mouseDown {
		bg = ui.Theme.SurfaceDown
	}

	ui.drawRect(w.bounds, bg)
	ui.drawBorder(w.bounds, 1, ui.Theme.PanelBorder)

	scale := float32(1)
	textW := ui.textWidth(w.label, scale)
	tx := w.bounds.X + (w.bounds.W-textW)*0.5
	ty := w.bounds.Y + (w.bounds.H-7*scale)*0.5
	ui.drawText(w.label, tx, ty, scale, ui.Theme.Text)
}

type Input struct {
	name        string
	value       *string
	placeholder string
	bounds      Rect
	onChange    func(string)
}

func (ui *UI) Input(id string, value *string) *Input {
	if w, ok := ui.byID[id].(*Input); ok {
		w.value = value
		return w
	}

	w := &Input{name: id, value: value}
	ui.add(w)
	return w
}

func (w *Input) id() string { return w.name }

func (w *Input) Rect(x, y, width, height float32) *Input {
	w.bounds = R(x, y, width, height)
	return w
}

func (w *Input) RectR(r Rect) *Input {
	w.bounds = r
	return w
}

func (w *Input) Placeholder(text string) *Input {
	w.placeholder = text
	return w
}

func (w *Input) OnChange(fn func(string)) *Input {
	w.onChange = fn
	return w
}

func (w *Input) draw(ui *UI) {
	focused := ui.focusID == w.name
	bg := ui.Theme.Input
	if focused {
		bg = ui.Theme.InputFocus
	}

	ui.drawRect(w.bounds, bg)
	border := ui.Theme.PanelBorder
	if focused {
		border = ui.Theme.Accent
	}
	ui.drawBorder(w.bounds, 1, border)

	text := ""
	color := ui.Theme.Text
	if w.value != nil {
		text = *w.value
	}
	if text == "" && w.placeholder != "" {
		text = w.placeholder
		color = ui.Theme.TextMuted
	}
	if focused {
		text += "_"
	}
	ui.drawText(text, w.bounds.X+8, w.bounds.Y+(w.bounds.H-7)*0.5, 1, color)
}

type Slider struct {
	name     string
	label    string
	value    *float32
	min, max float32
	bounds   Rect
	onChange func(float32)
}

func (ui *UI) Slider(id string, value *float32, min, max float32) *Slider {
	if w, ok := ui.byID[id].(*Slider); ok {
		w.value = value
		w.min = min
		w.max = max
		return w
	}

	w := &Slider{name: id, value: value, min: min, max: max}
	ui.add(w)
	return w
}

func (w *Slider) id() string { return w.name }

func (w *Slider) Rect(x, y, width, height float32) *Slider {
	w.bounds = R(x, y, width, height)
	return w
}

func (w *Slider) RectR(r Rect) *Slider {
	w.bounds = r
	return w
}

func (w *Slider) Label(text string) *Slider {
	w.label = text
	return w
}

func (w *Slider) OnChange(fn func(float32)) *Slider {
	w.onChange = fn
	return w
}

func (w *Slider) setFromMouse(x float32) func() {
	if w.value == nil || w.bounds.W <= 0 {
		return nil
	}
	t := clamp((x-w.bounds.X)/w.bounds.W, 0, 1)
	*w.value = lerp(w.min, w.max, t)
	if w.onChange != nil {
		value := *w.value
		return func() { w.onChange(value) }
	}
	return nil
}

func (w *Slider) draw(ui *UI) {
	value := w.min
	if w.value != nil {
		value = *w.value
	}
	t := float32(0)
	if w.max != w.min {
		t = clamp((value-w.min)/(w.max-w.min), 0, 1)
	}

	label := w.label
	if label != "" {
		label = fmt.Sprintf("%s %.2f", label, value)
		ui.drawText(label, w.bounds.X, w.bounds.Y, 1, ui.Theme.TextMuted)
	}

	track := R(w.bounds.X, w.bounds.Y+w.bounds.H*0.5-3, w.bounds.W, 6)
	fill := R(track.X, track.Y, track.W*t, track.H)
	knob := R(w.bounds.X+w.bounds.W*t-5, w.bounds.Y+w.bounds.H*0.5-8, 10, 16)

	ui.drawRect(track, ui.Theme.Track)
	ui.drawRect(fill, ui.Theme.Accent)
	ui.drawRect(knob, ui.Theme.AccentHover)
}

type Checkbox struct {
	name     string
	label    string
	value    *bool
	bounds   Rect
	onChange func(bool)
}

func (ui *UI) Checkbox(id, label string, value *bool) *Checkbox {
	if w, ok := ui.byID[id].(*Checkbox); ok {
		w.label = label
		w.value = value
		return w
	}

	w := &Checkbox{name: id, label: label, value: value}
	ui.add(w)
	return w
}

func (w *Checkbox) id() string { return w.name }

func (w *Checkbox) Rect(x, y, width, height float32) *Checkbox {
	w.bounds = R(x, y, width, height)
	return w
}

func (w *Checkbox) RectR(r Rect) *Checkbox {
	w.bounds = r
	return w
}

func (w *Checkbox) OnChange(fn func(bool)) *Checkbox {
	w.onChange = fn
	return w
}

func (w *Checkbox) draw(ui *UI) {
	box := R(w.bounds.X, w.bounds.Y+(w.bounds.H-14)*0.5, 14, 14)
	ui.drawRect(box, ui.Theme.Input)
	ui.drawBorder(box, 1, ui.Theme.PanelBorder)

	if w.value != nil && *w.value {
		ui.drawRect(box.Inset(3), ui.Theme.Accent)
	}

	ui.drawText(w.label, w.bounds.X+22, w.bounds.Y+(w.bounds.H-7)*0.5, 1, ui.Theme.Text)
}

type ProgressBar struct {
	name     string
	value    *float32
	min, max float32
	bounds   Rect
}

func (ui *UI) Progress(id string, value *float32, min, max float32) *ProgressBar {
	if w, ok := ui.byID[id].(*ProgressBar); ok {
		w.value = value
		w.min = min
		w.max = max
		return w
	}

	w := &ProgressBar{name: id, value: value, min: min, max: max}
	ui.add(w)
	return w
}

func (w *ProgressBar) id() string { return w.name }

func (w *ProgressBar) Rect(x, y, width, height float32) *ProgressBar {
	w.bounds = R(x, y, width, height)
	return w
}

func (w *ProgressBar) RectR(r Rect) *ProgressBar {
	w.bounds = r
	return w
}

func (w *ProgressBar) draw(ui *UI) {
	value := w.min
	if w.value != nil {
		value = *w.value
	}
	t := float32(0)
	if w.max != w.min {
		t = clamp((value-w.min)/(w.max-w.min), 0, 1)
	}
	ui.drawRect(w.bounds, ui.Theme.Track)
	ui.drawRect(R(w.bounds.X, w.bounds.Y, w.bounds.W*t, w.bounds.H), ui.Theme.Accent)
	ui.drawBorder(w.bounds, 1, ui.Theme.PanelBorder)
}
