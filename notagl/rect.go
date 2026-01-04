package notagl

import "NotaborEngine/notamath"

type Rect struct {
	center notamath.Po2
	W, H   float32
}

func (r Rect) AddToOrders(orders *[]DrawOrder) {
	hw := r.W / 2
	hh := r.H / 2

	c := r.center
	*orders = append(*orders, DrawOrder{
		Vertices: []notamath.Po2{
			{c.X - hw, c.Y - hh},
			{c.X + hw, c.Y - hh},
			{c.X + hw, c.Y + hh},

			{c.X - hw, c.Y - hh},
			{c.X + hw, c.Y + hh},
			{c.X - hw, c.Y + hh},
		}})
}
