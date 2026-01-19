// Package math provides fixed-point 3D mathematics for rendering.
//
// This package implements mathematics used by RuneScape 2 (build 317)
// rendering pipeline, including precomputed trigonometric tables,
// fixed-point arithmetic, and 3D transformations.
//
// Key Components:
//   - Vec3, Vec2: Fixed-point vector types
//   - Matrix: 3x3 transformation matrices
//   - TrigTable: Precomputed sine/cosine lookup tables
//
// Thread Safety: Safe for concurrent read access to trig tables
// Performance: O(1) for all operations (lookup tables)
package math

import (
	stdmath "math"
)

const (
	// RS2 Screen dimensions from DrawingArea.java
	RS2_SCREEN_WIDTH  = 512
	RS2_SCREEN_HEIGHT = 334
	RS2_CENTER_X      = 511 // centerX = bottomX - 1
	RS2_CENTER_Y      = 256 // centerY = bottomX / 2

	// Fixed-point math constants
	FIXED_SHIFT  = 16   // 16.16 format
	SCREEN_SHIFT = 9    // Screen projection shift (<< 9 = *512)
	TRIG_SIZE    = 2048 // 360° rotation resolution

	// Rendering constants
	NEAR_PLANE = 50   // Minimum render distance
	FAR_PLANE  = 3500 // Maximum render distance
)

// IVec3 defines interface for 3D vector operations.
//
// This interface provides the public contract for 3D vector operations
// using fixed-point arithmetic. Implementations maintain private fields
// and expose functionality through accessor methods.
type IVec3 interface {
	// Accessor methods
	X() int
	Y() int
	Z() int

	// Vector operations
	Add(other IVec3) IVec3
	Subtract(other IVec3) IVec3
	Scale(scalar int) IVec3
	Length() int
	LengthSquared() int
	Normalize() IVec3

	// Transformations
	Rotate(mat IMatrix) IVec3
}

// IVec2 defines interface for 2D screen coordinate operations.
type IVec2 interface {
	X() int
	Y() int
}

// IMatrix defines interface for 3x3 transformation matrix operations.
//
// Provides field accessor methods for fixed-point transformation matrices.
// Minimal interface design - only methods needed for current usage.
type IMatrix interface {
	M00() int
	M01() int
	M02() int
	M10() int
	M11() int
	M12() int
	M20() int
	M21() int
	M22() int
}

// ITrigTable defines interface for trigonometric table access.
//
// Provides access to precomputed sine/cosine lookup tables.
// Used for efficient rotation calculations without runtime trig functions.
type ITrigTable interface {
	Sin(index int) int
	Cos(index int) int
}

// Vec3 represents a 3D vector using fixed-point arithmetic.
//
// Uses 16.16 fixed-point format where 16 bits are integer
// and 16 bits are fractional. This provides precision while avoiding
// expensive floating-point operations.
type Vec3 struct {
	x, y, z int // Fixed-point 16.16 format
}

// Vec2 represents a 2D screen coordinate using fixed-point arithmetic.
type Vec2 struct {
	x, y int // Fixed-point 16.16 format
}

// Matrix represents a 3x3 transformation matrix using fixed-point math.
type Matrix struct {
	m00, m01, m02 int // Fixed-point 16.16 values
	m10, m11, m12 int
	m20, m21, m22 int
}

// TrigTable provides precomputed trigonometric lookup tables.
//
// These tables match the RS2 client's format with 2048 entries
// for full 360° rotation, using 65536 as the fixed-point scale.
type TrigTable struct {
	sin [2048]int // 65536 * sin(angle) for angle 0-2048
	cos [2048]int // 65536 * cos(angle) for angle 0-2048
}

// NewVec3 creates a new 3D vector.
//
// Returns IVec3 interface to hide implementation details.
func NewVec3(x, y, z int) IVec3 {
	return &Vec3{x: x, y: y, z: z}
}

// NewVec2 creates a new 2D vector.
//
// Returns IVec2 interface to hide implementation details.
func NewVec2(x, y int) IVec2 {
	return &Vec2{x: x, y: y}
}

// NewMatrix creates a new 3x3 identity matrix.
//
// Returns IMatrix interface to hide implementation details.
func NewMatrix() IMatrix {
	return &Matrix{
		m00: 65536, m01: 0, m02: 0,
		m10: 0, m11: 65536, m12: 0,
		m20: 0, m21: 0, m22: 65536,
	}
}

// NewTrigTable initializes and returns precomputed trigonometric tables.
//
// Returns ITrigTable interface to hide implementation details.
func NewTrigTable() ITrigTable {
	table := &TrigTable{}

	for i := 0; i < 2048; i++ {
		angle := float64(i) * 0.0030679615
		table.sin[i] = int(65536.0 * stdmath.Sin(angle))
		table.cos[i] = int(65536.0 * stdmath.Cos(angle))
	}

	return table
}
