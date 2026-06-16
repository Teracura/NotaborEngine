package main

import (
	"NotaborEngine/internal/notasdl"
	"NotaborEngine/internal/notatask"
)

type Window struct {
	handle *notasdl.Window
}

func (w *Window) LoadVisual(name string, path string, opts VisualOptions) (*Visual, error) {
	v, err := w.handle.LoadVisual(name, path, notasdl.VisualOptions(opts))
	if err != nil {
		return nil, err
	}
	return &Visual{handle: v}, nil
}

func (w *Window) Draw(alpha float32, cam any, entity *Entity) error {
	return w.handle.Draw(alpha, nil, entity.handle)
}

func (w *Window) Move(dx, dy int) {
	w.handle.Move(dx, dy)
}

func (w *Window) Close() {
	w.handle.ShouldClose = true
}

type WindowConfig struct {
	X, Y, W, H int
	Title      string
	Type       WindowType
	Resizable  bool
	TargetFPS  float32
	Loops      []*notatask.Loop
}

type WindowType = notasdl.WindowType

const (
	WindowFullscreen = notasdl.Fullscreen
	WindowWindowed   = notasdl.Windowed
)

type VisualOptions = notasdl.VisualOptions
type VisualMask = notasdl.VisualMask

const (
	MaskNone   = notasdl.MaskNone
	MaskCircle = notasdl.MaskCircle
)
