package notassets

import (
	"NotaborEngine/notacollision"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
)

type Entity struct {
	ID   string
	Name string

	Sprite  *Sprite
	Polygon notagl.Polygon

	Collider notacollision.Collider

	Active  bool
	Visible bool
}

// NewEntity creates a basic entity with a sprite
func NewEntity(id, name string, sprite *Sprite) *Entity {
	return &Entity{
		ID:      id,
		Name:    name,
		Sprite:  sprite,
		Active:  true,
		Visible: true,
	}
}

// NewEntityWithPolygon creates an entity with a polygon (for shapes without textures)
func NewEntityWithPolygon(id, name string, polygon notagl.Polygon) *Entity {
	return &Entity{
		ID:      id,
		Name:    name,
		Polygon: polygon,
		Active:  true,
		Visible: true,
	}
}

// SetPosition updates both sprite position and polygon transform
func (e *Entity) SetPosition(x, y float32) {
	if e.Sprite != nil {
		e.Sprite.X = int32(x)
		e.Sprite.Y = int32(y)
	}

	// Update polygon position if it exists
	if len(e.Polygon.Vertices) > 0 {
		e.Polygon.Transform.Position = notamath.Vec2{X: x, Y: y}
	}
}

// GetPosition returns the entity's position
func (e *Entity) GetPosition() (float32, float32) {
	if e.Sprite != nil {
		return float32(e.Sprite.X), float32(e.Sprite.Y)
	}

	if len(e.Polygon.Vertices) > 0 {
		return e.Polygon.Transform.Position.X, e.Polygon.Transform.Position.Y
	}

	return 0, 0
}
