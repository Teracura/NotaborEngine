package notagl

import "NotaborEngine/notamath"

type Rect struct {
	Center    notamath.Po2
	W, H      float32
	Transform notamath.Transform2D
}

func (r Rect) AddToOrders(orders *[]DrawOrder2D) {
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

	for i := range vertices {
		vertices[i] = r.Transform.TransformPoint(vertices[i])
	}

	*orders = append(*orders, DrawOrder2D{Vertices: vertices})
}
