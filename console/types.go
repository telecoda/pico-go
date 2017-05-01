package console

import (
	"github.com/telecoda/pico-go/api"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// This is the main control package for the virtual console
// the console package coordinates with all other console subsystems

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
	api.Config

	currentMode ModeType
	modes       map[ModeType]Mode
	hasQuit     bool

	cart Cartridge

	window   *sdl.Window
	renderer *sdl.Renderer

	palette
	font *ttf.Font
}

type Cartridge interface {
	Init()
	Render()
	Update()
}

type cartridge struct {
}

type Mode interface {
	Render() error
	Update() error
	PollEvents() error
	PixelBuffer
}

type PixelBuffer interface {
	Destroy()
	api.PicoGoAPI
	// Cls()
	// ClsColor(color api.Color)

	// Flip() error

	// PrintColorAt(str string, x, y int, color api.Color)
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

type palette map[int]rgba
