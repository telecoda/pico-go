package code

import (
	"github.com/telecoda/pico-go/console"
)

// This example demonstrates the drawing primatives

type cartridge struct {
	*console.BaseCartridge
}

// NewCart - initialise a struct implementing Cartridge interface
func NewCart() console.Cartridge {
	return &cartridge{
		BaseCartridge: console.NewBaseCart(),
	}
}

// Init - called once when cart is initialised
func (c *cartridge) Init() {
}

// Update -  called once every frame
func (c *cartridge) Update() {
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.ClsWithColor(console.PICO8_BLACK)
	c.PrintAtWithColor("DRAWING:", 50, 5, console.PICO8_WHITE)
	c.LineWithColor(0, 0, 128, 128, console.PICO8_WHITE)

	c.LineWithColor(0, 10, 128, 10, console.PICO8_WHITE)

	c.LineWithColor(10, 0, 10, 128, console.PICO8_WHITE)

	c.LineWithColor(30, 30, 60, 50, console.PICO8_WHITE)

	// c.LineWithColor(0, 12, 128, 12, console.PICO8_WHITE)
	// c.PrintAtWithColor("RECTS:", 10, 32, console.PICO8_WHITE)
	// c.Rect(45, 30, 55, 40)
	// c.Color(console.PICO8_GREEN)
	// c.RectFill(65, 30, 75, 40)
	// c.RectFillWithColor(85, 25, 105, 45, console.PICO8_RED)
	// c.PrintAtWithColor("CIRCLE:", 10, 55, console.PICO8_WHITE)
	// c.Circle(50, 57, 5)
	// c.Color(console.PICO8_BLUE)
	// c.CircleFill(70, 57, 5)
	// c.CircleFillWithColor(95, 57, 10, console.PICO8_BROWN)
	// c.PrintAtWithColor("LINES:", 10, 77, console.PICO8_WHITE)
	// c.Color(console.PICO8_LIGHT_GRAY)
	// c.Line(45, 77, 105, 77)
	// c.LineWithColor(45, 79, 105, 79, console.PICO8_YELLOW)
	// c.PrintAtWithColor("POINTS:", 10, 99, console.PICO8_WHITE)
	// c.PSet(50, 99)
	// c.PSetWithColor(70, 99, console.PICO8_PINK)
	// // get color of point // earlier rect
	// pointColor := c.PGet(85, 25)
	// c.PSetWithColor(95, 99, pointColor)
}
