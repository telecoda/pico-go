package console

import (
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

type Config struct {
	ConsoleWidth  uint
	ConsoleHeight uint
	WindowWidth   uint
	WindowHeight  uint
	FPS           uint8
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
	Cls()                 // Clear screen
	ClsColor(color Color) // Clear screen

	Flip() error // Copy graphics buffer to screen

	//Print(str string)                               // Print a string of characters to the screen
	//PrintAt(str string, x, y int)                   // Print a string of characters to the screen at position
	PrintColorAt(str string, x, y int, color Color) // Print a string of characters to the screen at position with color

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

type Console interface {
	LoadCart(path string) error
	Run() error
	Destroy()
	SetMode(newMode ModeType)
}

type console struct {
	Config

	currentMode   ModeType
	secondaryMode ModeType
	modes         map[ModeType]Mode
	hasQuit       bool

	cart Cartridge

	window   *sdl.Window
	renderer *sdl.Renderer

	palette
	font *ttf.Font
	logo *sdl.Surface
}

type Cartridge interface {
	Init()
	Render()
	Update()
}

type cartridge struct {
}

type Mode interface {
	Init() error
	Render() error
	Update() error
	HandleEvent(event sdl.Event) error
	PixelBuffer
}

type PixelBuffer interface {
	Destroy()
	PicoGoAPI
}

type mode struct {
	pixelBuffer
}

type pixelBuffer struct {
	cursor       pos
	pixelSurface *sdl.Surface // offscreen pixel buffer
	psRect       *sdl.Rect    // rect of pixelSurface
}

var title = "pico-go virtual games console"

type size struct {
	width  uint
	height uint
}

type pos struct {
	x uint
	y uint
}

type rgba struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type palette map[Color]rgba
