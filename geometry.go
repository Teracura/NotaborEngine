package main

import (
	"NotaborEngine/internal/notageometry"
	"NotaborEngine/internal/notamath"
)

// Polygon is pure geometry data (a set of 2D points forming a closed shape).
type Polygon struct {
	handle *notageometry.Polygon
}

// NewPolygon creates a polygon from the given points.
func NewPolygon(points []Po2) *Polygon {
	internal := make([]notamath.Po2, len(points))
	for i, p := range points {
		internal[i] = notamath.Po2(p)
	}
	return &Polygon{handle: &notageometry.Polygon{Points: internal}}
}

// CreateRectangle returns a rectangular polygon centred at origin with the given width and height.
func CreateRectangle(w, h float32) *Polygon {
	return &Polygon{handle: notageometry.CreateRectangle(w, h)}
}

// Fixate re-centres all polygon points so the centroid is at the origin.
func (p *Polygon) Fixate() {
	p.handle.Fixate()
}

// PointInTriangle checks if point p lies inside the triangle (a, b, c).
func PointInTriangle(p, a, b, c Po2) bool {
	return notageometry.PointInTriangle(
		notamath.Po2(p), notamath.Po2(a),
		notamath.Po2(b), notamath.Po2(c),
	)
}

// PolygonCentroid returns the centroid of the given points.
func PolygonCentroid(points []Po2) Po2 {
	internal := make([]notamath.Po2, len(points))
	for i, p := range points {
		internal[i] = notamath.Po2(p)
	}
	return Po2(notageometry.PolygonCentroid(internal))
}
