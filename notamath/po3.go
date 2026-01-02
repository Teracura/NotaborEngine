package notamath

import "fmt"

type Po3 struct {
	X, Y, Z float32
}

func (p Po3) Add(v Vec3) Po3 {
	return Po3{p.X + v.X, p.Y + v.Y, p.Z + v.Z}
}

func (p Po3) SubPo(q Po3) Vec3 {
	return Vec3{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

func (p Po3) SubVec(q Vec3) Vec3 {
	return Vec3{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

func (p Po3) DistanceSquared(q Po3) float32 {
	return p.SubPo(q).LenSquared()
}

func (p Po3) Distance(q Po3) float32 {
	return p.SubPo(q).Len()
}

func (p Po3) Equals(q Po3, eps float32) bool {
	return p.SubPo(q).LenSquared() <= eps*eps
}

func (p Po3) String() string {
	return fmt.Sprintf("Point3(%f, %f, %f)", p.X, p.Y, p.Z)
}
