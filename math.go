package main

import (
	"fmt"

	"NotaborEngine/internal/notamath"
)

// Vec2 is a 2-dimensional vector with X and Y float32 components.
type Vec2 notamath.Vec2

// NewVec2 returns a Vec2 with the given components.
func NewVec2(x, y float32) Vec2 {
	return Vec2(notamath.Vec2{X: x, Y: y})
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2(notamath.Vec2(v).Add(notamath.Vec2(o)))
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2(notamath.Vec2(v).Sub(notamath.Vec2(o)))
}

func (v Vec2) Mul(s float32) Vec2 {
	return Vec2(notamath.Vec2(v).Mul(s))
}

func (v Vec2) Div(s float32) Vec2 {
	return Vec2(notamath.Vec2(v).Div(s))
}

func (v Vec2) Dot(o Vec2) float32 {
	return notamath.Vec2(v).Dot(notamath.Vec2(o))
}

func (v Vec2) Cross(o Vec2) float32 {
	return notamath.Vec2(v).Cross(notamath.Vec2(o))
}

func (v Vec2) LenSquared() float32 {
	return notamath.Vec2(v).LenSquared()
}

func (v Vec2) Len() float32 {
	return notamath.Vec2(v).Len()
}

func (v Vec2) Distance(o Vec2) float32 {
	return notamath.Vec2(v).Distance(notamath.Vec2(o))
}

func (v Vec2) Neg() Vec2 {
	return Vec2(notamath.Vec2(v).Neg())
}

func (v Vec2) Normalize() Vec2 {
	return Vec2(notamath.Vec2(v).Normalize())
}

func (v Vec2) Perp() Vec2 {
	return Vec2(notamath.Vec2(v).Perp())
}

func (v Vec2) Lerp(a Vec2, t float32) Vec2 {
	return Vec2(notamath.Vec2(v).Lerp(notamath.Vec2(a), t))
}

func (v Vec2) Project(onto Vec2) Vec2 {
	return Vec2(notamath.Vec2(v).Project(notamath.Vec2(onto)))
}

func (v Vec2) Angle(o Vec2) float32 {
	return notamath.Vec2(v).Angle(notamath.Vec2(o))
}

func (v Vec2) Rotate(rad float32) Vec2 {
	return Vec2(notamath.Vec2(v).Rotate(rad))
}

func (v Vec2) String() string {
	return fmt.Sprintf("Vec2(%f, %f)", v.X, v.Y)
}

// Vec3 is a 3-dimensional vector with X, Y, Z float32 components.
type Vec3 notamath.Vec3

// NewVec3 returns a Vec3 with the given components.
func NewVec3(x, y, z float32) Vec3 {
	return Vec3(notamath.Vec3{X: x, Y: y, Z: z})
}

func (v Vec3) Add(o Vec3) Vec3 {
	return Vec3(notamath.Vec3(v).Add(notamath.Vec3(o)))
}

func (v Vec3) Sub(o Vec3) Vec3 {
	return Vec3(notamath.Vec3(v).Sub(notamath.Vec3(o)))
}

func (v Vec3) Mul(s float32) Vec3 {
	return Vec3(notamath.Vec3(v).Mul(s))
}

func (v Vec3) Div(s float32) Vec3 {
	return Vec3(notamath.Vec3(v).Div(s))
}

func (v Vec3) Neg() Vec3 {
	return Vec3(notamath.Vec3(v).Neg())
}

func (v Vec3) Dot(o Vec3) float32 {
	return notamath.Vec3(v).Dot(notamath.Vec3(o))
}

func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3(notamath.Vec3(v).Cross(notamath.Vec3(o)))
}

func (v Vec3) LenSquared() float32 {
	return notamath.Vec3(v).LenSquared()
}

func (v Vec3) Len() float32 {
	return notamath.Vec3(v).Len()
}

func (v Vec3) Normalize() Vec3 {
	return Vec3(notamath.Vec3(v).Normalize())
}

func (v Vec3) Distance(o Vec3) float32 {
	return notamath.Vec3(v).Distance(notamath.Vec3(o))
}

