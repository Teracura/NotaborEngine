package main

import (
	"NotaborEngine/internal/notacollision"
	"NotaborEngine/internal/notaentity"
	"NotaborEngine/internal/notageometry"
	"NotaborEngine/internal/notamath"
	"NotaborEngine/internal/notashader"
	"NotaborEngine/internal/notatexture"
)

// EntityManager manages game entities (create, remove, query) using
// an internal ECS architecture with sparse-set storage.
type EntityManager struct {
	handle *notaentity.EntityManager
}

// CreateEntity creates a new named entity and returns a handle to it.
func (em *EntityManager) CreateEntity(name string) *Entity {
	return &Entity{handle: em.handle.CreateEntity(name)}
}

// Flush applies all pending transform updates and syncs colliders.
// Call once per frame after moving/rotating/scaling entities.
func (em *EntityManager) Flush() {
	em.handle.Flush()
}

// GetEntity retrieves an entity by its name. Returns nil if not found or inactive.
func (em *EntityManager) GetEntity(name string) *Entity {
	e := em.handle.GetEntity(name)
	if e == nil {
		return nil
	}
	return &Entity{handle: e}
}

// GetEntities returns all active entities.
func (em *EntityManager) GetEntities() []*Entity {
	internal := em.handle.GetEntities()
	result := make([]*Entity, len(internal))
	for i, e := range internal {
		result[i] = &Entity{handle: e}
	}
	return result
}

// Remove deletes an entity by name and frees its ID for reuse.
func (em *EntityManager) Remove(name string) {
	em.handle.Remove(name)
}

// AddToCollisionGroup adds an entity to a named collision group for group-based
// collision solving.
func (em *EntityManager) AddToCollisionGroup(group string, e *Entity) {
	em.handle.AddToCollisionGroup(group, e.handle)
}

// SolveGroupCollision detects and records collisions between all entity pairs
// within the named collision group.
func (em *EntityManager) SolveGroupCollision(name string) {
	em.handle.SolveGroupCollision(name)
}

// Collides checks whether two entities are currently colliding (from the last
// SolveGroupCollision call).
func (em *EntityManager) Collides(a, b *Entity) bool {
	return em.handle.Collides(a.handle, b.handle)
}

// GetMTV returns the minimum translation vector to separate two colliding entities.
func (em *EntityManager) GetMTV(a, b *Entity) Vec2 {
	return Vec2(em.handle.GetMTV(a.handle, b.handle))
}

// Entity is a lightweight handle to a game object. Use the manager's
// fluent builder methods (WithSprite, WithCollider, etc.) to attach components.
type Entity struct {
	handle *notaentity.Entity
}

// WithVisual attaches a sprite (and optional material) to the entity.
func (e *Entity) WithVisual(v *Visual) *Entity {
	e.handle.WithVisual(v.handle)
	return e
}

// WithCollision attaches a collision profile to the entity.
func (e *Entity) WithCollision(p CollisionProfile) *Entity {
	e.handle.WithCollision(notaentity.CollisionProfile{Collider: p.handle})
	return e
}

// WithColor sets the colour tint of the entity.
func (e *Entity) WithColor(c Color) *Entity {
	e.handle.WithColor(wrapColor(c))
	return e
}

// Move translates the entity by a world-space delta (additive).
func (e *Entity) Move(delta Vec2) {
	e.handle.Move(notamath.Vec2(delta))
}

// Rotate rotates the entity by the given angle in radians (additive).
func (e *Entity) Rotate(rad float32) {
	e.handle.Rotate(rad)
}

// Scale scales the entity by a factor (multiplicative).
func (e *Entity) Scale(factor Vec2) {
	e.handle.Scale(notamath.Vec2(factor))
}

// Position returns the entity's current world-space position.
func (e *Entity) Position() Vec2 {
	return Vec2(e.handle.Position())
}

// Rotation returns the entity's current rotation in radians.
func (e *Entity) Rotation() float32 {
	return e.handle.Rotation()
}

// ScaleValue returns the entity's current scale factor.
func (e *Entity) ScaleValue() Vec2 {
	return Vec2(e.handle.ScaleValue())
}

// IsActive returns whether the entity is active (processed each frame).
func (e *Entity) IsActive() bool {
	return e.handle.IsActive()
}

// SetActive sets the active state of the entity.
func (e *Entity) SetActive(active bool) {
	e.handle.SetActive(active)
}

// IsVisible returns whether the entity is visible (rendered each frame).
func (e *Entity) IsVisible() bool {
	return e.handle.IsVisible()
}

// SetVisible sets the visibility state of the entity.
func (e *Entity) SetVisible(visible bool) {
	e.handle.SetVisible(visible)
}

// WithSprite attaches a sprite to the entity.
func (e *Entity) WithSprite(s *notatexture.Sprite) *Entity {
	e.handle.WithSprite(s)
	return e
}

// WithPolygon attaches a polygon geometry to the entity.
func (e *Entity) WithPolygon(p *notageometry.Polygon) *Entity {
	e.handle.WithPolygon(p)
	return e
}

// WithCollider attaches a collider for collision detection.
func (e *Entity) WithCollider(c *Collider) *Entity {
	e.handle.WithCollider(c.handle)
	return e
}

// WithShader attaches a custom shader to the entity.
func (e *Entity) WithShader(s *notashader.Shader) *Entity {
	e.handle.WithShader(s)
	return e
}

// WithMaterial attaches a material (shader instance) to the entity.
func (e *Entity) WithMaterial(m *notashader.Material) *Entity {
	e.handle.WithMaterial(m)
	return e
}

// WithCircle is a convenience method that attaches a circular collider with the given radius.
func (e *Entity) WithCircle(radius float32) *Entity {
	e.handle.WithCircle(radius)
	return e
}

// Visual bundles a Sprite and optional Material for rendering.
type Visual struct {
	handle *notaentity.Visual
}

// CollisionProfile describes how an entity participates in collision queries.
type CollisionProfile struct {
	handle notacollision.Collider
}

// CircleCollision creates a circular collision profile centred on the entity origin.
func CircleCollision(radius float32) CollisionProfile {
	p := notaentity.CircleCollision(radius)
	return CollisionProfile{handle: p.Collider}
}

// PolygonCollision creates a polygon collision profile from local-space points.
func PolygonCollision(points []Po2) CollisionProfile {
	internal := make([]notamath.Po2, len(points))
	for i, p := range points {
		internal[i] = notamath.Po2(p)
	}
	p := notaentity.PolygonCollision(internal)
	return CollisionProfile{handle: p.Collider}
}

// CustomCollision wraps an already-constructed Collider as a collision profile.
func CustomCollision(c *Collider) CollisionProfile {
	return CollisionProfile{handle: c.handle}
}
