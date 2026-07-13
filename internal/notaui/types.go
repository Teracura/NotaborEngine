package notaui

import (
	"NotaborEngine/internal/notacolor"
)

type Rect struct {
	X, Y, W, H float32
}

func R(x, y, w, h float32) Rect {
	return Rect{X: x, Y: y, W: w, H: h}
}

func (r Rect) Contains(x, y float32) bool {
	return x >= r.X && x <= r.X+r.W && y >= r.Y && y <= r.Y+r.H
}

func (r Rect) Inset(v float32) Rect {
	return Rect{X: r.X + v, Y: r.Y + v, W: r.W - v*2, H: r.H - v*2}
}

func (r Rect) IsEmpty() bool {
	return r.W <= 0 || r.H <= 0
}

// Widget is the interface all UI elements satisfy.
// Grid is also a Widget (it implements draw/bounds/setBounds).
type Widget interface {
	draw(*UI)
	id() string
	setBounds(Rect)
	bounds() Rect
	gridCol() int
	gridRow() int
	colSpan() int
	rowSpan() int
	setGridPos(col, row int)
	setColSpan(n int)
	setRowSpan(n int)
}

type gridMeta struct {
	col     int
	row     int
	colSpan int
	rowSpan int
}

func defaultGridMeta() gridMeta {
	return gridMeta{col: -1, row: -1, colSpan: 1, rowSpan: 1}
}

func (m *gridMeta) setGridPos(col, row int) { m.col, m.row = col, row }
func (m *gridMeta) setColSpan(n int)        { m.colSpan = n }
func (m *gridMeta) setRowSpan(n int)        { m.rowSpan = n }

type Theme struct {
	Panel        notacolor.Color
	PanelBorder  notacolor.Color
	Surface      notacolor.Color
	SurfaceHover notacolor.Color
	SurfaceDown  notacolor.Color
	Accent       notacolor.Color
	AccentHover  notacolor.Color
	Text         notacolor.Color
	TextMuted    notacolor.Color
	Input        notacolor.Color
	InputFocus   notacolor.Color
	Track        notacolor.Color
}

func DefaultTheme() Theme {
	return Theme{
		Panel:        notacolor.RGBA(0.06, 0.07, 0.08, 0.82),
		PanelBorder:  notacolor.RGBA(0.36, 0.42, 0.45, 0.85),
		Surface:      notacolor.RGBA(0.16, 0.19, 0.20, 0.92),
		SurfaceHover: notacolor.RGBA(0.22, 0.27, 0.29, 0.95),
		SurfaceDown:  notacolor.RGBA(0.10, 0.13, 0.14, 0.98),
		Accent:       notacolor.RGBA(0.10, 0.62, 0.55, 1),
		AccentHover:  notacolor.RGBA(0.16, 0.78, 0.68, 1),
		Text:         notacolor.RGBA(0.93, 0.96, 0.94, 1),
		TextMuted:    notacolor.RGBA(0.62, 0.68, 0.67, 1),
		Input:        notacolor.RGBA(0.03, 0.04, 0.045, 0.94),
		InputFocus:   notacolor.RGBA(0.08, 0.14, 0.14, 0.98),
		Track:        notacolor.RGBA(0.08, 0.10, 0.105, 0.96),
	}
}

func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}
