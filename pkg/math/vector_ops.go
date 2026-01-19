package math

// DotProduct calculates 3D vector dot product.
func DotProduct(a, b IVec3) int {
	return (a.X()*b.X() + a.Y()*b.Y() + a.Z()*b.Z()) >> 16
}

// CrossProduct calculates 3D vector cross product.
func CrossProduct(a, b IVec3) IVec3 {
	return NewVec3(
		(a.Y()*b.Z()-a.Z()*b.Y())>>16,
		(a.Z()*b.X()-a.X()*b.Z())>>16,
		(a.X()*b.Y()-a.Y()*b.X())>>16,
	)
}
