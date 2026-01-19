// Package renderer provides core software 3D rendering capabilities.
//
// This package implements a software 3D renderer matching RuneScape 2
// (build 317, 377 cache) rendering pipeline. It uses scanline
// rasterization with Gouraud shading, Painter's algorithm for depth sorting,
// and fixed-point arithmetic throughout.
//
// Key Components:
//   - Renderer: Main rendering orchestrator
//   - ScanlineRenderer: Triangle rasterization implementation
//   - RenderStats: Performance and debug metrics
//
// Thread Safety: Not safe for concurrent render calls during single frame
// Performance: Target 60fps at 512x384 resolution with ~1000 triangles
package renderer

import (
	"io"

	"solmud/pkg/math"
)

// IRenderer defines software rasterization operations.
//
// This interface provides the core contract for rendering 3D scenes using
// software algorithms (no hardware acceleration). It integrates with any
// framebuffer implementation via io.Writer interface.
type IRenderer interface {
	// RenderScene renders all visible models to the provided framebuffer writer.
	//
	// Parameters:
	//   - scene: Scene description with models, camera, and lighting
	//   - writer: Any io.Writer implementation (bytes.Buffer, os.File, etc.)
	//
	// Returns:
	//   - stats: Performance metrics for this render pass
	//   - error: Rendering error or invalid state
	//
	// Thread Safety: Not safe for concurrent calls
	RenderScene(scene *Scene, writer io.Writer) (*RenderStats, error)

	// ClearFramebuffer resets framebuffer to background color via writer.
	//
	// Parameters:
	//   - writer: Any io.Writer implementation
	//   - color: Background color in 0xAARRGGBB format
	//
	// Thread Safety: Not safe for concurrent calls
	ClearFramebuffer(writer io.Writer, color int) error

	// SetClipRegion restricts rendering to rectangular area.
	//
	// Parameters:
	//   - writer: Any io.Writer implementation (with seeking support)
	//   - region: Clipping rectangle in screen coordinates
	//
	// Thread Safety: Not safe for concurrent calls
	SetClipRegion(writer io.Seeker, region *Rect) error
}

// RenderStats provides performance metrics for a render pass.
type RenderStats struct {
	triangles_rendered int // Number of triangles drawn
	vertices_processed int // Number of vertices transformed
	pixels_filled      int // Number of pixels written to framebuffer
	culled_triangles   int // Number of triangles culled (backface/frustum)
	render_time_ms     int // Time taken for render pass
	frames_per_second  float64
}

// Rect represents a rectangular region for clipping.
type Rect struct {
	x      int // Left coordinate
	y      int // Top coordinate
	width  int // Width in pixels
	height int // Height in pixels
}

// Scene describes the complete 3D scene to render.
type Scene struct {
	Models     []*ModelInstance
	Camera     *Camera
	Lights     []*Light
	Background int   // Background color
	ClipRegion *Rect // Optional clip region (nil = full screen)
}

// ModelInstance represents a 3D model placed in the world.
type ModelInstance struct {
	ModelID  int
	Position math.Vec3 // World-space position (X, Y, Z)
	Rotation int       // Yaw rotation (0-2048, 2048 = 360°)
	Scale    int       // Uniform scale factor
	Visible  bool
	Distance int // Distance from camera for sorting
}

// Camera defines view and projection settings.
type Camera struct {
	Position       math.Vec3 // World-space position
	Yaw            int       // Horizontal rotation (0-2048)
	Pitch          int       // Vertical rotation (0-2048)
	ViewportWidth  int
	ViewportHeight int
	FovScale       int // Perspective scale (512 in RS2)
	NearPlane      int // Near culling (50 units)
	FarPlane       int // Far culling (3500 units)
}

// Light defines a directional or point light source.
type Light struct {
	Position  math.Vec3 // World-space position (for point lights)
	Direction math.Vec3 // Direction (for directional lights)
	Intensity int       // 0-255 brightness
	Color     int       // RGB color
}
