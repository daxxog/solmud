package math

// NewRotationMatrixX creates a 3x3 rotation matrix around X-axis.
func NewRotationMatrixX(trig ITrigTable, angle int) IMatrix {
	sin := trig.Sin(angle)
	cos := trig.Cos(angle)

	return &Matrix{
		m00: 65536, m01: 0, m02: 0,
		m10: 0, m11: cos, m12: -sin,
		m20: 0, m21: sin, m22: cos,
	}
}

// NewRotationMatrixY creates a 3x3 rotation matrix around Y-axis.
func NewRotationMatrixY(trig ITrigTable, angle int) IMatrix {
	sin := trig.Sin(angle)
	cos := trig.Cos(angle)

	return &Matrix{
		m00: cos, m01: 0, m02: sin,
		m10: 0, m11: 65536, m12: 0,
		m20: -sin, m21: 0, m22: cos,
	}
}

// NewRotationMatrixZ creates a 3x3 rotation matrix around Z-axis.
func NewRotationMatrixZ(trig ITrigTable, angle int) IMatrix {
	sin := trig.Sin(angle)
	cos := trig.Cos(angle)

	return &Matrix{
		m00: cos, m01: -sin, m02: 0,
		m10: sin, m11: cos, m12: 0,
		m20: 0, m21: 0, m22: 65536,
	}
}

// NewRotationMatrix creates a 3x3 rotation matrix around all three axes.
//
// Combines X, Y, and Z rotations by multiplying matrices in order Z * Y * X.
// Angles are in trig table index format (0-2048 for 360°).
func NewRotationMatrix(trig ITrigTable, angleX, angleY, angleZ int) IMatrix {
	matX := NewRotationMatrixX(trig, angleX)
	matY := NewRotationMatrixY(trig, angleY)
	matZ := NewRotationMatrixZ(trig, angleZ)

	return Multiply(Multiply(matZ, matY), matX)
}

// Multiply combines two 3x3 matrices.
func Multiply(a, b IMatrix) IMatrix {
	a00, a01, a02 := a.M00(), a.M01(), a.M02()
	a10, a11, a12 := a.M10(), a.M11(), a.M12()
	a20, a21, a22 := a.M20(), a.M21(), a.M22()

	b00, b01, b02 := b.M00(), b.M01(), b.M02()
	b10, b11, b12 := b.M10(), b.M11(), b.M12()
	b20, b21, b22 := b.M20(), b.M21(), b.M22()

	m00 := (a00*b00 + a01*b10 + a02*b20) >> 16
	m01 := (a00*b01 + a01*b11 + a02*b21) >> 16
	m02 := (a00*b02 + a01*b12 + a02*b22) >> 16

	m10 := (a10*b00 + a11*b10 + a12*b20) >> 16
	m11 := (a10*b01 + a11*b11 + a12*b21) >> 16
	m12 := (a10*b02 + a11*b12 + a12*b22) >> 16

	m20 := (a20*b00 + a21*b10 + a22*b20) >> 16
	m21 := (a20*b01 + a21*b11 + a22*b21) >> 16
	m22 := (a20*b02 + a21*b12 + a22*b22) >> 16

	return &Matrix{
		m00: m00, m01: m01, m02: m02,
		m10: m10, m11: m11, m12: m12,
		m20: m20, m21: m21, m22: m22,
	}
}

func (m *Matrix) M00() int {
	return m.m00
}

func (m *Matrix) M01() int {
	return m.m01
}

func (m *Matrix) M02() int {
	return m.m02
}

func (m *Matrix) M10() int {
	return m.m10
}

func (m *Matrix) M11() int {
	return m.m11
}

func (m *Matrix) M12() int {
	return m.m12
}

func (m *Matrix) M20() int {
	return m.m20
}

func (m *Matrix) M21() int {
	return m.m21
}

func (m *Matrix) M22() int {
	return m.m22
}
