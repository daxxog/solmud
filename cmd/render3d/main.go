package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	TODO_REMOVE_debug "solmud/pkg/debug"
	"solmud/pkg/framebuffer"
	"solmud/pkg/math"
	"solmud/pkg/renderer"
)

const (
	TODO_REMOVE_debug_enabled = false
)

const (
	screenWidth  = math.RS2_SCREEN_WIDTH
	screenHeight = math.RS2_SCREEN_HEIGHT
)

type Game struct {
	frame_buffer *framebuffer.EbitengineFramebuffer
	renderer     renderer.IRenderer
	camera       math.IVec3
	vertices     []math.IVec3
	rotationX    int
	rotationY    int
	rotationZ    int
	trig_table   math.ITrigTable
}

func NewGame(trig math.ITrigTable) *Game {
	renderer, _ := renderer.NewScanlineRenderer(trig)
	camera := math.NewVec3(0, 0, 0)

	game := &Game{
		frame_buffer: nil,
		renderer:     renderer,
		camera:       camera,
		vertices: []math.IVec3{
			math.NewVec3(-997, -610, 2000),
			math.NewVec3(-1411, -24, 2000),
			math.NewVec3(-583, 171, 2000),
		},
		rotationX:  0,
		rotationY:  0,
		rotationZ:  0,
		trig_table: trig,
	}
	return game
}

func (g *Game) Update() error {
	g.rotationX = (g.rotationX + 2) % math.TRIG_SIZE
	g.rotationY = (g.rotationY + 3) % math.TRIG_SIZE
	g.rotationZ = (g.rotationZ + 1) % math.TRIG_SIZE
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.frame_buffer == nil {
		g.frame_buffer = framebuffer.NewEbitengineFramebuffer(screen)
	}

	g.frame_buffer.Image().Fill(color.NRGBA{0, 0, 0, 255})

	if err := g.renderTriangle(); err != nil {
		TODO_REMOVE_debug.Printf("[ERROR] renderTriangle failed: %v\n", err)
	}
}

func (g *Game) renderTriangle() error {
	if TODO_REMOVE_debug_enabled {
		TODO_REMOVE_debug.Printf("[GEOMETRY] rotation: X=%d, Y=%d, Z=%d\n", g.rotationX, g.rotationY, g.rotationZ)
	}

	mat := math.NewRotationMatrix(g.trig_table, g.rotationX, g.rotationY, g.rotationZ)
	projected := make([]math.IVec2, 3)
	visible_count := 0

	for i, v := range g.vertices {
		rotated := v.Rotate(mat)
		screen, depth, visible := math.Project(rotated, g.camera, math.RS2_CENTER_X, math.RS2_CENTER_Y)

		if TODO_REMOVE_debug_enabled {
			TODO_REMOVE_debug.Printf("[GEOMETRY] vertex %d:\n", i)
			TODO_REMOVE_debug.Printf("[GEOMETRY]   world: (%6d, %6d, %6d)\n", v.X(), v.Y(), v.Z())
			TODO_REMOVE_debug.Printf("[GEOMETRY]   rotated: (%6d, %6d, %6d)\n", rotated.X(), rotated.Y(), rotated.Z())
			if visible {
				TODO_REMOVE_debug.Printf("[GEOMETRY]   screen: (%4d, %4d), depth=%4d\n", screen.X(), screen.Y(), depth)
			} else {
				TODO_REMOVE_debug.Printf("[GEOMETRY]   clipped (outside view)\n")
			}
		}

		if visible {
			projected[i] = screen
			visible_count++
		}
	}

	if visible_count != 3 {
		return fmt.Errorf("triangle has %d visible vertices (expected 3)", visible_count)
	}

	g.fillTriangle(projected[0], projected[1], projected[2], color.NRGBA{255, 0, 0, 255})
	return nil
}

