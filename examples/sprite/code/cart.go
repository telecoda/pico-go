package code

import (
	"fmt"

	"github.com/telecoda/pico-go/console"
)

// This example demonstrates the drawing primatives

type cartridge struct {
	cfg                 console.Config // holds details of console config
	console.PixelBuffer                // ref to console display

	// example vars below
	rot float64
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
}

// Update -  called once every frame
func (c *cartridge) Update() {
	c.rot += -4
	if c.rot > 360 {
		c.rot = 0
	}
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.ClsWithColor(console.WHITE)
	c.PrintAtWithColor("SPRITES:", 50, 5, console.BLACK)
	c.Line(0, 12, 128, 12)
	c.PrintAtWithColor("FLIPX: false", 10, 18, console.BLACK)
	c.PrintAtWithColor("FLIPY: false", 10, 26, console.BLACK)
	c.PrintAtWithColor(fmt.Sprintf("R: %d", int(c.rot)), 100, 22, console.BLACK)
	c.Sprite(2, 70, 16, 2, 2, 16, 16, c.rot, false, false)
	c.PrintAtWithColor("FLIPX: true", 10, 38, console.BLACK)
	c.PrintAtWithColor("FLIPY: false", 10, 46, console.BLACK)
	c.Sprite(0, 70, 36, 2, 2, 16, 16, 0, true, false)
	c.PrintAtWithColor("FLIPX: false", 10, 58, console.BLACK)
	c.PrintAtWithColor("FLIPY: true", 10, 66, console.BLACK)
	c.Sprite(0, 70, 56, 2, 2, 16, 16, 0, false, true)
	c.PrintAtWithColor("FLIPX: true", 10, 78, console.BLACK)
	c.PrintAtWithColor("FLIPY: true", 10, 86, console.BLACK)
	c.Sprite(0, 70, 76, 2, 2, 16, 16, 0, true, true)
	c.PrintAtWithColor("Scaled:", 10, 98, console.BLACK)
	c.Sprite(40, 40, 96, 4, 2, 64, 32, 0, false, false)
}
