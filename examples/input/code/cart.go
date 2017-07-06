package code

import (
	"github.com/telecoda/pico-go/console"
)

// Code must implement console.Cartridge interface

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
	c.ClsWithColor(console.BLUE)
}

// Update -  called once every frame
func (c *cartridge) Update() {
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls()
	c.PrintAtWithColor("INPUT:", 50, 5, console.WHITE)
	c.Line(0, 12, 128, 12)

	c.PrintAtWithColor("UP:", 20, 30, console.WHITE)
	if c.Btn(console.P1_BUTT_UP) {
		c.PrintAtWithColor("PRESSED", 48, 30, console.WHITE)
	}

	c.PrintAtWithColor("DOWN:", 20, 50, console.WHITE)
	if c.Btn(console.P1_BUTT_DOWN) {
		c.PrintAtWithColor("PRESSED", 48, 50, console.WHITE)
	}

	c.PrintAtWithColor("LEFT:", 20, 70, console.WHITE)
	if c.Btn(console.P1_BUTT_LEFT) {
		c.PrintAtWithColor("PRESSED", 48, 70, console.WHITE)
	}

	c.PrintAtWithColor("RIGHT:", 20, 90, console.WHITE)
	if c.Btn(console.P1_BUTT_RIGHT) {
		c.PrintAtWithColor("PRESSED", 48, 90, console.WHITE)
	}
}
