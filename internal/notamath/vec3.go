package notamath

import (
	"fmt"
	"math"
)

type Vec3 struct {
	X, Y, Z float32
}

func (v Vec3) Add(o Vec3) Vec3 {
	return Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Vec3) Sub(o Vec3) Vec3 {
	return Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v Vec3) Mul(s float32) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v Vec3) Div(s float32) Vec3 {
	return Vec3{v.X / s, v.Y / s, v.Z / s}
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func (v Vec3) Dot(o Vec3) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{
		X: v.Y*o.Z - v.Z*o.Y,
		Y: v.Z*o.X - v.X*o.Z,
		Z: v.X*o.Y - v.Y*o.X,
	}
}

func (v Vec3) LenSquared() float32 {
	return v.Dot(v)
}

func (v Vec3) Len() float32 {
	return float32(math.Sqrt(float64(v.LenSquared())))
}

func (v Vec3) Normalize() Vec3 {
	l := v.Len()
	if l == 0 {
		return Vec3{}
	}
	return v.Div(l)
}

func (v Vec3) Distance(o Vec3) float32 {
	return v.Sub(o).Len()
}

func (v Vec3) Lerp(to Vec3, t float32) Vec3 {
	return v.Add(to.Sub(v).Mul(t))
}

func (v Vec3) Project(onto Vec3) Vec3 {
	d := onto.LenSquared()
	if d == 0 {
		return Vec3{}
	}
	return onto.Mul(v.Dot(onto) / d)
}

// Angle returns the angle (reference angle) between this vector and the other vector in radians.
func (v Vec3) Angle(o Vec3) float32 {
	d := v.Dot(o)
	l := v.Len() * o.Len()
	if l == 0 {
		return 0
	}
	return float32(math.Acos(float64(d / l)))
}

// Rotate rotates the vector around the given axis vector by the given angle in radians.
func (v Vec3) Rotate(axis Vec3, angle float32) Vec3 {
	k := axis.Normalize()
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	return v.Mul(c).
		Add(k.Cross(v).Mul(s)).
		Add(k.Mul(k.Dot(v) * (1 - c)))
}

func (v Vec3) String() string {
	return fmt.Sprintf("Vector3(%f, %f, %f)", v.X, v.Y, v.Z)
}
