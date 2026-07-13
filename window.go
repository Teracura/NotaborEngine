package main

import (
	"NotaborEngine/internal/notageometry"
	"NotaborEngine/internal/notasdl"
	"NotaborEngine/internal/notashader"
	"NotaborEngine/internal/notatask"
	"NotaborEngine/internal/notatexture"
)

// Window represents an OS window with a GPU render surface.
// It handles drawing, input events, and window-state management.
type Window struct {
	handle *notasdl.Window
}

// LoadVisual loads a texture from disk and creates a Visual (sprite + material)
// using the given options (size, circle-mask, etc.).
func (w *Window) LoadVisual(name string, path string, opts VisualOptions) (*Visual, error) {
	v, err := w.handle.LoadVisual(name, path, notasdl.VisualOptions(opts))
	if err != nil {
		return nil, err
	}
	return &Visual{handle: v}, nil
}

// Draw renders entities through the window's camera view.
// Pass nil for cam to use the window's default camera.
func (w *Window) Draw(alpha float32, cam *Camera2D, entity *Entity) error {
	internalCam := (*notasdl.Camera2D)(nil)
	if cam != nil {
		internalCam = cam.handle
	}
	return w.handle.Draw(alpha, internalCam, entity.handle)
}

// Move translates the window by the given pixel delta (with boundary clamping).
func (w *Window) Move(dx, dy int) {
	w.handle.Move(dx, dy)
}

// Close marks the window for closure on the next frame.
func (w *Window) Close() {
	w.handle.ShouldClose = true
}

// LoadTexture loads an image from disk and caches it as a GPU texture.
func (w *Window) LoadTexture(name, path string) (*notatexture.Texture, error) {
	return w.handle.LoadTexture(name, path)
}

// GetTexture retrieves a previously loaded texture by name.
func (w *Window) GetTexture(name string) (*notatexture.Texture, error) {
	return w.handle.GetTexture(name)
}

// UnloadTexture releases a cached texture.
func (w *Window) UnloadTexture(name string) error {
	return w.handle.UnloadTexture(name)
}

// LoadShader compiles and caches a vertex/fragment shader pair.
func (w *Window) LoadShader(name, vertexPath, fragmentPath string) (*notashader.Shader, error) {
	return w.handle.LoadShader(name, vertexPath, fragmentPath)
}

// GetShader retrieves a previously loaded shader by name.
func (w *Window) GetShader(name string) (*notashader.Shader, error) {
	return w.handle.GetShader(name)
}

// ReloadShader recompiles a cached shader from its original source files.
func (w *Window) ReloadShader(name string) error {
	return w.handle.ReloadShader(name)
}

// UnloadShader releases a cached shader and its GPU resources.
func (w *Window) UnloadShader(name string) error {
	return w.handle.UnloadShader(name)
}

// CreateMaterial creates a new material from an existing shader.
func (w *Window) CreateMaterial(shader *notashader.Shader) *notashader.Material {
	return w.handle.CreateMaterial(shader)
}

// LoadMaterial compiles a shader pair and returns a material wrapping it.
func (w *Window) LoadMaterial(name, vertexPath, fragmentPath string) (*notashader.Material, error) {
	return w.handle.LoadMaterial(name, vertexPath, fragmentPath)
}

// CreateSprite creates a named sprite from an existing texture and polygon.
func (w *Window) CreateSprite(name string, texture *notatexture.Texture, polygon *notageometry.Polygon) (*notatexture.Sprite, error) {
	return w.handle.CreateSprite(name, texture, polygon)
}

// LoadSprite loads a texture and creates a sprite from it in one step.
func (w *Window) LoadSprite(name, texturePath string, polygon *notageometry.Polygon) (*notatexture.Sprite, error) {
	return w.handle.LoadSprite(name, texturePath, polygon)
}

// GetSprite retrieves a previously created sprite by name.
func (w *Window) GetSprite(name string) (*notatexture.Sprite, error) {
	return w.handle.GetSprite(name)
}

// RemoveSprite deletes a sprite (does not remove the underlying texture).
func (w *Window) RemoveSprite(name string) error {
	return w.handle.RemoveSprite(name)
}

// LoadCircleSprite is a convenience helper that loads a texture and creates
// a sprite with a circular shader mask of the given radius.
func (w *Window) LoadCircleSprite(name, texturePath string, radius float32) (*notatexture.Sprite, *notashader.Material, error) {
	return w.handle.LoadCircleSprite(name, texturePath, radius)
}

// GetSize returns the current window dimensions.
func (w *Window) GetSize() (width, height int) {
	return w.handle.GetSize()
}

// SetSize updates the window size at runtime.
func (w *Window) SetSize(width, height int) {
	w.handle.SetSize(width, height)
}

// GetPosition returns the window's current screen position.
func (w *Window) GetPosition() (x, y int) {
	return w.handle.GetPosition()
}

// SetPosition moves the window to the given screen coordinates (with boundary clamping).
func (w *Window) SetPosition(x, y int) {
	w.handle.SetPosition(x, y)
}

// SetVSync enables or disables vertical sync. When disabled, the driver may
// fall back to mailbox or immediate present modes.
func (w *Window) SetVSync(enabled bool) {
	w.handle.SetVSync(enabled)
}

// WindowConfig describes the initial properties of a window.
type WindowConfig struct {
	X, Y, W, H int
	Title      string
	Type       WindowType
	Resizable  bool
	TargetFPS  float32
	Loops      []*notatask.Loop
}

// WindowType specifies the display mode of a window.
type WindowType = notasdl.WindowType

const (
	WindowFullscreen = notasdl.Fullscreen
	WindowWindowed   = notasdl.Windowed
)

// VisualOptions controls sprite and material properties during visual creation.
type VisualOptions = notasdl.VisualOptions

// VisualMask selects a shader mask effect for the visual.
type VisualMask = notasdl.VisualMask

const (
	MaskNone   = notasdl.MaskNone
	MaskCircle = notasdl.MaskCircle
)
