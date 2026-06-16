package main

import (
	"NotaborEngine/internal/notaentity"
)

type EntityManager struct {
	handle *notaentity.EntityManager
}

func (em *EntityManager) CreateEntity(name string) *Entity {
	return &Entity{handle: em.handle.CreateEntity(name)}
}

func (em *EntityManager) Flush() {
	em.handle.Flush()
}

type Entity struct {
	handle *notaentity.Entity
}

func (e *Entity) WithVisual(v *Visual) *Entity {
	e.handle.WithVisual(v.handle)
	return e
}

func (e *Entity) WithCollision(p CollisionProfile) *Entity {
	e.handle.WithCollision(notaentity.CollisionProfile(p))
	return e
}

func (e *Entity) WithColor(c Color) *Entity {
	e.handle.WithColor(wrapColor(c))
	return e
}

func (e *Entity) Move(delta Vec2) {
	e.handle.Move(bridgeVec2(delta))
}

type Visual struct {
	handle *notaentity.Visual
}

type CollisionProfile notaentity.CollisionProfile

func CircleCollision(radius float32) CollisionProfile {
	return CollisionProfile(notaentity.CircleCollision(radius))
}
