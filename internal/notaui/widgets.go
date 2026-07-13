package notaui

import (
	"NotaborEngine/internal/notacolor"
	"fmt"
)

// ─── Text ────────────────────────────────────────────────────────────────

type Text struct {
	name   string
	rect   Rect
	scale  float32
	color  notacolor.Color
	text   string
	textFn func() string
	meta   gridMeta
}

func (ui *UI) Text(id, text string) *Text {
	if w, ok := ui.byID[id].(*Text); ok {
		w.text = text
		w.textFn = nil
		return w
	}
	w := &Text{name: id, text: text, scale: 1, color: ui.Theme.Text, meta: defaultGridMeta()}
	ui.add(w)
	return w
}

func (ui *UI) TextFunc(id string, text func() string) *Text {
	if w, ok := ui.byID[id].(*Text); ok {
		w.textFn = text
		return w
	}
	w := &Text{name: id, textFn: text, scale: 1, color: ui.Theme.Text, meta: defaultGridMeta()}
	ui.add(w)
	return w
}

func (w *Text) id() string          { return w.name }
func (w *Text) setBounds(r Rect)    { w.rect = r }
func (w *Text) bounds() Rect        { return w.rect }
func (w *Text) gridCol() int        { return w.meta.col }
func (w *Text) gridRow() int        { return w.meta.row }
func (w *Text) colSpan() int        { return w.meta.colSpan }
func (w *Text) rowSpan() int        { return w.meta.rowSpan }
func (w *Text) setGridPos(c, r int) { w.meta.setGridPos(c, r) }
func (w *Text) setColSpan(n int)    { w.meta.setColSpan(n) }
func (w *Text) setRowSpan(n int)    { w.meta.setRowSpan(n) }

