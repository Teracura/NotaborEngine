package main

import (
	"NotaborEngine/internal/notacore"
	"NotaborEngine/internal/notasdl"
)

// Engine is the main engine structure
type Engine struct {
	handle *notacore.Engine
}

// Settings represents the engine's configuration settings
type Settings struct {
	Vsync      bool
	Muted      bool
	SoundLevel float32
}

func wrapSettings(s *Settings) *notacore.Settings {
	return &notacore.Settings{
		Vsync:      s.Vsync,
		Muted:      s.Muted,
		SoundLevel: s.SoundLevel,
	}
}

// CreateEngine initializes a new engine instance with the given settings
func CreateEngine(settings *Settings) (*Engine, error) {
	h, err := notacore.CreateEngine(wrapSettings(settings))
	if err != nil {
		return nil, err
	}
	return &Engine{handle: h}, nil
}

// Run starts the engine main loop
func (e *Engine) Run() error {
	return e.handle.Run()
}

// Shutdown cleans up engine resources
func (e *Engine) Shutdown() {
	e.handle.Shutdown()
}

// CreateWindow creates a new window
func (e *Engine) CreateWindow(cfg *WindowConfig) (*Window, error) {
	internalCfg := &notasdl.WindowConfig{
		X:         cfg.X,
		Y:         cfg.Y,
		W:         cfg.W,
		H:         cfg.H,
		Title:     cfg.Title,
		Type:      cfg.Type,
		Resizable: cfg.Resizable,
		TargetFPS: cfg.TargetFPS,
		Loops:     cfg.Loops,
	}
	w, err := e.handle.CreateWindow(internalCfg)
	if err != nil {
		return nil, err
	}
	return &Window{handle: w}, nil
}

// Input returns the engine's input manager facade
func (e *Engine) Input() *InputManager {
	return &InputManager{handle: e.handle.Input}
}

// EntityManager returns the engine's entity manager facade
func (e *Engine) EntityManager() *EntityManager {
	return &EntityManager{handle: e.handle.EntityManager}
}
