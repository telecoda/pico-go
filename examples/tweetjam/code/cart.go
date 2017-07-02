package code

import (
	"math/rand"

	"github.com/telecoda/pico-go/console"
)

// Code must implement console.Cartridge interface

type cartridge struct {
	cfg                 console.Config // holds details of console config
	console.PixelBuffer                // ref to console display
	running             bool
}

// NewCart - initialise a struct implementing Cartridge interface
func NewCart() console.Cartridge {
	cart := &cartridge{
		cfg: console.DefaultConfig(),
	}

	cart.cfg.FPS = 30

	return cart
}

// GetConfig - return config need for Cart to run
func (c *cartridge) GetConfig() console.Config {
	return c.cfg
}

/* This is the original tweetcart code
s={}w=128 r=rnd for i=1,w do s[i]={}p=s[i]p[1]=r(w)end::a::cls()for i=1,w do p=s[i]pset(p[1],i,i%3+5)p[1]=(p[1]-i%3)%w end flip()goto a
*/

// Init - called once when cart is initialised
func (c *cartridge) Init(pb console.PixelBuffer) {
	// the Init method receives a PixelBuffer reference
	// hold onto this reference, this is the display that
	// your code will be drawing onto each frame
	c.PixelBuffer = pb
	c.running = true
}

// Update -  called once every frame
func (c *cartridge) Update() {
}

// Render - called once every frame
func (c *cartridge) Render() {
	c.ClsWithColor(console.BLACK)

	// init stars
	/*
		s={}
		w=128
		r=rnd
		for i=1,w do
			s[i]={}
			p=s[i]
			p[1]=r(w)
		end
	*/
	w := 128
	s := make([]int, w, w)

	for i := 0; i < w; i++ {
		s[i] = rand.Intn(w)
	}

	/*
			cls()
		for i=1,w do
			p=s[i]
			pset(p[1],i,i%3+5)
			p[1]=(p[1]-i%3)%w
		end
	*/
	for c.IsRunning() {
		c.Cls()
		for i := 0; i < w; i++ {
			c.PSetWithColor(s[i], i, console.Color(i%3+5))
			s[i] = (s[i] - (i % 3)) % w
			if s[i] < 0 {
				s[i] += w
			}
		}
		c.Flip()
	}

}

func (c *cartridge) IsRunning() bool {
	return c.running
}

func (c *cartridge) Stop() {
	c.running = false
}
