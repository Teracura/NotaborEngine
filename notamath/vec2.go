package notamath

import (
	"fmt"
	"math"
)

type Vec2 struct {
	X, Y float32
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{v.X + o.X, v.Y + o.Y}
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{v.X - o.X, v.Y - o.Y}
}

func (v Vec2) Mul(s float32) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

func (v Vec2) Div(s float32) Vec2 {
	return Vec2{v.X / s, v.Y / s}
}

func (v Vec2) Dot(o Vec2) float32 {
	return v.X*o.X + v.Y*o.Y
}

func (v Vec2) Cross(o Vec2) float32 {
	return v.X*o.Y - v.Y*o.X
}

func (v Vec2) LenSquared() float32 {
	return v.Dot(v)
}

func (v Vec2) Distance(o Vec2) float32 {
	return v.Sub(o).Len()
}

func (v Vec2) Len() float32 {
	return float32(math.Sqrt(float64(v.LenSquared())))
}

func (v Vec2) Neg() Vec2 {
	return Vec2{-v.X, -v.Y}
}

func (v Vec2) Normalize() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{}
	}
	return v.Mul(1 / l)
}

func (v Vec2) Perp() Vec2 {
	return Vec2{-v.Y, v.X}
}

func (v Vec2) Lerp(a Vec2, t float32) Vec2 {
	return v.Add(a.Sub(v).Mul(t))
}

func (v Vec2) Project(onto Vec2) Vec2 {
	d := onto.LenSquared()
	if d == 0 {
		return Vec2{}
	}
	return onto.Mul(v.Dot(onto) / d)
}

func (v Vec2) Angle(o Vec2) float32 {
	d := v.Dot(o)
	l := v.Len() * o.Len()
	if l == 0 {
		return 0
	}
	return float32(math.Acos(float64(d / l)))
}

func (v Vec2) Rotate(rad float32) Vec2 {
	c := float32(math.Cos(float64(rad)))
	s := float32(math.Sin(float64(rad)))
	return Vec2{
		v.X*c - v.Y*s,
		v.X*s + v.Y*c,
	}
}

func (v Vec2) String() string {
	return fmt.Sprintf("Vector2(%f, %f)", v.X, v.Y)
}