func (v Vec3) Lerp(to Vec3, t float32) Vec3 {
	return Vec3(notamath.Vec3(v).Lerp(notamath.Vec3(to), t))
}

func (v Vec3) Project(onto Vec3) Vec3 {
	return Vec3(notamath.Vec3(v).Project(notamath.Vec3(onto)))
}

func (v Vec3) Angle(o Vec3) float32 {
	return notamath.Vec3(v).Angle(notamath.Vec3(o))
}

func (v Vec3) Rotate(axis Vec3, angle float32) Vec3 {
	return Vec3(notamath.Vec3(v).Rotate(notamath.Vec3(axis), angle))
}

func (v Vec3) String() string {
	return fmt.Sprintf("Vec3(%f, %f, %f)", v.X, v.Y, v.Z)
}

// Po2 is a 2-dimensional point with X, Y float32 components.
type Po2 notamath.Po2

// NewPo2 returns a Po2 with the given coordinates.
func NewPo2(x, y float32) Po2 {
	return Po2(notamath.Po2{X: x, Y: y})
}

func (p Po2) Add(v Vec2) Po2 {
	return Po2(notamath.Po2(p).Add(notamath.Vec2(v)))
}

func (p Po2) Sub(q Po2) Vec2 {
	return Vec2(notamath.Po2(p).Sub(notamath.Po2(q)))
}

func (p Po2) DistanceSquared(q Po2) float32 {
	return notamath.Po2(p).DistanceSquared(notamath.Po2(q))
}

func (p Po2) Distance(q Po2) float32 {
	return notamath.Po2(p).Distance(notamath.Po2(q))
}

func (p Po2) Equals(q Po2, eps float32) bool {
	return notamath.Po2(p).Equals(notamath.Po2(q), eps)
}

func (p Po2) ToVec2() Vec2 {
	return Vec2(notamath.Po2(p).ToVec2())
}

func (p Po2) String() string {
	return fmt.Sprintf("Po2(%f, %f)", p.X, p.Y)
}

// Orient returns the signed area of the triangle (a, b, c).
// Positive means counter-clockwise.
func Orient(a, b, c Po2) float32 {
	return notamath.Orient(notamath.Po2(a), notamath.Po2(b), notamath.Po2(c))
}

// Po3 is a 3-dimensional point with X, Y, Z float32 components.
type Po3 notamath.Po3

// NewPo3 returns a Po3 with the given coordinates.
func NewPo3(x, y, z float32) Po3 {
	return Po3(notamath.Po3{X: x, Y: y, Z: z})
}

func (p Po3) Add(v Vec3) Po3 {
	return Po3(notamath.Po3(p).Add(notamath.Vec3(v)))
}

func (p Po3) Sub(q Po3) Vec3 {
	return Vec3(notamath.Po3(p).SubPo(notamath.Po3(q)))
}

func (p Po3) DistanceSquared(q Po3) float32 {
	return notamath.Po3(p).DistanceSquared(notamath.Po3(q))
}

func (p Po3) Distance(q Po3) float32 {
	return notamath.Po3(p).Distance(notamath.Po3(q))
}

func (p Po3) Equals(q Po3, eps float32) bool {
	return notamath.Po3(p).Equals(notamath.Po3(q), eps)
}

func (p Po3) String() string {
	return fmt.Sprintf("Po3(%f, %f, %f)", p.X, p.Y, p.Z)
}

// Mat3 is a 3x3 matrix stored in row-major order (9 elements).
type Mat3 struct {
	M [9]float32
}

// Mat3Identity returns the 3x3 identity matrix.
func Mat3Identity() Mat3 {
	return Mat3(notamath.Mat3Identity())
}

// Mat3Translation returns a translation matrix from a Vec2 offset.
func Mat3Translation(t Vec2) Mat3 {
	return Mat3(notamath.Mat3Translation(notamath.Vec2(t)))
}

