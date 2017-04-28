package console

import (
	"fmt"

	"github.com/telecoda/pico-go/config"
	"github.com/telecoda/pico-go/display"
	"github.com/telecoda/pico-go/events"
	"github.com/veandco/go-sdl2/sdl"
)

// This is the main control package for the virtual console
// the console package coordinates with all other console subsystems

type Console interface {
	Run() error
	Destroy()
}

type console struct {
	config.Config
	Display      display.Display
	EventHandler events.EventHandler
}

func New(cfg config.Config) (Console, error) {
	console := &console{
		Config: cfg,
	}

	// initialise display
	if console.Verbose {
		fmt.Printf("Initialising display\n")
	}

	display, err := display.New(cfg)
	if err != nil {
		return nil, err
	}

	console.Display = display

	// initialise event handler
	if console.Verbose {
		fmt.Printf("Initialising event handler\n")
	}

	handler, err := events.New()
	if err != nil {
		return nil, err
	}

	console.EventHandler = handler

	return console, nil
}

// Run is the main run loop
func (c *console) Run() error {
	// poll events
	for !c.EventHandler.HasQuit() {
		if err := c.EventHandler.PollEvents(); err != nil {
			return err
		}

		if err := c.Display.Render(); err != nil {
			return err
		}

		// wait for a bit
		sdl.Delay(16)
	}

	return nil
}

func (c *console) Destroy() {
	c.Display.Destroy()
}
