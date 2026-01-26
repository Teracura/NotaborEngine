package notacollision

import "NotaborEngine/notamath"

type ColliderType int

const (
	AABB ColliderType = iota
	Circle
)

type AABBCollider struct {
	Min notamath.Vec2
	Max notamath.Vec2
}

type Collider interface {
	Type() ColliderType
	AABB() AABBCollider
}

func Intersects(a, b Collider) bool {
	return false
}

func IsInside(a, b Collider) bool {
	return false
}
