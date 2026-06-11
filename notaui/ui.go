package notaui

import (
	"errors"
	"sync"

	"NotaborEngine/notacolor"
	"NotaborEngine/notacore"
	"NotaborEngine/notageometry"
	"NotaborEngine/notamath"
	"NotaborEngine/notasdl"
	"NotaborEngine/notashader"

	"github.com/Zyko0/go-sdl3/sdl"
)

type widget interface {
	draw(*UI)
	id() string
}

type UI struct {
	window   *notasdl.Window
	material *notashader.Material
	rect     *notageometry.Polygon

	Theme Theme

	mu        sync.Mutex
	widgets   []widget
	byID      map[string]widget
	mouseX    float32
	mouseY    float32
	mouseDown bool
	activeID  string
	focusID   string
}

func New(engine *notacore.Engine, window *notasdl.Window) (*UI, error) {
	if window == nil {
		return nil, errors.New("notaui: window is nil")
	}
	if window.Runtime == nil || window.Runtime.ShaderMgr == nil {
		return nil, errors.New("notaui: window runtime is not initialized")
	}

	material, err := window.BasicMaterial()
	if err != nil {
		return nil, err
	}

	ui := &UI{
		window:   window,
		material: material,
		rect:     notageometry.CreateRectangle(1, 1),
		Theme:    DefaultTheme(),
		byID:     make(map[string]widget),
	}

	if engine != nil && engine.Platform != nil {
		engine.Platform.SubscribeEvents(ui.HandleEvent)
	}

	return ui, nil
}

func (ui *UI) Draw() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	for _, w := range ui.widgets {
		w.draw(ui)
	}
}

func (ui *UI) HasKeyboardFocus() bool {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	return ui.focusID != ""
}

func (ui *UI) HandleEvent(event notasdl.Event) {
	if ui == nil || ui.window == nil || event.WindowID != uint32(ui.window.ID) {
		return
	}

	var callbacks []func()

	ui.mu.Lock()
	switch event.Type {
	case notasdl.EventMouseMove:
		ui.mouseX = event.X
		ui.mouseY = event.Y
		callbacks = append(callbacks, ui.dragActive(event.X, event.Y)...)

	case notasdl.EventMouseDown:
		if event.MouseBtn == notasdl.MouseButton(sdl.BUTTON_LEFT) {
			ui.mouseX = event.X
			ui.mouseY = event.Y
			ui.mouseDown = true
			ui.activeID = ""
			ui.focusID = ""

			if hit := ui.hit(event.X, event.Y); hit != nil {
				ui.activeID = hit.id()
				if _, ok := hit.(*Input); ok {
					ui.focusID = hit.id()
				}
				callbacks = append(callbacks, ui.dragActive(event.X, event.Y)...)
			}
		}

	case notasdl.EventMouseUp:
		if event.MouseBtn == notasdl.MouseButton(sdl.BUTTON_LEFT) {
			ui.mouseX = event.X
			ui.mouseY = event.Y
			ui.mouseDown = false
			if hit := ui.hit(event.X, event.Y); hit != nil && hit.id() == ui.activeID {
				callbacks = append(callbacks, ui.activate(hit)...)
			}
			ui.activeID = ""
		}

	case notasdl.EventKeyDown:
		callbacks = append(callbacks, ui.keyDown(event.Key)...)
	}
	ui.mu.Unlock()

	for _, callback := range callbacks {
		if callback != nil {
			callback()
		}
	}
}

func (ui *UI) add(w widget) {
	if _, exists := ui.byID[w.id()]; exists {
		return
	}
	ui.byID[w.id()] = w
	ui.widgets = append(ui.widgets, w)
}

func (ui *UI) hit(x, y float32) widget {
	for i := len(ui.widgets) - 1; i >= 0; i-- {
		switch w := ui.widgets[i].(type) {
		case *Button:
			if w.bounds.Contains(x, y) {
				return w
			}
		case *Input:
			if w.bounds.Contains(x, y) {
				return w
			}
		case *Slider:
			if w.bounds.Contains(x, y) {
				return w
			}
		case *Checkbox:
			if w.bounds.Contains(x, y) {
				return w
			}
		}
	}
	return nil
}

func (ui *UI) activate(w widget) []func() {
	switch v := w.(type) {
	case *Button:
		if v.onClick != nil {
			return []func(){v.onClick}
		}
	case *Checkbox:
		if v.value != nil {
			*v.value = !*v.value
			if v.onChange != nil {
				value := *v.value
				return []func(){func() { v.onChange(value) }}
			}
		}
	}
	return nil
}

func (ui *UI) dragActive(x, y float32) []func() {
	if ui.activeID == "" {
		return nil
	}
	if w, ok := ui.byID[ui.activeID].(*Slider); ok {
		if callback := w.setFromMouse(x); callback != nil {
			return []func(){callback}
		}
	}
	return nil
}

func (ui *UI) keyDown(key notasdl.Key) []func() {
	w, ok := ui.byID[ui.focusID].(*Input)
	if !ok || w.value == nil {
		return nil
	}

	changed := func() []func() {
		if w.onChange == nil {
			return nil
		}
		value := *w.value
		return []func(){func() { w.onChange(value) }}
	}

	switch sdl.Keycode(uint32(key)) {
	case sdl.K_BACKSPACE:
		if len(*w.value) > 0 {
			*w.value = (*w.value)[:len(*w.value)-1]
			return changed()
		}
	case sdl.K_RETURN, sdl.K_ESCAPE:
		ui.focusID = ""
	case sdl.K_SPACE:
		*w.value += " "
		return changed()
	default:
		r := rune(key)
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}
		if r >= 32 && r <= 126 {
			*w.value += string(r)
			return changed()
		}
	}

	return nil
}

func (ui *UI) drawRect(r Rect, color notacolor.Color) {
	if r.IsEmpty() || color.A <= 0 {
		return
	}

	w := float32(ui.window.Config.W)
	h := float32(ui.window.Config.H)
	if w <= 0 || h <= 0 {
		return
	}

	clipW := (r.W / w) * 2
	clipH := (r.H / h) * 2
	centerX := ((r.X + r.W*0.5) / w * 2) - 1
	centerY := 1 - ((r.Y + r.H*0.5) / h * 2)

	model := notamath.Mat3TRS(
		notamath.Vec2{X: centerX, Y: centerY},
		0,
		notamath.Vec2{X: clipW, Y: clipH},
	)

	ui.window.Runtime.Renderer.SubmitPolygon(ui.rect, model, color, nil, nil, ui.material)
}

func (ui *UI) drawBorder(r Rect, thickness float32, color notacolor.Color) {
	if thickness <= 0 {
		return
	}
	ui.drawRect(R(r.X, r.Y, r.W, thickness), color)
	ui.drawRect(R(r.X, r.Y+r.H-thickness, r.W, thickness), color)
	ui.drawRect(R(r.X, r.Y, thickness, r.H), color)
	ui.drawRect(R(r.X+r.W-thickness, r.Y, thickness, r.H), color)
}

func (ui *UI) drawText(text string, x, y, scale float32, color notacolor.Color) {
	if scale <= 0 {
		scale = 1
	}
	cx := x
	for _, r := range text {
		if r == '\n' {
			cx = x
			y += 8 * scale
			continue
		}
		ui.drawGlyph(r, cx, y, scale, color)
		cx += 6 * scale
	}
}

func (ui *UI) textWidth(text string, scale float32) float32 {
	if scale <= 0 {
		scale = 1
	}
	return float32(len([]rune(text))) * 6 * scale
}
