package main

import (
	"NotaborEngine/internal/notamath"
	"NotaborEngine/internal/notasdl"
)

// Camera2D controls how a window views world-space content: position, rotation, zoom,
// with support for smooth transitions.
type Camera2D struct {
	handle *notasdl.Camera2D
}

// NewCamera2D creates a 2D camera with identity zoom and zero position/rotation.
func NewCamera2D() *Camera2D {
	return &Camera2D{handle: notasdl.NewCamera2D()}
}

// Position returns the camera's world-space centre.
func (c *Camera2D) Position() Vec2 {
	return Vec2(c.handle.Position())
}

// SetPosition moves the camera to an absolute world-space position.
func (c *Camera2D) SetPosition(pos Vec2) {
	c.handle.SetPosition(notamath.Vec2(pos))
}

// Move translates the camera by a world-space delta.
func (c *Camera2D) Move(delta Vec2) {
	c.handle.Move(notamath.Vec2(delta))
}

// Rotation returns the camera rotation in radians.
func (c *Camera2D) Rotation() float32 {
	return c.handle.Rotation()
}

// SetRotation sets the camera rotation in radians.
func (c *Camera2D) SetRotation(rot float32) {
	c.handle.SetRotation(rot)
}

// Rotate adds a rotation delta in radians.
func (c *Camera2D) Rotate(delta float32) {
	c.handle.Rotate(delta)
}

// Zoom returns the camera zoom factor on each axis.
func (c *Camera2D) Zoom() Vec2 {
	return Vec2(c.handle.Zoom())
}

// SetZoom sets a uniform zoom factor on both axes.
func (c *Camera2D) SetZoom(zoom float32) {
	c.handle.SetZoom(zoom)
}

// SetZoomXY sets independent zoom factors for each axis.
func (c *Camera2D) SetZoomXY(x, y float32) {
	c.handle.SetZoomXY(x, y)
}

// ViewMatrix returns the affine matrix that transforms world space into camera view.
func (c *Camera2D) ViewMatrix() Mat3 {
	return Mat3(c.handle.ViewMatrix())
}
