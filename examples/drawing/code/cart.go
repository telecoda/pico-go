package code

import (
	"github.com/telecoda/pico-go/console"
)

// This example demonstrates the drawing primatives

type cartridge struct {
	cfg                 console.Config // holds details of console config
	console.PixelBuffer                // ref to console display

}

// NewCart - initialise a struct implementing Cartridge interface
func NewCart() console.Cartridge {
	return &cartridge{
		cfg: console.DefaultConfig(),
	}
}

// GetConfig - return config need for Cart to run
func (c *cartridge) GetConfig() console.Config {
	return c.cfg
}

// Init - called once when cart is initialised
func (c *cartridge) Init(pb console.PixelBuffer) {
	// the Init method receives a PixelBuffer reference
	// hold onto this reference, this is the display that
	// your code will be drawing onto each frame
	c.PixelBuffer = pb
	c.ClsWithColor(console.BLACK)
}

// Update -  called once every frame
func (c *cartridge) Update() {
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls()
	c.RectFillWithColor(10, 10, 20, 20, console.RED)
	c.RectFill(30, 30, 50, 50)
	c.Color(console.GREEN)
	c.Rect(60, 60, 80, 100)
	c.LineWithColor(0, 60, 70, 85, console.BLUE)
	c.Color(console.LIGHT_GRAY)
	c.Line(0, 65, 75, 90)
	c.PSetWithColor(90, 90, console.BROWN)
	c.PSet(100, 90)
	// get color of rect
	rectColor := c.PGet(60, 60)
	c.PrintAtWithColor("Here!", 50, 10, rectColor)
	c.Circle(64, 64, 30)
	c.CircleWithColor(64, 64, 40, console.RED)
	c.CircleFill(64, 64, 20)
	c.CircleFillWithColor(64, 64, 10, console.YELLOW)
}
