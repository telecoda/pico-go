package console

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type cli struct {
	console *console
	PixelBuffer
	*cursor
}

var event sdl.Event
var joysticks [16]*sdl.Joystick

func newCLIMode(c *console) Mode {
	cli := &cli{

		console: c,
	}
	pb, _ := newPixelBuffer(c.Config)
	cursor := newCursor(pb, RED)

	cli.PixelBuffer = pb
	cli.cursor = cursor
	cursor.x = 1
	cursor.y = 8
	return cli
}

func (c *cli) HandleEvent(event sdl.Event) error {
	switch t := event.(type) {
	case *sdl.KeyDownEvent:
		switch t.Keysym.Sym {
		case sdl.K_RIGHT:
			fmt.Printf("Switching to code editor\n")
			c.console.SetMode(CODE_EDITOR)
		}
	default:
		fmt.Printf("Some event: %#v \n", event)
	}

	return nil
}

func (c *cli) Init() error {
	// get native pixel buffer
	c.PixelBuffer.ClsColor(BLACK)
	pb := c.PixelBuffer.(*pixelBuffer)

	logoRect := &sdl.Rect{X: 0, Y: 0, W: _logoWidth, H: _logoHeight}
	screenRect := &sdl.Rect{X: 0, Y: 0, W: _logoWidth, H: _logoHeight}
	_console.logo.Blit(logoRect, pb.pixelSurface, screenRect)

	title := fmt.Sprintf("PICO-GO %s", _version)
	c.PixelBuffer.PrintColorAt(title, 0, 24, LIGHT_GRAY)
	c.PixelBuffer.PrintColorAt("(C) 2017 @TELECODA", 0, 32, LIGHT_GRAY)
	c.PixelBuffer.PrintColorAt("TYPE HELP FOR HELP", 0, 48, LIGHT_GRAY)

	c.PixelBuffer.PrintColorAt(">", 0, 64, WHITE)

	return nil
}

func (c *cli) Update() error {
	return nil
}

func (c *cli) Render() error {
	c.cursor.render()
	return nil
}
