package main

import (
	"NotaborEngine/internal/notacolor"
)

// Color represents an RGBA color using normalized float32 values [0.0, 1.0].
type Color struct {
	R, G, B, A float32
}

// Built-in color constants matching the internal engine palette.
var (
	White       = Color{1, 1, 1, 1}
	Black       = Color{0, 0, 0, 1}
	Red         = Color{1, 0, 0, 1}
	Green       = Color{0, 1, 0, 1}
	Blue        = Color{0, 0, 1, 1}
	Magenta     = Color{1, 0, 1, 1}
	Yellow      = Color{1, 1, 0, 1}
	Cyan        = Color{0, 1, 1, 1}
	Gray        = Color{0.5, 0.5, 0.5, 1}
	Silver      = Color{0.75, 0.75, 0.75, 1}
	Maroon      = Color{0.5, 0, 0, 1}
	Olive       = Color{0.5, 0.5, 0, 1}
	Navy        = Color{0, 0, 0.5, 1}
	Purple      = Color{0.5, 0, 0.5, 1}
	Teal        = Color{0, 0.5, 0.5, 1}
	Orange      = Color{1, 0.5, 0, 1}
	Transparent = Color{0, 0, 0, 0}
)

// RGBA returns a new clamped color with explicit red, green, blue, and alpha channels.
func RGBA(r, g, b, a float32) Color {
	return fromInternal(notacolor.RGBA(r, g, b, a))
}

// RGB returns a new clamped color with full opacity.
func RGB(r, g, b float32) Color {
	return fromInternal(notacolor.RGB(r, g, b))
}

// FromHex parses a hex color string (e.g. "#FF0000" or "FF0000FF") into a public Color structure.
func FromHex(hex string) (Color, error) {
	internalColor, err := notacolor.FromHex(hex)
	if err != nil {
		return Color{}, err
	}
	return fromInternal(internalColor), nil
}

// WithAlpha returns a copy of the color with a modified alpha value.
func (c Color) WithAlpha(a float32) Color {
	return fromInternal(toInternal(c).WithAlpha(a))
}

// Clamp returns a copy of the color with all values constrained between 0.0 and 1.0.
func (c Color) Clamp() Color {
	return fromInternal(toInternal(c).Clamp())
}

// ToVec4 extracts the color channels into a standard 4-element array.
func (c Color) ToVec4() [4]float32 {
	return c.toInternal().ToVec4()
}

// Lerp linearly interpolates between this color and another based on alpha variant t [0.0, 1.0].
func (c Color) Lerp(to Color, t float32) Color {
	return fromInternal(toInternal(c).Lerp(toInternal(to), t))
}

// --- Internal Conversion Utilities --- (TEMPORARY)

// toInternal translates a public facade Color structure into its internal implementation counterpart.
func (c Color) toInternal() notacolor.Color {
	return notacolor.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// toInternal package level helper function
func toInternal(c Color) notacolor.Color {
	return notacolor.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// fromInternal translates an internal engine color back into the clean public facade representation.
func fromInternal(c notacolor.Color) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: c.A}
}
