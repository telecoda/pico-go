package console

import (
	"image"
	"time"
)

type cursor struct {
	pos
	color Color
	pb    *pixelBuffer
	rect  *image.Rectangle
	on    bool
	speed time.Duration
}

func newCursor(pb PixelBuffer, color Color) *cursor {
	c := &cursor{
		color: color,
		pb:    pb.getPixelBuffer(),
		rect:  &image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: _console.Config.fontWidth, Y: _console.Config.fontHeight}},
		on:    true,
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
	if c.on {
		c.pb.RectFillWithColor(c.pos.x*_console.Config.fontWidth, c.pos.y*_console.Config.fontHeight, _console.Config.fontWidth, _console.Config.fontHeight-1, c.color)
	}
}
