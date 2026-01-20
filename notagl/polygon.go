package notagl

import (
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
)

type Polygon struct {
	Vertices  []notamath.Po2
	Transform notamath.Transform2D
	Color     notashader.Color   // Fallback / Single color
	Colors    []notashader.Color // Gradient colors (one per vertex)
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

	verts := make([]Vertex2D, len(p.Vertices))
	useGradient := len(p.Colors) == len(p.Vertices)

	for i, v := range p.Vertices {
		c := p.Color
		if useGradient {
			c = p.Colors[i]
		}

		verts[i] = Vertex2D{
			Pos:   mat.TransformPo2(v),
			Color: c,
		}
	}

	*orders = append(*orders, DrawOrder2D{
		Vertices: verts,
	})
}

func (p *Polygon) SetVerticalGradient(top, bottom notashader.Color) {
	p.Colors = make([]notashader.Color, len(p.Vertices))

	minY, maxY := p.Vertices[0].Y, p.Vertices[0].Y
	for _, v := range p.Vertices {
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	rangeY := maxY - minY
	for i, v := range p.Vertices {
		t := (v.Y - minY) / rangeY
		p.Colors[i] = bottom.Lerp(top, t)
	}
}

func (p *Polygon) SetColor(c notashader.Color) {
	p.Color = c
}

func (p *Polygon) SetHorizontalGradient(left, right notashader.Color) {
	p.Colors = make([]notashader.Color, len(p.Vertices))
	for i, v := range p.Vertices {
		p.Colors[i] = left.Lerp(right, v.X/p.Vertices[len(p.Vertices)-1].X)
	}
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

func CreateCircle(center notamath.Po2, radius float32) Polygon {
	size := radius * 2
	rect := CreateRectangle(center, size, size)
	return rect
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
		X: cx * inv,
		Y: cy * inv,
	}
}
