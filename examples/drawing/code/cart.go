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
	//c.ClsWithColor(console.RED)
}

// Update -  called once every frame
func (c *cartridge) Update() {
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.ClsWithColor(console.BLACK)
	c.PSet(0, 0)
	c.PrintAtWithColor("DRAWING:", 50, 5, console.WHITE)
	c.Line(0, 12, 128, 12)
	c.PrintAtWithColor("RECTS:", 10, 32, console.WHITE)
	c.Rect(45, 30, 55, 40)
	c.Color(console.GREEN)
	c.RectFill(65, 30, 75, 40)
	c.RectFillWithColor(85, 25, 105, 45, console.RED)
	c.PrintAtWithColor("CIRCLE:", 10, 55, console.WHITE)
	c.Circle(50, 57, 5)
	c.Color(console.BLUE)
	c.CircleFill(70, 57, 5)
	c.CircleFillWithColor(95, 57, 10, console.BROWN)
	c.PrintAtWithColor("LINES:", 10, 77, console.WHITE)
	c.Color(console.LIGHT_GRAY)
	c.Line(45, 77, 105, 77)
	c.LineWithColor(45, 79, 105, 79, console.YELLOW)
	c.PrintAtWithColor("POINTS:", 10, 99, console.WHITE)
	c.PSet(50, 99)
	c.PSetWithColor(70, 99, console.PINK)
	// get color of point // earlier rect
	pointColor := c.PGet(85, 25)
	c.PSetWithColor(95, 99, pointColor)
}
