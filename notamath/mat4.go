package notamath

import "math"

type Mat4 struct {
	M [16]float32
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
		1, 0, 0, t.X,
		0, 1, 0, t.Y,
		0, 0, 1, t.Z,
		0, 0, 0, 1,
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

	return Mat4{M: [16]float32{
		c + k.X*k.X*t, k.X*k.Y*t - k.Z*s, k.X*k.Z*t + k.Y*s, 0,
		k.Y*k.X*t + k.Z*s, c + k.Y*k.Y*t, k.Y*k.Z*t - k.X*s, 0,
		k.Z*k.X*t - k.Y*s, k.Z*k.Y*t + k.X*s, c + k.Z*k.Z*t, 0,
		0, 0, 0, 1,
	}}
}

func (m Mat4) Mul(b Mat4) Mat4 {
	var r Mat4

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			r.M[row*4+col] =
				m.M[row*4+0]*b.M[0*4+col] +
					m.M[row*4+1]*b.M[1*4+col] +
					m.M[row*4+2]*b.M[2*4+col] +
					m.M[row*4+3]*b.M[3*4+col]
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

	if !bottomRow0001Check(m, 1e-5) || !bottomRow0001Check(b, 1e-5) {
		return m.Mul(b)
	}

	mTransOnly := m.isTranslationOnly()
	bTransOnly := b.isTranslationOnly()

	if mTransOnly && bTransOnly {
		return Mat4Translation(Vec3{
			X: m.M[3] + b.M[3],
			Y: m.M[7] + b.M[7],
			Z: m.M[11] + b.M[11],
		})
	}

	if mTransOnly {
		result := b
		result.M[3] += m.M[3]
		result.M[7] += m.M[7]
		result.M[11] += m.M[11]
		return result
	}

	if bTransOnly {
		result := m
		result.M[3] += b.M[3]
		result.M[7] += b.M[7]
		result.M[11] += b.M[11]
		return result
	}

	mScaleOnly := m.isScaleOnly()
	bScaleOnly := b.isScaleOnly()

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
			m.M[0] * sx, m.M[1] * sx, m.M[2] * sx, m.M[3],
			m.M[4] * sy, m.M[5] * sy, m.M[6] * sy, m.M[7],
			m.M[8] * sz, m.M[9] * sz, m.M[10] * sz, m.M[11],
			m.M[12], m.M[13], m.M[14], m.M[15],
		}}
	}

	return m.Mul(b)
}

func (m Mat4) isTranslationOnly() bool {
	const eps = 1e-5
	if !offDiagonalZeroCheck(m, eps) {
		return false
	}
	if !floatEqual(m.M[0], 1, eps) ||
		!floatEqual(m.M[5], 1, eps) ||
		!floatEqual(m.M[10], 1, eps) {
		return false
	}
	return true
}

func (m Mat4) isScaleOnly() bool {
	const eps = 1e-5

	if !offDiagonalZeroCheck(m, eps) {
		return false
	}

	// Translation must be zero (last column except M[15])
	if !floatEqual(m.M[3], 0, eps) ||
		!floatEqual(m.M[7], 0, eps) ||
		!floatEqual(m.M[11], 0, eps) {
		return false
	}
	return true
}

func bottomRow0001Check(m Mat4, eps float32) bool {
	if !floatEqual(m.M[12], 0, eps) ||
		!floatEqual(m.M[13], 0, eps) ||
		!floatEqual(m.M[14], 0, eps) ||
		!floatEqual(m.M[15], 1, eps) {
		return false
	}
	return true
}

func offDiagonalZeroCheck(m Mat4, eps float32) bool {
	if !floatEqual(m.M[1], 0, eps) ||
		!floatEqual(m.M[2], 0, eps) ||
		!floatEqual(m.M[4], 0, eps) ||
		!floatEqual(m.M[6], 0, eps) ||
		!floatEqual(m.M[8], 0, eps) ||
		!floatEqual(m.M[9], 0, eps) {
		return false
	}
	return true
}

func (m Mat4) TransformPo3(p Po3) Po3 {
	return Po3{
		X: m.M[0]*p.X + m.M[1]*p.Y + m.M[2]*p.Z + m.M[3],
		Y: m.M[4]*p.X + m.M[5]*p.Y + m.M[6]*p.Z + m.M[7],
		Z: m.M[8]*p.X + m.M[9]*p.Y + m.M[10]*p.Z + m.M[11],
	}
}

func (m Mat4) TransformVec3(v Vec3) Vec3 {
	return Vec3{
		X: m.M[0]*v.X + m.M[1]*v.Y + m.M[2]*v.Z,
		Y: m.M[4]*v.X + m.M[5]*v.Y + m.M[6]*v.Z,
		Z: m.M[8]*v.X + m.M[9]*v.Y + m.M[10]*v.Z,
	}
}

