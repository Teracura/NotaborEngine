package notamath

import "math"

type Transform2D struct {
	Position Vec2
	Rotation float32 // radians
	Scale    Vec2

	// previous (for interpolation)
	prevPosition Vec2
	prevRotation float32
	prevScale    Vec2

	Dirty  bool //mutated without mapping to a matrix
	matrix Mat3
}

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

func (t *Transform2D) dirtify() {
	t.Dirty = true
}

func (t *Transform2D) TransformPoint(p Po2) Po2 {
	return t.Matrix().TransformPo2(p)
}

func (t *Transform2D) TransformVector(v Vec2) Vec2 {
	return t.Matrix().TransformVec2(v)
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

// Snapshot saves the previous transform values
func (t *Transform2D) Snapshot() {
	t.prevPosition = t.Position
	t.prevRotation = t.Rotation
	t.prevScale = t.Scale
}

func (t *Transform2D) InterpolatedMatrix(alpha float32) Mat3 {
	pos := t.prevPosition.Lerp(t.Position, alpha)
	scale := t.prevScale.Lerp(t.Scale, alpha)
	rot := lerpAngle(t.prevRotation, t.Rotation, alpha)

	return Mat3TRS(pos, rot, scale)
}

func lerpAngle(a, b, t float32) float32 {
	d := b - a
	for d > math.Pi {
		d -= 2 * math.Pi
	}
	for d < -math.Pi {
		d += 2 * math.Pi
	}
	return a + d*t
}
