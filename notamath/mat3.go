package notamath

import (
	"fmt"
	"math"
)

type Mat3 struct {
	M [9]float32
}

type Transform2D struct {
	Position Vec2
	Rotation float32 // radians
	Scale    Vec2
	Dirty    bool //mutated without mapping to a matrix
	matrix   Mat3
}

func Mat3Identity() Mat3 {
	return Mat3{M: [9]float32{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}}
}

func Mat3Translation(t Vec2) Mat3 {
	return Mat3{M: [9]float32{
		1, 0, 0,
		0, 1, 0,
		t.X, t.Y, 1,
	}}
}

func Mat3Scale(s Vec2) Mat3 {
	return Mat3{M: [9]float32{
		s.X, 0, 0,
		0, s.Y, 0,
		0, 0, 1,
	}}
}

func Mat3Rotation(rad float32) Mat3 {
	c := float32(math.Cos(float64(rad)))
	s := float32(math.Sin(float64(rad)))

	return Mat3{M: [9]float32{
		c, s, 0,
		-s, c, 0,
		0, 0, 1,
	}}
}

func Mat3Shear(kx float32, ky float32) Mat3 {
	return Mat3{M: [9]float32{
		1, kx, 0,
		ky, 1, 0,
		0, 0, 1,
	}}
}

// Mat3TRS provides translation, rotation, and scale in one matrix
func Mat3TRS(pos Vec2, rot float32, scale Vec2) Mat3 {
	c := float32(math.Cos(float64(rot)))
	s := float32(math.Sin(float64(rot)))

	return Mat3{M: [9]float32{
		scale.X * c, scale.X * s, 0,
		-scale.Y * s, scale.Y * c, 0,
		pos.X, pos.Y, 1,
	}}
}

// NewTransform2D creates new transform parameters
func NewTransform2D() Transform2D {
	return Transform2D{
		Scale:  Vec2{1, 1},
		Dirty:  true,
		matrix: Mat3Identity(),
	}
}

func (t *Transform2D) SetPosition(p Vec2) {
	defer t.dirtify()
	t.Position = p
}

func (t *Transform2D) SetRotation(r float32) {
	defer t.dirtify()
	t.Rotation = r
}

func (t *Transform2D) SetScale(s Vec2) {
	defer t.dirtify()
	t.Scale = s
}

func (t *Transform2D) Matrix() Mat3 {
	if t.Dirty {
		t.matrix = Mat3TRS(t.Position, t.Rotation, t.Scale)
		t.Dirty = false
	}
	return t.matrix
}

// WorldMatrix returns the world matrix of this transform
func (t *Transform2D) WorldMatrix(parent *Transform2D) Mat3 {
	if parent == nil {
		return t.Matrix()
	}
	return parent.Matrix().Mul(t.Matrix())
}

func (t *Transform2D) dirtify() {
	t.Dirty = true
}

func (m Mat3) Mul(b Mat3) Mat3 {
	return Mat3{M: [9]float32{
		m.M[0]*b.M[0] + m.M[3]*b.M[1] + m.M[6]*b.M[2],
		m.M[1]*b.M[0] + m.M[4]*b.M[1] + m.M[7]*b.M[2],
		m.M[2]*b.M[0] + m.M[5]*b.M[1] + m.M[8]*b.M[2],

		m.M[0]*b.M[3] + m.M[3]*b.M[4] + m.M[6]*b.M[5],
		m.M[1]*b.M[3] + m.M[4]*b.M[4] + m.M[7]*b.M[5],
		m.M[2]*b.M[3] + m.M[5]*b.M[4] + m.M[8]*b.M[5],

		m.M[0]*b.M[6] + m.M[3]*b.M[7] + m.M[6]*b.M[8],
		m.M[1]*b.M[6] + m.M[4]*b.M[7] + m.M[7]*b.M[8],
		m.M[2]*b.M[6] + m.M[5]*b.M[7] + m.M[8]*b.M[8],
	}}
}

func (t *Transform2D) TransformPoint(p Po2) Po2 {
	return t.Matrix().TransformPo2(p)
}

func (t *Transform2D) TransformVector(v Vec2) Vec2 {
	return t.Matrix().TransformVec(v)
}

func (t *Transform2D) TranslateBy(delta Vec2) {
	t.Position = t.Position.Add(delta)
	t.dirtify()
}

func (t *Transform2D) RotateBy(delta float32) {
	t.Rotation += delta
	t.dirtify()
}

func (t *Transform2D) ScaleBy(factor Vec2) {
	t.Scale = Vec2{t.Scale.X * factor.X, t.Scale.Y * factor.Y}
	t.dirtify()
}
func (m Mat3) TransformPo2(p Po2) Po2 {
	return Po2{
		X: m.M[0]*p.X + m.M[3]*p.Y + m.M[6],
		Y: m.M[1]*p.X + m.M[4]*p.Y + m.M[7],
	}
}

func (m Mat3) TransformVec(v Vec2) Vec2 {
	return Vec2{
		X: m.M[0]*v.X + m.M[3]*v.Y,
		Y: m.M[1]*v.X + m.M[4]*v.Y,
	}
}

func (m Mat3) Transpose() Mat3 {
	return Mat3{M: [9]float32{
		m.M[0], m.M[3], m.M[6],
		m.M[1], m.M[4], m.M[7],
		m.M[2], m.M[5], m.M[8],
	}}
}

func (m Mat3) InverseAffine() Mat3 {
	a, b, c := m.M[0], m.M[1], m.M[6]
	d, e, f := m.M[3], m.M[4], m.M[7]

	// Determinant
	det := m.Det()
	if det == 0 {
		return Mat3Identity() // fallback
	}
	invDet := 1 / det

	return Mat3{M: [9]float32{
		e * invDet, -b * invDet, 0,
		-d * invDet, a * invDet, 0,
		(d*f - e*c) * invDet, (b*c - a*f) * invDet, 1,
	}}
}

func (m Mat3) Det() float32 {
	a, b := m.M[0], m.M[1]
	c, d := m.M[3], m.M[4]
	return a*d - b*c
}

func (m Mat3) String() string {
	return fmt.Sprintf("Matrix3x3(%f, %f, %f\n"+
		", %f, %f, %f\n"+
		", %f, %f, %f)", m.M[0], m.M[1], m.M[2], m.M[3], m.M[4], m.M[5], m.M[6], m.M[7], m.M[8])
}

// Mat3TRSS provides translation, rotation, scale, and shear in one matrix
func Mat3TRSS(pos Vec2, rot float32, scale Vec2, shearX, shearY float32) Mat3 {
	trs := Mat3TRS(pos, rot, scale)
	return Mat3Shear(shearX, shearY).Mul(trs)
}
