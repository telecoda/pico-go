package code

import (
	"fmt"

	"github.com/telecoda/pico-go/console"
)

// Code must implement console.Cartridge interface

type cartridge struct {
	*console.BaseCartridge

	counter int // just used in demo code
	x       int
	y       int
	speedY  int
}

// NewCart - initialise a struct implementing Cartridge interface
func NewCart() console.Cartridge {
	return &cartridge{
		BaseCartridge: console.NewBaseCart(),
	}
}

// Init - called once when cart is initialised
func (c *cartridge) Init() {
	c.counter = 0
	c.x = 40
	c.y = 0
	c.speedY = 2
}

// Update -  called once every frame
func (c *cartridge) Update() {
	c.y += c.speedY
	if c.y < 0 {
		c.y = 0
		c.speedY = -c.speedY
	}

	if c.y > c.GetConfig().ConsoleHeight {
		c.y = c.GetConfig().ConsoleHeight
		c.speedY = -c.speedY
	}

	c.counter++
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.ClsWithColor(c.GetConfig().BgColor)
	c.PrintAtWithColor(fmt.Sprintf("counter:%d", c.counter), c.x, c.y, c.GetConfig().FgColor)
}