func (w *Text) At(x, y float32) *Text {
	w.rect.X, w.rect.Y = x, y
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

func (w *Text) ColSpan(n int) *Text { w.meta.colSpan = n; return w }
func (w *Text) RowSpan(n int) *Text { w.meta.rowSpan = n; return w }

func (w *Text) draw(ui *UI) {
	text := w.text
	if w.textFn != nil {
		text = w.textFn()
	}
	ui.drawText(text, w.rect.X, w.rect.Y, w.scale, w.color)
}

// ─── Panel ───────────────────────────────────────────────────────────────

type Panel struct {
	name   string
	rect   Rect
	fill   notacolor.Color
	border notacolor.Color
	meta   gridMeta
}

func (ui *UI) Panel(id string) *Panel {
	if w, ok := ui.byID[id].(*Panel); ok {
		return w
	}
	w := &Panel{name: id, fill: ui.Theme.Panel, border: ui.Theme.PanelBorder, meta: defaultGridMeta()}
	ui.add(w)
	return w
}

func (w *Panel) id() string          { return w.name }
func (w *Panel) setBounds(r Rect)    { w.rect = r }
func (w *Panel) bounds() Rect        { return w.rect }
func (w *Panel) gridCol() int        { return w.meta.col }
func (w *Panel) gridRow() int        { return w.meta.row }
func (w *Panel) colSpan() int        { return w.meta.colSpan }
func (w *Panel) rowSpan() int        { return w.meta.rowSpan }
func (w *Panel) setGridPos(c, r int) { w.meta.setGridPos(c, r) }
func (w *Panel) setColSpan(n int)    { w.meta.setColSpan(n) }
func (w *Panel) setRowSpan(n int)    { w.meta.setRowSpan(n) }

func (w *Panel) Rect(x, y, width, height float32) *Panel {
	w.rect = R(x, y, width, height)
	return w
}
func (w *Panel) RectR(r Rect) *Panel { w.rect = r; return w }
func (w *Panel) Colors(fill, border notacolor.Color) *Panel {
	w.fill = fill
	w.border = border
	return w
}

func (w *Panel) ColSpan(n int) *Panel { w.meta.colSpan = n; return w }
func (w *Panel) RowSpan(n int) *Panel { w.meta.rowSpan = n; return w }

func (w *Panel) draw(ui *UI) {
	ui.drawRect(w.rect, w.fill)
	ui.drawBorder(w.rect, 1, w.border)
}

// ─── Button ──────────────────────────────────────────────────────────────

type Button struct {
	name    string
	label   string
	rect    Rect
	onClick func()
	meta    gridMeta
}

func (ui *UI) Button(id, label string) *Button {
	if w, ok := ui.byID[id].(*Button); ok {
		w.label = label
		return w
	}
	w := &Button{name: id, label: label, meta: defaultGridMeta()}
	ui.add(w)
	return w
}

func (w *Button) id() string          { return w.name }
func (w *Button) setBounds(r Rect)    { w.rect = r }
func (w *Button) bounds() Rect        { return w.rect }
func (w *Button) gridCol() int        { return w.meta.col }
func (w *Button) gridRow() int        { return w.meta.row }
func (w *Button) colSpan() int        { return w.meta.colSpan }
func (w *Button) rowSpan() int        { return w.meta.rowSpan }
func (w *Button) setGridPos(c, r int) { w.meta.setGridPos(c, r) }
func (w *Button) setColSpan(n int)    { w.meta.setColSpan(n) }
func (w *Button) setRowSpan(n int)    { w.meta.setRowSpan(n) }

func (w *Button) Rect(x, y, width, height float32) *Button {
	w.rect = R(x, y, width, height)
	return w
}
func (w *Button) RectR(r Rect) *Button { w.rect = r; return w }
func (w *Button) OnClick(fn func()) *Button {
	w.onClick = fn
	return w
}

func (w *Button) ColSpan(n int) *Button { w.meta.colSpan = n; return w }
func (w *Button) RowSpan(n int) *Button { w.meta.rowSpan = n; return w }

func (w *Button) draw(ui *UI) {
	bg := ui.Theme.Surface
	if w.rect.Contains(ui.mouseX, ui.mouseY) {
		bg = ui.Theme.SurfaceHover
	}
	if ui.activeID == w.name && ui.mouseDown {
		bg = ui.Theme.SurfaceDown
	}

	ui.drawRect(w.rect, bg)
	ui.drawBorder(w.rect, 1, ui.Theme.PanelBorder)

	scale := float32(1)
	textW := ui.textWidth(w.label, scale)
	tx := w.rect.X + (w.rect.W-textW)*0.5
	ty := w.rect.Y + (w.rect.H-7*scale)*0.5
	ui.drawText(w.label, tx, ty, scale, ui.Theme.Text)
}

// ─── Input ───────────────────────────────────────────────────────────────

type Input struct {
	name        string
	value       *string
	placeholder string
	rect        Rect
	onChange    func(string)
	meta        gridMeta
}

func (ui *UI) Input(id string, value *string) *Input {
	if w, ok := ui.byID[id].(*Input); ok {
		w.value = value
		return w
	}
	w := &Input{name: id, value: value, meta: defaultGridMeta()}
	ui.add(w)
	return w
}

func (w *Input) id() string          { return w.name }
func (w *Input) setBounds(r Rect)    { w.rect = r }
func (w *Input) bounds() Rect        { return w.rect }
func (w *Input) gridCol() int        { return w.meta.col }
func (w *Input) gridRow() int        { return w.meta.row }
func (w *Input) colSpan() int        { return w.meta.colSpan }
func (w *Input) rowSpan() int        { return w.meta.rowSpan }
func (w *Input) setGridPos(c, r int) { w.meta.setGridPos(c, r) }
func (w *Input) setColSpan(n int)    { w.meta.setColSpan(n) }
func (w *Input) setRowSpan(n int)    { w.meta.setRowSpan(n) }

func (w *Input) Rect(x, y, width, height float32) *Input {
	w.rect = R(x, y, width, height)
	return w
}
func (w *Input) RectR(r Rect) *Input { w.rect = r; return w }
func (w *Input) Placeholder(text string) *Input {
	w.placeholder = text
	return w
}
func (w *Input) OnChange(fn func(string)) *Input {
	w.onChange = fn
	return w
}

func (w *Input) ColSpan(n int) *Input { w.meta.colSpan = n; return w }
func (w *Input) RowSpan(n int) *Input { w.meta.rowSpan = n; return w }

func (w *Input) draw(ui *UI) {
	focused := ui.focusID == w.name
	bg := ui.Theme.Input
	if focused {
		bg = ui.Theme.InputFocus
	}

	ui.drawRect(w.rect, bg)
	border := ui.Theme.PanelBorder
	if focused {
		border = ui.Theme.Accent
	}
	ui.drawBorder(w.rect, 1, border)

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
	ui.drawText(text, w.rect.X+8, w.rect.Y+(w.rect.H-7)*0.5, 1, color)
}

// ─── Slider ──────────────────────────────────────────────────────────────

type Slider struct {
	name     string
	label    string
	value    *float32
	min, max float32
	rect     Rect
	onChange func(float32)
	meta     gridMeta
}

func (ui *UI) Slider(id string, value *float32, min, max float32) *Slider {
	if w, ok := ui.byID[id].(*Slider); ok {
		w.value = value
		w.min = min
		w.max = max
		return w
	}
	w := &Slider{name: id, value: value, min: min, max: max, meta: defaultGridMeta()}
	ui.add(w)
	return w
}

func (w *Slider) id() string          { return w.name }
func (w *Slider) setBounds(r Rect)    { w.rect = r }
func (w *Slider) bounds() Rect        { return w.rect }
func (w *Slider) gridCol() int        { return w.meta.col }
func (w *Slider) gridRow() int        { return w.meta.row }
func (w *Slider) colSpan() int        { return w.meta.colSpan }
func (w *Slider) rowSpan() int        { return w.meta.rowSpan }
func (w *Slider) setGridPos(c, r int) { w.meta.setGridPos(c, r) }
func (w *Slider) setColSpan(n int)    { w.meta.setColSpan(n) }
func (w *Slider) setRowSpan(n int)    { w.meta.setRowSpan(n) }

func (w *Slider) Rect(x, y, width, height float32) *Slider {
	w.rect = R(x, y, width, height)
	return w
}
func (w *Slider) RectR(r Rect) *Slider { w.rect = r; return w }
func (w *Slider) Label(text string) *Slider {
	w.label = text
	return w
}
func (w *Slider) OnChange(fn func(float32)) *Slider {
	w.onChange = fn
	return w
}

func (w *Slider) ColSpan(n int) *Slider { w.meta.colSpan = n; return w }
func (w *Slider) RowSpan(n int) *Slider { w.meta.rowSpan = n; return w }

func (w *Slider) setFromMouse(x float32) func() {
	if w.value == nil || w.rect.W <= 0 {
		return nil
	}
	t := clamp((x-w.rect.X)/w.rect.W, 0, 1)
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
		ui.drawText(label, w.rect.X, w.rect.Y, 1, ui.Theme.TextMuted)
	}

	track := R(w.rect.X, w.rect.Y+w.rect.H*0.5-3, w.rect.W, 6)
	fill := R(track.X, track.Y, track.W*t, track.H)
	knob := R(w.rect.X+w.rect.W*t-5, w.rect.Y+w.rect.H*0.5-8, 10, 16)

	ui.drawRect(track, ui.Theme.Track)
	ui.drawRect(fill, ui.Theme.Accent)
	ui.drawRect(knob, ui.Theme.AccentHover)
}

