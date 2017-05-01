package console

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type cli struct {
	console *console
	PixelBuffer
}

var event sdl.Event
var joysticks [16]*sdl.Joystick

func newCLIMode(c *console) Mode {
	cli := &cli{
		console: c,
	}
	pb, _ := newPixelBuffer(c.Config)

	cli.PixelBuffer = pb
	return cli
}

func (c *cli) PollEvents() error {
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			c.console.hasQuit = true
		case *sdl.KeyDownEvent:
			switch t.Keysym.Sym {
			case sdl.K_RIGHT:
				fmt.Printf("Switching to code editor\n")
				c.console.SetMode(CODE_EDITOR)
			}
		default:
			fmt.Printf("Some event: %#v \n", event)
		}
	}

	return nil

}

func (c *cli) Update() error {
	return nil
}

func (c *cli) Render() error {
	c.PixelBuffer.ClsColor(8)

	c.PixelBuffer.PrintColorAt("CLI Print Test", 10, 10, 7)
	return nil
}
