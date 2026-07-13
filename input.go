package main

import (
	"NotaborEngine/internal/notacore"
)

// InputContext holds per-frame input state for querying key/mouse/gamepad.
//
// Use engine.Input().GetContext() to obtain the active context.
type InputContext struct {
	handle *notacore.InputContext
}

// InputSignal represents a named, bindable input action (e.g. "jump", "move-left").
// Query its state each frame with Pressed / Held / Released / Idle.
type InputSignal struct {
	handle *notacore.InputSignal
}

// Held returns true while the input is physically held down.
func (s *InputSignal) Held() bool {
	return s.handle.Held()
}

// Pressed returns true on the exact frame the input became active.
func (s *InputSignal) Pressed() bool {
	return s.handle.Pressed()
}

// Released returns true on the exact frame the input became inactive.
func (s *InputSignal) Released() bool {
	return s.handle.Released()
}

// Idle returns true when the input is not currently active.
func (s *InputSignal) Idle() bool {
	return s.handle.Idle()
}

// InputManager owns input contexts and signal registrations.
// Call Start to begin polling, Stop to pause.
type InputManager struct {
	handle *notacore.InputManager
}

// Start begins the input-polling loop at the given frequency (Hz).
func (im *InputManager) Start(Hz float32) {
	im.handle.Start(Hz)
}

// Stop halts the input-polling loop.
func (im *InputManager) Stop() {
	im.handle.Stop()
}

// GetContext returns the default input context.
func (im *InputManager) GetContext() *InputContext {
	return &InputContext{handle: im.handle.GetContext()}
}

// Get retrieves a previously bound signal by name. Returns nil if not found.
func (im *InputManager) Get(name string) *InputSignal {
	s := im.handle.Get(name)
	if s == nil {
		return nil
	}
	return &InputSignal{handle: s}
}

// InputKey represents a physical input (keyboard key, mouse button, gamepad button/axis).
type InputKey notacore.StateInput

const (
	KeyW = InputKey(notacore.KeyW)
	KeyA = InputKey(notacore.KeyA)
	KeyS = InputKey(notacore.KeyS)
	KeyD = InputKey(notacore.KeyD)
	KeyQ = InputKey(notacore.KeyQ)
	KeyE = InputKey(notacore.KeyE)

	MouseLeft   = InputKey(notacore.MouseLeft)
	MouseRight  = InputKey(notacore.MouseRight)
	MouseMiddle = InputKey(notacore.MouseMiddle)
)

// Input creates a simple named input signal bound to a single key.
func Input(name string, key InputKey, ctx *InputContext) *InputSignal {
	s := notacore.Input(name, notacore.StateInput(key), ctx.handle)
	return &InputSignal{handle: s}
}

// InputCombo creates an input signal that activates when all given keys are held.
func InputCombo(name string, ctx *InputContext, keys ...InputKey) *InputSignal {
	internalKeys := make([]notacore.StateInput, len(keys))
	for i, k := range keys {
		internalKeys[i] = notacore.StateInput(k)
	}
	s := notacore.InputCombo(name, ctx.handle, internalKeys...)
	return &InputSignal{handle: s}
}
