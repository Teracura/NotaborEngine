package notassets

import (
	"NotaborEngine/notacollision"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"math"
)

// Entity is the central unit: drawable, collidable, movable
type Entity struct {
	ID   string
	Name string

	Sprite   *Sprite
	Polygon  *notagl.Polygon
	Collider notacollision.Collider

	Active  bool
	Visible bool
}

// NewEntity creates a basic empty entity
func NewEntity(id, name string) *Entity {
	return &Entity{
		ID:      id,
		Name:    name,
		Active:  true,
		Visible: true,
	}
}

func (e *Entity) SetSprite(s *Sprite)                  { e.Sprite = s }
func (e *Entity) SetPolygon(p *notagl.Polygon)         { e.Polygon = p }
func (e *Entity) SetCollider(c notacollision.Collider) { e.Collider = c }

func (e *Entity) Move(delta notamath.Vec2) {
	if !e.Active {
		return
	}

	if e.Sprite != nil {
		e.Sprite.X += delta.X
		e.Sprite.Y += delta.Y
	}

	if e.Polygon != nil && len(e.Polygon.Vertices) > 0 {
		e.Polygon.Transform.TranslateBy(delta)
	}

	if e.Collider != nil {
		e.Collider.Move(delta)
	}
}

// Rotate rotates polygon (optional for sprite)
func (e *Entity) Rotate(rad float32) {
	if !e.Active {
		return
	}

	if e.Polygon != nil {
		e.Polygon.Transform.RotateBy(rad)
	}

	if e.Collider != nil {
		e.Collider.Rotate(rad)
	}

	if e.Sprite != nil {
		e.Sprite.X += e.Sprite.X * float32(math.Cos(float64(rad)))
		e.Sprite.Y += e.Sprite.Y * float32(math.Sin(float64(rad)))
	}
}

func (e *Entity) Draw(renderer *notagl.Renderer2D) {
	if !e.Visible || !e.Active {
		return
	}

	// Draw polygon first
	if e.Polygon != nil && len(e.Polygon.Vertices) > 0 {
		renderer.Submit(e.Polygon)
	}

	// Draw sprite on top
	if e.Sprite != nil && e.Sprite.Texture != nil {
		quad := notagl.CreateTextureQuad(
			notamath.Po2{X: e.Sprite.X, Y: e.Sprite.Y},
			float32(e.Sprite.srcWidth), float32(e.Sprite.srcHeight),
		)
		quad.Transform.Position = notamath.Vec2{X: e.Sprite.X, Y: e.Sprite.Y}
		renderer.Submit(&quad)
	}
}

// CollidesWith returns true if this entity collides with another
func (e *Entity) CollidesWith(other *Entity) bool {
	if e.Collider == nil || other.Collider == nil {
		return false
	}
	return notacollision.Intersects(e.Collider, other.Collider)
}
