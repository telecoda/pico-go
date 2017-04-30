package console

import "github.com/telecoda/pico-go/api"

// Cls - clears pixel buffer
func (c *console) Cls() {
	c.Display.Cls()
}

// ClsColor - fill pixel buffer with a set color
func (c *console) ClsColor(colorId api.Color) {
	c.Display.ClsColor(colorId)
}

// Flip - copy offscreen buffer to onscreen buffer
func (c *console) Flip() {
	c.Display.Flip()
}