// ─── Checkbox ────────────────────────────────────────────────────────────

type Checkbox struct {
	name     string
	label    string
	value    *bool
	rect     Rect
	onChange func(bool)
	meta     gridMeta
}

func (ui *UI) Checkbox(id, label string, value *bool) *Checkbox {
	if w, ok := ui.byID[id].(*Checkbox); ok {
		w.label = label
		w.value = value
		return w
	}
	w := &Checkbox{name: id, label: label, value: value, meta: defaultGridMeta()}
	ui.add(w)
	return w
}

func (w *Checkbox) id() string          { return w.name }
func (w *Checkbox) setBounds(r Rect)    { w.rect = r }
func (w *Checkbox) bounds() Rect        { return w.rect }
func (w *Checkbox) gridCol() int        { return w.meta.col }
func (w *Checkbox) gridRow() int        { return w.meta.row }
func (w *Checkbox) colSpan() int        { return w.meta.colSpan }
func (w *Checkbox) rowSpan() int        { return w.meta.rowSpan }
func (w *Checkbox) setGridPos(c, r int) { w.meta.setGridPos(c, r) }
func (w *Checkbox) setColSpan(n int)    { w.meta.setColSpan(n) }
func (w *Checkbox) setRowSpan(n int)    { w.meta.setRowSpan(n) }

