package console

import (
	"github.com/telecoda/pico-go/api"
	"github.com/telecoda/pico-go/display"
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
	api.PicoGoAPI
}

type console struct {
	api.Config
	display.Display

	currentMode ModeType
	modes       map[ModeType]Mode
	hasQuit     bool

	cart Cartridge
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
}

type mode struct{}
