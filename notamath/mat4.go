package notamath

import "math"

type Mat4 struct {
	M [16]float32
}

type Transform3D struct {
	Position Vec3
	Rotation Vec3 // Reference vector for rotation
	Scale    Vec3
	Dirty    bool // true if matrix needs to be recomputed
	matrix   Mat4 // cached TRS matrix
}

func Mat4Identity() Mat4 {
	return Mat4{M: [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}}
}

func Mat4Translation(t Vec3) Mat4 {
	return Mat4{M: [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		t.X, t.Y, t.Z, 1,
	}}
}

func Mat4Scale(s Vec3) Mat4 {
	return Mat4{M: [16]float32{
		s.X, 0, 0, 0,
		0, s.Y, 0, 0,
		0, 0, s.Z, 0,
		0, 0, 0, 1,
	}}
}

func Mat4RotationAxisAngle(axis Vec3, angle float32) Mat4 {
	k := axis.Normalize()
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))
	t := 1 - c

	x, y, z := k.X, k.Y, k.Z

	return Mat4{M: [16]float32{
		t*x*x + c, t*x*y + s*z, t*x*z - s*y, 0,
		t*x*y - s*z, t*y*y + c, t*y*z + s*x, 0,
		t*x*z + s*y, t*y*z - s*x, t*z*z + c, 0,
		0, 0, 0, 1,
	}}
}

func (m Mat4) Mul(b Mat4) Mat4 {
	var r Mat4
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			r.M[col*4+row] =
				m.M[0*4+row]*b.M[col*4+0] +
					m.M[1*4+row]*b.M[col*4+1] +
					m.M[2*4+row]*b.M[col*4+2] +
					m.M[3*4+row]*b.M[col*4+3]
		}
	}
	return r
}

// SmartMul checks common special cases which take less computational time on average
func (m Mat4) SmartMul(b Mat4) Mat4 {
	if m == Mat4Identity() {
		return b
	}
	if b == Mat4Identity() {
		return m
	}

	mTransOnly := m.IsTranslationOnly()
	bTransOnly := b.IsTranslationOnly()

	if mTransOnly && bTransOnly {
		return Mat4Translation(Vec3{
			X: m.M[12] + b.M[12],
			Y: m.M[13] + b.M[13],
			Z: m.M[14] + b.M[14],
		})
	}
	if mTransOnly {
		result := b
		result.M[12] += m.M[12]
		result.M[13] += m.M[13]
		result.M[14] += m.M[14]
		return result
	}
	if bTransOnly {
		result := m
		result.M[12] += b.M[12]
		result.M[13] += b.M[13]
		result.M[14] += b.M[14]
		return result
	}

	mScaleOnly := m.IsScaleOnly()
	bScaleOnly := b.IsScaleOnly()

	if mScaleOnly && bScaleOnly {
		return Mat4Scale(Vec3{
			X: m.M[0] * b.M[0],
			Y: m.M[5] * b.M[5],
			Z: m.M[10] * b.M[10],
		})
	}
	if mScaleOnly {
		sx, sy, sz := m.M[0], m.M[5], m.M[10]
		return Mat4{M: [16]float32{
			sx * b.M[0], sx * b.M[1], sx * b.M[2], sx * b.M[3],
			sy * b.M[4], sy * b.M[5], sy * b.M[6], sy * b.M[7],
			sz * b.M[8], sz * b.M[9], sz * b.M[10], sz * b.M[11],
			b.M[12], b.M[13], b.M[14], b.M[15],
		}}
	}
	if bScaleOnly {
		sx, sy, sz := b.M[0], b.M[5], b.M[10]
		return Mat4{M: [16]float32{
			m.M[0] * sx, m.M[1] * sy, m.M[2] * sz, m.M[3],
			m.M[4] * sx, m.M[5] * sy, m.M[6] * sz, m.M[7],
			m.M[8] * sx, m.M[9] * sy, m.M[10] * sz, m.M[11],
			m.M[12] * sx, m.M[13] * sy, m.M[14] * sz, m.M[15],
		}}
	}
	return m.Mul(b)
}

func (m Mat4) TransformPo3(p Po3) Po3 {
	return Po3{
		X: m.M[0]*p.X + m.M[4]*p.Y + m.M[8]*p.Z + m.M[12],
		Y: m.M[1]*p.X + m.M[5]*p.Y + m.M[9]*p.Z + m.M[13],
		Z: m.M[2]*p.X + m.M[6]*p.Y + m.M[10]*p.Z + m.M[14],
	}
}

func (m Mat4) TransformVec3(v Vec3) Vec3 {
	return Vec3{
		X: m.M[0]*v.X + m.M[4]*v.Y + m.M[8]*v.Z,
		Y: m.M[1]*v.X + m.M[5]*v.Y + m.M[9]*v.Z,
		Z: m.M[2]*v.X + m.M[6]*v.Y + m.M[10]*v.Z,
	}
}