// Mat3Scale returns a scale matrix from a Vec2 factor.
func Mat3Scale(s Vec2) Mat3 {
	return Mat3(notamath.Mat3Scale(notamath.Vec2(s)))
}

// Mat3Rotation returns a rotation matrix by the given angle in radians.
func Mat3Rotation(rad float32) Mat3 {
	return Mat3(notamath.Mat3Rotation(rad))
}

// Mat3Shear returns a shear matrix with the given factors.
func Mat3Shear(kx, ky float32) Mat3 {
	return Mat3(notamath.Mat3Shear(kx, ky))
}

// Mat3TRS returns a combined translation-rotation-scale matrix (applied in that order).
func Mat3TRS(pos Vec2, rot float32, scale Vec2) Mat3 {
	return Mat3(notamath.Mat3TRS(notamath.Vec2(pos), rot, notamath.Vec2(scale)))
}

func (m Mat3) Mul(b Mat3) Mat3 {
	return Mat3(notamath.Mat3(m).Mul(notamath.Mat3(b)))
}

func (m Mat3) TransformPo2(p Po2) Po2 {
	return Po2(notamath.Mat3(m).TransformPo2(notamath.Po2(p)))
}

func (m Mat3) TransformVec2(v Vec2) Vec2 {
	return Vec2(notamath.Mat3(m).TransformVec2(notamath.Vec2(v)))
}

func (m Mat3) Transpose() Mat3 {
	return Mat3(notamath.Mat3(m).Transpose())
}

func (m Mat3) Det() float32 {
	return notamath.Mat3(m).Det()
}

func (m Mat3) InverseAffine() Mat3 {
	return Mat3(notamath.Mat3(m).InverseAffine())
}

func (m Mat3) String() string {
	return fmt.Sprintf(
		"[%f %f %f\n %f %f %f\n %f %f %f]",
		m.M[0], m.M[1], m.M[2],
		m.M[3], m.M[4], m.M[5],
		m.M[6], m.M[7], m.M[8],
	)
}

// Transform2D encapsulates position, rotation, and scale with lazy matrix evaluation
// and frame interpolation support.
type Transform2D struct {
	handle *notamath.Transform2D
}

// NewTransform2D creates a Transform2D with identity scale and zero position/rotation.
func NewTransform2D() *Transform2D {
	t := notamath.NewTransform2D()
	return &Transform2D{handle: &t}
}

func (t *Transform2D) SetPosition(p Vec2) {
	t.handle.SetPosition(notamath.Vec2(p))
}

func (t *Transform2D) SetRotation(r float32) {
	t.handle.SetRotation(r)
}

func (t *Transform2D) SetScale(s Vec2) {
	t.handle.SetScale(notamath.Vec2(s))
}

func (t *Transform2D) Matrix() Mat3 {
	return Mat3(t.handle.Matrix())
}

func (t *Transform2D) TransformPoint(p Po2) Po2 {
	return Po2(t.handle.TransformPoint(notamath.Po2(p)))
}

func (t *Transform2D) TransformVector(v Vec2) Vec2 {
	return Vec2(t.handle.TransformVector(notamath.Vec2(v)))
}

func (t *Transform2D) TranslateBy(delta Vec2) {
	t.handle.TranslateBy(notamath.Vec2(delta))
}

func (t *Transform2D) RotateBy(delta float32) {
	t.handle.RotateBy(delta)
}

func (t *Transform2D) ScaleBy(factor Vec2) {
	t.handle.ScaleBy(notamath.Vec2(factor))
}

func (t *Transform2D) Snapshot() {
	t.handle.Snapshot()
}

func (t *Transform2D) InterpolatedMatrix(alpha float32) Mat3 {
	return Mat3(t.handle.InterpolatedMatrix(alpha))
}

// bridgeVec2 converts a public Vec2 to an internal notamath.Vec2.
func bridgeVec2(v Vec2) notamath.Vec2 {
	return notamath.Vec2(v)
}

// unbridgeVec2 converts an internal notamath.Vec2 to a public Vec2.
func unbridgeVec2(v notamath.Vec2) Vec2 {
	return Vec2(v)
}
