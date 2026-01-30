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
	Polygon *notagl.Polygon

	Collider notacollision.Collider

	Active  bool
	Visible bool
}

func NewEntity(id, name string) *Entity {
	return &Entity{
		ID:      id,
		Name:    name,
		Active:  true,
		Visible: true,
	}
}

func (e *Entity) SetPolygon(Polygon *notagl.Polygon) {
	e.Polygon = Polygon
}

func (e *Entity) SetCollider(Collider notacollision.Collider) {
	e.Collider = Collider
}

func (e *Entity) SetSprite(sprite *Sprite) {
	e.Sprite = sprite
}

func (e *Entity) Move(delta notamath.Vec2) {
	if e.Sprite != nil {
		e.Sprite.X += delta.X
		e.Sprite.Y += delta.Y
	}

	// Update polygon position if it exists
	if len(e.Polygon.Vertices) > 0 {
		e.Polygon.Transform.TranslateBy(delta)
	}

	if e.Collider != nil {
		e.Collider.Move(delta)
	}
}
