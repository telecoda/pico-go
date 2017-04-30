package console

import (
	"github.com/telecoda/pico-go/config"
	"github.com/telecoda/pico-go/display"
	"github.com/telecoda/pico-go/events"
)

// This is the main control package for the virtual console
// the console package coordinates with all other console subsystems

type Console interface {
	LoadCart(path string) error
	Run() error
	Destroy()
}

type console struct {
	config.Config
	Display      display.Display
	EventHandler events.EventHandler

	cart Cartridge
}

type Cartridge interface {
	Init()
	Render()
	Update()
}

type cartridge struct {
}
