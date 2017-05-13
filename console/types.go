package console

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Config struct {
	ConsoleWidth  int
	ConsoleHeight int
	WindowWidth   int
	WindowHeight  int
	FPS           int
	Verbose       bool
}

type Color int

/*
	This package tries to replicate the pico8 API as closely as possible
	During development I will be trying to implement more an more of the API
	To achieve feature parity with pico8
	Documented extensively here http://pico-8.wikia.com/wiki/Category:API
*/

type PicoGoAPI interface {
	Cls()                       // Clear screen
	ClsWithColor(colorID Color) // Clear screen

	Color(colorID Color) // Set drawing color (colour!!!)
	Cursor(x, y int)     // Set text cursor

	Flip() error // Copy graphics buffer to screen

	GetCursor() pos
	Print(str string) // Print a string of characters to the screen at default pos
	//PrintAt(str string, x, y int)                   // Print a string of characters to the screen at position
	PrintAtWithColor(str string, x, y int, colorID Color) // Print a string of characters to the screen at position with color
	//RectFill(x0, y0, x1, y1 int)
	RectFillWithColor(x0, y0, x1, y1 int, colorID Color)
	ScrollUpLine()
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
	PicoGoAPI
}

var title = "pico-go virtual games console"

type size struct {
	width  int
	height int
}

type rgba struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type palette map[Color]rgba
