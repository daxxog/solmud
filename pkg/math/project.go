package math

import (
	"solmud/pkg/assert"
)

// Project transforms world coordinate to screen coordinate using RS2 algorithm.
//
// RS2 screen projection formula: screenX = (worldX << 9) / depth + centerX
// Uses 9-bit shift for screen projection (<< 9 = *512) matching original client.
//
// Parameters:
//   - world: 3D point in world space
//   - camera: Camera position
//   - center_x: Viewport center X (DrawingArea.centerX = 511)
//   - center_y: Viewport center Y (DrawingArea.centerY = 256)
//
// Returns:
//   - screen: 2D screen coordinate
//   - depth: Z distance from camera
//   - visible: False if behind camera or clipped
func Project(world IVec3, camera IVec3, center_x, center_y int) (IVec2, int, bool) {
	rel_x := world.X() - camera.X()
	rel_y := world.Y() - camera.Y()
	rel_z := world.Z() - camera.Z()

	assert.NewAssert().True(rel_z != 0, "relative Z cannot be zero (division by zero)")

	if rel_z < NEAR_PLANE || rel_z > FAR_PLANE {
		return NewVec2(0, 0), rel_z, false
	}

	screen_x := ((rel_x << SCREEN_SHIFT) / rel_z) + center_x
	screen_y := ((rel_y << SCREEN_SHIFT) / rel_z) + center_y

	if screen_x < 0 || screen_x >= RS2_SCREEN_WIDTH {
		return NewVec2(0, 0), rel_z, false
	}
	if screen_y < 0 || screen_y >= RS2_SCREEN_HEIGHT {
		return NewVec2(0, 0), rel_z, false
	}

	return NewVec2(screen_x, screen_y), rel_z, true
}
