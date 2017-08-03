package console

import (
	"image"
	"time"
)

type cursor struct {
	pos
	color   Color
	surface *image.RGBA
	rect    *image.Rectangle
	on      bool
	speed   time.Duration
}

func newCursor(pb PixelBuffer, color Color) *cursor {
	c := &cursor{
		color:   color,
		surface: pb.getPixelBuffer().pixelSurface,
		rect:    &image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: _console.Config.fontWidth, Y: _console.Config.fontHeight}},
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
	// draw at pos
	// TODO
	//	dRect := &sdl.Rect{X: (int32(c.pos.x * _console.Config.fontWidth)), Y: (int32(c.pos.y * _console.Config.fontHeight)), W: int32(_console.Config.fontWidth), H: int32(_console.Config.fontHeight - 1)}
	if c.on {
		// TODO
		//c.surface.FillRect(dRect, uint32(c.color))
	}
}