func Mat4TRS(pos Vec3, axis Vec3, angle float32, scale Vec3) Mat4 {
	t := Mat4Translation(pos)
	r := Mat4RotationAxisAngle(axis, angle)
	s := Mat4Scale(scale)
	return t.SmartMul(r).SmartMul(s)
}

func floatEqual(a, b, eps float32) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= eps
}

// Mat4Perspective creates a perspective projection matrix (shrinks or expands based on distance from object)
func Mat4Perspective(fovY, aspect, near, far float32) Mat4 {
	f := float32(1.0 / math.Tan(float64(fovY*0.5)))

	return Mat4{M: [16]float32{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / (near - far), (2 * far * near) / (near - far),
		0, 0, -1, 0,
	}}
}

// Mat4LookAt creates a view matrix that looks at the center point from the eye point
func Mat4LookAt(eye Vec3, center Vec3, up Vec3) Mat4 {
	f := center.Sub(eye).Normalize()
	s := f.Cross(up).Normalize()
	u := s.Cross(f)

	return Mat4{M: [16]float32{
		s.X, s.Y, s.Z, -s.Dot(eye),
		u.X, u.Y, u.Z, -u.Dot(eye),
		-f.X, -f.Y, -f.Z, f.Dot(eye),
		0, 0, 0, 1,
	}}
}

// Mat4Ortho creates an orthographic projection matrix, used for 2D rendering
func Mat4Ortho(left, right, bottom, top, near, far float32) Mat4 {
	return Mat4{M: [16]float32{
		2 / (right - left), 0, 0, -(right + left) / (right - left),
		0, 2 / (top - bottom), 0, -(top + bottom) / (top - bottom),
		0, 0, -2 / (far - near), -(far + near) / (far - near),
		0, 0, 0, 1,
	}}
}

// InverseAffine returns the inverse of the matrix, ignoring translation
func (m Mat4) InverseAffine() Mat4 {
	// Extract 3x3 linear part
	a00, a01, a02 := m.M[0], m.M[1], m.M[2]
	a10, a11, a12 := m.M[4], m.M[5], m.M[6]
	a20, a21, a22 := m.M[8], m.M[9], m.M[10]

	det := a00*(a11*a22-a12*a21) -
		a01*(a10*a22-a12*a20) +
		a02*(a10*a21-a11*a20)

	if det == 0 {
		return Mat4Identity()
	}

	invDet := 1 / det

	// Inverse 3x3
	r00 := (a11*a22 - a12*a21) * invDet
	r01 := (a02*a21 - a01*a22) * invDet
	r02 := (a01*a12 - a02*a11) * invDet

	r10 := (a12*a20 - a10*a22) * invDet
	r11 := (a00*a22 - a02*a20) * invDet
	r12 := (a02*a10 - a00*a12) * invDet

	r20 := (a10*a21 - a11*a20) * invDet
	r21 := (a01*a20 - a00*a21) * invDet
	r22 := (a00*a11 - a01*a10) * invDet

	// Correct row-major translation
	tx, ty, tz := m.M[3], m.M[7], m.M[11]

	itx := -(r00*tx + r01*ty + r02*tz)
	ity := -(r10*tx + r11*ty + r12*tz)
	itz := -(r20*tx + r21*ty + r22*tz)

	return Mat4{M: [16]float32{
		r00, r01, r02, itx,
		r10, r11, r12, ity,
		r20, r21, r22, itz,
		0, 0, 0, 1,
	}}
}

// NormalMatrix returns the inverse transpose of the upper-left 3x3 part of the matrix, ignoring translation
func (m Mat4) NormalMatrix() Mat3 {
	a00, a01, a02 := m.M[0], m.M[1], m.M[2]
	a10, a11, a12 := m.M[4], m.M[5], m.M[6]
	a20, a21, a22 := m.M[8], m.M[9], m.M[10]

	det := a00*(a11*a22-a12*a21) -
		a01*(a10*a22-a12*a20) +
		a02*(a10*a21-a11*a20)

	if det == 0 {
		return Mat3Identity()
	}

	invDet := 1 / det

	r00 := (a11*a22 - a12*a21) * invDet
	r01 := (a02*a21 - a01*a22) * invDet
	r02 := (a01*a12 - a02*a11) * invDet

	r10 := (a12*a20 - a10*a22) * invDet
	r11 := (a00*a22 - a02*a20) * invDet
	r12 := (a02*a10 - a00*a12) * invDet

	r20 := (a10*a21 - a11*a20) * invDet
	r21 := (a01*a20 - a00*a21) * invDet
	r22 := (a00*a11 - a01*a10) * invDet

	// transpose (inverse-transpose)
	return Mat3{M: [9]float32{
		r00, r10, r20,
		r01, r11, r21,
		r02, r12, r22,
	}}
}
