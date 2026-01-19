// Package framebuffer provides display buffer integration with Ebitengine.
//
// This package implements io.Writer and io.Seeker interfaces on top of
// Ebitengine's Image type, allowing software renderer output
// to be written directly to the display.
package framebuffer

import (
	"io"
	"sync"

	ebiten "github.com/hajimehoshi/ebiten/v2"
)

// IFramebuffer defines display buffer operations.
type IFramebuffer interface {
	GetWriter() io.Writer
}

// EbitengineFramebuffer integrates with Ebitengine v2 image system.
type EbitengineFramebuffer struct {
	image  *ebiten.Image
	buf    []byte
	width  int
	height int
	mutex  sync.RWMutex
}

// NewEbitengineFramebuffer wraps existing Ebitengine image.
func NewEbitengineFramebuffer(image *ebiten.Image) *EbitengineFramebuffer {
	bounds := image.Bounds()
	return &EbitengineFramebuffer{
		image:  image,
		buf:    make([]byte, bounds.Dx()*bounds.Dy()*4),
		width:  bounds.Dx(),
		height: bounds.Dy(),
	}
}

func (ef *EbitengineFramebuffer) Write(p []byte) (int, error) {
	ef.mutex.Lock()
	defer ef.mutex.Unlock()

	copy(ef.buf, p)

	return len(p), nil
}

func (ef *EbitengineFramebuffer) GetWriter() io.Writer {
	return ef
}

func (ef *EbitengineFramebuffer) Width() int {
	return ef.width
}

func (ef *EbitengineFramebuffer) Height() int {
	return ef.height
}

func (ef *EbitengineFramebuffer) Image() *ebiten.Image {
	return ef.image
}
