package code

import (
	"fmt"

	"github.com/telecoda/pico-go/console"
)

// Code must implement console.Cartridge interface

type cartridge struct {
	cfg                 console.Config // holds details of console config
	console.PixelBuffer                // ref to console display

	counter int // just used in demo code
	x       int
	y       int
	speedX  int
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
	c.ClsWithColor(console.BLUE)

	c.counter = 0
	c.x = 50
	c.y = 50
	c.speedX = 2
}

// Update -  called once every frame
func (c *cartridge) Update() {
	c.x += c.speedX
	if c.x < 0 {
		c.x = 0
		c.speedX = -c.speedX
	}

	if c.x > c.cfg.ConsoleWidth {
		c.x = c.cfg.ConsoleWidth
		c.speedX = -c.speedX
	}

	c.counter++
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.Cls()
	c.PrintAtWithColor(fmt.Sprintf("c:%d", c.counter), c.x, c.y, console.RED)
}
