package notagl

import "NotaborEngine/notamath"

type Polygon struct {
	Vertices []notamath.Po2
}

func (p Polygon) AddToOrders(orders *[]DrawOrder2D) {
	*orders = append(*orders, DrawOrder2D{
		Vertices: p.Vertices})
}
