// Package main provides the public facade for the NotaborEngine.
// All public types wrap internal implementations, keeping the engine
// internals hidden from users while exposing a clean, documented API.
package main

import (
	"NotaborEngine/internal/notacore"
	"NotaborEngine/internal/notasdl"
)

// Engine is the top-level game engine instance. Create one with CreateEngine,
// configure windows and entities, then call Run to start the main loop.
type Engine struct {
	handle *notacore.Engine
}

// Settings control engine-wide behaviour: vsync, audio mute, and master volume.
type Settings struct {
	Vsync      bool    // Locks TargetFPS to monitor's refresh rate
	Muted      bool    // Completely disables sound output
	SoundLevel float32 // Master volume (0.0 = silent, 1.0 = full)
}

func wrapSettings(s *Settings) *notacore.Settings {
	return &notacore.Settings{
		Vsync:      s.Vsync,
		Muted:      s.Muted,
		SoundLevel: s.SoundLevel,
	}
}

// CreateEngine initialises a new engine with the given settings.
// It sets up the platform layer (SDL), audio, input, and entity systems.
func CreateEngine(settings *Settings) (*Engine, error) {
	h, err := notacore.CreateEngine(wrapSettings(settings))
	if err != nil {
		return nil, err
	}
	return &Engine{handle: h}, nil
}

// Run starts the engine main loop and blocks until all windows are closed.
func (e *Engine) Run() error {
	return e.handle.Run()
}

// Shutdown cleans up all engine resources (windows, GPU, audio, input).
func (e *Engine) Shutdown() {
	e.handle.Shutdown()
}

// CreateWindow creates a new window on screen with the given configuration.
// The window is automatically registered with the engine's render loop.
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

// Input returns the engine's input manager for configuring input signals.
func (e *Engine) Input() *InputManager {
	return &InputManager{handle: e.handle.Input}
}

// EntityManager returns the engine's entity manager for creating and
// managing game entities (ECS).
func (e *Engine) EntityManager() *EntityManager {
	return &EntityManager{handle: e.handle.EntityManager}
}

// SoundManager returns the engine's sound manager for playing audio.
func (e *Engine) SoundManager() *SoundManager {
	return &SoundManager{handle: e.handle.SoundManager}
}

// GetSettings returns a copy of the engine's current settings.
func (e *Engine) GetSettings() Settings {
	s := e.handle.GetSettings()
	return Settings{
		Vsync:      s.Vsync,
		Muted:      s.Muted,
		SoundLevel: s.SoundLevel,
	}
}

// ChangeSettings applies new settings to the engine at runtime.
func (e *Engine) ChangeSettings(s *Settings) {
	e.handle.ChangeSettings(wrapSettings(s))
}
