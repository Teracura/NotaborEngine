package notagl

import (
	"NotaborEngine/notamath"
)

type Rect struct {
	Center    notamath.Po2
	W, H      float32
	Transform notamath.Transform2D
}

func (r *Rect) AddToOrders(orders *[]DrawOrder2D, alpha float32) {
	hw := r.W / 2
	hh := r.H / 2

	c := r.Center
	vertices := []notamath.Po2{
		{c.X - hw, c.Y - hh},
		{c.X + hw, c.Y - hh},
		{c.X + hw, c.Y + hh},

		{c.X - hw, c.Y - hh},
		{c.X + hw, c.Y + hh},
		{c.X - hw, c.Y + hh},
	}

	mat := r.Transform.InterpolatedMatrix(alpha)

	for i := range vertices {
		vertices[i] = mat.TransformPo2(vertices[i])
	}

	*orders = append(*orders, DrawOrder2D{Vertices: vertices})
}
