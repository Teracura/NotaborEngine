package notagl

import (
	"NotaborEngine/notamath"
	"math"
)

type Polygon struct {
	Vertices  []notamath.Po2
	Transform notamath.Transform2D
}

// Fixate Adjusts points according to center point
func (p *Polygon) Fixate() {
	center := PolygonCentroid(p.Vertices)

	local := make([]notamath.Po2, len(p.Vertices))
	for i, v := range p.Vertices {
		local[i] = notamath.Po2{
			X: v.X - center.X,
			Y: v.Y - center.Y,
		}
	}

	p.Vertices = local
	p.Transform.Position = notamath.Vec2{
		X: center.X,
		Y: center.Y,
	}
}

func (p *Polygon) AddToOrders(orders *[]DrawOrder2D, alpha float32) {
	mat := p.Transform.InterpolatedMatrix(alpha)
	verts := make([]notamath.Po2, len(p.Vertices))
	for i, v := range p.Vertices {
		verts[i] = mat.TransformPo2(v)
	}

	*orders = append(*orders, DrawOrder2D{
		Vertices: verts,
	})
}

func CreateRectangle(center notamath.Po2, w, h float32) Polygon {
	hw := w / 2
	hh := h / 2
	rect := Polygon{Vertices: []notamath.Po2{
		{-hw, -hh},
		{+hw, -hh},
		{+hw, +hh},
		{-hw, +hh},
	},
		Transform: notamath.NewTransform2D()}
	rect.Transform.Position = notamath.Vec2{X: center.X, Y: center.Y}
	return rect
}

func CreateCircle(origin notamath.Po2, radius float32, segments int) Polygon {
	vertices := make([]notamath.Po2, segments)
	for i := 0; i < segments; i++ {
		vertices[i] = origin.Add(notamath.Vec2{
			X: float32(math.Cos(float64(i) * 2 * math.Pi / float64(segments))),
			Y: float32(math.Sin(float64(i) * 2 * math.Pi / float64(segments))),
		}.Mul(radius))
	}
	return Polygon{Vertices: vertices, Transform: notamath.NewTransform2D()}
}

func IsCCW(poly []notamath.Po2) bool {
	var area float32
	for i := 0; i < len(poly); i++ {
		a := poly[i]
		b := poly[(i+1)%len(poly)]
		area += (b.X - a.X) * (b.Y + a.Y)
	}
	return area < 0
}

func PointInTriangle(p, a, b, c notamath.Po2) bool {
	o1 := notamath.Orient(a, b, p)
	o2 := notamath.Orient(b, c, p)
	o3 := notamath.Orient(c, a, p)

	hasNeg := (o1 < 0) || (o2 < 0) || (o3 < 0)
	hasPos := (o1 > 0) || (o2 > 0) || (o3 > 0)

	return !(hasNeg && hasPos)
}

func IsEar(prev, curr, next notamath.Po2, poly []notamath.Po2) bool {
	// Must be convex (CCW polygon)
	if notamath.Orient(prev, curr, next) <= 0 {
		return false
	}

	for _, p := range poly {
		if p == prev || p == curr || p == next {
			continue
		}
		if PointInTriangle(p, prev, curr, next) {
			return false
		}
	}
	return true
}

func PolygonCentroid(poly []notamath.Po2) notamath.Po2 {
	var cx, cy, area float32

	n := len(poly)
	for i := 0; i < n; i++ {
		p0 := poly[i]
		p1 := poly[(i+1)%n]

		cross := p0.X*p1.Y - p1.X*p0.Y
		area += cross
		cx += (p0.X + p1.X) * cross
		cy += (p0.Y + p1.Y) * cross
	}

	area *= 0.5
	if area == 0 {
		return notamath.Po2{} // degenerate polygon
	}

	inv := 1.0 / (6.0 * area)
	return notamath.Po2{
		X: cx * float32(inv),
		Y: cy * float32(inv),
	}
}
