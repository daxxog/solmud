// Package math provides fixed-point 3D mathematics for rendering.
//
// This package implements the mathematics used by RuneScape 2 (build 317)
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

// Vec3 represents a 3D vector using fixed-point arithmetic.
//
// Uses 16.16 fixed-point format where 16 bits are integer
// and 16 bits are fractional. This provides precision while avoiding
// expensive floating-point operations.
type Vec3 struct {
	x int // Fixed-point 16.16 format
	y int // Fixed-point 16.16 format
	z int // Fixed-point 16.16 format
}

// Vec2 represents a 2D screen coordinate using fixed-point arithmetic.
type Vec2 struct {
	x int // Fixed-point 16.16 format
	y int // Fixed-point 16.16 format
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
func NewVec3(x, y, z int) Vec3 {
	return Vec3{x: x, y: y, z: z}
}

// NewVec2 creates a new 2D vector.
func NewVec2(x, y int) Vec2 {
	return Vec2{x: x, y: y}
}

// NewMatrix creates a new 3x3 identity matrix.
func NewMatrix() Matrix {
	return Matrix{
		m00: 65536, m01: 0, m02: 0,
		m10: 0, m11: 65536, m12: 0,
		m20: 0, m21: 0, m22: 65536,
	}
}

// NewTrigTable initializes and returns precomputed trigonometric tables.
//
// Returns a TrigTable with 2048 entries for sine and cosine.
func NewTrigTable() *TrigTable {
	table := &TrigTable{}

	// PI / 1024 for each index step (2048 entries = 2π)
	// 0.0030679615 rad per step
	for i := 0; i < 2048; i++ {
		angle := float64(i) * 0.0030679615
		table.sin[i] = int(65536.0 * stdmath.Sin(angle))
		table.cos[i] = int(65536.0 * stdmath.Cos(angle))
	}

	return table
}

// Rotate applies rotation matrix to vector.
func (v Vec3) Rotate(mat *Matrix) Vec3 {
	return Vec3{
		x: (v.x*mat.m00 + v.y*mat.m01 + v.z*mat.m02) >> 16,
		y: (v.x*mat.m10 + v.y*mat.m11 + v.z*mat.m12) >> 16,
		z: (v.x*mat.m20 + v.y*mat.m21 + v.z*mat.m22) >> 16,
	}
}

// Project transforms world coordinate to screen coordinate.
//
// Parameters:
//   - world: 3D point in world space
//   - camera: Camera position and orientation
//   - screen_width: Viewport width in pixels
//   - screen_height: Viewport height in pixels
//
// Returns:
//   - screen: 2D screen coordinate
//   - depth: Z distance from camera
//   - visible: False if behind camera or clipped
func Project(world Vec3, camera *Vec3, screen_width, screen_height int) (screen Vec2, depth int, visible bool) {
	// Calculate relative position
	rel_x := world.x - camera.x
	rel_y := world.y - camera.y
	rel_z := world.z - camera.z

	// Near plane culling
	if rel_z < 50 || rel_z > 3500 {
		return Vec2{}, rel_z, false
	}

	// Perspective projection
	fov_scale := 512 // RS2 perspective scale factor
	screen_x := screen_width/2 + (rel_x*fov_scale)/rel_z
	screen_y := screen_height/2 + (rel_y*fov_scale)/rel_z

	// Frustum culling
	if screen_x < 0 || screen_x >= screen_width {
		return Vec2{}, rel_z, false
	}
	if screen_y < 0 || screen_y >= screen_height {
		return Vec2{}, rel_z, false
	}

	return Vec2{x: screen_x, y: screen_y}, rel_z, true
}

// DotProduct calculates 3D vector dot product.
func DotProduct(a, b Vec3) int {
	return (a.x*b.x + a.y*b.y + a.z*b.z) >> 16
}

// CrossProduct calculates 3D vector cross product.
func CrossProduct(a, b Vec3) Vec3 {
	return Vec3{
		x: (a.y*b.z - a.z*b.y) >> 16,
		y: (a.z*b.x - a.x*b.z) >> 16,
		z: (a.x*b.y - a.y*b.x) >> 16,
	}
}

// Add adds two vectors.
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{
		x: v.x + other.x,
		y: v.y + other.y,
		z: v.z + other.z,
	}
}

// Subtract subtracts other vector from this vector.
func (v Vec3) Subtract(other Vec3) Vec3 {
	return Vec3{
		x: v.x - other.x,
		y: v.y - other.y,
		z: v.z - other.z,
	}
}

// Scale multiplies vector by scalar.
func (v Vec3) Scale(scalar int) Vec3 {
	return Vec3{
		x: (v.x * scalar) >> 16,
		y: (v.y * scalar) >> 16,
		z: (v.z * scalar) >> 16,
	}
}

// LengthSquared returns squared length (avoids sqrt for comparisons).
func (v Vec3) LengthSquared() int {
	return (v.x*v.x + v.y*v.y + v.z*v.z) >> 16
}

// Length returns vector length (uses sqrt).
func (v Vec3) Length() int {
	return int(sqrt(float64(v.LengthSquared())))
}

func sqrt(x float64) float64 {
	// Simple integer sqrt approximation
	// In production, use math.Sqrt from standard library
	z := 1.0
	for i := 1; i < 20; i++ {
		z = (z + x/z) / 2.0
	}
	return z
}
