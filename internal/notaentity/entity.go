package notaentity

import (
	"NotaborEngine/internal/notacollision"
	"NotaborEngine/internal/notacolor"
	"NotaborEngine/internal/notageometry"
	"NotaborEngine/internal/notamath"
	"NotaborEngine/internal/notarender"
	"NotaborEngine/internal/notashader"
	"NotaborEngine/internal/notatexture"
)

// EntityID is a unique identifier for an entity.
type EntityID uint32

// Entity represents an entity in the ECS system.
// It is a lightweight handle that references components stored in the EntityManager.
type Entity struct {
	ID      EntityID
	manager *EntityManager
}

// Visual is a reusable visual bundle containing sprite and material.
type Visual struct {
	Sprite   *notatexture.Sprite
	Material *notashader.Material
}

// CollisionProfile describes how an entity should participate in collision queries.
type CollisionProfile struct {
	Collider notacollision.Collider
}

// NewVisual creates a reusable visual bundle from a sprite and optional material.
func NewVisual(sprite *notatexture.Sprite, material *notashader.Material) *Visual {
	return &Visual{
		Sprite:   sprite,
		Material: material,
	}
}

// CircleCollision creates a circular collision profile centered on the entity origin.
func CircleCollision(radius float32) CollisionProfile {
	return CollisionProfile{
		Collider: notacollision.NewCircleCollider(notamath.Po2{}, radius),
	}
}

// PolygonCollision creates a polygon collision profile from local-space points.
func PolygonCollision(points []notamath.Po2) CollisionProfile {
	return CollisionProfile{
		Collider: notacollision.NewPolygonCollider(points),
	}
}

// CustomCollision wraps an already constructed collider as a collision profile.
func CustomCollision(collider notacollision.Collider) CollisionProfile {
	return CollisionProfile{
		Collider: collider,
	}
}

// newEntity creates a new entity with the given ID and manager reference.
func newEntity(id EntityID, manager *EntityManager) Entity {
	return Entity{
		ID:      id,
		manager: manager,
	}
}

// WithSprite assigns a sprite to the entity.
func (e *Entity) WithSprite(s *notatexture.Sprite) *Entity {
	e.manager.SetSprite(e.ID, s)
	return e
}

// WithPolygon assigns a polygon geometry to the entity.
func (e *Entity) WithPolygon(p *notageometry.Polygon) *Entity {
	e.manager.SetPolygon(e.ID, p)
	return e
}

// WithCollider assigns a collider to the entity.
func (e *Entity) WithCollider(c notacollision.Collider) *Entity {
	e.manager.SetCollider(e.ID, c)
	return e
}

// WithShader assigns a custom shader to the entity.
func (e *Entity) WithShader(s *notashader.Shader) *Entity {
	e.manager.SetShader(e.ID, s)
	return e
}

// WithMaterial assigns a material to the entity.
func (e *Entity) WithMaterial(m *notashader.Material) *Entity {
	e.manager.SetMaterial(e.ID, m)
	return e
}

// WithVisual assigns a sprite and optional material to the entity in one call.
func (e *Entity) WithVisual(v *Visual) *Entity {
	if v == nil {
		return e
	}
	if v.Sprite != nil {
		e.WithSprite(v.Sprite)
	}
	if v.Material != nil {
		e.WithMaterial(v.Material)
	}
	return e
}

// WithCollision assigns a collision profile to the entity in one call.
func (e *Entity) WithCollision(profile CollisionProfile) *Entity {
	if profile.Collider == nil {
		return e
	}
	return e.WithCollider(profile.Collider)
}

// WithCircle is a compatibility helper that assigns a circular collider directly.
func (e *Entity) WithCircle(radius float32) *Entity {
	return e.WithCollider(notacollision.NewCircleCollider(notamath.Po2{}, radius))
}

// WithCircleSprite is a compatibility helper that assigns a sprite, material, and inferred circle collider together.
func (e *Entity) WithCircleSprite(sprite *notatexture.Sprite, material *notashader.Material) *Entity {
	e.WithSprite(sprite)
	e.WithMaterial(material)

	if sprite != nil && sprite.Polygon != nil {
		e.WithCircle(circleRadiusFromPolygon(sprite.Polygon))
	}

	return e
}

