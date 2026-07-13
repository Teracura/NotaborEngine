package main

import (
	"NotaborEngine/internal/notacollision"
	"NotaborEngine/internal/notamath"
)

// Collider wraps a collision shape (circle or polygon) for intersection testing.
type Collider struct {
	handle notacollision.Collider
}

// NewCircleCollider creates a circular collider with the given local-space center and radius.
func NewCircleCollider(center Po2, radius float32) *Collider {
	return &Collider{
		handle: notacollision.NewCircleCollider(notamath.Po2(center), radius),
	}
}

// NewPolygonCollider creates a polygon collider from a set of local-space points.
func NewPolygonCollider(points []Po2) *Collider {
	internal := make([]notamath.Po2, len(points))
	for i, p := range points {
		internal[i] = notamath.Po2(p)
	}
	return &Collider{
		handle: notacollision.NewPolygonCollider(internal),
	}
}

// Intersects checks whether two colliders overlap and returns the minimum
// translation vector (MTV) to separate them.
func Intersects(a, b *Collider) (bool, Vec2) {
	ok, mtv := notacollision.Intersects(a.handle, b.handle)
	return ok, Vec2(mtv)
}

// SetMaximumMTVTravelDistance sets the maximum distance the MTV can move
// colliders per frame. Default is 1.0. Adjust with care.
func SetMaximumMTVTravelDistance(amount float32) {
	notacollision.SetMaximumMTVTravelDistance(amount)
}
