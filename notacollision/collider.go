package notacollision

import (
	"NotaborEngine/notamath"
	"math"
)

type AABBCollider struct {
	Min notamath.Vec2
	Max notamath.Vec2
}

type Collider interface {
	AABB() AABBCollider
	Move(delta notamath.Vec2)
	Rotate(delta float32)
}

type CircleCollider struct {
	Center notamath.Po2
	Radius float32
}

type PolygonCollider struct {
	Vertices []notamath.Po2
}

func (c *CircleCollider) AABB() AABBCollider {
	return AABBCollider{
		Min: notamath.Vec2{
			X: c.Center.X - c.Radius,
			Y: c.Center.Y - c.Radius,
		},
		Max: notamath.Vec2{
			X: c.Center.X + c.Radius,
			Y: c.Center.Y + c.Radius,
		},
	}
}

func (c *CircleCollider) Move(delta notamath.Vec2) {
	c.Center = c.Center.Add(delta)
}

func (c *CircleCollider) Rotate(delta float32) {
	return //rotating doesn't do anything for circles
}

func (p *PolygonCollider) AABB() AABBCollider {
	if len(p.Vertices) == 0 {
		return AABBCollider{}
	}

	minX := p.Vertices[0].X
	minY := p.Vertices[0].Y
	maxX := p.Vertices[0].X
	maxY := p.Vertices[0].Y

	for i := 1; i < len(p.Vertices); i++ {
		v := p.Vertices[i]

		if v.X < minX {
			minX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	return AABBCollider{
		Min: notamath.Vec2{X: minX, Y: minY},
		Max: notamath.Vec2{X: maxX, Y: maxY},
	}
}

func (p *PolygonCollider) Move(delta notamath.Vec2) {
	for _, vert := range p.Vertices {
		vert = vert.Add(delta)
	}
}

func (p *PolygonCollider) Rotate(delta float32) {
	for _, vert := range p.Vertices {
		vert.X += vert.X * float32(math.Cos(float64(delta)))
		vert.Y += vert.Y * float32(math.Sin(float64(delta)))
	}
}

func BroadPhase(a, b Collider) bool {
	return AABBIntersects(a.AABB(), b.AABB())
}

func AABBIntersects(a, b AABBCollider) bool {
	return a.Min.X <= b.Max.X &&
		a.Max.X >= b.Min.X &&
		a.Min.Y <= b.Max.Y &&
		a.Max.Y >= b.Min.Y
}

func Intersects(a, b Collider) bool {
	if !BroadPhase(a, b) {
		return false
	}

	switch a := a.(type) {
	case *CircleCollider:
		switch b := b.(type) {
		case *CircleCollider:
			return circleVsCircle(a, b)
		case *PolygonCollider:
			return circleVsPolygon(a, b)
		}
	case *PolygonCollider:
		switch b := b.(type) {
		case *CircleCollider:
			return circleVsPolygon(b, a)
		case *PolygonCollider:
			return polygonVsPolygon(a, b)
		}
	}

	return false
}

func circleVsCircle(a, b *CircleCollider) bool {
	dx := a.Center.X - b.Center.X
	dy := a.Center.Y - b.Center.Y
	r := a.Radius + b.Radius

	return dx*dx+dy*dy <= r*r
}

func polygonVsPolygon(a, b *PolygonCollider) bool {
	nA := len(a.Vertices)
	nB := len(b.Vertices)

	// 1. Edge vs edge
	for i := 0; i < nA; i++ {
		a1 := a.Vertices[i]
		a2 := a.Vertices[(i+1)%nA]

		for j := 0; j < nB; j++ {
			b1 := b.Vertices[j]
			b2 := b.Vertices[(j+1)%nB]

			if segmentsIntersect(a1, a2, b1, b2) {
				return true
			}
		}
	}

	if pointInPolygon(a.Vertices[0], b.Vertices) {
		return true
	}

	if pointInPolygon(b.Vertices[0], a.Vertices) {
		return true
	}

	return false
}

func circleVsPolygon(c *CircleCollider, p *PolygonCollider) bool {
	center := notamath.Po2{X: c.Center.X, Y: c.Center.Y}
	r2 := c.Radius * c.Radius
	n := len(p.Vertices)

	for i := 0; i < n; i++ {
		a := p.Vertices[i]
		b := p.Vertices[(i+1)%n]

		closest := closestPointOnSegment(a, b, center)
		if center.DistanceSquared(closest) <= r2 {
			return true
		}
	}

	if pointInPolygon(center, p.Vertices) {
		return true
	}

	return false
}

//HELPERS

const epsilon float32 = 1e-6

func segmentsIntersect(p1, p2, q1, q2 notamath.Po2) bool {
	o1 := notamath.Orient(p1, p2, q1)
	o2 := notamath.Orient(p1, p2, q2)
	o3 := notamath.Orient(q1, q2, p1)
	o4 := notamath.Orient(q1, q2, p2)

	// Proper intersection
	if o1*o2 < 0 && o3*o4 < 0 {
		return true
	}

	// Collinear cases
	if almostZero(o1) && onSegment(p1, p2, q1) {
		return true
	}
	if almostZero(o2) && onSegment(p1, p2, q2) {
		return true
	}
	if almostZero(o3) && onSegment(q1, q2, p1) {
		return true
	}
	if almostZero(o4) && onSegment(q1, q2, p2) {
		return true
	}

	return false
}

func almostZero(v float32) bool {
	if v < 0 {
		return -v < epsilon
	}
	return v < epsilon
}

func onSegment(a, b, p notamath.Po2) bool {
	return p.X >= min(a.X, b.X)-epsilon &&
		p.X <= max(a.X, b.X)+epsilon &&
		p.Y >= min(a.Y, b.Y)-epsilon &&
		p.Y <= max(a.Y, b.Y)+epsilon
}

func pointInPolygon(point notamath.Po2, poly []notamath.Po2) bool {
	inside := false
	n := len(poly)

	for i := 0; i < n; i++ {
		j := (i + n - 1) % n

		pi := poly[i]
		pj := poly[j]

		intersect := ((pi.Y > point.Y) != (pj.Y > point.Y)) &&
			(point.X < (pj.X-pi.X)*(point.Y-pi.Y)/(pj.Y-pi.Y)+pi.X)

		if intersect {
			inside = !inside
		}
	}

	return inside
}

func closestPointOnSegment(a, b notamath.Po2, p notamath.Po2) notamath.Po2 {
	ab := b.Sub(a)
	ap := p.Sub(a)

	t := ap.Dot(ab) / ab.LenSquared()

	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	return notamath.Po2{
		X: a.X + ab.X*t,
		Y: a.Y + ab.Y*t,
	}
}
