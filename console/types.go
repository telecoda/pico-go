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

type PicoGraphicsAPI interface {
	Clearer
	Drawer
	Paletter
	Printer
	Spriter
}

type PicoInputAPI interface {
	Btn(id int) bool
}

type Clearer interface {
	Cls()                       // Clear screen
	ClsWithColor(colorID Color) // Clear screen with color
}

type Drawer interface {
	Color(colorID Color) // Set drawing color (colour!!!)
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
	PaletteReset()
	PaletteCopy() Paletter
	GetColorID(rgba rgba) Color
	GetRGBA(color Color) (rgba, uint32)
	GetSDLColors() []sdl.Color
	MapColor(fromColor Color, toColor Color) error
	SetTransparent(color Color, enabled bool) error
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
	systemSprite(n, x, y, w, h, dw, dh int, rot float64, flipX, flipY bool)
}

type ModeType int

const (
	CLI ModeType = iota
	CODE_EDITOR
	SPRITE_EDITOR
	MAP_EDITOR
	SFX_EDITOR
	MUSIC_EDITOR
	RUNTIME
)

type ConsoleType string

const (
	PICO8       = "pico8"
	TIC80       = "tic80"
	ZX_SPECTRUM = "zxspectrum"
	CBM64       = "cbm64"
)

var ConsoleTypes = map[ConsoleType]string{
	PICO8:       "PICO8",
	TIC80:       "TIC80",
	ZX_SPECTRUM: "ZX_SPECTRUM",
	CBM64:       "CBM64",
}

const TOTAL_COLORS = 256

type Configger interface {
	GetConfig() Config
}

type Cartridge interface {
	// BaseCartridge methods already implemented
	Configger
	initPb(pb PixelBuffer)
	IsRunning() bool
	Stop()
	PicoInputAPI
	// User implemented methods below
	Init()
	Render()
	Update()
}

type Runtime interface {
	Mode
	PicoInputAPI
	LoadCart(cart Cartridge) error
}

type PixelBuffer interface {
	Flip() error // Copy graphics buffer to screen
	Destroy()
	GetFrame() *sdl.Surface
	PicoGraphicsAPI
	getPixelBuffer() *pixelBuffer
}

var title = "pico-go virtual games console"

type size struct {
	width  int
	height int
}
