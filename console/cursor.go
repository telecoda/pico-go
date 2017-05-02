package console

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type cursor struct {
	pos
	color   Color
	surface *sdl.Surface
	rect    *sdl.Rect
	on      bool
	speed   time.Duration
}

func newCursor(pb PixelBuffer, color Color) *cursor {
	c := &cursor{
		color:   color,
		surface: pb.(*pixelBuffer).pixelSurface,
		rect:    &sdl.Rect{X: 0, Y: 0, W: _charWidth, H: _charHeight},
		on:      true,
	}

	ticker := time.NewTicker(_cursorFlash)

	go c.cursorFlasher(ticker)

	return c
}

func (c *cursor) cursorFlasher(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			c.on = !c.on // toggle flash
		}
	}
}

// draw cursor on surface
func (c *cursor) render() {
	var color uint32
	// draw at pos
	dRect := &sdl.Rect{X: (int32(c.pos.x) * _charWidth) + 1, Y: (int32(c.pos.y) * _charHeight) - 1, W: _charWidth - 2, H: _charHeight - 2}
	if c.on {
		_, color = _console.palette.getRGBA(c.color)
	} else {
		_, color = _console.palette.getRGBA(BLACK)
	}
	c.surface.FillRect(dRect, color)
}
