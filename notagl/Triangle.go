package notagl

import "NotaborEngine/notamath"

type Triangle struct {
	vertices [3]notamath.Po2
}

func (t Triangle) AddToOrders(orders *[]DrawOrder) {
	*orders = append(*orders, DrawOrder{
		Vertices: []notamath.Po2{
			{t.vertices[0].X, t.vertices[0].Y},
			{t.vertices[1].X, t.vertices[1].Y},
			{t.vertices[2].X, t.vertices[2].Y},
		}})
}
