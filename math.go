package main

import (
	"NotaborEngine/internal/notamath"
)

type Vec2 struct {
	X, Y float32
}

func (v Vec2) Mul(s float32) Vec2 {
	return Vec2{X: v.X * s, Y: v.Y * s}
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{X: v.X + o.X, Y: v.Y + o.Y}
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{X: v.X - o.X, Y: v.Y - o.Y}
}

// Internal bridge - strictly not named toInternal/fromInternal
func bridgeVec2(v Vec2) notamath.Vec2 {
	return notamath.Vec2{X: v.X, Y: v.Y}
}

func unbridgeVec2(v notamath.Vec2) Vec2 {
	return Vec2{X: v.X, Y: v.Y}
}