// WithColor assigns a color to the entity.
func (e *Entity) WithColor(c notacolor.Color) *Entity {
	e.manager.SetColor(e.ID, c)
	return e
}

// Move moves an entity by a vector amount. Movement is additively applied.
func (e *Entity) Move(delta notamath.Vec2) {
	if !e.manager.IsActive(e.ID) {
		return
	}
	e.manager.SubmitMove(e.ID, delta)
}

// Rotate rotates an entity by an amount (radians). Rotation is additively applied.
func (e *Entity) Rotate(rad float32) {
	if !e.manager.IsActive(e.ID) {
		return
	}
	e.manager.SubmitRotation(e.ID, rad)
}

// Scale scales an entity by a factor. Scaling is multiplicatively applied.
func (e *Entity) Scale(factor notamath.Vec2) {
	if !e.manager.IsActive(e.ID) {
		return
	}
	e.manager.SubmitScale(e.ID, factor)
}

// Position gets the current position of the entity.
func (e *Entity) Position() notamath.Vec2 {
	return e.manager.getPosition(e.ID)
}

// Rotation gets the current rotation of the entity (radians).
func (e *Entity) Rotation() float32 {
	return e.manager.getRotation(e.ID)
}

// ScaleValue gets the current relative scale of the entity.
func (e *Entity) ScaleValue() notamath.Vec2 {
	return e.manager.getScale(e.ID)
}

// IsActive returns whether the entity is active.
func (e *Entity) IsActive() bool {
	return e.manager.IsActive(e.ID)
}

// SetActive sets the active state of the entity.
func (e *Entity) SetActive(active bool) {
	e.manager.SetActive(e.ID, active)
}

// IsVisible returns whether the entity is visible.
func (e *Entity) IsVisible() bool {
	return e.manager.IsVisible(e.ID)
}

// SetVisible sets the visibility state of the entity.
func (e *Entity) SetVisible(visible bool) {
	e.manager.SetVisible(e.ID, visible)
}

// Draw sends a draw request to the renderer.
func (e *Entity) Draw(renderer *notarender.Renderer, alpha float32) error {
	return e.DrawWithView(renderer, notamath.Mat3Identity(), alpha)
}

// DrawWithView sends a draw request to the renderer using the provided view matrix.
func (e *Entity) DrawWithView(renderer *notarender.Renderer, view notamath.Mat3, alpha float32) error {
	if !e.manager.IsVisible(e.ID) || !e.manager.IsActive(e.ID) {
		return nil
	}

	pos := e.manager.getPosition(e.ID)
	scale := e.manager.getScale(e.ID)
	rot := e.manager.getRotation(e.ID)

	model := view.Mul(notamath.Mat3TRS(pos, rot, scale))

	color := e.manager.GetColor(e.ID)

	sprite := e.manager.GetSprite(e.ID)
	if sprite != nil && sprite.Polygon != nil {
		renderer.SubmitPolygon(sprite.Polygon, model, color, sprite.Texture, e.manager.GetShader(e.ID), e.manager.GetMaterial(e.ID))
		return nil
	}

	poly := e.manager.GetPolygon(e.ID)
	if poly != nil {
		renderer.SubmitPolygon(poly, model, color, nil, e.manager.GetShader(e.ID), e.manager.GetMaterial(e.ID))
	}

	return nil
}

// GetId returns the entity's ID as a string (for backward compatibility).
func (e *Entity) GetId() string {
	return e.manager.GetEntityName(e.ID)
}

func circleRadiusFromPolygon(poly *notageometry.Polygon) float32 {
	if poly == nil || len(poly.Points) == 0 {
		return 0
	}

	minX, maxX := poly.Points[0].X, poly.Points[0].X
	minY, maxY := poly.Points[0].Y, poly.Points[0].Y
	for _, p := range poly.Points[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	width := maxX - minX
	height := maxY - minY
	if width < height {
		return width * 0.5
	}
	return height * 0.5
}