func (w *Checkbox) Rect(x, y, width, height float32) *Checkbox {
	w.rect = R(x, y, width, height)
	return w
}
func (w *Checkbox) RectR(r Rect) *Checkbox { w.rect = r; return w }
func (w *Checkbox) OnChange(fn func(bool)) *Checkbox {
	w.onChange = fn
	return w
}

func (w *Checkbox) ColSpan(n int) *Checkbox { w.meta.colSpan = n; return w }
func (w *Checkbox) RowSpan(n int) *Checkbox { w.meta.rowSpan = n; return w }

func (w *Checkbox) draw(ui *UI) {
	box := R(w.rect.X, w.rect.Y+(w.rect.H-14)*0.5, 14, 14)
	ui.drawRect(box, ui.Theme.Input)
	ui.drawBorder(box, 1, ui.Theme.PanelBorder)

	if w.value != nil && *w.value {
		ui.drawRect(box.Inset(3), ui.Theme.Accent)
	}

	ui.drawText(w.label, w.rect.X+22, w.rect.Y+(w.rect.H-7)*0.5, 1, ui.Theme.Text)
}

// ─── ProgressBar ─────────────────────────────────────────────────────────

type ProgressBar struct {
	name     string
	value    *float32
	min, max float32
	rect     Rect
	meta     gridMeta
}

func (ui *UI) Progress(id string, value *float32, min, max float32) *ProgressBar {
	if w, ok := ui.byID[id].(*ProgressBar); ok {
		w.value = value
		w.min = min
		w.max = max
		return w
	}
	w := &ProgressBar{name: id, value: value, min: min, max: max, meta: defaultGridMeta()}
	ui.add(w)
	return w
}

func (w *ProgressBar) id() string          { return w.name }
func (w *ProgressBar) setBounds(r Rect)    { w.rect = r }
func (w *ProgressBar) bounds() Rect        { return w.rect }
func (w *ProgressBar) gridCol() int        { return w.meta.col }
func (w *ProgressBar) gridRow() int        { return w.meta.row }
func (w *ProgressBar) colSpan() int        { return w.meta.colSpan }
func (w *ProgressBar) rowSpan() int        { return w.meta.rowSpan }
func (w *ProgressBar) setGridPos(c, r int) { w.meta.setGridPos(c, r) }
func (w *ProgressBar) setColSpan(n int)    { w.meta.setColSpan(n) }
func (w *ProgressBar) setRowSpan(n int)    { w.meta.setRowSpan(n) }

func (w *ProgressBar) Rect(x, y, width, height float32) *ProgressBar {
	w.rect = R(x, y, width, height)
	return w
}
func (w *ProgressBar) RectR(r Rect) *ProgressBar { w.rect = r; return w }

func (w *ProgressBar) ColSpan(n int) *ProgressBar { w.meta.colSpan = n; return w }
func (w *ProgressBar) RowSpan(n int) *ProgressBar { w.meta.rowSpan = n; return w }

func (w *ProgressBar) draw(ui *UI) {
	value := w.min
	if w.value != nil {
		value = *w.value
	}
	t := float32(0)
	if w.max != w.min {
		t = clamp((value-w.min)/(w.max-w.min), 0, 1)
	}
	ui.drawRect(w.rect, ui.Theme.Track)
	ui.drawRect(R(w.rect.X, w.rect.Y, w.rect.W*t, w.rect.H), ui.Theme.Accent)
	ui.drawBorder(w.rect, 1, ui.Theme.PanelBorder)
}