func Mat4TRS(pos Vec3, axis Vec3, angle float32, scale Vec3) Mat4 {
	t := Mat4Translation(pos)
	r := Mat4RotationAxisAngle(axis, angle)
	s := Mat4Scale(scale)
	return t.SmartMul(r).SmartMul(s)
}

// Mat4Perspective creates a perspective projection matrix (shrinks or expands based on distance from object)
func Mat4Perspective(fovY, aspect, near, far float32) Mat4 {
	f := float32(1.0 / math.Tan(float64(fovY*0.5)))

	return Mat4{M: [16]float32{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / (near - far), -1,
		0, 0, (2 * far * near) / (near - far), 0,
	}}
}

// Mat4LookAt creates a view matrix that looks at the center point from the eye point
func Mat4LookAt(eye Vec3, center Po3, up Vec3) Mat4 {
	f := center.SubVec(eye).Normalize()
	s := f.Cross(up).Normalize()
	u := s.Cross(f)

	return Mat4{M: [16]float32{
		s.X, u.X, -f.X, 0,
		s.Y, u.Y, -f.Y, 0,
		s.Z, u.Z, -f.Z, 0,
		-s.Dot(eye), -u.Dot(eye), f.Dot(eye), 1,
	}}
}

// Mat4Ortho creates an orthographic projection matrix, used for 2D rendering
func Mat4Ortho(left, right, bottom, top, near, far float32) Mat4 {
	return Mat4{M: [16]float32{
		2 / (right - left), 0, 0, 0,
		0, 2 / (top - bottom), 0, 0,
		0, 0, -2 / (far - near), 0,
		-(right + left) / (right - left),
		-(top + bottom) / (top - bottom),
		-(far + near) / (far - near),
		1,
	}}
}

// InverseAffine returns the inverse of the matrix, ignoring translation
func (m Mat4) InverseAffine() Mat4 {
	// extract linear 3x3 part
	a00, a01, a02 := m.M[0], m.M[1], m.M[2]
	a10, a11, a12 := m.M[4], m.M[5], m.M[6]
	a20, a21, a22 := m.M[8], m.M[9], m.M[10]

	// determinant
	det := a00*(a11*a22-a12*a21) - a01*(a10*a22-a12*a20) + a02*(a10*a21-a11*a20)
	if det == 0 {
		return Mat4Identity()
	}
	invDet := 1 / det

	// invert linear part
	r00 := (a11*a22 - a12*a21) * invDet
	r01 := (a02*a21 - a01*a22) * invDet
	r02 := (a01*a12 - a02*a11) * invDet

	r10 := (a12*a20 - a10*a22) * invDet
	r11 := (a00*a22 - a02*a20) * invDet
	r12 := (a02*a10 - a00*a12) * invDet

	r20 := (a10*a21 - a11*a20) * invDet
	r21 := (a01*a20 - a00*a21) * invDet
	r22 := (a00*a11 - a01*a10) * invDet

	// invert translation
	t := Vec3{m.M[12], m.M[13], m.M[14]}
	ti := Vec3{
		-(r00*t.X + r10*t.Y + r20*t.Z),
		-(r01*t.X + r11*t.Y + r21*t.Z),
		-(r02*t.X + r12*t.Y + r22*t.Z),
	}

	return Mat4{M: [16]float32{
		r00, r01, r02, 0,
		r10, r11, r12, 0,
		r20, r21, r22, 0,
		ti.X, ti.Y, ti.Z, 1,
	}}
}

// NormalMatrix returns the inverse transpose of the upper-left 3x3 part of the matrix, ignoring translation
func (m Mat4) NormalMatrix() Mat3 {
	// extract 3x3 linear part
	a00, a01, a02 := m.M[0], m.M[1], m.M[2]
	a10, a11, a12 := m.M[4], m.M[5], m.M[6]
	a20, a21, a22 := m.M[8], m.M[9], m.M[10]

	// determinant
	det := a00*(a11*a22-a12*a21) - a01*(a10*a22-a12*a20) + a02*(a10*a21-a11*a20)
	if det == 0 {
		return Mat3Identity()
	}
	invDet := 1 / det

	// invert
	r00 := (a11*a22 - a12*a21) * invDet
	r01 := (a02*a21 - a01*a22) * invDet
	r02 := (a01*a12 - a02*a11) * invDet

	r10 := (a12*a20 - a10*a22) * invDet
	r11 := (a00*a22 - a02*a20) * invDet
	r12 := (a02*a10 - a00*a12) * invDet

	r20 := (a10*a21 - a11*a20) * invDet
	r21 := (a01*a20 - a00*a21) * invDet
	r22 := (a00*a11 - a01*a10) * invDet

	// transpose for normal matrix
	return Mat3{M: [9]float32{
		r00, r10, r20,
		r01, r11, r21,
		r02, r12, r22,
	}}
}

