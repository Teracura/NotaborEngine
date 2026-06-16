package main

import (
	"NotaborEngine/internal/notacore"
)

// InputContext is an opaque handle to the input state
type InputContext struct {
	handle *notacore.InputContext
}

// InputSignal represents a bindable input action
type InputSignal struct {
	handle *notacore.InputSignal
}

func (s *InputSignal) Held() bool {
	return s.handle.Held()
}

func (s *InputSignal) Pressed() bool {
	return s.handle.Pressed()
}

func (s *InputSignal) Released() bool {
	return s.handle.Released()
}

func (s *InputSignal) Idle() bool {
	return s.handle.Idle()
}

// InputManager manages engine inputs
type InputManager struct {
	handle *notacore.InputManager
}

func (im *InputManager) Start(Hz float32) {
	im.handle.Start(Hz)
}

func (im *InputManager) Stop() {
	im.handle.Stop()
}

func (im *InputManager) GetContext() *InputContext {
	return &InputContext{handle: im.handle.GetContext()}
}

func (im *InputManager) Get(name string) *InputSignal {
	s := im.handle.Get(name)
	if s == nil {
		return nil
	}
	return &InputSignal{handle: s}
}

// InputKey represents a physical input (key, mouse button, etc.)
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

// Input creates a simple input signal
func Input(name string, key InputKey, ctx *InputContext) *InputSignal {
	s := notacore.Input(name, notacore.StateInput(key), ctx.handle)
	return &InputSignal{handle: s}
}

// InputCombo creates an input signal from multiple keys
func InputCombo(name string, ctx *InputContext, keys ...InputKey) *InputSignal {
	internalKeys := make([]notacore.StateInput, len(keys))
	for i, k := range keys {
		internalKeys[i] = notacore.StateInput(k)
	}
	s := notacore.InputCombo(name, ctx.handle, internalKeys...)
	return &InputSignal{handle: s}
}
