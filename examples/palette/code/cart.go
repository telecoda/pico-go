package code

import (
	"math/rand"

	"github.com/telecoda/pico-go/console"
)

// Code must implement console.Cartridge interface

type cartridge struct {
	*console.BaseCartridge

	// example vars below
	mapAnim      bool
	frameCount   int
	totalFrames  int
	currentColor int
}

// NewCart - initialise a struct implementing Cartridge interface
func NewCart() console.Cartridge {
	return &cartridge{
		BaseCartridge: console.NewBaseCart(console.Pico8Config()),
	}
}

// Init - called once when cart is initialised
func (c *cartridge) Init() {
	// the Init method receives a PixelBuffer reference
	// hold onto this reference, this is the display that
	// your code will be drawing onto each frame
	c.frameCount = 0
	c.totalFrames = 25
	c.currentColor = 0
	c.mapAnim = false
}

// Update -  called once every frame
func (c *cartridge) Update() {
	c.frameCount++
	if c.frameCount > c.totalFrames {
		// trigger update

		if c.mapAnim {
			c.MapColor(console.Color(c.currentColor), console.PICO8_RED)
			c.MapColor(console.Color(15), console.Color(rand.Intn(16)))
		} else {
			c.SetTransparent(console.Color(c.currentColor), true)
		}
		// reset counters
		c.frameCount = 0
		c.currentColor++
		if c.currentColor > 16 {
			// all colors have been swapped reset
			c.currentColor = 0
			c.PaletteReset()
			c.mapAnim = !c.mapAnim
		}
	}
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls()
	c.RectFillWithColor(0, 0, 32, 32, 0)
	c.RectFillWithColor(32, 0, 64, 32, 1)
	c.RectFillWithColor(64, 0, 96, 32, 2)
	c.RectFillWithColor(96, 0, 128, 32, 3)

	c.RectFillWithColor(0, 32, 32, 64, 4)
	c.RectFillWithColor(32, 32, 64, 64, 5)
	c.RectFillWithColor(64, 32, 96, 64, 6)
	c.RectFillWithColor(96, 32, 128, 64, 7)

	c.RectFillWithColor(0, 64, 32, 96, 8)
	c.RectFillWithColor(32, 64, 64, 96, 9)
	c.RectFillWithColor(64, 64, 96, 96, 10)
	c.RectFillWithColor(96, 64, 128, 96, 11)

	c.RectFillWithColor(0, 96, 32, 128, 12)
	c.RectFillWithColor(32, 96, 64, 128, 13)
	c.RectFillWithColor(64, 96, 96, 128, 14)
	c.RectFillWithColor(96, 96, 128, 128, 15)

	c.PrintAtWithColor("PALETTE:", 46, 5, 15)
	c.Line(0, 12, 128, 12)

	if c.mapAnim {
		c.PrintAtWithColor("COLORS CAN BE SWAPPED.", 20, 20, 15)
	} else {
		c.PrintAtWithColor("COLORS CAN BE TRANSPARENT.", 12, 20, 15)
	}
}