func (m Mat4) IsTranslationOnly() bool {
	const eps = 1e-5

	if !floatEqual(m.M[0], 1, eps) || !floatEqual(m.M[5], 1, eps) || !floatEqual(m.M[10], 1, eps) {
		return false
	}
	if !floatEqual(m.M[1], 0, eps) || !floatEqual(m.M[2], 0, eps) ||
		!floatEqual(m.M[4], 0, eps) || !floatEqual(m.M[6], 0, eps) ||
		!floatEqual(m.M[8], 0, eps) || !floatEqual(m.M[9], 0, eps) {
		return false
	}

	if !floatEqual(m.M[15], 1, eps) {
		return false
	}

	return true
}

func (m Mat4) IsScaleOnly() bool {
	const eps = 1e-5

	if !floatEqual(m.M[1], 0, eps) || !floatEqual(m.M[2], 0, eps) ||
		!floatEqual(m.M[4], 0, eps) || !floatEqual(m.M[6], 0, eps) ||
		!floatEqual(m.M[8], 0, eps) || !floatEqual(m.M[9], 0, eps) {
		return false
	}

	if !floatEqual(m.M[12], 0, eps) || !floatEqual(m.M[13], 0, eps) || !floatEqual(m.M[14], 0, eps) {
		return false
	}

	if !floatEqual(m.M[15], 1, eps) {
		return false
	}

	return true
}

func (m Mat4) IsRotationOnly() bool {
	const eps = 1e-3

	if !floatEqual(m.M[12], 0, eps) || !floatEqual(m.M[13], 0, eps) || !floatEqual(m.M[14], 0, eps) {
		return false
	}

	if !floatEqual(m.M[15], 1, eps) {
		return false
	}

	if !floatEqual(m.M[0]*m.M[0]+m.M[1]*m.M[1]+m.M[2]*m.M[2], 1, eps) ||
		!floatEqual(m.M[4]*m.M[4]+m.M[5]*m.M[5]+m.M[6]*m.M[6], 1, eps) ||
		!floatEqual(m.M[8]*m.M[8]+m.M[9]*m.M[9]+m.M[10]*m.M[10], 1, eps) {
		return false
	}

	if !floatEqual(m.M[0]*m.M[4]+m.M[1]*m.M[5]+m.M[2]*m.M[6], 0, eps) ||
		!floatEqual(m.M[0]*m.M[8]+m.M[1]*m.M[9]+m.M[2]*m.M[10], 0, eps) ||
		!floatEqual(m.M[4]*m.M[8]+m.M[5]*m.M[9]+m.M[6]*m.M[10], 0, eps) {
		return false
	}

	return true
}

func floatEqual(a, b, eps float32) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= eps
}

func Mat4RotationXYZ(rx, ry, rz float32) Mat4 {
	cx, sx := float32(math.Cos(float64(rx))), float32(math.Sin(float64(rx)))
	cy, sy := float32(math.Cos(float64(ry))), float32(math.Sin(float64(ry)))
	cz, sz := float32(math.Cos(float64(rz))), float32(math.Sin(float64(rz)))

	return Mat4{M: [16]float32{
		cy * cz, cx*sz + sx*sy*cz, sx*sz - cx*sy*cz, 0,
		-cy * sz, cx*cz - sx*sy*sz, sx*cz + cx*sy*sz, 0,
		sy, -sx * cy, cx * cy, 0,
		0, 0, 0, 1,
	}}
}

func NewTransform3D() Transform3D {
	return Transform3D{
		Scale:  Vec3{1, 1, 1},
		Dirty:  true,
		matrix: Mat4Identity(),
	}
}

func (t *Transform3D) SetPosition(p Vec3) {
	t.Position = p
	t.Dirty = true
}

func (t *Transform3D) SetRotation(r Vec3) {
	t.Rotation = r
	t.Dirty = true
}

func (t *Transform3D) SetScale(s Vec3) {
	t.Scale = s
	t.Dirty = true
}

func (t *Transform3D) Matrix() Mat4 {
	if !t.Dirty {
		return t.matrix
	}

	tr := Mat4Translation(t.Position)

	rot := Mat4RotationXYZ(t.Rotation.X, t.Rotation.Y, t.Rotation.Z)

	sc := Mat4Scale(t.Scale)

	// Combine TRS
	t.matrix = tr.SmartMul(rot).SmartMul(sc)
	t.Dirty = false
	return t.matrix
}

func (t *Transform3D) TransformPo3(p Po3) Po3 {
	return t.Matrix().TransformPo3(p)
}

func (t *Transform3D) TransformVec3(v Vec3) Vec3 {
	return t.Matrix().TransformVec3(v)
}

func (t *Transform3D) TranslateBy(delta Vec3) {
	t.Position = t.Position.Add(delta)
	t.Dirty = true
}

func (t *Transform3D) RotateBy(delta Vec3) {
	t.Rotation = t.Rotation.Add(delta)
	t.Dirty = true
}

func (t *Transform3D) ScaleBy(factor Vec3) {
	t.Scale = Vec3{t.Scale.X * factor.X, t.Scale.Y * factor.Y, t.Scale.Z * factor.Z}
	t.Dirty = true
}

func (t *Transform3D) WorldMatrix(parent *Transform3D) Mat4 {
	if parent == nil {
		return t.Matrix()
	}
	return parent.Matrix().SmartMul(t.Matrix())
}
