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
	c.ClsWithColor(console.WHITE)
	c.Sprite(1, 10, 10, 2, 2)
}
