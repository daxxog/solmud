package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"solmud/pkg/debug"
	"solmud/pkg/framebuffer"
	"solmud/pkg/math"
	"solmud/pkg/renderer"
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
	rotation     int
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
			math.NewVec3(0, -50, 1000),
			math.NewVec3(-50, 50, 1000),
			math.NewVec3(50, 50, 1000),
		},
		rotation:   0,
		trig_table: trig,
	}
	return game
}

func (g *Game) Update() error {
	g.rotation = (g.rotation + 1) % math.TRIG_SIZE
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.frame_buffer == nil {
		g.frame_buffer = framebuffer.NewEbitengineFramebuffer(screen)
	}

	g.frame_buffer.Image().Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})

	g.renderTriangle()
}

func (g *Game) renderTriangle() error {
	debug.Printf("[DEBUG] renderTriangle ENTRY - camera: X=%d, Y=%d, Z=%d\n",
		g.camera.X(), g.camera.Y(), g.camera.Z())

	mat := math.NewMatrix()
	projected := make([]math.IVec2, 3)
	visible_count := 0

	for i, v := range g.vertices {
		rotated := v.Rotate(mat)
		screen, _, visible := math.Project(rotated, g.camera, math.RS2_CENTER_X, math.RS2_CENTER_Y)
		if visible {
			projected[i] = screen
			visible_count++
		}
	}

	debug.Printf("[DEBUG]   Visible vertices: %d/3\n", visible_count)

	if visible_count != 3 {
		return fmt.Errorf("triangle has %d visible vertices (expected 3)", visible_count)
	}

	g.fillTriangle(projected[0], projected[1], projected[2], color.NRGBA{255, 255, 255, 255})

	debug.Printf("[DEBUG] renderTriangle END\n")
	return nil
}

func (g *Game) fillTriangle(p0, p1, p2 math.IVec2, c color.NRGBA) {
	x0, y0 := p0.X(), p0.Y()
	x1, y1 := p1.X(), p1.Y()
	x2, y2 := p2.X(), p2.Y()

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

	if dy10 > 0 {
		slope01 := ((x1 - x0) << 16) / dy10
		slope02 := ((x2 - x0) << 16) / dy20
		for y := y0; y <= y1; y++ {
			x_start := (x0 << 16) + ((y - y0) * slope01 >> 16)
			x_end := (x0 << 16) + ((y - y0) * slope02 >> 16)
			g.scanline(x_start>>16, x_end>>16, y, c)
		}
	}

	if dy21 > 0 {
		slope12 := ((x2 - x1) << 16) / dy21
		slope02 := ((x2 - x0) << 16) / dy20
		for y := y1; y <= y2; y++ {
			x_start := (x1 << 16) + ((y - y1) * slope12 >> 16)
			x_end := (x0 << 16) + ((y - y0) * slope02 >> 16)
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