func (g *Game) fillTriangle(p0, p1, p2 math.IVec2, c color.NRGBA) {
	x0, y0 := p0.X(), p0.Y()
	x1, y1 := p1.X(), p1.Y()
	x2, y2 := p2.X(), p2.Y()

	if TODO_REMOVE_debug_enabled {
		TODO_REMOVE_debug.Printf("[RASTERIZE] triangle before sort:\n")
		TODO_REMOVE_debug.Printf("[RASTERIZE]   p0: (%4d, %4d)\n", x0, y0)
		TODO_REMOVE_debug.Printf("[RASTERIZE]   p1: (%4d, %4d)\n", x1, y1)
		TODO_REMOVE_debug.Printf("[RASTERIZE]   p2: (%4d, %4d)\n", x2, y2)
	}

	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	if y0 > y2 {
		x0, y0, x2, y2 = x2, y2, x0, y0
	}
	if y1 > y2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dy10 := y1 - y0
	dy20 := y2 - y0
	dy21 := y2 - y1

	if TODO_REMOVE_debug_enabled {
		TODO_REMOVE_debug.Printf("[RASTERIZE] triangle after sort:\n")
		TODO_REMOVE_debug.Printf("[RASTERIZE]   top: (%4d, %4d)\n", x0, y0)
		TODO_REMOVE_debug.Printf("[RASTERIZE]   mid: (%4d, %4d)\n", x1, y1)
		TODO_REMOVE_debug.Printf("[RASTERIZE]   bot: (%4d, %4d)\n", x2, y2)
		TODO_REMOVE_debug.Printf("[RASTERIZE]   bounds: Y[%d,%d], dy10=%d, dy20=%d, dy21=%d\n", y0, y2, dy10, dy20, dy21)
	}

	if dy10 > 0 {
		if TODO_REMOVE_debug_enabled {
			TODO_REMOVE_debug.Printf("[RASTERIZE] upper triangle: Y=%d..%d\n", y0, y1)
		}
		slope01 := ((x1 - x0) << 16) / dy10
		slope02 := ((x2 - x0) << 16) / dy20
		for y := y0; y <= y1; y++ {
			x_start := (x0 << 16) + ((y - y0) * slope01 >> 16)
			x_end := (x0 << 16) + ((y - y0) * slope02 >> 16)

			if TODO_REMOVE_debug_enabled && (y == y0 || y == y1) {
				TODO_REMOVE_debug.Printf("[RASTERIZE]   Y=%d: X=%d..%d\n", y, x_start>>16, x_end>>16)
			}

			g.scanline(x_start>>16, x_end>>16, y, c)
		}
	}

	if dy21 > 0 {
		if TODO_REMOVE_debug_enabled {
			TODO_REMOVE_debug.Printf("[RASTERIZE] lower triangle: Y=%d..%d\n", y1, y2)
		}
		slope12 := ((x2 - x1) << 16) / dy21
		slope02 := ((x2 - x0) << 16) / dy20
		for y := y1; y <= y2; y++ {
			x_start := (x1 << 16) + ((y - y1) * slope12 >> 16)
			x_end := (x0 << 16) + ((y - y0) * slope02 >> 16)

			if TODO_REMOVE_debug_enabled && (y == y1 || y == y2) {
				TODO_REMOVE_debug.Printf("[RASTERIZE]   Y=%d: X=%d..%d\n", y, x_start>>16, x_end>>16)
			}

			g.scanline(x_start>>16, x_end>>16, y, c)
		}
	}
}

func (g *Game) scanline(x_start, x_end, y int, c color.NRGBA) {
	if x_start > x_end {
		x_start, x_end = x_end, x_start
	}

	if y < 0 || y >= screenHeight {
		return
	}

	if x_start < 0 {
		x_start = 0
	}
	if x_end >= screenWidth {
		x_end = screenWidth - 1
	}

	img := g.frame_buffer.Image()

	for x := x_start; x <= x_end; x++ {
		img.Set(x, y, c)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("RS2 317 - Triangle Test")
	trig := math.NewTrigTable()
	game := NewGame(trig)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
