package math

import (
	stdmath "math"
)

func (v *Vec3) X() int {
	return v.x
}

func (v *Vec3) Y() int {
	return v.y
}

func (v *Vec3) Z() int {
	return v.z
}

func (v *Vec3) Add(other IVec3) IVec3 {
	return NewVec3(
		v.x+other.X(),
		v.y+other.Y(),
		v.z+other.Z(),
	)
}

func (v *Vec3) Subtract(other IVec3) IVec3 {
	return NewVec3(
		v.x-other.X(),
		v.y-other.Y(),
		v.z-other.Z(),
	)
}

func (v *Vec3) Scale(scalar int) IVec3 {
	return NewVec3(
		(v.x*scalar)>>16,
		(v.y*scalar)>>16,
		(v.z*scalar)>>16,
	)
}

func (v *Vec3) LengthSquared() int {
	return (v.x*v.x + v.y*v.y + v.z*v.z) >> 16
}

func (v *Vec3) Length() int {
	return int(stdmath.Sqrt(float64(v.LengthSquared())))
}

func (v *Vec3) Normalize() IVec3 {
	length := v.Length()
	if length == 0 {
		return v
	}
	scale := (65536 << 16) / length
	return v.Scale(scale)
}

func (v *Vec3) Rotate(mat IMatrix) IVec3 {
	return v.rotateInternal(mat.(*Matrix))
}

func (v *Vec3) rotateInternal(mat *Matrix) IVec3 {
	return NewVec3(
		(v.x*mat.m00+v.y*mat.m01+v.z*mat.m02)>>16,
		(v.x*mat.m10+v.y*mat.m11+v.z*mat.m12)>>16,
		(v.x*mat.m20+v.y*mat.m21+v.z*mat.m22)>>16,
	)
}
