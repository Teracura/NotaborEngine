package notagl

import (
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
)

type Mesh struct {
	Vertices  []notamath.Po3
	Transform notamath.Transform3D
	Color     notashader.Color
	Colors    []notashader.Color
}

func (m *Mesh) Fixate() {
	center := meshCentroid(m.Vertices)

	local := make([]notamath.Po3, len(m.Vertices))
	for i, v := range m.Vertices {
		local[i] = notamath.Po3{
			X: v.X - center.X,
			Y: v.Y - center.Y,
			Z: v.Z - center.Z,
		}
	}

	m.Vertices = local
	m.Transform.Position = notamath.Vec3(center)
}
func (m *mesh) AddToOrders(orders *[]DrawOrder3D, alpha float32) {
	mat := m.Transform.InterpolatedMatrix(alpha, 0)

	verts := make([]Vertex3D, len(m.Vertices))
	useGradient := len(m.Colors) == len(m.Vertices)

	for i, v := range m.Vertices {
		c := m.Color
		if useGradient {
			c = m.Colors[i]
		}

		verts[i] = Vertex3D{
			Pos:   mat.TransformPo3(v),
			Color: c,
		}
	}

	*orders = append(*orders, DrawOrder3D{
		Vertices: verts,
	})
}
func (m *Mesh) SetHorizontalGradient(
	left, right notashader.Color,
) {
	m.Colors = make([]notashader.Color, len(m.Vertices))

	minX, maxX := m.Vertices[0].X, m.Vertices[0].X
	for _, v := range m.Vertices {
		if v.X < minX {
			minX = v.X
		}
		if v.X > maxX {
			maxX = v.X
		}
	}

	rx := maxX - minX
	if rx == 0 {
		return
	}

	for i, v := range m.Vertices {
		t := (v.X - minX) / rx
		m.Colors[i] = left.Lerp(right, t)
	}
}

func (m *Mesh) SetVerticalGradient(
	bottom, top notashader.Color,
) {
	if len(m.Colors) != len(m.Vertices) {
		m.Colors = make([]notashader.Color, len(m.Vertices))
	}

	minY, maxY := m.Vertices[0].Y, m.Vertices[0].Y
	for _, v := range m.Vertices {
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	ry := maxY - minY
	if ry == 0 {
		return
	}

	for i, v := range m.Vertices {
		t := (v.Y - minY) / ry
		m.Colors[i] = bottom.Lerp(top, t)
	}
}

func (m *Mesh) SetDepthGradient(
	near, far notashader.Color,
) {
	if len(m.Colors) != len(m.Vertices) {
		m.Colors = make([]notashader.Color, len(m.Vertices))
	}

	minZ, maxZ := m.Vertices[0].Z, m.Vertices[0].Z
	for _, v := range m.Vertices {
		if v.Z < minZ {
			minZ = v.Z
		}
		if v.Z > maxZ {
			maxZ = v.Z
		}
	}

	rz := maxZ - minZ
	if rz == 0 {
		return
	}

	for i, v := range m.Vertices {
		t := (v.Z - minZ) / rz
		m.Colors[i] = near.Lerp(far, t)
	}
}

func meshCentroid(verts []notamath.Po3) notamath.Po3 {
	var cx, cy, cz float32
	for _, v := range verts {
		cx += v.X
		cy += v.Y
		cz += v.Z
	}

	n := float32(len(verts))
	return notamath.Po3{
		X: cx / n,
		Y: cy / n,
		Z: cz / n,
	}
}
