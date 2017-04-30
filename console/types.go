package console

import (
	"github.com/telecoda/pico-go/api"
	"github.com/telecoda/pico-go/display"
	"github.com/telecoda/pico-go/events"
)

// This is the main control package for the virtual console
// the console package coordinates with all other console subsystems

type Mode int

const (
	CLI Mode = iota
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
	SetMode(mode Mode)
	api.PicoGoAPI
}

type console struct {
	api.Config
	display.Display
	EventHandler events.EventHandler

	mode Mode

	cart Cartridge
}

type Cartridge interface {
	Init()
	Render()
	Update()
}

type cartridge struct {
}
