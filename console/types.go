package console

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Color int

/*
	This package tries to replicate the pico8 API as closely as possible
	During development I will be trying to implement more an more of the API
	To achieve feature parity with pico8
	Documented extensively here http://pico-8.wikia.com/wiki/Category:API
*/

type PicoGoAPI interface {
	Flip() error // Copy graphics buffer to screen
	Clearer
	Drawer
	Paletter
	Printer
	Spriter
}

type Clearer interface {
	Cls()                       // Clear screen
	ClsWithColor(colorID Color) // Clear screen with color
}

type Drawer interface {
	// drawing primatives
	Circle(x, y, r int)
	CircleWithColor(x, y, r int, colorID Color)
	CircleFill(x, y, r int)
	CircleFillWithColor(x, y, r int, colorID Color)
	Line(x0, y0, x1, y1 int)
	LineWithColor(x0, y0, x1, y1 int, colorID Color)
	PGet(x, y int) Color
	PSet(x, y int)
	PSetWithColor(x, y int, colorID Color)
	Rect(x0, y0, x1, y1 int)
	RectWithColor(x0, y0, x1, y1 int, colorID Color)
	RectFill(x0, y0, x1, y1 int)
	RectFillWithColor(x0, y0, x1, y1 int, colorID Color)
}

type Paletter interface {
	Color(colorID Color) // Set drawing color (colour!!!)
}

type Printer interface {
	// Text/Printing
	Cursor(x, y int) // Set text cursor
	GetCursor() pos
	Print(str string)                                     // Print a string of characters to the screen at default pos
	PrintAt(str string, x, y int)                         // Print a string of characters to the screen at position
	PrintAtWithColor(str string, x, y int, colorID Color) // Print a string of characters to the screen at position with color
	ScrollUpLine()
}

type Spriter interface {
	// TODO lots of params!! needs a bit of overloading love
	Sprite(n, x, y, w, h, dw, dh int, rot float64, flipX, flipY bool)
}

type ModeType int

const (
	CLI ModeType = iota
	CODE_EDITOR
	RUNTIME
	MAP_EDITOR
	SFX_EDITOR
	MUSIC_EDITOR
)

const TOTAL_COLORS = 256

type Configger interface {
	GetConfig() Config
}

type Cartridge interface {
	Configger
	Init(pb PixelBuffer)
	Render()
	Update()
}

type Mode interface {
	Init() error
	Render() error
	Update() error
	HandleEvent(event sdl.Event) error
	PixelBuffer
}

type Runtime interface {
	Mode
	LoadCart(cart Cartridge) error
}

type PixelBuffer interface {
	Destroy()
	GetFrame() *sdl.Surface
	PicoGoAPI
}

var title = "pico-go virtual games console"

type size struct {
	width  int
	height int
}
