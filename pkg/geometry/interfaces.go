// Package geometry provides 3D model loading and parsing capabilities.
//
// This package handles parsing of RuneScape 2 (build 317) model files
// from 377 cache format. It implements delta-encoded vertex
// compression and face indexing system from original client.
//
// Key Components:
//   - ModelLoader: Parses .dat model files from any io.Reader
//   - Model: In-memory representation of 3D geometry
//   - Vertex, Face, TexCoord: Primitive geometry types
//
// Thread Safety: Safe for concurrent loading from different readers
// Performance: O(n) where n is model vertex count
package geometry

import (
	"io"

	"solmud/pkg/math"
)

// IModelLoader defines interface for loading 3D models from RS2 cache.
//
// This interface supports loading from any data source (files, memory,
// network, etc.) by accepting the standard io.Reader interface.
type IModelLoader interface {
	// LoadModel reads model data from any io.Reader implementation.
	//
	// Parameters:
	//   - reader: Any source implementing io.Reader (os.File, bytes.Reader, etc.)
	//   - expected_size: Optional size hint for buffer pre-allocation
	//
	// Returns:
	//   - *Model: Parsed 3D geometry
	//   - error: Parse error, I/O error, or invalid data
	//
	// Thread Safety: Safe to call concurrently with different readers
	LoadModel(reader io.Reader, expected_size int64) (*Model, error)

	// LoadModelFromFile loads a model by ID from cache directory.
	//
	// Parameters:
	//   - cache_path: Path to cache directory (e.g., "mopar/cache/")
	//   - model_id: Unique model identifier from 377 cache
	//
	// Returns:
	//   - *Model: Cached or newly loaded model
	//   - error: File not found, cache read error, or parse error
	LoadModelFromFile(cache_path string, model_id int) (*Model, error)

	// ReleaseModel frees resources for a loaded model.
	//
	// Parameters:
	//   - model: Model instance to release
	//
	// Returns:
	//   - error: Release error (model in use, cache error)
	//
	// Thread Safety: Safe to call concurrently for different models
	ReleaseModel(model *Model) error
}

// Model represents a 3D mesh with vertices, faces, and texture data.
//
// This structure mirrors the RS2 build 317 Model.java format,
// using fixed-point arithmetic and face-based geometry.
type Model struct {
	// Basic counts
	vertex_count  int
	face_count    int
	texture_count int

	// Geometry data
	vertices       []Vertex
	faces          []Face
	texture_coords []TexCoord

	// Rendering properties
	face_colors   [][3]int // Per-face Gouraud shading colors
	face_alpha    []int    // Per-face transparency (0-255)
	face_priority []int    // Per-face rendering priority

	// Bounds for culling
	bounds       Bounds
	radius       int // Horizontal radius from model center
	model_height int // Vertical height
}

// Vertex represents a 3D point with position and optional attributes.
//
// Uses fixed-point 16.16 format for all coordinates.
type Vertex struct {
	x     int  // X position (East/West)
	y     int  // Y position (Up/Down, positive is up)
	z     int  // Z position (North/South)
	alpha byte // Vertex alpha (-1 if none)
}

// Face represents a triangle face with indices and rendering properties.
type Face struct {
	vertex_indices [3]int // Indices into vertices array
	texture_index  int    // Texture ID (or -1 for untextured)
	face_type      int    // Type: 0=flat, 1=textured, 2=alpha, 3=priority
	color_index    int    // Index into color palette
	priority       int    // Depth sorting priority
}

// TexCoord represents texture coordinates for vertex mapping.
type TexCoord struct {
	u int // U coordinate (0-65535)
	v int // V coordinate (0-65535)
}

// Bounds represents an axis-aligned bounding box in world space.
type Bounds struct {
	min math.Vec3 // Minimum extents
	max math.Vec3 // Maximum extents
}
